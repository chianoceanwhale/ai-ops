package crondtask

import (
	"fmt"
	"strconv"

	"ai-ops/configs/agent"
	"os/exec"
	"runtime"
	"sync"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	ServiceTask Task
)

var (
	// 定时任务调度管理器
	serviceCron *cron.Cron

	// 同一任务是否有实例处于运行中
	runInstance Instance

	// 任务计数-正在运行的任务
	taskCount TaskCount

	// 并发队列, 限制同时运行的任务数量
	concurrencyQueue ConcurrencyQueue
)

// 并发队列
type ConcurrencyQueue struct {
	queue chan struct{}
}

func (cq *ConcurrencyQueue) Add() {
	cq.queue <- struct{}{}
}

func (cq *ConcurrencyQueue) Done() {
	<-cq.queue
}

// 任务计数
type TaskCount struct {
	wg   sync.WaitGroup
	exit chan struct{}
}

func (tc *TaskCount) Add() {
	tc.wg.Add(1)
}

func (tc *TaskCount) Done() {
	tc.wg.Done()
}

func (tc *TaskCount) Exit() {
	tc.wg.Done()
	<-tc.exit
}

func (tc *TaskCount) Wait() {
	tc.Add()
	tc.wg.Wait()
	close(tc.exit)
}

// 任务ID作为Key
type Instance struct {
	m sync.Map
}

// 是否有任务处于运行中
func (i *Instance) has(key int) bool {
	_, ok := i.m.Load(key)

	return ok
}

func (i *Instance) add(key int) {
	i.m.Store(key, struct{}{})
}

func (i *Instance) done(key int) {
	i.m.Delete(key)
}

type Task struct {
	m sync.Map
}

func (t *Task) get(key string) cron.EntryID {
	v, ok := t.m.Load(key)
	if !ok {
		return 0
	}

	return v.(cron.EntryID)
}

func (t *Task) add(key string, entry cron.EntryID) {
	t.m.Store(key, entry)
}

func (t *Task) done(key string) {
	t.m.Delete(key)
}

type TaskResult struct {
	Result string
	Err    error
}

//invoke task run
func (task Task) Run(taskModel OpsCronDetail) {
	go createJob(taskModel)()
}

//handle interface
type Handler interface {
	Run(taskModel OpsCronDetail) (string, error)
}

// SHELL调用执行任务
type ShellHandler struct{}

func (h *ShellHandler) Run(taskModel OpsCronDetail) (result string, err error) {
	cmd := exec.Command("sh", "-c", taskModel.CronContent)
	out, err := cmd.Output()
	return string(out), err
}

//init crond
func (task Task) Initialize() {
	serviceCron = cron.New(cron.WithSeconds())
	serviceCron.Start()
	concurrencyQueue = ConcurrencyQueue{queue: make(chan struct{}, 1)}
	taskCount = TaskCount{sync.WaitGroup{}, make(chan struct{})}
	go taskCount.Wait()
	log.Info("开始初始化定时任务")
	taskNum := 0

	taskList := ActiveList()
	//遍历所有激活任务并添加到定时任务中
	for _, item := range taskList {
		task.Add(item)
		taskNum++
	}
	log.Infof("定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
}

// PanicToError Panic转换为error
func PanicToError(f func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf(PanicTrace(e))
		}
	}()
	f()
	return
}

// PanicTrace panic调用链跟踪
func PanicTrace(err interface{}) string {
	stackBuf := make([]byte, 4096)
	n := runtime.Stack(stackBuf, false)

	return fmt.Sprintf("panic: %v %s", err, stackBuf[:n])
}

// 删除任务
func (task Task) Remove(bindFlag string) {
	entryId := task.get(bindFlag)
	serviceCron.Remove(entryId)
}

// 添加任务
func (task Task) Add(taskModel OpsCronDetail) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Task Add panic, task: %v", taskModel)
		}
	}()
	taskFunc := createJob(taskModel)
	if taskFunc == nil {
		log.Error("创建任务处理Job失败,不支持的任务协议#", taskModel)
		return
	}

	log.Info("-----taskModel.CronTime-----------", taskModel.CronTime)

	entryId, err := serviceCron.AddFunc(taskModel.CronTime, taskFunc)
	if err != nil {
		log.Error("添加任务到调度器失败#", err)
	}
	task.add(taskModel.BindFlag, entryId)
}

func createJob(taskModel OpsCronDetail) cron.FuncJob {
	var handler Handler = nil
	handler = new(ShellHandler)
	taskFunc := func() {
		taskCount.Add()
		defer taskCount.Done()

		log.Debugf("任务命令-%s", taskModel)

		concurrencyQueue.Add()
		defer concurrencyQueue.Done()

		log.Infof("开始执行任务#%s#命令-%s", taskModel.CronName, taskModel.CronContent)
		taskResult := execJob(handler, taskModel)
		log.Infof("任务完成#%s#命令-%s", taskModel.CronName, taskModel.CronContent)
		log.Info("任务执行结果为：%s", taskResult)
	}

	return taskFunc
}

// 执行具体任务
func execJob(handler Handler, taskModel OpsCronDetail) TaskResult {
	defer func() {
		if err := recover(); err != nil {
			log.Error("panic#service/task.go:execJob#", err)
		}
	}()
	// 默认只运行任务一次
	var output string
	var err error

	output, err = handler.Run(taskModel)
	if err == nil {
		return TaskResult{Result: output, Err: err}
	}

	return TaskResult{Result: output, Err: err}
}

// 添加任务到定时器
func AddTaskToTimer(bindFlag string) {
	task, flag := GetDetailInfoByBindFlag(bindFlag)
	if flag {
		ServiceTask.RemoveAndAdd(task)
	}
}

//add task after delete task
func (task Task) RemoveAndAdd(taskModel OpsCronDetail) {
	task.Remove(taskModel.BindFlag)
	task.Add(taskModel)
}

//active job task
func ActiveList() []OpsCronDetail {
	var taskList []OpsCronDetail
	for _, v := range agent.TaskCrondData {
		tmpData := v.([]interface{})
		for _, j := range tmpData {
			innerTmpData := j.(map[string]interface{})
			switch innerTmpData["start_flag"].(float64) {
			case 1:
				log.Info("----------%v task loading----------", innerTmpData["cron_name"])
				taskList = append(taskList, OpsCronDetailInterfaceTransferStruct(innerTmpData))
			case 2:
				log.Info("----------%v Task Forbding----------", innerTmpData["cron_name"])
			default:
				log.Info("----------%v Load Task Failed,Can't match type----------", innerTmpData["cron_name"])
			}
		}
	}
	return taskList
}

//Get detail info by bind_flag param
func GetDetailInfoByBindFlag(bindFlag string) (OpsCronDetail, bool) {
	var taskList OpsCronDetail
	var Flag bool
	for _, v := range agent.TaskCrondData {
		tmpData := v.([]interface{})
		for _, j := range tmpData {
			innerTmpData := j.(map[string]interface{})
			if innerTmpData["start_flag"].(string) == bindFlag {
				taskList = OpsCronDetailInterfaceTransferStruct(innerTmpData)
				Flag = true
			}
		}
	}
	return taskList, Flag
}

//data reload and flush to memory
func ReloadDataFlushToMemory(changeData OpsCronDetail, operatorType int) map[string]interface{} {
	for k, _ := range agent.TaskCrondData {
		tmpData := agent.TaskCrondData[k].([]interface{})
		switch operatorType {
		case 1:
			//add data
			tmpData = append(tmpData, changeData)
			agent.TaskCrondData[k] = tmpData
			ServiceTask.Add(changeData)

		case 2:
			//update data
			var equelFlag int
			for i, j := range tmpData {
				innerTmpData := j.(map[string]interface{})
				if innerTmpData["start_flag"].(string) == changeData.BindFlag {
					equelFlag = i
					retTmpData := DeleteElementFromSlice(tmpData, equelFlag)
					retTmpData = append(retTmpData, changeData)
					agent.TaskCrondData[k] = retTmpData
					ServiceTask.RemoveAndAdd(changeData)
					break
				}
			}
		case 3:
			//delete data
			var equelFlag int
			for i, j := range tmpData {
				innerTmpData := j.(map[string]interface{})
				if innerTmpData["start_flag"].(string) == changeData.BindFlag {
					// ...
					equelFlag = i
					retTmpData := DeleteElementFromSlice(tmpData, equelFlag)
					agent.TaskCrondData[k] = retTmpData
					ServiceTask.Remove(changeData.BindFlag)
					break
				}
			}
		case 4:
			var equelFlag int
			for i, j := range tmpData {
				innerTmpData := j.(map[string]interface{})
				if innerTmpData["start_flag"].(string) == changeData.BindFlag {
					equelFlag = i
					retTmpData := DeleteElementFromSlice(tmpData, equelFlag)
					retTmpData = append(retTmpData, changeData)
					agent.TaskCrondData[k] = retTmpData
					ServiceTask.RemoveAndAdd(changeData)
					break
				}
			}
			if equelFlag == 0 {
				ReloadDataFlushToMemory(changeData, 1)
			}
		}
	}
	return agent.TaskCrondData
}

//delete element from slice
func DeleteElementFromSlice(orgSlice []interface{}, index int) []interface{} {
	return append(orgSlice[:index], orgSlice[index+1:]...)
}

//crond detail
type OpsCronDetail struct {
	BindFlag    string `json:"bind_flag" `   //关联环境标识
	CronName    string `json:"cron_name"`    //任务名
	CronTime    string `json:"cron_time" `   //定时时间
	AddUser     string `json:"add_user" `    //添加人
	CronContent string `json:"cron_content"` //定时任务内容
	StartFlag   int    `json:"start_flag"`   //'任务是否启用[1-启用;2-禁用],默认为1'
}

type OpsCrontJsonData struct {
	CrondData []OpsCronDetail `json:"crond_data"`
}

//OpsCronDetail interface transfer struct
func OpsCronDetailInterfaceTransferStruct(taskOrgObj interface{}) OpsCronDetail {
	var taskObj OpsCronDetail
	tmpData := taskOrgObj.(map[string]interface{})
	taskObj.BindFlag = tmpData["bind_flag"].(string)
	taskObj.CronName = tmpData["cron_name"].(string)
	taskObj.CronTime = tmpData["cron_time"].(string)
	taskObj.AddUser = tmpData["add_user"].(string)
	taskObj.CronContent = tmpData["cron_content"].(string)
	startFlag := fmt.Sprintf("%f", tmpData["start_flag"].(float64))
	taskObj.StartFlag, _ = strconv.Atoi(startFlag)
	return taskObj
}

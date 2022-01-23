package agentservice

import (
	"ai-ops/configs/agent"
	"ai-ops/configs/task"
	"ai-ops/internal/agent/protos"
	"ai-ops/tools/utils"
	"context"

	"github.com/gookit/goutil/jsonutil"
)

func (this *AgentService) CrontaskSync(ctx context.Context, req *protos.CronCommandRequest) (*protos.AgentResponse, error) {
	var data crondtask.OpsCronDetail
	// var data crond.OpsCronTask
	var err error
	data.AddUser = req.AddUser
	data.CronName = req.CronName
	data.CronTime = req.CronTime
	data.CronContent = req.CronContent
	data.BindFlag = req.BindFlag
	data.StartFlag = int(req.StartFlag)
	filePath := utils.GetFileAbsPath(agent.DataConfig.Dir, agent.DataConfig.FileName)
	dataMapInterface := crondtask.ReloadDataFlushToMemory(data, int(req.OperatorType))
	err = jsonutil.WriteFile(filePath, dataMapInterface)
	if err != nil {
		return &protos.AgentResponse{Data: string(req.CronName)}, err
	}
	return &protos.AgentResponse{Data: string(req.CronName)}, nil
}

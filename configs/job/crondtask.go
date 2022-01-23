package job

import (
	"ai-ops/internal/tools/utils"

	"github.com/gookit/goutil/jsonutil"
)

var TaskCrondData map[string]interface{}

//load json file
func LoadJSONFile() map[string]interface{} {
	filePath := utils.GetFileAbsPath(DataConfig.Dir, DataConfig.FileName)
	var taskCrond map[string]interface{}
	jsonutil.ReadFile(filePath, &taskCrond)
	return taskCrond
}

//init  task json file data
func InitTaskCrondData() {
	filePath := utils.GetFileAbsPath(DataConfig.Dir, DataConfig.FileName)
	var TaskCrondData map[string]interface{}
	jsonutil.ReadFile(filePath, TaskCrondData)
}

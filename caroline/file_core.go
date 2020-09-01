package caroline

import (
	"FrontEndOJudger/models"
	"FrontEndOJudger/pkg/setting"
	"fmt"
	"io/ioutil"
	"os"
)

/**
写入本地磁盘 ./test_chamber/{creator}/{submitid}
*/
func WriteSubmitToFile(labSubmit *models.LabSubmit) (string, string) {
	testChamberDirName := fmt.Sprintf("%s/%s/%s/%d/", setting.JudgerSetting.TestChamberBaseDir, setting.JudgerSetting.TestChamberDir, labSubmit.Creator, labSubmit.ID)
	testChamberFileName := fmt.Sprintf("%sindex.html", testChamberDirName)
	testChamberUrlName := fmt.Sprintf("%s/%d", labSubmit.Creator, labSubmit.ID)
	// 检查是否存在
	_, err := os.Stat(testChamberFileName)
	if err == nil || os.IsExist(err) {
		return testChamberFileName, ""
	}
	os.MkdirAll(testChamberDirName, 0777)
	ioutil.WriteFile(testChamberFileName, []byte(labSubmit.SubmitData), 0777)
	return testChamberFileName, testChamberUrlName
}

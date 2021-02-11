package caroline

import (
	"FrontEndOJudger/models"
	"FrontEndOJudger/pkg/setting"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

/**
写入本地磁盘 ./test_chamber/{submitid}
*/
func WriteSubmitToFile(labSubmit *models.LabSubmit) (string, string, error) {
	testChamberDirName := fmt.Sprintf("%s/%d/", setting.JudgerSetting.TestChamberBaseDir, labSubmit.ID)
	testChamberFileName := fmt.Sprintf("%sindex.html", testChamberDirName)
	testChamberUrlName := fmt.Sprintf("%d", labSubmit.ID)
	// 检查是否存在
	_, err := os.Stat(testChamberFileName)
	if err == nil || os.IsExist(err) {
		return testChamberFileName, "", nil
	}
	err = os.MkdirAll(testChamberDirName, 0777)
	if err != nil {
		log.Printf("Can not mkdir [%s] error [%s]", testChamberDirName, err)
	}
	return testChamberFileName, testChamberUrlName, ioutil.WriteFile(testChamberFileName, []byte(labSubmit.SubmitData), 0777)
}

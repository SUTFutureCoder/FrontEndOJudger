package caroline

import (
	"FrontEndOJudger/models"
	"FrontEndOJudger/pkg/setting"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func Judge() {
	// Step1 determine scan db or listen mq
	ch := make(chan *models.LabSubmit, setting.JudgerSetting.JudgerSum)

	// scan db
	labSubmits, err := models.GetSubmitByStatus(models.LABSUBMITSTATUS_PENDING, setting.JudgerSetting.JudgerSum)
	if err != nil {
		log.Printf("Get LabSubmits Error:%v", err)
		return
	}

	for i := 0; i < setting.JudgerSetting.JudgerSum; i++ {
		go JudgeQueue(ch)
	}

	for _, labSubmit := range labSubmits {
		ch <- labSubmit
	}

}

func JudgeQueue(ch chan *models.LabSubmit) {
	v := <-ch
	log.Printf("Start judge submitId[%d]", v.ID)
	submitRet, err := JudgeSubmit(v.ID)
	if submitRet != nil {
		log.Printf("End judge submitId[%d] with status[%d] err[%#v]", submitRet.ID, submitRet.Status, err)
		return
	}
	log.Printf("[ERROR] End judge submitId[%d] with empty result err[%#v]", v.ID, err)
}

func JudgeSubmit(submitId uint64) (*models.LabSubmit, error) {

	// 获取lab_id
	labSubmit, err := models.GetSubmitById(submitId)
	if err != nil {
		log.Printf("")
		return labSubmit, err
	}

	if labSubmit == nil {
		return nil, errors.New("labSubmit is nil")
	}

	// 如果非法实验室、脏数据，直接判定失败
	if labSubmit.LabID == 0 {
		log.Printf("find the lab id of submit is 0, maybe dirty data, submitId[%d]", submitId)
		return labSubmit, updateSubmitStatus(submitId, models.LABSUBMITSTATUS_PENDING, models.LABSUBMITSTATUS_SYSTEM_ERROR, labSubmit)
	}

	// 获取case信息
	testcaseIds, err := models.GetLabTestcaseMapByLabId(labSubmit.LabID)
	if len(testcaseIds) == 0 {
		// 如果没有testcase直接判定AC
		log.Printf("find the length of testcase Ids is empty, labId[%d], submitId[%d]", labSubmit.LabID, submitId)
		return labSubmit, updateSubmitStatus(submitId, models.LABSUBMITSTATUS_PENDING, models.LABSUBMITSTATUS_NO_TESTCASE, labSubmit)
	}

	// 获取testcase详情
	testcases, err := models.GetTestcaseByIds(testcaseIds)

	// 变更状态 乐观锁
	rows, err := models.UpdateSubmitStatusResult(submitId, models.LABSUBMITSTATUS_PENDING, models.LABSUBMITSTATUS_JUDING, "")
	if err != nil {
		log.Printf("change submit status error submitId:%d fromStatus:%d toStatus:%d rows:%d err:%v", submitId, models.LABSUBMITSTATUS_PENDING, models.LABSUBMITSTATUS_JUDING, rows, err)
		return labSubmit, errors.New(fmt.Sprintf("change submit status error submitId:%d fromStatus:%d toStatus:%d rows:%d err:%v", submitId, models.LABSUBMITSTATUS_PENDING, models.LABSUBMITSTATUS_JUDING, rows, err))
	}

	// 执行测试用例
	_, testChamberUrlName, err := WriteSubmitToFile(labSubmit)
	if err != nil {
		return labSubmit, err
	}

	// 记录执行时间

	// 实际执行
	//testResults := ExecCaroline("file://"+testChamberFileName, testcases, submitId)
	testResults := ExecCaroline(fmt.Sprintf("%s:%s/%s", setting.JudgerSetting.TestChamberAddr, setting.JudgerSetting.TestChamberPort, testChamberUrlName), testcases, submitId)

	// 获取测试结果 更新结果
	labSubmit.Status = models.LABSUBMITSTATUS_ACCEPTED
	for _, v := range testResults {
		// 如果有错误按照最后一个为准
		if v.Status != models.LABSUBMITSTATUS_ACCEPTED {
			labSubmit.Status = v.Status
		}
	}
	jsonByte, err := json.Marshal(testResults)
	labSubmit.SubmitResult = string(jsonByte)

	// 更新结果
	return labSubmit, updateSubmitStatus(submitId, models.LABSUBMITSTATUS_JUDING, labSubmit.Status, labSubmit)
}

func updateSubmitStatus(submitId uint64, fromStatus, toStatus int, labSubmit *models.LabSubmit) error {
	_, err := models.UpdateSubmitStatusResult(submitId, fromStatus, toStatus, labSubmit.SubmitResult)
	if err != nil {
		log.Printf("change submit status and result error submitId:%d fromStatus:%d toStatus:%d submitresult:%v err:%v", submitId, fromStatus, labSubmit.Status, labSubmit.SubmitResult, err)
		return errors.New(fmt.Sprintf("change submit status error submitId:%d fromStatus:%d toStatus:%d submitresult:%v err:%v", submitId, fromStatus, labSubmit.Status, labSubmit.SubmitResult, err))
	}
	log.Println(labSubmit.SubmitResult)
	return nil
}

package caroline

import (
	"FrontEndOJudger/models"
	"FrontEndOJudger/pkg/setting"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
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
	fmt.Println(time.Now().UnixNano() / 1e6)
	_ = JudgeSubmit(v.ID)
	fmt.Println(time.Now().UnixNano() / 1e6)
}

func JudgeSubmit(submitId uint64) error {

	// 获取lab_id
	labSubmit, err := models.GetSubmitById(submitId)
	if err != nil {
		log.Printf("")
		return err
	}

	if labSubmit == nil || labSubmit.LabID == 0 {
		return nil
	}

	// 获取case信息
	testcaseIds, err := models.GetLabTestcaseMapByLabId(labSubmit.LabID)
	if len(testcaseIds) == 0 {
		return nil
	}

	// 获取testcase详情
	testcases, err := models.GetTestcaseByIds(testcaseIds)

	// 变更状态 乐观锁
	rows, err := models.UpdateSubmitStatusResult(submitId, models.LABSUBMITSTATUS_PENDING, models.LABSUBMITSTATUS_JUDING, "")
	if rows != 1 || err != nil {
		log.Printf("change submit status error submitId:%d fromStatus:%d toStatus:%d rows:%d err:%v", submitId, models.LABSUBMITSTATUS_PENDING, models.LABSUBMITSTATUS_JUDING, rows, err)
		return errors.New(fmt.Sprintf("change submit status error submitId:%d fromStatus:%d toStatus:%d rows:%d err:%v", submitId, models.LABSUBMITSTATUS_PENDING, models.LABSUBMITSTATUS_JUDING, rows, err))
	}

	// 执行测试用例
	testChamberFileName := WriteSubmitToFile(labSubmit)
	testResults := ExecCaroline("file://"+testChamberFileName, testcases, submitId)

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
	rows, err = models.UpdateSubmitStatusResult(submitId, models.LABSUBMITSTATUS_JUDING, labSubmit.Status, labSubmit.SubmitResult)
	if rows != 1 || err != nil {
		log.Printf("change submit status and result error submitId:%d fromStatus:%d toStatus:%d submitresult:%v err:%v", submitId, models.LABSUBMITSTATUS_PENDING, models.LABSUBMITSTATUS_JUDING, labSubmit.SubmitResult, err)
		return errors.New(fmt.Sprintf("change submit status error submitId:%d fromStatus:%d toStatus:%d submitresult:%v err:%v", submitId, models.LABSUBMITSTATUS_PENDING, models.LABSUBMITSTATUS_JUDING, labSubmit.SubmitResult, err))
	}
	log.Println(labSubmit.SubmitResult)
	return nil
}

package caroline

import (
	"FrontEndOJudger/models"
	"FrontEndOJudger/pkg/setting"
	"fmt"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

func Judge() {
	wg.Add(1)
	// Step1 determine scan db or listen mq
	ch := make(chan *models.LabSubmit, setting.JudgerSetting.JudgerSum)

	// scan db
	labSubmits, err := models.GetSubmitByStatus(models.LABSUBMITSTATUS_PENDING, setting.JudgerSetting.JudgerSum * 2)
	if err != nil {
		log.Printf("Get LabSubmits Error:%v", err)
		return
	}

	for _, labSubmit := range labSubmits {
		ch <- labSubmit
		ch <- labSubmit
		ch <- labSubmit
		ch <- labSubmit
	}

	for i := 0; i < setting.JudgerSetting.JudgerSum; i++ {
		go JudgeQueue(ch)
	}
	wg.Wait()
}

func JudgeQueue(ch chan *models.LabSubmit) {
	v := <- ch
	fmt.Println((time.Now().UnixNano() / 1e6))
	JudgeSubmit(v.ID)
	fmt.Println((time.Now().UnixNano() / 1e6))
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
	//rows, err := models.UpdateSubmitStatus(submitId, models.LABSUBMITSTATUS_PENDING, models.LABSUBMITSTATUS_JUDING)
	//if rows != 1 || err != nil {
	//	log.Printf("change submit status error submitId:%d fromStatus:%d toStatus:%d err:%v", submitId, models.LABSUBMITSTATUS_PENDING, models.LABSUBMITSTATUS_JUDING, err)
	//	return errors.New(fmt.Sprintf("change submit status error submitId:%d fromStatus:%d toStatus:%d err:%v", submitId, models.LABSUBMITSTATUS_PENDING, models.LABSUBMITSTATUS_JUDING, err))
	//}

	// 执行测试用例
	testChamberFileName := WriteSubmitToFile(labSubmit)
	testResults := ExecCaroline("file://"+testChamberFileName, testcases, submitId)

	// 获取测试结果
	_ = testResults
	return nil
}

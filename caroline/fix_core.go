package caroline

import (
	"FrontEndOJudger/models"
	"FrontEndOJudger/pkg/setting"
	"log"
)

func FixExpiredJudgingSubmits() {
	// STEP1 get expired submits
	labSubmit := &models.LabSubmit{}
	labSubmits, err := labSubmit.GetExpiredJudgingSubmits(setting.JudgerSetting.JudgerSum)

	if err != nil {
		log.Printf("[ERROR] failed get expired juding submits err:[%v]", err)
		return
	}

	// STEP2 change status to pending
	for _, labSubmit := range labSubmits {
		labSubmit.UpdateStatusResult(labSubmit.Status, models.LABSUBMITSTATUS_PENDING, labSubmit.SubmitResult)
	}

	// FIN wait for judge main process
}
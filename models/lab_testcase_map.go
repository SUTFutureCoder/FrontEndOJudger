package models

import (
	"database/sql"
	"strconv"
)

// LabTestcaseMap 实验室、测试用例关联表
type LabTestcaseMap struct {
	Model
	// LabID 实验室id
	LabID int64 `json:"lab_id"`
	// TestcaseID 测试用例id
	TestcaseID int64 `json:"testcase_id"`
}

func GetLabTestcaseMapByLabId(labId uint64) ([]interface{}, error) {
	var testcaseIds []interface{}
	stmt, err := DB.Prepare("SELECT testcase_id FROM lab_testcase_map WHERE lab_id = ? AND status = 1")
	rows, err := stmt.Query(
		&labId,
	)
	defer rows.Close()
	for rows.Next() {
		var testcaseId int
		rows.Scan(&testcaseId)
		testcaseIds = append(testcaseIds, strconv.Itoa(testcaseId))
	}
	return testcaseIds, err
}

func (labTestCaseMap *LabTestcaseMap) Insert(tx *sql.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare("INSERT INTO lab_testcase_map (lab_id, testcase_id, creator, create_time) VALUES (?,?,?,?)")
	defer stmt.Close()
	result, err := stmt.Exec(
		labTestCaseMap.LabID,
		labTestCaseMap.TestcaseID,
		labTestCaseMap.Creator,
		labTestCaseMap.CreateTime,
	)
	return result, err

}

package models

// LabSubmit 提交表
type LabSubmit struct {
	Model
	// LabID 实验室id
	LabID uint64 `json:"lab_id"`
	// SubmitData 提交内容
	SubmitData string `json:"submit_data"`
	// SubmitResult 提交结果
	SubmitResult string `json:"submit_result"`
}

/**
使用标准ACM OnlineJudget状态
Pending:在线评测系统正忙，需要等待一段时间才能评测你的代码。
Pending Rejudge:测试数据更新了，现在评测系统需要重新评判你的代码。
Compiling:评测系统正在编译你的程序。
Judging Test #<Test Data Number>:评测系统现在正在测试你的程序。
Accepted:你的程序通过了所有的测试点。
Presentation Error(PE):你输出的格式与测试数据的标准格式有差别。请规范检查空行、空格和制表符。
Wrong Answer(WA):你的程序输出的结果与正确答案不同。仅通过样例并不代表这是正确答案。
Time Limit Exceeded(TLE):你的程序花费的时间超过了时间限制（一个标程1000ms）。试着优化算法。
Memory Limit Exceeded(MLE):你的程序花费的内存超过了内存限制（一般为64MB或128MB）。
Output Limit Exceeded(OLE):你的程序输出了超过标准答案两倍的字符。则一般是死循环所致。
Runtime Error(RE):你的程序发生了运行时错误。有可能是数组越界，指针错误或除以0。
Compile Error(CE):编译器发现了源代码的语法错误，以至于无法产生可执行文件。可以查看错误信息。
Compile OK:比赛结束前不能知道分数，仅显示编译是否通过。这是编译通过。
Test:OJ网站管理员功能，可以测试运行。
Other Error:你的程序出现了其它错误。
System Error(SE):评测系统出现了问题。请向管理员汇报。
*/
const (
	EMPTY = iota
	LABSUBMITSTATUS_PENDING
	LABSUBMITSTATUS_ERROR
	LABSUBMITSTATUS_COMPILING
	LABSUBMITSTATUS_JUDING
	LABSUBMITSTATUS_ACCEPTED
	LABSUBMITSTATUS_PRESENTATION_ERROR
	LABSUBMITSTATUS_WRONG_ANSWER
	LABSUBMITSTATUS_TIME_LIMIT_EXCEEDED
	LABSUBMITSTATUS_MEMORY_LIMIT_EXCEEDED
	LABSUBMITSTATUS_OUPUT_LIMIT_EXCEED
	LABSUBMITSTATUS_RUNTIME_ERROR
	LABSUBMITSTATUS_COMPILE_ERROR
	LABSUBMITSTATUS_COMPILE_OK
	LABSUBMITSTATUS_TEST
	LABSUBMITSTATUS_OTHER_ERROR
	LABSUBMITSTATUS_SYSTEM_ERROR
)

func GetSubmitById(submitId uint64) (*LabSubmit, error) {
	stmt, err := DB.Prepare("SELECT id, lab_id, submit_data, submit_result, status, creator, create_time, update_time FROM lab_submit WHERE id = ? AND status = 1")
	row := stmt.QueryRow(&submitId)
	labSubmit := new(LabSubmit)
	row.Scan(&labSubmit.ID, &labSubmit.LabID, &labSubmit.SubmitData, &labSubmit.SubmitResult, &labSubmit.Status, &labSubmit.Creator, &labSubmit.CreateTime, &labSubmit.UpdateTime)
	return labSubmit, err
}

func GetSubmitByStatus(status, size int) ([]*LabSubmit, error) {
	stmt, err := DB.Prepare("SELECT id, lab_id, submit_data, submit_result, status, creator, create_time, update_time FROM lab_submit WHERE status = ? LIMIT ?")
	rows, err := stmt.Query(
			&status,
			&size,
		)
	defer rows.Close()
	var labSubmits []*LabSubmit
	for rows.Next() {
		var labSubmit LabSubmit
		rows.Scan(
				&labSubmit.ID,
				&labSubmit.LabID,
				&labSubmit.SubmitData,
				&labSubmit.SubmitResult,
				&labSubmit.Status,
				&labSubmit.Creator,
				&labSubmit.CreateTime,
				&labSubmit.UpdateTime,
			)
		labSubmits = append(labSubmits, &labSubmit)
	}
	return labSubmits, err
}

func UpdateSubmitStatus(submitId uint64, fromStaus, toStatus int) (int64, error) {
	stmt, err := DB.Prepare("UPDATE lab_submit SET status=? WHERE status=? AND id=?")
	ret, err := stmt.Exec(&toStatus, &fromStaus, &submitId)
	rowsAffected, err := ret.RowsAffected()
	return rowsAffected, err
}

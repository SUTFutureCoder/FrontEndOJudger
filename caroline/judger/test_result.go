package judger


// 测试结果构造
type TestResult struct {
	Id             uint64
	TestCaseId     int
	TestCaseInput  string
	SubmitOutput   string
	TestcaseOutput string
	Status         int
	Err            string
}

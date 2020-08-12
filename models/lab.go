package models

// Lab 实验室表
type Lab struct {
	Model
	// LabName 实验室名称
	LabName string `json:"lab_name"`
	// LabDesc 实验室描述
	LabDesc string `json:"lab_desc"`
	// LabType 实验室类型
	LabType int8 `json:"lab_type"`
	// LabSample 实验室样例或地址
	LabSample string `json:"lab_sample"`
}

const (
	LABTYPE_HTML = iota
	LABTYPE_CSS
	LABTYPE_JS
	LABTYPE_VUE
	LABTYPE_COMPLEX
	LABTYPE_PRD
	LABTYPE_IMITATE
	LABTYPE_SECURITY
	LABTYPE_OTHER
)

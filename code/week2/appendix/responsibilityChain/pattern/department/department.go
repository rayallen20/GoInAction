package department

import "responsibilityChain/pattern"

type Department interface {
	Execute(patient *pattern.Patient) // Execute 用于让病人执行当前部门的处理
	SetNext(department Department)    // SetNext 设置下一个部门
}

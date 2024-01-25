package department

import (
	"fmt"
	"responsibilityChain/pattern"
)

// Medical 药房部门
type Medical struct {
	next Department // next 药房部门处理病人完成后,即将处理病人的下一个部门
}

// Execute 本方法用于执行药房部门处理病人的过程 并调用下一个处理病人部门的执行逻辑
func (m *Medical) Execute(patient *pattern.Patient) {
	if patient.MedicineDone {
		fmt.Printf("Medical already given to patient\n")
		if m.next != nil {
			m.next.Execute(patient)
			return
		}
	}

	fmt.Printf("Medical giving to patient\n")
	patient.MedicineDone = true
	if m.next != nil {
		m.next.Execute(patient)
	}

	return
}

// SetNext 设置下一个处理病人的部门
func (m *Medical) SetNext(next Department) {
	m.next = next
}

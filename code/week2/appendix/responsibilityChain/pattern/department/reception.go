package department

import (
	"fmt"
	"responsibilityChain/pattern"
)

// Reception 挂号部门
type Reception struct {
	next Department // next 挂号部门处理病人完成后,即将处理病人的下一个部门
}

// Execute 本方法用于执行挂号部门处理病人的过程 并调用下一个处理病人部门的执行逻辑
func (r *Reception) Execute(patient *pattern.Patient) {
	if patient.RegistrationDone {
		fmt.Printf("patient registration already done\n")
		if r.next != nil {
			r.next.Execute(patient)
		}
		return
	}

	fmt.Printf("Reception registering patient now\n")
	patient.RegistrationDone = true
	if r.next != nil {
		r.next.Execute(patient)
	}
	return
}

// SetNext 设置下一个处理病人的部门
func (r *Reception) SetNext(department Department) {
	r.next = department
}

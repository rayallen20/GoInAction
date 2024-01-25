package department

import (
	"fmt"
	"responsibilityChain/pattern"
)

// Doctor 医生部门
type Doctor struct {
	next Department // next 医生部门处理病人完成后,即将处理病人的下一个部门
}

// Execute 本方法用于执行医生部门处理病人的过程 并调用下一个处理病人部门的执行逻辑
func (d *Doctor) Execute(patient *pattern.Patient) {
	if patient.DoctorCheckUpDone {
		fmt.Printf("Doctor checkup already done\n")
		if d.next != nil {
			d.next.Execute(patient)
			return
		}
	}

	fmt.Printf("Doctor checking patient\n")
	patient.DoctorCheckUpDone = true
	if d.next != nil {
		d.next.Execute(patient)
	}
	return
}

// SetNext 设置下一个处理病人的部门
func (d *Doctor) SetNext(next Department) {
	d.next = next
}

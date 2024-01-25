package department

import (
	"fmt"
	"responsibilityChain/pattern"
)

// Cashier 收费部门
type Cashier struct {
	next Department // next 收费部门处理病人完成后,即将处理病人的下一个部门
}

// Execute 本方法用于执行收费部门处理病人的过程 并调用下一个处理病人部门的执行逻辑
func (c *Cashier) Execute(patient *pattern.Patient) {
	if patient.PaymentDone {
		fmt.Printf("Payment Done\n")

		if c.next != nil {
			c.next.Execute(patient)
			return
		}
	}

	fmt.Printf("Cashier getting money from patient patient\n")
	patient.PaymentDone = true
	if c.next != nil {
		c.next.Execute(patient)
	}

	return
}

// SetNext 设置下一个处理病人的部门
func (c *Cashier) SetNext(department Department) {
	c.next = department
}

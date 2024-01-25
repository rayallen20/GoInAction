package main

import (
	"responsibilityChain/pattern"
	"responsibilityChain/pattern/department"
)

func main() {
	patient := &pattern.Patient{
		Name: "abc",
	}

	departmentChain := initDepartmentChain()
	departmentChain.Execute(patient)
}

// initDepartmentChain 本函数用于初始化部门链并返回位于链首部的第1个部门
// Tips: 注意部门初始化的顺序和部门处理病人的顺序是相反的
func initDepartmentChain() department.Department {
	// 初始化收费部门
	cashier := &department.Cashier{}

	// 初始化药房部门
	medical := &department.Medical{}
	medical.SetNext(cashier)

	// 初始化医生部门
	doctor := &department.Doctor{}
	doctor.SetNext(medical)

	// 初始化挂号部门
	reception := &department.Reception{}
	reception.SetNext(doctor)

	return reception
}

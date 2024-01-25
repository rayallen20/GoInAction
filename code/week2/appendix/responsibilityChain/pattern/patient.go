package pattern

// Patient 病人结构体
type Patient struct {
	Name              string // 姓名
	RegistrationDone  bool   // 是否完成挂号
	DoctorCheckUpDone bool   // 是否完成医生检查
	MedicineDone      bool   // 是否完成取药
	PaymentDone       bool   // 是否完成缴费
}

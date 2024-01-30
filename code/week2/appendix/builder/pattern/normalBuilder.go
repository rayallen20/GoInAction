package pattern

// NormalBuilder 普通房屋生成器
type NormalBuilder struct {
	windowType string // windowType 普通房屋的窗户类型
	doorType   string // doorType 普通房屋的门的类型
	floor      int    // floor 普通房屋的房屋楼层数
}

// newNormalBuilder NormalBuilder 的构造函数
func newNormalBuilder() *NormalBuilder {
	return &NormalBuilder{}
}

// setWindowType 本方法用于设置普通房屋的窗户类型
func (b *NormalBuilder) setWindowType() {
	b.windowType = "木窗户"
}

// setDoorType 本方法用于设置普通房屋的门类型
func (b *NormalBuilder) setDoorType() {
	b.doorType = "木门"
}

// setFloor 本方法用于设置普通房屋的楼层数
func (b *NormalBuilder) setFloor() {
	b.floor = 2
}

// getHouse 本方法用于创建普通房屋
func (b *NormalBuilder) getHouse() House {
	return House{
		windowType: b.windowType,
		doorType:   b.doorType,
		floor:      b.floor,
	}
}

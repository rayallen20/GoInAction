package pattern

// IglooBuilder 冰屋生成器
type IglooBuilder struct {
	windowType string // windowType 冰屋的窗户类型
	doorType   string // doorType 冰屋的门的类型
	floor      int    // floor 冰屋的房屋楼层数
}

// newIglooBuilder IglooBuilder 的构造函数
func newIglooBuilder() *IglooBuilder {
	return &IglooBuilder{}
}

// setWindowType 本方法用于设置冰屋的窗户类型
func (b *IglooBuilder) setWindowType() {
	b.windowType = "雪窗户"
}

// setDoorType 本方法用于设置冰屋的门类型
func (b *IglooBuilder) setDoorType() {
	b.doorType = "雪门"
}

// setFloor 本方法用于设置冰屋的楼层数
func (b *IglooBuilder) setFloor() {
	b.floor = 1
}

// getHouse 本方法用于创建冰屋
func (b *IglooBuilder) getHouse() House {
	return House{
		windowType: b.windowType,
		doorType:   b.doorType,
		floor:      b.floor,
	}
}

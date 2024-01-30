package pattern

// House 房屋结构体 即最终对外提供的产品
type House struct {
	windowType string // windowType 窗户类型
	doorType   string // doorType 门的类型
	floor      int    // floor 房屋楼层数
}

// GetWindowType 本方法用于获取房屋的窗户类型
func (h House) GetWindowType() string {
	return h.windowType
}

// GetDoorType 本方法用于获取房屋的门的类型
func (h House) GetDoorType() string {
	return h.doorType
}

// GetFloor 本方法用于获取房屋的楼层数
func (h House) GetFloor() int {
	return h.floor
}

package pattern

// IBuilder 生成器接口 本接口用于定义创建房屋的必要过程
type IBuilder interface {
	setWindowType()  // setWindowType 本方法用于设置窗户类型
	setDoorType()    // setDoorType 本方法用于设置门的类型
	setFloor()       // setFloor 本方法用于设置房屋楼层数
	getHouse() House // getHouse 本方法用于创建并返回房屋
}

// GetBuilder 本函数用于根据给定的生成器类型 创建具体生成器
func GetBuilder(builderType string) IBuilder {
	switch builderType {
	case "normal":
		return newNormalBuilder()
	case "igloo":
		return newIglooBuilder()
	default:
		return nil
	}
}

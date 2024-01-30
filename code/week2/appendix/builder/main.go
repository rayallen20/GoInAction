package main

import (
	"builder/pattern"
	"fmt"
)

func main() {
	// 普通房屋生成器实例
	normalBuilder := pattern.GetBuilder("normal")
	// 冰屋生成器实例
	iglooBuilder := pattern.GetBuilder("igloo")

	// 管理者实例
	director := pattern.NewDirector(normalBuilder)
	// 生成普通房屋
	normalHouse := director.BuildHouse()

	fmt.Printf("普通房屋的窗户类型: %s\n", normalHouse.GetWindowType())
	fmt.Printf("普通房屋的门的类型: %s\n", normalHouse.GetDoorType())
	fmt.Printf("普通房屋的楼层数: %d\n", normalHouse.GetFloor())

	// 重置管理者实例中的生成器实例
	director.SetBuilder(iglooBuilder)
	// 生成冰屋
	iglooHouse := director.BuildHouse()

	fmt.Printf("冰屋的窗户类型: %s\n", iglooHouse.GetWindowType())
	fmt.Printf("冰屋的门的类型: %s\n", iglooHouse.GetDoorType())
	fmt.Printf("冰屋的楼层数: %d\n", iglooHouse.GetFloor())
}

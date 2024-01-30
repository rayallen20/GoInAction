package pattern

// Director 管理者结构体 本结构体用于编排生成器的生成过程 并最终生成产品
// 本结构体的存在使得客户端不必关心生成器的生成过程
type Director struct {
	builder IBuilder // builder 生成器接口的实现
}

// NewDirector Director 的构造函数
func NewDirector(builder IBuilder) *Director {
	return &Director{
		builder: builder,
	}
}

// SetBuilder 本方法用于为管理者结构体设置具体生成器
// 本方法存在的意义在于可以在运行时动态的改变管理者实例中的生成器实例
func (d *Director) SetBuilder(builder IBuilder) {
	d.builder = builder
}

// BuildHouse 本方法用于编排具体生成器的生成过程 并最终生成产品
func (d *Director) BuildHouse() House {
	d.builder.setWindowType()
	d.builder.setDoorType()
	d.builder.setFloor()
	return d.builder.getHouse()
}

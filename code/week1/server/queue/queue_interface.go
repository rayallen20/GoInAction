package queue

import "context"

// 1. 应该使用泛型
// 2. 并发队列 有 阻塞 和 非阻塞形态 (没有控制超时时长的参数或方法)
// 3. 通过返回值判定是否超时
type Queue interface {
	// Append 追加一个元素到队列尾部
	Append(ctx context.Context, v any) error
	// Prepend 追加一个元素到队列头部
	Prepend(ctx context.Context, v any) error
	// Pop 弹出队列头部的元素 (因为我在设计一个MQ 所以没有从队列尾部取元素这个操作我认为)
	// 将返回的元素从MQ中删除掉
	Pop(ctx context.Context) (any, error)
	// Random 随机返回队列中的一个元素
	// Random 认为队列中的所有元素都是平等的 各元素分摊负载均衡
	Random(ctx context.Context) (any, error)
	// Peek 打印队列头部的元素 (只打印不弹出)
	// 返回的元素没有从MQ中删除
	// 打印 这个"行为"太抽象 无法确定它最终打印到 console/file/JSON等 还不如返回一个队头的副本
	Peek() (any, error)
	// IsEmpty 判断队列是否为空
	IsEmpty() bool
	// Size 返回队列的大小
	Size() int
	// Cap 返回队列的容量
	// 1. 基于链表的实现 没有Cap()这个行为
	// 2. 把Cap()还是定义在接口上 但链表的实现返回Size()即可
	// 3. 基于链表的允许外部来传一个最大容量 表示链表内最多允许的节点个数
	// 4. 基于链表的Cap()没上限
	Cap() int
	// expandCap 扩容 这个操作只能MQ内部完成
	// "扩容"不是所有的MQ都有的行为 只有具备"扩容"能力的MQ才需要实现这个方法
	// 接口应该是"通用"的行为
	expandCap()

	// Traversal 遍历
	// f 每个元素都调用这个函数
	Traversal(ctx context.Context, f func(any) bool)
}

//type Queue interface {
//	Enqueue(v any) error
//	Dequeue() (v any, err error)
//	Traversal(ctx context.Context, f func(any) error) error
//}

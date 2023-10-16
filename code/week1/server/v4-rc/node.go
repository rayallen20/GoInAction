package v4_rc

// node 路由树中的节点
type node struct {
	// path 路由路径
	path string
	// children 子节点 key为子节点的路由路径 value为路径对应子节点
	children map[string]*node
	// HandleFunc 路由对应的处理函数
	HandleFunc
}

// findOrCreate 本方法用于根据给定的path值 在当前节点的子节点中查找path为给定path值的节点
// 找到则返回 未找到则创建
func (n *node) findOrCreate(segment string) *node {
	if n.children == nil {
		n.children = make(map[string]*node)
	}

	target, exist := n.children[segment]
	if !exist {
		// 当前节点的子节点映射中不存在目标子节点 则创建目标子节点 将子节点加入当前节点的子节点映射后返回
		target = &node{
			path: segment,
		}
		n.children[segment] = target
		return target
	}

	// 当前节点的子节点映射中存在目标子节点 则直接返回
	return target
}

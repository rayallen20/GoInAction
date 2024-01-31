package web

import (
	"fmt"
	"strings"
)

// router 路由森林 用于支持对路由树的操作
type router struct {
	// trees 路由森林 按HTTP动词组织路由树
	// 该map中 key为HTTP动词 value为路由树的根节点
	// 即: 每个HTTP动词对应一棵路由树 指向每棵路由树的根节点
	trees map[string]*node
}

// newRouter 创建路由森林
func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// addRoute 注册路由到路由森林中的路由树上
func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	if path == "" {
		panic("web: 路由不能为空字符串")
	}

	if path[0] != '/' {
		panic("web: 路由必须以 '/' 开头")
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路由不能以 '/' 结尾")
	}

	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	if path == "/" {
		if root.HandleFunc != nil {
			panic("web: 路由冲突,重复注册路由 [/] ")
		}
		root.HandleFunc = handleFunc

		// 记录节点的全路由
		root.fullRoute = path
		return
	}

	path = strings.TrimLeft(path, "/")
	segments := strings.Split(path, "/")

	target := root
	for _, segment := range segments {
		if segment == "" {
			panic("web: 路由中不得包含连续的'/'")
		}

		child := target.childOrCreate(segment)
		target = child
	}

	if target.HandleFunc != nil {
		panic(fmt.Sprintf("web: 路由冲突,重复注册路由 [%s] ", path))
	}

	target.HandleFunc = handleFunc

	// 记录节点的全路由
	target.fullRoute = path
}

// findRoute 根据给定的HTTP方法和路由路径,在路由森林中查找对应的节点
// 若该节点为参数路径节点,则不仅返回该节点,还返回参数名和参数值
// 否则,仅返回该节点
func (r *router) findRoute(method string, path string) (*matchNode, bool) {
	targetMatchNode := &matchNode{}
	root, ok := r.trees[method]
	// 给定的HTTP动词在路由森林中不存在对应的路由树,则直接返回false
	if !ok {
		return nil, false
	}

	// 对根节点做特殊处理
	if path == "/" {
		targetMatchNode.node = root
		return targetMatchNode, true
	}

	// 给定的HTTP动词在路由森林中存在对应的路由树,则在该路由树中查找对应的节点
	// 去掉前导和后置的"/"
	path = strings.Trim(path, "/")
	segments := strings.Split(path, "/")

	// Tips: 同样的 这里我认为用target作为变量名表现力更强
	target := root

	for _, segment := range segments {
		child, isParamChild, found := target.childOf(segment)
		// 如果在当前节点的子节点映射中没有找到对应的子节点,则直接返回
		if !found {
			return nil, false
		}

		// 若当前节点为参数节点,则将参数名和参数值保存到targetMatchNode中
		if isParamChild {
			// 参数名是形如 :id 的格式, 因此需要去掉前导的:
			name := child.path[1:]
			// 参数值就是当前路由路径中的路由段
			value := segment
			targetMatchNode.addPathParams(name, value)
		}

		// 如果在当前节点的子节点映射中找到了对应的子节点,则继续在该子节点中查找
		target = child
	}

	// 如果找到了对应的节点,则返回该节点
	// Tips: 此处有2种设计 一种是用标量表示是否找到了子节点
	// Tips: 另一种是 return target, target.HandleFunc != nil
	// Tips: 这种返回就表示找到了子节点且子节点必然有对应的业务处理函数
	// 此处我倾向用第1种设计 因为方法名叫findRoute,表示是否找到节点的意思.而非表示是否找到了一个有对应的业务处理函数的节点
	targetMatchNode.node = target
	return targetMatchNode, true
}

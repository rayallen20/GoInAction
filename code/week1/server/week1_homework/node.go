package week1_homework

import (
	"fmt"
	"regexp"
	"strings"
)

// nodeType 节点类型
// 0: 静态路由
// 1: 正则路由
// 2: 参数路由
// 3: 通配符路由
type nodeType int

const (
	// nodeTypeStatic 静态路由
	nodeTypeStatic = iota
	// nodeTypeReg 正则路由
	nodeTypeReg
	// nodeTypeParam 参数路由
	nodeTypeParam
	// nodeTypeAny 通配符路由
	nodeTypeAny
)

// node 路由树的节点
type node struct {
	// typ 节点类型
	typ nodeType

	// path 当前节点的路径
	path string

	// children 子路由路径到子节点的映射
	children map[string]*node

	// wildcardChild 通配符子节点
	wildcardChild *node

	// paramChild 参数子节点
	paramChild *node

	// regChild 正则子节点
	regChild *node

	// regExp 编译好的正则表达式
	regExp *regexp.Regexp

	// HandleFunc 路由对应的业务逻辑
	HandleFunc
}

// childOrCreate 本方法用于在节点上获取给定的子节点,如果给定的子节点不存在则创建
func (n *node) childOrCreate(segment string) *node {
	if n.isRegChild(segment) {
		err := n.createRegChild(segment)
		if err != nil {
			panic(err.Error())
		}
		return n.regChild
	}

	if n.isParamChild(segment) {
		err := n.createParamChild(segment)
		if err != nil {
			panic(err.Error())
		}
		return n.paramChild
	}

	if n.isWildcardChild(segment) {
		err := n.createWildcardChild(segment)
		if err != nil {
			panic(err.Error())
		}
		return n.wildcardChild
	}

	n.createStaticChild(segment)
	return n.children[segment]
}

// isRegChild 判断当前节点是否为正则路由子节点
// 判断依据:
// 1. 以":"开头
// 2. 包含"("
// 3. 以")"结尾
func (n *node) isRegChild(segment string) bool {
	return strings.HasPrefix(segment, ":") && strings.Contains(segment, "(") && strings.HasSuffix(segment, ")")
}

// createRegChild 创建正则路由子节点
// 创建依据:
// 1. 若同路由段上的参数路由子节点已存在 则不允许再创建
// 2. 若同路由段上的正则路由子节点已存在 则不允许再创建
// 3. 若同路由段上的通配符路由子节点已存在 则不允许再创建
// 4. 若无法从路由段中提取出正则表达式 则不允许创建
func (n *node) createRegChild(segment string) error {
	if n.regChild != nil {
		msg := fmt.Sprintf("web: 非法路由,已有正则子节点 %s .不允许同时注册多个正则子节点", n.regChild.regExp.String())
		return fmt.Errorf(msg)
	}

	if n.paramChild != nil {
		msg := fmt.Sprintf("web: 非法路由,已有参数子节点 %s .不允许同时注册正则子节点与参数子节点", n.paramChild.path)
		return fmt.Errorf(msg)
	}

	if n.wildcardChild != nil {
		msg := fmt.Sprintf("web: 非法路由,已有通配符子节点 %s .不允许同时注册通配符子节点与正则子节点", n.wildcardChild.path)
		return fmt.Errorf(msg)
	}

	regStr, err := n.extractReg(segment)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	re, err := regexp.Compile(regStr)
	if err != nil {
		msg := fmt.Sprintf("web: 非法路由,无法编译正则表达式.路由段 %s", segment)
		return fmt.Errorf(msg)
	}

	n.regChild = &node{
		path:   segment,
		typ:    nodeTypeReg,
		regExp: re,
	}

	return nil
}

// isParamChild 判断当前节点是否为参数路由子节点
// 判断依据:
// 1. 以":"开头 且不以")"结尾
func (n *node) isParamChild(segment string) bool {
	return strings.HasPrefix(segment, ":") && !strings.HasSuffix(segment, ")")
}

// createParamChild 创建参数路由子节点
// 创建依据:
// 1. 若同路由段上的正则路由子节点已存在 则不允许再创建
// 2. 若同路由段上的参数路由子节点已存在 则不允许再创建
// 3. 若同路由段上的通配符路由子节点已存在 则不允许再创建
func (n *node) createParamChild(segment string) error {
	if n.regChild != nil {
		msg := fmt.Sprintf("web: 非法路由,已有正则子节点 %s .不允许同时注册正则子节点与参数子节点", n.regChild.regExp.String())
		return fmt.Errorf(msg)
	}

	if n.paramChild != nil {
		msg := fmt.Sprintf("web: 非法路由,已有参数子节点 %s .不允许同时注册多个参数子节点", n.paramChild.path)
		return fmt.Errorf(msg)
	}

	if n.wildcardChild != nil {
		msg := fmt.Sprintf("web: 非法路由,已有通配符子节点 %s .不允许同时注册通配符子节点与参数子节点", n.wildcardChild.path)
		return fmt.Errorf(msg)
	}

	n.paramChild = &node{
		path: segment,
		typ:  nodeTypeParam,
	}
	return nil
}

// isWildcardChild 判断当前节点是否为通配符路由子节点
// 判断依据:
// 1. 路由段为"*"
func (n *node) isWildcardChild(segment string) bool {
	return segment == "*"
}

// createWildcardChild 创建通配符路由子节点
// 创建依据:
// 1. 若同路由段上的正则路由子节点已存在 则不允许再创建
// 2. 若同路由段上的参数路由子节点已存在 则不允许再创建
// 3. 若同路由段上的通配符路由子节点已存在 则不允许再创建
func (n *node) createWildcardChild(segment string) error {
	if n.regChild != nil {
		msg := fmt.Sprintf("web: 非法路由,已有正则子节点 %s .不允许同时注册通配符子节点与正则子节点", n.regChild.regExp.String())
		return fmt.Errorf(msg)
	}

	if n.paramChild != nil {
		msg := fmt.Sprintf("web: 非法路由,已有参数子节点 %s .不允许同时注册通配符子节点与参数子节点", n.paramChild.path)
		return fmt.Errorf(msg)
	}

	if n.wildcardChild != nil {
		msg := fmt.Sprintf("web: 非法路由,已有通配符子节点 %s .不允许同时注册多个通配符子节点", n.wildcardChild.path)
		return fmt.Errorf(msg)
	}

	n.wildcardChild = &node{
		path: segment,
		typ:  nodeTypeAny,
	}

	return nil
}

// createStaticChild 创建静态路由子节点
func (n *node) createStaticChild(segment string) {
	if n.children == nil {
		n.children = map[string]*node{}
	}

	if _, ok := n.children[segment]; ok {
		return
	}

	res := &node{
		path: segment,
		typ:  nodeTypeStatic,
	}
	n.children[segment] = res

	return
}

// childOf 根据给定的path在当前节点的子节点映射中查找对应的子节点
// 查找优先级: 静态路由 > 正则路由 > 参数路由 > 通配符路由
func (n *node) childOf(path string) (child *node, found bool) {
	// 当前节点的子节点映射为空 则尝试查找 正则子节点 或 参数路由子节点 或 通配符子节点
	if n.children == nil {
		// 优先尝试查找正则子节点
		if n.regChild != nil && n.regChild.regExp.MatchString(path) {
			return n.regChild, true
		}

		// 再尝试查找参数路由子节点
		if n.paramChild != nil {
			return n.paramChild, true
		}

		// 最后尝试查找通配符子节点
		return n.wildcardChild, n.wildcardChild != nil
	}

	// 在子当前节点的节点映射中查找对应的子节点 若未找到同样尝试返回当前节点的参数子节点
	// 若参数子节点为空 则尝试返回当前节点的通配符子节点
	child, found = n.children[path]
	if !found {
		if n.paramChild != nil {
			return n.paramChild, true
		}
		return n.wildcardChild, n.wildcardChild != nil
	}

	// 找到了对应的子节点 则返回该子节点
	return child, found
}

// extractReg 从给定的路由段中提取出正则表达式
func (n *node) extractReg(segment string) (regStr string, err error) {
	startIndex := strings.Index(segment, "(")
	endIndex := strings.LastIndex(segment, ")")
	if startIndex != -1 && endIndex != -1 && endIndex > startIndex+1 {
		return segment[startIndex+1 : endIndex], nil
	}

	return "", fmt.Errorf("web: 非法路由,正则路由格式错误.路由段 %s", segment)
}

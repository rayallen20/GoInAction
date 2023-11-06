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
	// 如果路径为参数 则查找当前节点的参数子节点 或创建一个当前节点的参数子节点 并返回
	// Tips: 判断参数路由的条件: 以":"开头 且不以")"结尾
	if strings.HasPrefix(segment, ":") && !strings.HasSuffix(segment, ")") {
		// 若当前节点存在通配符子节点 则不允许注册参数子节点
		if n.wildcardChild != nil {
			panic("web: 非法路由,已有通配符路由.不允许同时注册通配符路由和参数路由")
		}

		// 若当前节点的参数子节点不为空 说明当前节点已被注册了一个参数子节点 不允许再注册参数子节点
		if n.paramChild != nil {
			msg := fmt.Sprintf("web: 路由冲突,参数路由冲突.已存在路由 %s", n.paramChild.path)
			panic(msg)
		}

		n.paramChild = &node{
			path: segment,
			typ:  nodeTypeParam,
		}
		return n.paramChild
	}

	// 如果路径为正则 则查找当前节点的正则子节点 或创建一个当前节点的正则子节点 并返回
	// Tips: 判断正则路由的条件: 以":"开头 且包含"(" 且以")"结尾
	if strings.HasPrefix(segment, ":") && strings.Contains(segment, "(") && strings.HasSuffix(segment, ")") {
		// 若当前节点存在正则子节点 则不允许注册正则子节点
		if n.regChild != nil {
			msg := fmt.Sprintf("web: 非法路由,已有正则路由 %s .不允许同时注册多个正则路由", n.regChild.regExp.String())
			panic(msg)
		}

		// 若当前节点存在通配符子节点 则不允许注册正则子节点
		if n.wildcardChild != nil {
			panic("web: 非法路由,已有通配符路由.不允许同时注册通配符路由和正则路由")
		}

		// 提取正则表达式
		regStr, err := n.extractReg(segment)
		if err != nil {
			panic(err)
		}

		// 编译正则表达式
		re, err := regexp.Compile(regStr)
		if err != nil {
			msg := fmt.Sprintf("web: 非法路由,无法编译正则表达式.路由段 %s", segment)
			panic(msg)
		}

		// 创建正则子节点
		n.regChild = &node{
			path:   segment,
			typ:    nodeTypeReg,
			regExp: re,
		}
		return n.regChild
	}

	// 若路径为通配符 则查找当前节点的通配符子节点 或创建一个当前节点的通配符子节点 并返回
	if segment == "*" {
		// 若当前节点存在参数子节点 则不允许注册通配符子节点
		if n.paramChild != nil {
			panic("web: 非法路由,已有参数路由.不允许同时注册通配符路由和参数路由")
		}

		if n.wildcardChild == nil {
			n.wildcardChild = &node{
				path: segment,
				typ:  nodeTypeAny,
			}
		}
		return n.wildcardChild
	}

	// 如果当前节点的子节点映射为空 则创建一个子节点映射
	if n.children == nil {
		n.children = map[string]*node{}
	}

	res, ok := n.children[segment]
	// 如果没有找到子节点,则创建一个子节点;否则返回找到的子节点
	if !ok {
		res = &node{
			path: segment,
			typ:  nodeTypeStatic,
		}
		n.children[segment] = res
	}
	return res
}

// childOf 根据给定的path在当前节点的子节点映射中查找对应的子节点(即:匹配到了静态路由)
// 若未在子节点映射中找到对应子节点 则先尝试返回当前节点的正则子节点(即:匹配到了正则路由) 若正则子节点为空
// 则尝试返回当前节点的参数路由子节点(即:匹配到了参数路由) 若参数路由子节点为空
// 则尝试返回当前节点的通配符子节点(即:匹配到了通配符路由)
// 优先级: 静态路由 > 正则路由 > 参数路由 > 通配符路由
// child: 查找到的子节点
// isParamChild: 查找到的子节点是否为参数路由子节点
// found: 是否找到了对应的子节点
func (n *node) childOf(path string) (child *node, isParamChild bool, found bool) {
	// 当前节点的子节点映射为空 则有可能匹配到 正则子节点 或 参数路由子节点 或 通配符子节点
	// 此处优先查找正则子节点 再查找参数路由子节点 最后查找通配符子节点
	// 因为正则子节点更具体 所以正则子节点的优先级高于参数路由子节点
	// 此处优先查找参数路由子节点 因为参数路由子节点更具体 所以参数路由的优先级高于通配符路由
	if n.children == nil {
		// 如果当前节点的正则子节点不为空 则检测参数是否匹配正则表达式 若匹配则返回当前节点的正则子节点
		if n.regChild != nil {
			if n.regChild.regExp.MatchString(path) {
				return n.regChild, false, true
			}
		}

		// 如果当前节点的参数子节点不为空 则尝试返回当前节点的参数子节点
		if n.paramChild != nil {
			return n.paramChild, true, true
		}

		// 如果当前节点的参数子节点为空 则尝试返回当前节点的通配符子节点
		return n.wildcardChild, false, n.wildcardChild != nil
	}

	// 在子当前节点的节点映射中查找对应的子节点 若未找到同样尝试返回当前节点的参数子节点
	// 若参数子节点为空 则尝试返回当前节点的通配符子节点
	child, found = n.children[path]
	if !found {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.wildcardChild, false, n.wildcardChild != nil
	}

	// 找到了对应的子节点 则返回该子节点
	return child, false, found
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

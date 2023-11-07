package week1_homework

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"regexp"
	"testing"
)

// TestNode 测试路由树节点
// 由于此处我们要测试的是路由树的结构,因此不需要在测试路由树节点中添加路由处理函数
// 调用addRoute时写死一个HandleFunc即可
type TestNode struct {
	method string
	path   string
}

// TestRouter_addRoute 测试路由注册功能的结果是否符合预期
func TestRouter_addRoute(t *testing.T) {
	// step1. 构造路由树
	testRoutes := []TestNode{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
	}

	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}

	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandleFunc)
	}

	// step2. 验证路由树 断言二者是否相等
	wantRouter := &router{
		trees: map[string]*node{
			// GET方法路由树
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"user": {
						path: "user",
						children: map[string]*node{
							"home": &node{
								path:     "home",
								children: nil,
								// 注意路由是/user/home 因此只有最深层的节点才有handleFunc
								// /user和/ 都是没有handleFunc的
								HandleFunc: mockHandleFunc,
							},
						},
						HandleFunc: mockHandleFunc,
					},
					"order": &node{
						path: "order",
						children: map[string]*node{
							"detail": &node{
								path:       "detail",
								children:   nil,
								HandleFunc: mockHandleFunc,
							},
						},
						HandleFunc: nil,
					},
				},
				HandleFunc: mockHandleFunc,
			},

			// POST方法路由树
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"order": &node{
						path: "order",
						children: map[string]*node{
							"create": &node{
								path:       "create",
								children:   nil,
								HandleFunc: mockHandleFunc,
							},
						},
						HandleFunc: nil,
					},
					"login": &node{
						path:       "login",
						children:   nil,
						HandleFunc: mockHandleFunc,
					},
				},
				HandleFunc: nil,
			},
		},
	}

	// HandleFunc类型是方法,方法不可比较,因此只能比较两个路由树的结构是否相等
	// assert.Equal(t, wantRouter, r)

	msg, ok := wantRouter.equal(&r)
	assert.True(t, ok, msg)
}

// equal 比较两个路由森林是否相等
// msg: 两个路由森林不相等时的错误信息
// ok: 两个路由森林是否相等
func (r *router) equal(target *router) (msg string, ok bool) {
	// 如果目标路由森林为nil 则不相等
	if target == nil {
		return fmt.Sprintf("目标路由森林为nil"), false
	}

	// 如果两个路由森林中的路由树数量不同 则不相等
	rTreesNum := len(r.trees)
	yTreesNum := len(target.trees)
	if rTreesNum != yTreesNum {
		return fmt.Sprintf("路由森林中的路由树数量不相等,源路由森林有 %d 棵路由树, 目标路由森林有 %d 棵路由树", rTreesNum, yTreesNum), false
	}

	for method, tree := range r.trees {
		dstTree, ok := target.trees[method]

		// 如果目标router中没有对应HTTP方法的路由树 则不相等
		if !ok {
			return fmt.Sprintf("目标 router 中没有HTTP方法 %s的路由树", method), false
		}

		// 比对两棵路由树的结构是否相等
		msg, equal := tree.equal(dstTree)
		if !equal {
			return method + "-" + msg, false
		}
	}
	return "", true
}

// equal 比较两棵路由树是否相等
// msg: 两棵路由树不相等时的错误信息
// ok: 两棵路由树是否相等
func (n *node) equal(target *node) (msg string, ok bool) {
	// 如果目标节点为nil 则不相等
	if target == nil {
		return fmt.Sprintf("目标节点为nil"), false
	}

	// 如果两个节点的类型不相等 则不相等
	if n.typ != target.typ {
		return fmt.Sprintf("两个节点的类型不相等,源节点的类型为 %d,目标节点的类型为 %d", n.typ, target.typ), false
	}

	// 如果两个节点的path不相等 则不相等
	if n.path != target.path {
		return fmt.Sprintf("两个节点的path不相等,源节点的path为 %s,目标节点的path为 %s", n.path, target.path), false
	}

	// 若两个节点的子节点数量不相等 则不相等
	nChildrenNum := len(n.children)
	yChildrenNum := len(target.children)
	if nChildrenNum != yChildrenNum {
		return fmt.Sprintf("两个节点的子节点数量不相等,源节点的子节点数量为 %d,目标节点的子节点数量为 %d", nChildrenNum, yChildrenNum), false
	}

	// 若两个节点的参数子节点不相等 则不相等
	if n.paramChild != nil {
		if target.paramChild == nil {
			return fmt.Sprintf("目标节点的参数子节点为空"), false
		}
		_, paramIsEqual := n.paramChild.equal(target.paramChild)
		if !paramIsEqual {
			return fmt.Sprintf("两个节点的参数子节点不相等"), false
		}
	}

	// 若两个节点的正则子节点不相等 则不相等
	if n.regChild != nil {
		if target.regChild == nil {
			return fmt.Sprintf("目标节点的正则子节点为空"), false
		}

		// 比对两个正则子节点的正则表达式是否相等
		if n.regChild.regExp.String() != target.regChild.regExp.String() {
			return fmt.Sprintf("两个节点的正则子节点的正则表达式不相等, 期待的正则表达式为: %s, 实际的正则表达式为: %s", n.regChild.regExp.String(), target.regChild.regExp.String()), false
		}

		_, regIsEqual := n.regChild.equal(target.regChild)
		if !regIsEqual {
			return fmt.Sprintf("两个节点的正则子节点不相等"), false
		}
	}

	// 若两个节点的通配符子节点不相等 则不相等
	if n.wildcardChild != nil {
		if target.wildcardChild == nil {
			return fmt.Sprintf("目标节点的通配符子节点为空"), false
		}
		_, wildcardIsEqual := n.wildcardChild.equal(target.wildcardChild)
		if !wildcardIsEqual {
			return fmt.Sprintf("两个节点的通配符子节点不相等"), false
		}
	}

	// 若两个节点的handleFunc类型不同 则不相等
	nHandler := reflect.ValueOf(n.HandleFunc)
	yHandler := reflect.ValueOf(target.HandleFunc)
	if nHandler != yHandler {
		return fmt.Sprintf("%s节点的handleFunc不相等,源节点的handleFunc为 %v,目标节点的handleFunc为 %v", n.path, nHandler.Type().String(), yHandler.Type().String()), false
	}

	// 比对两个节点的子节点映射是否相等
	for path, child := range n.children {
		dstChild, ok := target.children[path]
		// 如果源节点的子节点中 存在目标节点没有的子节点 则不相等
		if !ok {
			return fmt.Sprintf("目标节点的子节点中没有path为 %s 的子节点", path), false
		}

		// 比对两个子节点是否相等
		msg, equal := child.equal(dstChild)
		if !equal {
			return msg, false
		}
	}

	return "", true
}

// equal 比较两个matchNode是否相等
// msg: 两个matchNode不相等时的错误信息
// ok: 两个matchNode是否相等
func (m *matchNode) equal(target *matchNode) (msg string, ok bool) {
	// 比对两个matchNode的node是否相等
	msg, equal := m.node.equal(target.node)
	if !equal {
		return msg, false
	}

	// 比对两个matchNode的pathParams是否相等
	if len(m.pathParams) != len(target.pathParams) {
		return fmt.Sprintf("两个matchNode的pathParams长度不相等"), false
	}

	for name, value := range m.pathParams {
		dstValue, ok := target.pathParams[name]
		if !ok {
			return fmt.Sprintf("目标matchNode的pathParams中没有name为 %s 的pathParam", name), false
		}

		if value != dstValue {
			return fmt.Sprintf("两个matchNode的pathParams中name为 %s 的pathParam的值不相等", name), false
		}
	}

	return "", true
}

// TestRouter_addRoute_Illegal_Case 测试路由注册功能的非法用例
func TestRouter_addRoute_Illegal_Case(t *testing.T) {
	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	// 为测试路由冲突 先注册路由
	r.addRoute(http.MethodGet, "/", mockHandleFunc)
	r.addRoute(http.MethodGet, "/user", mockHandleFunc)

	// step1. 断言路由注册功能的非法用例
	// 1.1 测试路由为空字符串
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "", mockHandleFunc)
	}, "web: 路由不能为空字符串")

	// 1.2 测试路由不以"/"开头
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "login", mockHandleFunc)
	}, "web: 路由必须以 '/' 开头")

	// 1.3 测试路由以"/"结尾
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/login/", mockHandleFunc)
	}, "web: 路由不能以 '/' 结尾")

	// 1.4 测试路由中包含连续的"/"
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/login///", mockHandleFunc)
	}, "web: 路由中不得包含连续的'/'")

	// 1.5 测试路由重复注册
	// a. 根节点重复注册
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandleFunc)
	}, "web: 路由冲突,重复注册路由 [/] ")

	// b. 普通节点重复注册
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/user", mockHandleFunc)
	}, "web: 路由冲突,重复注册路由 [/user] ")
}

// TestRouter_findRoute 测试路由查找功能
func TestRouter_findRoute(t *testing.T) {
	// step1. 构造路由树
	testRoutes := []TestNode{
		// GET方法路由树
		TestNode{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		TestNode{
			method: http.MethodGet,
			path:   "/",
		},
	}

	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}

	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandleFunc)
	}

	// step2. 构造测试用例
	testCases := []struct {
		name      string
		method    string
		path      string
		isFound   bool
		matchNode *matchNode
	}{
		// 测试HTTP动词不存在的用例
		{
			name:      "method not found",
			method:    http.MethodDelete,
			path:      "/user",
			isFound:   false,
			matchNode: nil,
		},

		// 测试完全命中的用例
		{
			name:    "order detail",
			method:  http.MethodGet,
			path:    "/order/detail",
			isFound: true,
			matchNode: &matchNode{
				node: &node{
					path:       "detail",
					children:   nil,
					HandleFunc: mockHandleFunc,
				},
			},
		},

		// 测试命中了节点但节点的HandleFunc为nil的情况
		{
			name:    "order",
			method:  http.MethodGet,
			path:    "/order",
			isFound: true,
			matchNode: &matchNode{
				node: &node{
					path: "order",
					children: map[string]*node{
						"detail": &node{
							path:       "detail",
							children:   nil,
							HandleFunc: mockHandleFunc,
						},
					},
					HandleFunc: nil,
				},
			},
		},

		// 测试根节点
		{
			name:    "",
			method:  http.MethodGet,
			path:    "/",
			isFound: true,
			matchNode: &matchNode{
				node: &node{
					path: "/",
					children: map[string]*node{
						"order": &node{
							path: "order",
							children: map[string]*node{
								"detail": &node{
									path:       "detail",
									children:   nil,
									HandleFunc: mockHandleFunc,
								},
							},
							HandleFunc: nil,
						},
					},
					HandleFunc: mockHandleFunc,
				},
			},
		},

		// 测试路由不存在的用例
		{
			name:      "path not found",
			method:    http.MethodGet,
			path:      "/user",
			isFound:   false,
			matchNode: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			foundNode, found := r.findRoute(testCase.method, testCase.path)
			// Tips: testCase.isFound是期望的结果,而found是实际的结果
			assert.Equal(t, testCase.isFound, found)

			// 没有找到路由就不用继续比较了
			if !found {
				return
			}

			// 此处和之前的测试一样 不能直接用assert.Equal()比较 因为HandleFunc不可比
			// 所以要用封装的node.equal()方法比较
			msg, found := testCase.matchNode.equal(foundNode)
			assert.True(t, found, msg)
		})
	}
}

// TestRouter_wildcard 测试通配符路由的注册功能
func TestRouter_wildcard(t *testing.T) {
	// step1. 构造路由树
	testRoutes := []TestNode{
		// 普通节点的通配符子节点测试用例
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		// 根节点的通配符子节点测试用例
		{
			method: http.MethodGet,
			path:   "/*",
		},
		// 通配符子节点的通配符子节点测试用例
		{
			method: http.MethodGet,
			path:   "/*/*",
		},
		// 通配符子节点的普通子节点测试用例
		{
			method: http.MethodGet,
			path:   "/*/*/order",
		},
		// 通配符子节点的普通子节点的通配符子节点
		{
			method: http.MethodGet,
			path:   "/*/*/order/*",
		},
	}

	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}

	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandleFunc)
	}

	// step2. 验证路由树 断言二者是否相等
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"order": {
						path:     "order",
						children: nil,
						wildcardChild: &node{
							path:          "*",
							children:      nil,
							wildcardChild: nil,
							HandleFunc:    mockHandleFunc,
						},
						HandleFunc: nil,
					},
				},
				wildcardChild: &node{
					path:     "*",
					children: nil,
					wildcardChild: &node{
						path: "*",
						children: map[string]*node{
							"order": {
								path:     "order",
								children: nil,
								wildcardChild: &node{
									path:          "*",
									children:      nil,
									wildcardChild: nil,
									HandleFunc:    mockHandleFunc,
								},
								HandleFunc: mockHandleFunc,
							},
						},
						wildcardChild: nil,
						HandleFunc:    mockHandleFunc,
					},
					HandleFunc: mockHandleFunc,
				},
				HandleFunc: nil,
			},
		},
	}

	msg, ok := wantRouter.equal(&r)
	assert.True(t, ok, msg)
}

// TestRouter_findRoute_wildcard 测试针对通配符路由的查找功能
func TestRouter_findRoute_wildcard(t *testing.T) {
	// step1. 构造路由树
	testRoutes := []TestNode{
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
	}

	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}

	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandleFunc)
	}

	// step2. 构造测试用例
	testCases := []struct {
		name      string
		method    string
		path      string
		isFound   bool
		matchNode *matchNode
	}{
		// 普通节点的通配符子节点测试用例
		{
			name:    "order wildcard",
			method:  http.MethodGet,
			path:    "/order/abc",
			isFound: true,
			matchNode: &matchNode{
				node: &node{
					path:          "*",
					typ:           nodeTypeAny,
					children:      nil,
					wildcardChild: nil,
					HandleFunc:    mockHandleFunc,
				},
			},
		},
		// 普通节点下普通子节点和通配符子节点共存的测试用例
		{
			name:    "order detail",
			method:  http.MethodGet,
			path:    "/order/detail",
			isFound: true,
			matchNode: &matchNode{
				node: &node{
					path:          "detail",
					typ:           nodeTypeStatic,
					children:      nil,
					wildcardChild: nil,
					HandleFunc:    mockHandleFunc,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			foundNode, found := r.findRoute(testCase.method, testCase.path)
			assert.Equal(t, testCase.isFound, found)

			if !found {
				return
			}

			msg, found := testCase.matchNode.equal(foundNode)
			assert.True(t, found, msg)
		})
	}
}

// TestRouter_addParamRoute 测试注册参数路由的结果是否符合预期
func TestRouter_addParamRoute(t *testing.T) {
	// step1. 构造路由树
	testRoutes := []TestNode{
		{
			method: http.MethodGet,
			path:   "/order/detail/:id",
		},
	}

	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}

	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandleFunc)
	}

	// step2. 验证路由树 断言二者是否相等
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:          "detail",
								children:      nil,
								wildcardChild: nil,
								paramChild: &node{
									path:          ":id",
									typ:           nodeTypeParam,
									children:      nil,
									wildcardChild: nil,
									paramChild:    nil,
									HandleFunc:    mockHandleFunc,
								},
								HandleFunc: nil,
							},
						},
						wildcardChild: nil,
						paramChild:    nil,
						HandleFunc:    nil,
					},
				},
				wildcardChild: nil,
				paramChild:    nil,
				HandleFunc:    nil,
			},
		},
	}

	msg, ok := wantRouter.equal(&r)
	assert.True(t, ok, msg)
}

// TestRouter_findRoute_param 测试针对参数路由的查找功能
func TestRouter_findRoute_param(t *testing.T) {
	// step1. 构造路由树
	testRoutes := []TestNode{
		{
			method: http.MethodGet,
			path:   "/order/detail/:id",
		},
	}

	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}

	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandleFunc)
	}

	// step2. 构造测试用例
	testCases := []struct {
		name      string
		method    string
		path      string
		isFound   bool
		matchNode *matchNode
	}{
		// 普通节点的参数路由子节点测试用例
		{
			name:    "order detail id",
			method:  http.MethodGet,
			path:    "/order/detail/123",
			isFound: true,
			matchNode: &matchNode{
				node: &node{
					path:          ":id",
					typ:           nodeTypeParam,
					children:      nil,
					wildcardChild: nil,
					paramChild:    nil,
					HandleFunc:    mockHandleFunc,
				},
				pathParams: map[string]string{
					"id": "123",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			foundNode, found := r.findRoute(testCase.method, testCase.path)
			assert.Equal(t, testCase.isFound, found)

			if !found {
				return
			}

			msg, found := testCase.matchNode.equal(foundNode)
			assert.True(t, found, msg)
		})
	}
}

// TestRouter_findRoute_param_and_reg_coexist 测试针对注册参数路由时,已有正则路由的情况
func TestRouter_findRoute_param_and_reg_coexist(t *testing.T) {
	// step1. 注册有冲突的路由
	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	r.addRoute(http.MethodGet, "/order/detail/:name(.+)", mockHandleFunc)

	// step2. 断言非法用例
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/order/detail/:id", mockHandleFunc)
	}, "web: 非法路由,已有正则子节点 .+ .不允许同时注册正则子节点与参数子节点")
}

// TestRouter_findRoute_param_and_param_coexist 测试针对注册参数路由时,已有参数路由的情况
func TestRouter_findRoute_param_and_param_coexist(t *testing.T) {
	// step1. 注册有冲突的路由
	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	r.addRoute(http.MethodGet, "/order/detail/:id", mockHandleFunc)

	// step2. 断言非法用例
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/order/detail/:name", mockHandleFunc)
	}, "web: 非法路由,已有参数子节点 :id .不允许同时注册多个参数子节点")
}

// TestRouter_findRoute_param_and_wildcard_coexist 测试针对注册参数路由时,已有通配符路由的情况
func TestRouter_findRoute_param_and_wildcard_coexist(t *testing.T) {
	// step1. 注册有冲突的路由
	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	r.addRoute(http.MethodGet, "/order/detail/*", mockHandleFunc)

	// step2. 断言非法用例
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/order/detail/:id", mockHandleFunc)
	}, "web: 非法路由,已有通配符子节点 * .不允许同时注册通配符子节点与参数子节点")
}

// TestRouter_findRoute_wildcard_and_reg_coexist 测试针对注册通配符路由时,已有正则路由的情况
func TestRouter_findRoute_wildcard_and_reg_coexist(t *testing.T) {
	// step1. 注册有冲突的路由
	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	r.addRoute(http.MethodGet, "/order/detail/:id([0-9]+)", mockHandleFunc)

	// step2. 断言非法用例
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/order/detail/*", mockHandleFunc)
	}, "web: 非法路由,已有正则子节点 [0-9]+ .不允许同时注册通配符子节点与正则子节点")
}

// TestRouter_findRoute_wildcard_and_param_coexist 测试针对注册通配符路由时,已有参数路由的情况
func TestRouter_findRoute_wildcard_and_param_coexist(t *testing.T) {
	// step1. 注册有冲突的路由
	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	r.addRoute(http.MethodGet, "/order/detail/:id", mockHandleFunc)

	// step2. 断言非法用例
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/order/detail/*", mockHandleFunc)
	}, "web: 非法路由,已有参数路由.不允许同时注册通配符路由和参数路由")
}

// TestRouter_findRoute_wildcard_and_wildcard_coexist 测试针对注册通配符路由时,已有通配符路由的情况
func TestRouter_findRoute_wildcard_and_wildcard_coexist(t *testing.T) {
	// step1. 注册有冲突的路由
	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	r.addRoute(http.MethodGet, "/order/detail/*", mockHandleFunc)

	// step2. 断言非法用例
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/order/detail/*", mockHandleFunc)
	}, "web: 非法路由,已有通配符子节点 * .不允许同时注册多个通配符子节点")
}

// TestRouter_findRoute_multistage_wildcard 测试在注册路由时,通配符出现在末尾的情况下,通配符匹配多段路由的情况
func TestRouter_findRoute_multistage_wildcard(t *testing.T) {
	// step1. 注册路由
	testRoutes := []struct {
		method string // method HTTP方法
		path   string // path 注册的路由字面量
	}{
		{
			method: http.MethodGet,
			path:   "/user/detail/*",
		},
		{
			method: http.MethodGet,
			path:   "/user/*/create",
		},
		{
			method: http.MethodGet,
			path:   "/*/order/show",
		},
	}

	mockHandle := func(ctx *Context) {}
	r := newRouter()
	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandle)
	}

	// step2. 构造测试用例
	testCases := []struct {
		name      string // name 测试用例名称
		method    string // method HTTP方法
		path      string // path 待匹配的路由字面量
		isFound   bool
		matchNode *matchNode
	}{
		{
			name:    "multi stage wildcard",
			method:  http.MethodGet,
			path:    "/user/detail/id/login",
			isFound: true,
			matchNode: &matchNode{
				node: &node{
					path:          "*",
					typ:           nodeTypeAny,
					children:      nil,
					wildcardChild: nil,
					paramChild:    nil,
					HandleFunc:    mockHandle,
				},
				pathParams: nil,
			},
		},
		{
			name:    "wildcard in middle",
			method:  http.MethodGet,
			path:    "/user/id/create",
			isFound: true,
			matchNode: &matchNode{
				node: &node{
					path:          "create",
					typ:           nodeTypeStatic,
					children:      nil,
					wildcardChild: nil,
					paramChild:    nil,
					HandleFunc:    mockHandle,
				},
				pathParams: nil,
			},
		},
		{
			name:    "wildcard in head",
			method:  http.MethodGet,
			path:    "/employee/order/show",
			isFound: true,
			matchNode: &matchNode{
				node: &node{
					path:          "show",
					typ:           nodeTypeStatic,
					children:      nil,
					wildcardChild: nil,
					paramChild:    nil,
					HandleFunc:    mockHandle,
				},
				pathParams: nil,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			foundNode, found := r.findRoute(testCase.method, testCase.path)
			assert.Equal(t, testCase.isFound, found)

			if !found {
				return
			}

			msg, found := testCase.matchNode.equal(foundNode)
			assert.True(t, found, msg)
		})
	}
}

// TestRouter_addRoute_reg 测试注册正则路由的结果是否符合预期
func TestRouter_addRoute_reg(t *testing.T) {
	// step1. 构造路由树
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/reg/:id([0-9]+)",
		},
		{
			method: http.MethodGet,
			path:   "/:name(^.+$)/abc",
		},
	}

	mockHandle := func(ctx *Context) {}
	r := newRouter()
	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandle)
	}

	// step2. 验证路由树
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				typ:  nodeTypeStatic,
				path: "/",
				children: map[string]*node{
					"reg": &node{
						typ:           nodeTypeStatic,
						path:          "reg",
						children:      nil,
						wildcardChild: nil,
						paramChild:    nil,
						regChild: &node{
							typ:           nodeTypeReg,
							path:          ":id([0-9]+)",
							children:      nil,
							wildcardChild: nil,
							paramChild:    nil,
							regChild:      nil,
							regExp:        regexp.MustCompile("[0-9]+"),
							HandleFunc:    mockHandle,
						},
						regExp:     nil,
						HandleFunc: nil,
					},
				},
				wildcardChild: nil,
				paramChild:    nil,
				regChild: &node{
					typ:  nodeTypeReg,
					path: ":name(^.+$)",
					children: map[string]*node{
						"abc": &node{
							typ:           nodeTypeStatic,
							path:          "abc",
							children:      nil,
							wildcardChild: nil,
							paramChild:    nil,
							regChild:      nil,
							regExp:        nil,
							HandleFunc:    mockHandle,
						},
					},
					wildcardChild: nil,
					paramChild:    nil,
					regChild:      nil,
					regExp:        regexp.MustCompile("^.+$"),
					HandleFunc:    nil,
				},
				regExp:     nil,
				HandleFunc: nil,
			},
		},
	}

	msg, ok := wantRouter.equal(&r)
	assert.True(t, ok, msg)
}

// TestRouter_findRoute_reg 测试针对正则路由的查找功能
func TestRouter_findRoute_reg(t *testing.T) {
	// step1. 构造路由树
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/reg/:id([0-9]+)",
		},
		{
			method: http.MethodGet,
			path:   "/user/:name(^.+$)",
		},
	}

	mockHandle := func(ctx *Context) {}
	r := newRouter()
	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandle)
	}

	// step2. 构造测试用例
	testCases := []struct {
		name      string
		method    string
		path      string
		isFound   bool
		matchNode *matchNode
	}{
		{
			name:    "reg in tail",
			method:  http.MethodGet,
			path:    "/reg/123",
			isFound: true,
			matchNode: &matchNode{
				node: &node{
					typ:           nodeTypeReg,
					path:          ":id([0-9]+)",
					children:      nil,
					wildcardChild: nil,
					paramChild:    nil,
					regChild:      nil,
					regExp:        regexp.MustCompile("[0-9]+"),
					HandleFunc:    mockHandle,
				},
				pathParams: map[string]string{
					"id": "123",
				},
			},
		},
		{
			name:    "reg in middle",
			method:  http.MethodGet,
			path:    "/user/peter",
			isFound: true,
			matchNode: &matchNode{
				node: &node{
					typ:           nodeTypeReg,
					path:          ":name(^.+$)",
					children:      nil,
					wildcardChild: nil,
					paramChild:    nil,
					regChild:      nil,
					regExp:        regexp.MustCompile("^.+$"),
					HandleFunc:    mockHandle,
				},
				pathParams: map[string]string{
					"name": "peter",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			foundNode, found := r.findRoute(testCase.method, testCase.path)
			assert.Equal(t, testCase.isFound, found)

			if !found {
				return
			}

			msg, found := testCase.matchNode.equal(foundNode)
			assert.True(t, found, msg)
		})
	}
}

// TestRouter_findRoute_reg_and_reg_coexist 测试针对注册正则路由时,已有正则路由的情况
func TestRouter_findRoute_reg_and_reg_coexist(t *testing.T) {
	// step1. 注册有冲突的路由
	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	r.addRoute(http.MethodGet, "/order/detail/:id([0-9]+)", mockHandleFunc)

	// step2. 断言非法用例
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/order/detail/:name(.+)", mockHandleFunc)
	}, "web: 非法路由,已有正则子节点 [0-9]+ .不允许同时注册通配符子节点与正则子节点")
}

// TestRouter_findRoute_reg_and_param_coexist 测试针对注册正则路由时,已有参数路由的情况
func TestRouter_findRoute_reg_and_param_coexist(t *testing.T) {
	// step1. 注册有冲突的路由
	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	r.addRoute(http.MethodGet, "/order/detail/:id", mockHandleFunc)

	// step2. 断言非法用例
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/order/detail/:name(.+)", mockHandleFunc)
	}, "web: 非法路由,已有参数子节点 :id .不允许同时注册通配符子节点与参数子节点")
}

// TestRouter_findRoute_reg_and_wildcard_coexist 测试针对注册正则路由时,已有通配符路由的情况
func TestRouter_findRoute_reg_and_wildcard_coexist(t *testing.T) {
	// step1. 注册有冲突的路由
	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	r.addRoute(http.MethodGet, "/order/detail/*", mockHandleFunc)

	// step2. 断言非法用例
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/order/detail/:name(.+)", mockHandleFunc)
	}, "web: 非法路由,已有通配符子节点 * .不允许同时注册多个通配符子节点")
}

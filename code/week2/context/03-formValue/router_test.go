package formValue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
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
			msg, found := testCase.matchNode.node.equal(foundNode.node)
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

			msg, found := testCase.matchNode.node.equal(foundNode.node)
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

			msg, found := testCase.matchNode.node.equal(foundNode.node)
			assert.True(t, found, msg)
		})
	}
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
	}, "web: 非法路由,已有通配符路由.不允许同时注册通配符路由和参数路由")
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

// TestRouter_findRoute_same_param_coexist 测试针对参数路由时,已有同名参数路由的情况
func TestRouter_findRoute_same_param_coexist(t *testing.T) {
	// step1. 注册有冲突的路由
	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	r.addRoute(http.MethodGet, "/order/detail/:id", mockHandleFunc)

	// step2. 断言非法用例
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/order/detail/:name", mockHandleFunc)
	}, "web: 路由冲突,参数路由冲突.已存在路由 id")
}

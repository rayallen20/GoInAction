# WEEK1作业

## PART1. 增强通配符匹配

### 1.1 功能需求

现在的通配符匹配只能支持1段,即:`/user/id/*`可以匹配到`/user/id/login`和`/user/id/detail`等,无法匹配到`/user/id/login/detail`这种多段的路由.

需要支持多段匹配,即`/user/id/*`可以匹配到`/user/id/login/detail`这种多段的路由.

但如果通配符出现在中间,例如`/user/*/detail`这种情况,则只匹配1段.即:`/user/*/detail`只能匹配到`/user/login/detail`或`/user/id/detail`,不能匹配到`/user/login/id/detail`.

### 1.2 测试用例

```go
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

			msg, found := testCase.matchNode.node.equal(foundNode.node)
			assert.True(t, found, msg)
		})
	}
}
```

### 1.3 实现

#### 1.3.1 实现思路

在查找路由树时:

- step1. 若未匹配到路由,则**尝试**在路由树中查找**节点类型为叶子节点**的通配符子节点
  - 叶子节点需满足的条件:无任何类型的子节点

#### 1.3.2 实现代码

```go
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
			// 若未匹配到节点 且 当前节点为通配符节点
			// 且 当前节点为叶子节点 则返回该叶子节点
			if target.path == "*" && target.children == nil && target.paramChild == nil && target.wildcardChild == nil {
				targetMatchNode.node = target
				return targetMatchNode, true
			}

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
```

## PART2. 实现正则匹配

### 2.1 功能需求

#### 2.1.1 基本功能

在路由注册时,支持正则表达式作为路由路径的一部分,例如:

- 注册路由: `/user/:id(^[0-9]+$)`,则该路由可以匹配到`/user/123`或`/user/456`
- 注册路由: `/order/:name(.+)`,则该路由可以匹配到`/order/abc`或`/order/def`

#### 2.1.2 复合场景

即:同时注册2个正则路由.例如:

- 注册路由1: `/user/:id(^[0-9]+$)`
- 注册路由2: `/user/:name(.+)`

则注册路由2时需要判断是否与路由1冲突,如果冲突,则直接panic.

理由:因为在查找路由时,如果请求路径为`/user/123`,则无法判断该请求路径应该匹配到路由1还是路由2.换言之:**开发者不应该设计这种路由**.

#### 2.1.3 路由优先级

静态匹配 > 正则匹配 > 路径参数 > 通配符匹配

这种优先级基本符合最左最精准原则.之所以正则匹配的优先级比路径参数的优先级高,是因为路径参数是正则匹配的一个特例,即`:id(.+)`(换言之,即`:参数名(任意字符出现任意多次)`).因此,其他的正则匹配只会比路径参数更精准,或者说更具象.

### 2.2 创建正则子节点

#### 2.2.1 测试用例

```go
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
```

#### 2.2.2 实现

##### a. 修改`node`结构体

**注:这一步参考了老师给的结构**

```go
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
```

##### b. 修改`childOrCreate()`方法

- step1. 判断当前给定的路由段是否表示正则路由段
  - 判断条件: 路由段以`:`开头 且包含`(` 且以`)`结尾,即为正则路由段
- step2. 如果是正则路由段,则提取出正则表达式,并编译正则表达式

```go
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

// extractReg 从给定的路由段中提取出正则表达式
func (n *node) extractReg(segment string) (regStr string, err error) {
	startIndex := strings.Index(segment, "(")
	endIndex := strings.LastIndex(segment, ")")
	if startIndex != -1 && endIndex != -1 && endIndex > startIndex+1 {
		return segment[startIndex+1 : endIndex], nil
	}

	return "", fmt.Errorf("web: 非法路由,正则路由格式错误.路由段 %s", segment)
}
```

##### c. 修改`node.equal()`方法

- 新增了对正则子节点的比对
- 新增了对节点类型的比对

```go
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
```

### 2.3 匹配正则子节点

#### 2.3.1 测试用例

```go
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

			msg, found := testCase.matchNode.node.equal(foundNode.node)
			assert.True(t, found, msg)

			// 比对参数是否相同
			assert.Equal(t, testCase.matchNode.pathParams, foundNode.pathParams)
		})
	}
}
```

#### 2.3.2 实现

只需在`childOf()`方法中,在查找参数子节点和通配符子节点之前,先查找正则子节点即可.

```go
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
```

## PART3. 整理代码

结合老师给出的`node`结构体,需要对代码做如下调整:

### 3.1 注册路由

现在仅仅是对非法的路由(例如路由是空字符串、路由中有多个连续的`/`等情况)做了处理,但对路由冲突的处理不够全面.

注册路由时,判断是否冲突的依据:同一个路由段上只能注册路径参数、通配符路由和正则路由中的一个

### 3.2 查找路由

`node.childOf()`方法不再返回`isParamChild`字段,因为`node`新增了`typ`字段.

因此在`router.findRoute()`方法中,只需对`node.childOf()`方法返回的`node`的`typ`字段做判断即可得知是否需要加参数

### 3.3 测试用例

`matchNode`结构体也应该有一个`equal()`方法,该方法除了判断`node`字段是否相等外,还应该判断`pathParams`字段是否相等

### 3.4 实现

#### 3.4.1 注册路由

```go
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
```

#### 3.4.2 查找路由

- `node.childOf()`: 不再返回`isParamChild`字段

```go
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
```

- `router.findRoute()`: 对`node.childOf()`方法返回的`node`的`typ`字段做判断,判断是否需要加参数

```go
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
		child, found := target.childOf(segment)
		if !found {
			// 若未匹配到节点 且 当前节点为通配符节点
			// 且 当前节点为叶子节点(这意味着注册时通配符子节点是最后一段路由段) 则返回该叶子节点
			if target.typ == nodeTypeAny && target.children == nil && target.paramChild == nil && target.wildcardChild == nil {
				targetMatchNode.node = target
				return targetMatchNode, true
			}

			return nil, false
		}

		// 若当前节点为正则节点,则将参数名和参数值保存到targetMatchNode中
		if child.typ == nodeTypeReg {
			// 参数名是 :id(正则表达式) 中的id 此处判断为从`:`开始到`(`结束的字符串
			name := child.path[1:strings.Index(child.path, "(")]
			// 参数值就是当前路由路径中的路由段
			value := segment
			targetMatchNode.addPathParams(name, value)
		}

		// 若当前节点为参数节点,则将参数名和参数值保存到targetMatchNode中
		if child.typ == nodeTypeParam {
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
```

#### 3.4.3 测试用例

##### a. 新增`matchNode.equal()`方法

```go
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
```

##### b. 针对参数路由的路由冲突测试用例

```go
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
```

##### c. 针对通配符路由的路由冲突测试用例

```go
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
```

##### d. 针对正则路由的路由冲突测试用例

```go
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
```
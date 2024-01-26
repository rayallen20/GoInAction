package web

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_Context(t *testing.T) {
	server := &HTTPServer{router: newRouter()}

	handleFunc := func(c *Context) {
		id, err := c.PathValue("id").AsInt64()
		if err != nil {
			c.Resp.WriteHeader(http.StatusBadRequest)
			_, _ = c.Resp.Write([]byte("id输入不正确: " + err.Error()))
		}

		_, _ = c.Resp.Write([]byte(fmt.Sprintf("id: %d", id)))
	}

	server.GET("/order/:id", handleFunc)
	_ = server.Start(":8091")
}

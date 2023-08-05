// handler/handler.go

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	server_ctxt *gin.Engine
}

func NewApiHandler(router *gin.Engine) *ApiHandler {
	return &ApiHandler{
		server_ctxt: router,
	}
}

func (api *ApiHandler) RegisterApiHandlers() (int, error) {
	api.server_ctxt.GET("/api", api.Index)

	// go api.PublishRandomImages()

	return 1, nil
}

func (api *ApiHandler) Index(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello World")
}

// func (api *ApiHandler) PublishRandomImages() {

// }

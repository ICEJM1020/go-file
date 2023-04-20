package davserver

import (
	"go-file/common"
	"go-file/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

var handler *webdav.Handler

func InitWebDAV() {
	handler = &webdav.Handler{
		Prefix:     "/webdav",
		FileSystem: webdav.Dir(common.VideoServePath),
		LockSystem: webdav.NewMemLS(),
	}
}

func GetFromAPI(dav *gin.RouterGroup) {
	dav.Use(middleware.WebDAVAuth())
	{
		dav.Any("/*path", ServeWebDAV)
		dav.Any("", ServeWebDAV)
		dav.Handle("PROPFIND", "/*path", ServeWebDAV)
		dav.Handle("PROPFIND", "", ServeWebDAV)
		dav.Handle("MKCOL", "/*path", ServeWebDAV)
		dav.Handle("LOCK", "/*path", ServeWebDAV)
		dav.Handle("UNLOCK", "/*path", ServeWebDAV)
		dav.Handle("PROPPATCH", "/*path", ServeWebDAV)
		dav.Handle("COPY", "/*path", ServeWebDAV)
		dav.Handle("MOVE", "/*path", ServeWebDAV)
	}
}

func ServeWebDAV(c *gin.Context) {
	// user := c.MustGet("user".)
	// ctx := context.WithValue(c.Request.Context(), "user")
	handler.ServeHTTP(c.Writer, c.Request)
}

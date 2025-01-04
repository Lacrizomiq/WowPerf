package authboss

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinRouter adapts Gin's router to Authboss's router interface
type GinRouter struct {
	engine *gin.Engine
	group  *gin.RouterGroup
}

// NewRouter creates a new router that adapts Gin to Authboss
func NewRouter(engine *gin.Engine, mountPath string) *GinRouter {
	return &GinRouter{
		engine: engine,
		group:  engine.Group(mountPath),
	}
}

// Get registers a GET route
func (r *GinRouter) Get(path string, handler http.Handler) {
	r.group.GET(path, wrapHandler(handler))
}

// Post registers a POST route
func (r *GinRouter) Post(path string, handler http.Handler) {
	r.group.POST(path, wrapHandler(handler))
}

// Delete registers a DELETE route
func (r *GinRouter) Delete(path string, handler http.Handler) {
	r.group.DELETE(path, wrapHandler(handler))
}

// ServeHTTP implements http.Handler
func (r *GinRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}

// wrapHandler converts an http.Handler to a gin.HandlerFunc
func wrapHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

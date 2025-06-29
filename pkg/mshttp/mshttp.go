package mshttp

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/mushanyux/MSIMServer/pkg/cache"
	"github.com/mushanyux/MSIMServer/pkg/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type UserRole string

const (
	Admin      UserRole = "admin"
	SuperAdmin UserRole = "superAdmin"
)

type MSHttp struct {
	r    *gin.Engine
	pool sync.Pool
}

func New() *MSHttp {
	l := &MSHttp{
		r:    gin.New(),
		pool: sync.Pool{},
	}
	l.r.Use(gin.Recovery())
	l.pool.New = func() interface{} { return allocateContext() }
	return l
}

func allocateContext() *Context { return &Context{Context: nil, lg: log.NewTLog("context")} }

type Context struct {
	*gin.Context
	lg log.Log
}

func (c *Context) reset() {
	c.Context = nil
}

func (c *Context) ResponseError(err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"msg":    err.Error(),
		"status": http.StatusBadRequest,
	})
}

func (c *Context) ResponseErrorf(msg string, err error) {
	if err != nil {
		c.lg.Error(msg, zap.Error(err), zap.String("path", c.FullPath()))
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"msg":    msg,
		"status": http.StatusBadRequest,
	})
}

func (c *Context) ResponseErrorWithStatus(err error, status int) {
	c.JSON(http.StatusBadRequest, gin.H{
		"msg":    err.Error(),
		"status": status,
	})
}

func (c *Context) GetPage() (pageIndex int64, pageSize int64) {
	pageIndex, _ = strconv.ParseInt(c.Query("page_index"), 10, 64)
	pageSize, _ = strconv.ParseInt(c.Query("page_size"), 10, 64)
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 15
	}
	return
}

func (c *Context) ResponseOK() {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

func (c *Context) Response(data interface{}) {
	c.JSON(http.StatusOK, data)
}

func (c *Context) ResponseWithStatus(status int, data interface{}) {
	c.JSON(status, data)
}

func (c *Context) GetLoginUID() string {
	return c.MustGet("uid").(string)
}

func (c *Context) GetAppID() string {
	return c.GetHeader("appid")
}

func (c *Context) GetLoginName() string {
	return c.MustGet("name").(string)
}

func (c *Context) GetLoginRole() string {
	return c.GetString("role")
}

func (c *Context) GetSpanContext() opentracing.SpanContext {
	return c.MustGet("spanContext").(opentracing.SpanContext)
}

func (c *Context) CheckLoginRole() error {
	role := c.GetLoginRole()
	if role == "" {
		return errors.New("登录用户角色错误")
	}
	if role != string(Admin) && role != string(SuperAdmin) {
		return errors.New("该用户无权限")
	}
	return nil
}

func (c *Context) CheckLoginRoleIsSuperAdmin() error {
	role := c.GetLoginRole()
	if role == "" {
		return errors.New("登录用户角色错误")
	}
	if role != string(SuperAdmin) {
		return errors.New("该用户无权限")
	}
	return nil
}

type HandlerFunc func(c *Context)

func (l *MSHttp) MSHttpHandlerFunc(handlerFunc HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		hc := l.pool.Get().(*Context)
		hc.reset()
		hc.Context = c
		defer l.pool.Put(hc)

		handlerFunc(hc)
	}
}

func (l *MSHttp) Run(addr ...string) error {
	return l.r.Run(addr...)
}

func (l *MSHttp) RunTLS(addr, certFile, keyFile string) error {
	return l.r.RunTLS(addr, certFile, keyFile)
}

func (l *MSHttp) POST(relativePath string, handlers ...HandlerFunc) {
	l.r.POST(relativePath, l.handlersToGinHandleFuncs(handlers)...)
}

func (l *MSHttp) GET(relativePath string, handlers ...HandlerFunc) {
	l.r.GET(relativePath, l.handlersToGinHandleFuncs(handlers)...)
}

func (l *MSHttp) Any(relativePath string, handlers ...HandlerFunc) {
	l.r.Any(relativePath, l.handlersToGinHandleFuncs(handlers)...)
}

func (l *MSHttp) Static(relativePath string, root string) {
	l.r.Static(relativePath, root)
}

func (l *MSHttp) LoadHTMLGlob(pattern string) {
	l.r.LoadHTMLGlob(pattern)
}

func (l *MSHttp) UseGin(handlers ...gin.HandlerFunc) {
	l.r.Use(handlers...)
}

func (l *MSHttp) Use(handlers ...HandlerFunc) {
	l.r.Use(l.handlersToGinHandleFuncs(handlers)...)
}

func (l *MSHttp) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	l.r.ServeHTTP(w, req)
}

func (l *MSHttp) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return newRouterGroup(l.r.Group(relativePath, l.handlersToGinHandleFuncs(handlers)...), l)
}

func (l *MSHttp) HandleContext(c *Context) {
	l.r.HandleContext(c.Context)
}

func (l *MSHttp) handlersToGinHandleFuncs(handlers []HandlerFunc) []gin.HandlerFunc {
	newHandlers := make([]gin.HandlerFunc, 0, len(handlers))
	if handlers != nil {
		for _, handler := range handlers {
			newHandlers = append(newHandlers, l.MSHttpHandlerFunc(handler))
		}
	}
	return newHandlers
}

func (l *MSHttp) AuthMiddleware(cache cache.Cache, tokenPrefix string) HandlerFunc {
	return func(c *Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "token不能为空，请先登录！",
			})
			return
		}
		uidAndName := GetLoginUID(token, tokenPrefix, cache)
		if strings.TrimSpace(uidAndName) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "请先登录！",
			})
			return
		}
		uidAndNames := strings.Split(uidAndName, "@")
		if len(uidAndNames) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "token有误！",
			})
			return
		}
		c.Set("uid", uidAndNames[0])
		c.Set("name", uidAndNames[1])
		if len(uidAndNames) > 2 {
			c.Set("role", uidAndNames[2])
		}
		c.Next()
	}
}

func GetLoginUID(token string, tokenPrefix string, cache cache.Cache) string {
	uid, err := cache.Get(tokenPrefix + token)
	if err != nil {
		return ""
	}
	return uid
}

type RouterGroup struct {
	*gin.RouterGroup
	L *MSHttp
}

func newRouterGroup(g *gin.RouterGroup, l *MSHttp) *RouterGroup {
	return &RouterGroup{RouterGroup: g, L: l}
}

func (r *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.POST(relativePath, r.L.handlersToGinHandleFuncs(handlers)...)
}
func (r *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.GET(relativePath, r.L.handlersToGinHandleFuncs(handlers)...)
}
func (r *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.DELETE(relativePath, r.L.handlersToGinHandleFuncs(handlers)...)
}
func (r *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.PUT(relativePath, r.L.handlersToGinHandleFuncs(handlers)...)
}

func CORSMiddleware() HandlerFunc {
	return func(c *Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, token, accept, origin, Cache-Control, X-Requested-With, appid, noncestr, sign, timestamp")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

package router

import (
	"ai-ops/common"
	"ai-ops/common/errcode"
	"ai-ops/configs"
	_ "ai-ops/docs"
	"ai-ops/middleware"
	"ai-ops/pkg/apis/admin"
	"ai-ops/pkg/apis/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {

	r := gin.New()
	middleware.InitMiddleware(r)
	g := r.Group("")

	g.GET("/check.do", healthCheck)

	adminRouter(g)

	// swagger；注意：生产环境可以注释掉
	sysSwaggerRouter(g)

	return r
}

//健康检查
func healthCheck(c *gin.Context) {
	log.Println("health check success")
	c.String(http.StatusOK, "Ok!")
}

func overviewRouter(r *gin.RouterGroup) {
	v1 := r.Group("/api/v1/overview")
	{
		v1.GET("/grafana", overview.GetGrafana)
	}
}

func adminRouter(r *gin.RouterGroup) {
	v1 := r.Group("/api/v1/admin")
	{
		v1.GET("/department/list", admin.ListDepartment)
		v1.GET("/user/info", admin.QueryUserInfo)
		v1.GET("/user/info/by-account", admin.QueryUserByAccount)
	}
}

func sysSwaggerRouter(r *gin.RouterGroup) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

//用户ID转UserName
func HandleUserName(c *gin.Context) {
	userName := manage.GetUserNameById(common.UCUserId)
	log.Infof("RequestUser Userid:%s,UserName:%s", common.UCUserId, userName)
	if userName == "" {
		c.Abort()
		panic("CustomError#" + strconv.Itoa(errcode.ERROR_USERNAME_PERMISSION_DENY) + "#" + errcode.ErrorCodeMsg(errcode.ERROR_USERNAME_PERMISSION_DENY))
	} else {
		common.UCUserName = userName
		c.SetCookie("Username", userName, 3600, "/", config.ApplicationConfig.Domain, false, false)
	}

	c.Next()
}

//操作记录
func HandleAddLog(c *gin.Context) {
	c.Request.Header.Add(common.FLAG_LOGGING, common.FLAG_ENABLE)
}

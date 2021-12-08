package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetCorsRouters(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "lang", "token", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
}

//func clickToken() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		//请求处理
//		token := c.GetHeader("token")
//		if token == "" {
//			ResponseErr(c, hcode.UserNotLoginErr)
//			c.Abort()
//			return
//		}
//		t, err := auth.GetAuth().CheckUserToken(token)
//		if err != nil {
//			ResponseErr(c, err)
//			c.Abort()
//			return
//		}
//		c.Set("uid", t.UID)
//		c.Next()
//		//处理后
//	}
//}
//
//func clickAdminToken() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		//请求处理
//		token := c.GetHeader("token")
//		if token == "" {
//			ResponseErr(c, hcode.UserNotLoginErr)
//			c.Abort()
//			return
//		}
//		t, err := auth.GetAuth().CheckAdminToken(token)
//		if err != nil {
//			ResponseErr(c, err)
//			c.Abort()
//			return
//		}
//		c.Set("adminId", t)
//		c.Next()
//		//处理后
//	}
//}

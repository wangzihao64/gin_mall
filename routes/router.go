package routes

import (
	api "gin_mall/api/v1"
	"gin_mall/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.StaticFS("/static", http.Dir("static"))
	v1 := r.Group("/api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})
		//用户操作
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)
		authed := v1.Group("/") //需要登录保护
		authed.Use(middleware.JWT())
		{
			//用户操作
			authed.PUT("user", api.UserUpdate)
			authed.PUT("avatar", api.UploadAvatar)

			//轮播图
			v1.GET("carousels", api.ListCarousel)
			authed.POST("user/sending-email", api.SendEmail)
			authed.POST("user/valid-email", api.ValidEmail)

			//显示金额
			authed.POST("money", api.ShowMoney)

			//商品操作
			authed.POST("product", api.CreateProduct)
		}
	}
	return r
}

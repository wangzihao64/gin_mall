package v1

import (
	"gin_mall/pkg/util"
	"gin_mall/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建商品
func CreateProduct(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files"]
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	createProductService := service.ProductService{}
	if err := c.ShouldBind(&createProductService); err != nil {
		res := createProductService.Create(c.Request.Context(), claim.ID, files)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrousObj.Infoln(err)
	}
}

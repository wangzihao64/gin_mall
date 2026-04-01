package service

import (
	"context"
	"database/sql/driver"
	"gin_mall/dao"
	"gin_mall/model"
	"gin_mall/pkg/e"
	"gin_mall/serizlizer"
	"mime/multipart"
)

type ProductService struct {
	Id            uint   `json:"id" form:"id"`
	Name          string `json:"name" form:"name"`
	Category      uint   `json:"category_id" form:"category_id"`
	Title         string `json:"title" form:"title"`
	Info          string `json:"info" form:"info"`
	ImgPath       string `json:"img_path" form:"img_path"`
	Price         string `json:"price" form:"price"`
	DiscountPrice string `json:"discount_price" form:"discount_price"`
	Onsale        bool   `json:"onsale" form:"onsale"`
	Num           int    `json:"num" form:"num"`
	model.BasePage
}

func (service *ProductService) Create(ctx context.Context, uId uint, files []*multipart.FileHeader) serizlizer.Response {
	var boss *model.User
	var err error
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	boss, _ := userDao.GetUserById(uId)
	//以第一章作为封面图
	tmp, _ := files[0].Open()
	path, err := UploadProductToLocalStatic(tmp, uId, service.Name)
}

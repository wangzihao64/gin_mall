package service

import (
	"context"
	"gin_mall/dao"
	"gin_mall/pkg/e"
	"gin_mall/pkg/util"
	"gin_mall/serizlizer"
)

type CarouselService struct {
}

func (service *CarouselService) List(ctx context.Context) serizlizer.Response {
	carouselDao := dao.NewCarouselDao(ctx)
	code := e.Success
	carousels, err := carouselDao.ListCarousel()
	if err != nil {
		util.LogrousObj.Infoln("err", err)
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serizlizer.BuildListResponse(serizlizer.BuildCarousels(carousels), uint(len(carousels)))
}

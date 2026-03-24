package dao

import (
	"context"
	"gin_mall/model"

	"gorm.io/gorm"
)

type CarouselDao struct {
	*gorm.DB
}

func NewCarouselDao(ctx context.Context) *CarouselDao {
	return &CarouselDao{NewDBClient(ctx)}
}

func NewCarouselDaoByDB(db *gorm.DB) *CarouselDao {
	return &CarouselDao{db}
}

// ListCarousel
func (dao *CarouselDao) ListCarousel() (Carousel []model.Carousel, err error) {
	err = dao.DB.Model(&model.Carousel{}).Find(&Carousel).Error
	return
}

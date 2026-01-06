package service

import (
	"context"
	"gin_mall/dao"
	"gin_mall/model"
	"gin_mall/pkg/e"
	"gin_mall/pkg/util"
	"gin_mall/serizlizer"
)

type UserService struct {
	NickName string `json:"nick_name" form:"nick_name"`
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
	Key      string `json:"key" form:"key"` //密钥
}

func (service UserService) Register(ctx context.Context) serizlizer.Response {
	var user model.User
	code := e.Success
	if service.Key == "" || len(service.Key) != 16 {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "密钥长度不足",
		}
	}
	// 10000 --->  密文存储对称加密操作
	util.Encrypt.SetKey(service.Key)
	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if exist {
		code = e.ErrorExistUser
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	user = model.User{
		Username: service.UserName,
		NickName: service.NickName,
		Status:   model.ActiveString,
		Avatar:   "avatar.JPG",
		Money:    util.Encrypt.AesEncoding("10000"), //初始金额的加密
	}
	//密码加密
	if err = user.SetPassword(service.Password); err != nil {
		code = e.ErrorFailEncryption
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//创建用户
	err = userDao.CreateUser(&user)
	if err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
func (service *UserService) Login(ctx context.Context) serizlizer.Response {
	var user *model.User
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	//db := dao.NewDBClient(ctx)

	//err := db.Model(&model.User{}).Where("username = ?", service.UserName).First(&user).Error
	//if err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return serizlizer.Response{}
	//	}
	//	return serizlizer.Response{}
	//}

	//判断用户是否存在
	user, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if !exist || err != nil {
		code = e.ErrorExistUserNotFound
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户不存在，请先注册",
		}
	}
	//校验密码
	if user.CheckPassword(service.Password) == false {
		code = e.ErrorNotCompare
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密码错误，请重新登陆",
		}
	}
	//http是一个无状态的协议（认证,token）
	//token签发
	token, err := util.GenerateToken(user.ID, service.UserName, 0)
	if err != nil {
		code = e.ErrorAuthToken
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "token error",
		}
	}
	return serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data: serizlizer.TokenData{
			Token: token,
			User:  serizlizer.BuildUser(user),
		},
	}
}

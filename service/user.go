package service

import (
	"context"
	"gin_mall/conf"
	"gin_mall/dao"
	"gin_mall/model"
	"gin_mall/pkg/e"
	"gin_mall/pkg/util"
	"gin_mall/serizlizer"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"gopkg.in/mail.v2"
)

type UserService struct {
	NickName string `json:"nick_name" form:"nick_name"`
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
	Key      string `json:"key" form:"key"` //密钥
}

type SendEmailService struct {
	Email         string `json:"email" form:"email"`
	Password      string `json:"password" form:"password"`
	OperationType uint   `json:"operation_type" form:"operation_type"`
	//1.绑定邮箱 2.解绑邮箱 3.改密码
}
type ValidEmailService struct {
}

type ShowMoneyService struct {
	Key string `json:"key" form:"key"`
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

// Update用户信息修改
func (service *UserService) Update(ctx context.Context, uId uint) serizlizer.Response {
	var user *model.User
	var err error
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	user, err = userDao.GetUserById(uId)
	//修改昵称nickname
	if service.NickName != "" {
		user.NickName = service.NickName
	}
	err = userDao.UpdateUserById(uId, user)
	if err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serizlizer.BuildUser(user),
	}
}

// 头像更新
func (service *UserService) Post(ctx context.Context, uId uint, file multipart.File, fileSize int64) serizlizer.Response {
	code := e.Success
	var user *model.User
	var err error
	UserDao := dao.NewUserDao(ctx)
	user, err = UserDao.GetUserById(uId)
	if err != nil {
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//保存图片到本地函数
	path, err := UploadAvatarToLocalStatic(file, uId, user.Username)
	if err != nil {
		code = e.ErrorUploadFail
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	user.Avatar = path
	err = UserDao.UpdateUserById(uId, user)
	if err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serizlizer.BuildUser(user),
	}
}

// 发送邮箱
func (service *SendEmailService) Send(ctx context.Context, uId uint) serizlizer.Response {
	code := e.Success
	var address string
	var notice *model.Notice //绑定邮箱，修改密码 模版通知
	token, err := util.GenerateEmailToken(uId, service.OperationType, service.Email, service.Password)
	if err != nil {
		code = e.ErrorAuthToken
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	noticeDao := dao.NewNoticeDao(ctx)
	notice, err = noticeDao.GetNoticeById(service.OperationType)
	if err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	address = conf.VaildEmail + token //发送方
	mailStr := notice.Text
	mailText := strings.Replace(mailStr, "email", address, -1)
	log.Println(address)
	log.Println(mailText)
	m := mail.NewMessage()
	m.SetHeader("From", conf.SmtpEmail)
	m.SetHeader("To", service.Email)
	m.SetHeader("Subject", "FanOne")
	m.SetBody("text/html", mailText)
	d := mail.NewDialer(conf.SmtpHost, 465, conf.SmtpEmail, conf.SmtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err = d.DialAndSend(m); err != nil {
		code = e.ErrorSendEmail
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *ValidEmailService) Valid(ctx context.Context, token string) serizlizer.Response {
	var userId uint
	var email string
	var password string
	var operationType uint
	code := e.Success
	//验证token
	if token == "" {
		code = e.InvalidParams
	} else {
		claims, err := util.ParseEmailToken(token)
		if err != nil {
			code = e.ErrorAuthToken
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ErrorAuthCheckTokenTimeout
		} else {
			userId = claims.UserID
			email = claims.Email
			password = claims.Password
			operationType = claims.OperationType
		}
	}
	if code != e.Success {
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//获取该用户的信息
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(userId)
	if err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if operationType == 1 {
		//绑定邮箱
		user.Email = email
	} else if operationType == 2 {
		user.Email = ""
	} else if operationType == 3 {
		err = user.SetPassword(password)
		if err != nil {
			code = e.Error
			return serizlizer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
	}
	err = userDao.UpdateUserById(userId, user)
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
		Data:   serizlizer.BuildUser(user),
	}
}

// 展示用户金额
func (service *ShowMoneyService) Show(ctx context.Context, uId uint) serizlizer.Response {
	code := e.Success
	UserDao := dao.NewUserDao(ctx)
	user, err := UserDao.GetUserById(uId)
	if err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serizlizer.Response{
		Status: code,
		Data:   serizlizer.BuildMoney(user, service.Key),
		Msg:    e.GetMsg(code),
	}
}

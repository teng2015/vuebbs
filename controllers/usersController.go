package controllers

import (
	"gopkg.in/baa.v1"
	"../db"
	"../dto/returnStatus"
	"../entity"
	"../utls/sessionUtls"
	"golang.org/x/tools/go/gcimporter15/testdata"
)

type usersController struct{}

var UsersController = usersController{}

/**
 通过板块id获取列表
 */
func (this *usersController) GetCount(b *baa.Context) {
	userCount := db.UserInfoDAL.GetCount()

	b.JSON(200, returnStatus.ReturnStatus{returnStatus.Ok, userCount})
}

/**
 获取当前用户
 */
func (this *usersController) GetCurrent(b *baa.Context) {
	user := sessionUtls.SessionUtls.GetCurrentUser(b)

	b.JSON(200, returnStatus.ReturnStatus{returnStatus.Ok, user})
}


/**
 注册用户信息
 */
func (this *usersController) Regist(b *baa.Context) {
	user := &entity.UserInfo{}
	user.LoginName = b.Req.PostFormValue("LoginName")
	user.LoginPwd = b.Req.PostFormValue("LoginPwd")
	user.EmailAddr = b.Req.PostFormValue("EmailAddr")
	user.NickName = b.Req.PostFormValue("NickName")


	//判断参数
	if user.CheckRegistFiled() {
		b.JSON(200, returnStatus.ReturnStatus{Status:returnStatus.ParamError, Data:nil })
	} else if db.UserInfoDAL.EmailOrLoginNameExist(*user) {
		//查看用户的登录名或者邮箱是否存在
		b.JSON(200, returnStatus.ReturnStatus{Status:returnStatus.UserExist, Data:nil })
	} else {
		//执行添加操作
		db.UserInfoDAL.Regist(user)
		if user.Id <= 0 {
			b.JSON(200, returnStatus.ReturnStatus{Status:returnStatus.Error, Data:nil })    //系统错误
		}

		//添加到session
		sessionUtls.SessionUtls.SetCurrentUser(user.Id, b)

		//返回用户信息
		b.JSON(200, returnStatus.ReturnStatus{Status:returnStatus.Ok, Data:user})
	}
}

//用户登录方法
func (this *usersController) Login(b *baa.Context) {
	user := &entity.UserInfo{}
	user.LoginName = b.Req.PostFormValue("LoginName")
	user.LoginPwd = b.Req.PostFormValue("LoginPwd")
	user.EmailAddr = b.Req.PostFormValue("EmailAddr")

	if (user.LoginName == "" && user.EmailAddr == "") || user.LoginPwd == "" {
		b.JSON(200, returnStatus.ReturnStatus{Status:returnStatus.ParamError})
	} else {
		userRet, err := db.UserInfoDAL.Login(user.LoginName, user.EmailAddr, user.LoginPwd)

		//登录过程出现错误
		if err != nil {
			b.JSON(200, returnStatus.ReturnStatus{Status:returnStatus.Error, Data:err })
			return
		}

		//用户未找到
		if userRet == nil || userRet.Id <= 0 {
			b.JSON(200, returnStatus.ReturnStatus{Status:returnStatus.UserNotExist })
			return
		}

		//登录成功  存入用户数据
		sessionUtls.SessionUtls.SetCurrentUser(user.Id, b)

		b.JSON(200, returnStatus.ReturnStatus{Status:returnStatus.Ok, Data:userRet})
	}
}
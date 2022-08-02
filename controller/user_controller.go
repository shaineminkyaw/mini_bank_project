package controller

import (
	"crypto/rsa"
	"fmt"
	"math/rand"
	"miniproject/config"
	"miniproject/ds"
	"miniproject/dto"
	"miniproject/model"
	"miniproject/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mazen160/go-random"
	"gorm.io/gorm"
)

type userController struct {
	H       *Handler
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

func NewUserController(h *Handler) *userController {
	return &userController{
		H:       h,
		Private: config.PrivateKey,
		Public:  config.PublicKey,
	}
}

func (ctr *userController) Register() {
	h := ctr.H
	group := h.Router.Group("api/user")
	group.GET("/all", ctr.list)
	group.POST("verifyCode", ctr.verifyCode)
	group.POST("/register", ctr.register)
	group.POST("/login", ctr.login)

	group.Use(AuthMiddleware())
	group.POST("/logout", ctr.logout)
	group.POST("/updateEmail", ctr.updateEmail)
	group.POST("/updatePassword", ctr.updatePassword)
	// h.Router.Use(midd)
}

//@@@user list
func (ctr *userController) list(c *gin.Context) {
	//
	resp := dto.ResObj{}
	user := &model.User{}
	users, err := ctr.H.UserService.FindAll(user)
	if err != nil {
		resp.ErrCode = 500
		resp.ErrMsg = "Something went wrong"
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.ErrCode = 0
	resp.ErrMsg = "success"
	resp.Data = users
	c.JSON(http.StatusOK, resp)
}

//@@@@verify code
func (ctr *userController) verifyCode(c *gin.Context) {
	//

	res := dto.ResObj{}
	req := dto.ReqVerifyCode{}
	vUser := &model.VerifyCode{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.ErrCode = 400
		res.ErrMsg = "wrong param :" + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	number := time.Now().UnixNano()
	rand.Seed(number)
	code := fmt.Sprintf("%v%v%v%v", rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10))
	expTime := time.Now().Add(time.Minute * 30)
	mail := req.Email
	cUser := &model.VerifyCode{
		Email:      mail,
		Code:       code,
		ExpireTime: expTime,
	}

	err = ds.DB.Model(&model.VerifyCode{}).Where("email = ?", mail).First(&vUser).Error
	if err == gorm.ErrRecordNotFound {
		err = ds.DB.Model(&model.VerifyCode{}).Create(&cUser).Error
		if err != nil {
			res.ErrCode = 403
			res.ErrMsg = "wrong on create verify code" + err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
	} else {
		data := ds.DB.Model(&model.VerifyCode{}).Where("email = ?", req.Email)
		err = data.Updates(&cUser).Error
		if err != nil {
			res.ErrCode = 403
			res.ErrMsg = "wrong on updates verify code" + err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		res.ErrCode = 408
		res.ErrMsg = "Something went wrong"
		c.JSON(http.StatusOK, res)
		return
	}
	res.ErrCode = 0
	res.ErrMsg = "success"
	res.Data = gin.H{
		"code": code,
	}
	c.JSON(http.StatusOK, res)

}

//@@@@register
func (ctr *userController) register(c *gin.Context) {
	//
	h := ctr.H
	res := dto.ResObj{}
	req := dto.ReqUserRegister{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.ErrCode = 403
		res.ErrMsg = "wrong param :" + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	//@@@validate verifycode
	vUser, err := ctr.H.UserService.GetVerifyUser(req.Email)
	if err != nil {
		res.ErrCode = 407
		res.ErrMsg = "wrong on verify code " + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	if vUser == nil {
		res.ErrCode = 408
		res.ErrMsg = "nothing verify code " + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	if req.VerifyCode != vUser.Code || time.Now().Unix() > vUser.ExpireTime.UnixNano() {
		res.ErrCode = 409
		res.ErrMsg = "verify code needs update" + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	//hash password
	userPassword, err := util.HashPassword(req.Password)
	if err != nil {
		res.ErrCode = 8000
		res.ErrMsg = "error on hash password :" + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	charset := "12345678"
	length := 8
	usr, _ := random.Random(length, charset, true)
	userName := "U" + usr
	currency := "USD"
	bankCard, err := util.GetBankCardNumber(req.City)
	if err != nil {
		res.ErrCode = 700
		res.ErrMsg = "error on bank card generate"
		c.JSON(http.StatusOK, res)
		return
	}
	user := &model.User{
		Username:       userName,
		Password:       userPassword,
		Email:          req.Email,
		RegisterIP:     c.ClientIP(),
		LastLoginIP:    c.ClientIP(),
		NationID:       req.NationID,
		BankCardNumber: bankCard,
		City:           req.City,
		Balance:        0,
		Currency:       currency,
		Type:           req.Type,
	}
	_, err = h.UserService.FindByEmail(req.Email)
	if err == gorm.ErrRecordNotFound {
		err := ds.DB.Model(&model.User{}).Create(&user).Error
		if err != nil {
			res.ErrCode = 400
			res.ErrMsg = "error on user create :" + err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
		err = util.SaveUserBankCard(req.City, user.Id, bankCard)
		if err != nil {
			res.ErrCode = 800
			res.ErrMsg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
	} else {
		res.ErrCode = 407
		res.ErrMsg = "User Already exist"
		c.JSON(http.StatusOK, res)
		return
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		res.ErrCode = 403
		res.ErrMsg = "Something went wrong " + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	res.ErrCode = 0
	res.ErrMsg = "success"
	res.Data = gin.H{
		"token": "",
		"data":  user,
	}
	c.JSON(http.StatusOK, res)

}

//@@@login
func (ctr *userController) login(c *gin.Context) {
	//
	res := dto.ResObj{}
	req := dto.ReqLogin{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.ErrCode = 403
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	user := &model.User{}
	err = ds.DB.Model(&model.User{}).Where("email = ?", req.Email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		res.ErrCode = 409
		res.ErrMsg = "User not found"
		c.JSON(http.StatusOK, res)
		return
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		res.ErrCode = 500
		res.ErrMsg = "Something went wrong"
		c.JSON(http.StatusOK, res)
		return
	}
	ok, err := util.ValidateHashedPassword(user.Password, req.Password)
	if err != nil {
		res.ErrCode = 420
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	if !ok {
		res.ErrCode = 430
		res.ErrMsg = "Password not match"
		c.JSON(http.StatusOK, res)
		return
	}
	token, err := ctr.H.TokenService.NewTokenPair(user, "")
	if err != nil {
		res.ErrCode = 407
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	res.ErrCode = 0
	res.ErrMsg = "success"
	res.Data = gin.H{
		"data":  user,
		"token": token.AccessToken.SS,
	}
	c.JSON(http.StatusOK, res)
}

//@@@logout
func (ctr *userController) logout(c *gin.Context) {
	//
	res := dto.ResObj{}
	token := &model.UserToken{}
	user := c.MustGet("user").(*model.User)
	db := ds.DB.Model(&model.UserToken{}).Where("uid = ?", user.Id)
	err := db.Delete(&token).Error
	if err != nil {
		res.ErrCode = 409
		res.ErrMsg = "something went wrong"
		c.JSON(http.StatusOK, res)
		return
	}
	res.ErrCode = 0
	res.ErrMsg = "success"
	c.JSON(http.StatusOK, res)

}

//@@@update Email
func (ctr *userController) updateEmail(c *gin.Context) {
	//
	res := dto.ResObj{}
	type reqNewEmail struct {
		Email string `json:"email" form:"email" binding:"required"`
	}
	req := reqNewEmail{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.ErrCode = 403
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	loginUser := c.MustGet("user").(*model.User)
	db := ds.DB.Model(&model.User{})
	db = db.Where("id = ?", loginUser.Id)
	err = db.Update("email", req.Email).Error
	if err != nil {
		res.ErrCode = 409
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	res.ErrCode = 0
	res.ErrMsg = "success"
	c.JSON(http.StatusOK, res)
}

// @@@ update Password
func (ctr *userController) updatePassword(c *gin.Context) {
	//
	res := dto.ResObj{}
	user := c.MustGet("user").(*model.User)
	type reqUpdatePassword struct {
		Password string `json:"password" form:"password" binding:"required,min=6"`
	}
	req := reqUpdatePassword{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.ErrCode = 403
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	pass, err := util.HashPassword(req.Password)
	if err != nil {
		res.ErrCode = 607
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	db := ds.DB.Model(&model.User{})
	db = db.Where("id = ?", user.Id)
	err = db.Update("password", pass).Error
	if err != nil {
		res.ErrCode = 608
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	res.ErrCode = 0
	res.ErrMsg = "success"
	c.JSON(http.StatusOK, res)

}

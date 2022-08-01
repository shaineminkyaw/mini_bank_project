package controller

import (
	"miniproject/config"
	"miniproject/ds"
	"miniproject/dto"
	"miniproject/model"
	"miniproject/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HeaderAuth struct {
	Header string `header:"Authorization"`
}

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		res := dto.ResObj{}
		h := &HeaderAuth{}
		err := c.ShouldBindHeader(&h)
		if err != nil {
			res.ErrCode = 409
			res.ErrMsg = "Error on header binding"
			c.JSON(http.StatusOK, res)
			return
		}
		if h.Header == "" {
			res.ErrCode = 403
			res.ErrMsg = "Unprivileged to vie this page"
			c.Abort()
			c.JSON(http.StatusOK, res)
			return
		}
		token, err := util.ValidateAccessToken(h.Header, config.PublicKey)
		if err != nil {
			res.ErrCode = 403
			res.ErrMsg = "not match token"
			c.Abort()
			c.JSON(http.StatusOK, res)
			return
		}

		var user *model.User
		db := ds.DB.Model(&model.User{})
		db = db.Where("id = ?", token.UserID)
		err = db.First(&user).Error
		if err != nil {
			res.ErrCode = 400
			res.ErrMsg = "not found user"
			c.Abort()
			c.JSON(http.StatusOK, res)
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

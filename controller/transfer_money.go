package controller

import (
	"fmt"
	"miniproject/ds"
	"miniproject/dto"
	"miniproject/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type transferMoney struct {
	H *Handler
}

func NewTransferMoney(h *Handler) *transferMoney {
	return &transferMoney{
		H: h,
	}
}

func (ctr *transferMoney) Register() {
	h := ctr.H

	group := h.Router.Group("api/money")
	group.POST("/transfer", ctr.transfer)
	group.POST("/deposit", ctr.deposit)
	group.POST("/withdrawl", ctr.withdrawl)
}

type ReqTransferMoney struct {
	FromEmail string  `json:"from_email" form:"from_email" binding:"required"`
	ToEmail   string  `json:"to_email" form:"to_email" binding:"required"`
	Amount    float64 `json:"amount" form:"amount" binding:"required"`
}

type RespTransferMoney struct {
	FromID    uint64    `json:"transfer_id"`
	FromEmail string    `json:"transfer_email"`
	ToID      uint64    `json:"receive_id"`
	ToEmail   string    `json:"receive_email"`
	Amount    float64   `json:"transfer_amount"`
	CreatedAT time.Time `json:"created_at"`
}

func (ctr *transferMoney) transfer(c *gin.Context) {
	//
	res := dto.ResObj{}
	req := &ReqTransferMoney{}
	var transfer, receive float64

	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.ErrCode = 403
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	tx := ds.DB.Begin()
	user := &model.User{}

	err = tx.Model(&model.User{}).Where("email = ?", req.FromEmail).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		res.ErrCode = 400
		res.ErrMsg = "wrong mail by transfer"
		c.JSON(http.StatusOK, res)
		return
	} else {
		if req.Amount >= user.Balance {
			res.ErrCode = 401
			res.ErrMsg = fmt.Sprintf("transfer balance is more than your balance %v", user.Balance)
			c.JSON(http.StatusOK, res)
			return
		}
		transfer = user.Balance - req.Amount
		err = tx.Model(&model.User{}).Where("email = ?", req.FromEmail).Update("balance", transfer).Error
		if err != nil {
			res.ErrCode = 402
			res.ErrMsg = err.Error()
			c.JSON(http.StatusOK, res)
			tx.Rollback()
			return
		}
		from := &model.MoneyEntries{
			AccountID:    user.Id,
			AccountEmail: req.FromEmail,
			Amount:       (-req.Amount),
			Type:         3,
		}
		err = tx.Model(&model.MoneyEntries{}).Create(&from).Error
		if err != nil {
			res.ErrCode = 403
			res.ErrMsg = err.Error()
			c.JSON(http.StatusOK, res)
			tx.Rollback()
			return
		}

	}

	//
	user1 := &model.User{}
	err = tx.Model(&model.User{}).Where("email = ?", req.ToEmail).First(&user1).Error
	if err == gorm.ErrRecordNotFound {
		res.ErrCode = 405
		res.ErrMsg = "wrong mail by receive"
		c.JSON(http.StatusOK, res)
		return
	} else {
		receive = user1.Balance + req.Amount
		err = tx.Model(&model.User{}).Where("email = ?", req.ToEmail).Update("balance", receive).Error
		if err != nil {
			res.ErrCode = 406
			res.ErrMsg = err.Error()
			c.JSON(http.StatusOK, res)
			tx.Rollback()
			return
		}
		to := &model.MoneyEntries{
			AccountID:    user1.Id,
			AccountEmail: user1.Email,
			Amount:       req.Amount,
			Type:         4,
		}
		err = tx.Model(&model.MoneyEntries{}).Create(&to).Error
		if err != nil {
			res.ErrCode = 407
			res.ErrMsg = err.Error()
			c.JSON(http.StatusOK, res)
			tx.Rollback()
			return
		}
	}
	//
	transaction := &model.MoneyTransfer{
		FromID:    user.Id,
		FromEmail: user.Email,
		ToID:      user1.Id,
		ToEmail:   user1.Email,
		Amount:    req.Amount,
	}
	err = tx.Model(&model.MoneyTransfer{}).Create(&transaction).Error
	if err != nil {
		res.ErrCode = 408
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	if err != nil {
		res.ErrCode = 900
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		res.ErrCode = 700
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	response := &RespTransferMoney{
		FromID:    user.Id,
		FromEmail: user.Email,
		ToID:      user1.Id,
		ToEmail:   user1.Email,
		Amount:    req.Amount,
		CreatedAT: time.Now(),
	}

	res.ErrCode = 0
	res.ErrMsg = "success"
	res.Data = response
	c.JSON(http.StatusOK, res)
}

type ReqDepositMoney struct {
	Email  string  `json:"email" form:"email" binding:"required"`
	Amount float64 `json:"amount" form:"amount" binding:"required,min=10"`
}
type ResponseDeposit struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Owner     string    `json:"owner"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func (ctr *transferMoney) deposit(c *gin.Context) {
	//
	res := dto.ResObj{}
	req := ReqDepositMoney{}
	var value float64
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.ErrCode = 403
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	tx := ds.DB.Begin()
	user := &model.User{}
	err = tx.Model(&model.User{}).Where("email = ?", req.Email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		res.ErrCode = 700
		res.ErrMsg = "email wrong"
		c.JSON(http.StatusOK, res)
		tx.Rollback()
		return
	} else {

		value = user.Balance + req.Amount
		err = tx.Model(&model.User{}).Where("email = ?", req.Email).Update("balance", value).Error
		if err != nil {
			res.ErrCode = 450
			res.ErrMsg = err.Error()
			c.JSON(http.StatusOK, res)
			tx.Rollback()
			return
		}
		userMoney := &model.UserMoney{
			Uid:    user.Id,
			Amount: req.Amount,
			Type:   1,
		}
		err = tx.Model(&model.UserMoney{}).Create(&userMoney).Error
		if err != nil {
			res.ErrCode = 900
			res.ErrMsg = err.Error()
			c.JSON(http.StatusOK, res)
			tx.Rollback()
			return
		}
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		res.ErrCode = 600
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	err = tx.Commit().Error
	if err != nil {
		res.ErrCode = 800
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		tx.Rollback()
		return
	}
	deposit := &ResponseDeposit{
		ID:        user.Id,
		Email:     user.Email,
		Owner:     user.Username,
		Balance:   value,
		CreatedAt: user.CreatedAt,
	}
	res.ErrCode = 0
	res.ErrMsg = "success"
	res.Data = deposit
	c.JSON(http.StatusOK, res)

}

type ReqWithdrawlMoney struct {
	Email  string  `json:"email" form:"email" binding:"required"`
	Amount float64 `json:"amount" form:"amount" binding:"required,min=10"`
}
type ResponseWithdrawl struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Owner     string    `json:"owner"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func (ctr *transferMoney) withdrawl(c *gin.Context) {
	//
	var value float64
	res := dto.ResObj{}
	req := ReqWithdrawlMoney{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.ErrCode = 403
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	tx := ds.DB.Begin()
	user := &model.User{}
	err = tx.Model(&model.User{}).Where("email = ?", req.Email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		res.ErrCode = 700
		res.ErrMsg = "email wrong"
		c.JSON(http.StatusOK, res)
		tx.Rollback()
		return
	} else {
		if req.Amount >= user.Balance {
			res.ErrCode = 760
			res.ErrMsg = fmt.Sprintf("withdrwal amount is more than your balance %v", user.Balance)
			c.JSON(http.StatusOK, res)
			return
		}
		value = user.Balance - req.Amount
		err = tx.Model(&model.User{}).Where("email = ?", req.Email).Update("balance", value).Error
		if err != nil {
			res.ErrCode = 450
			res.ErrMsg = err.Error()
			c.JSON(http.StatusOK, res)
			tx.Rollback()
			return
		}
		userMoney := &model.UserMoney{
			Uid:    user.Id,
			Amount: (-req.Amount),
			Type:   2,
		}
		err = tx.Model(&model.UserMoney{}).Create(&userMoney).Error
		if err != nil {
			res.ErrCode = 900
			res.ErrMsg = err.Error()
			c.JSON(http.StatusOK, res)
			tx.Rollback()
			return
		}
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		res.ErrCode = 600
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	err = tx.Commit().Error
	if err != nil {
		res.ErrCode = 800
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK, res)
		tx.Rollback()
		return
	}
	withdrawl := &ResponseDeposit{
		ID:        user.Id,
		Email:     user.Email,
		Owner:     user.Username,
		Balance:   value,
		CreatedAt: user.CreatedAt,
	}
	res.ErrCode = 0
	res.ErrMsg = "success"
	res.Data = withdrawl
	c.JSON(http.StatusOK, res)

}

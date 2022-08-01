package dto

// response object
type ResObj struct {
	ErrCode uint64      `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    interface{} `json:"data,omitempty"`
}

// req verifycode

type ReqVerifyCode struct {
	Email string `json:"email" form:"email" binding:"required"`
}

//register request
type ReqUserRegister struct {
	Email      string `json:"email" form:"email" binding:"required"`
	Password   string `json:"password" form:"password" binding:"required,min=6"`
	VerifyCode string `json:"verify_code" form:"verify_code" binding:"required"`
	NationID   string `json:"nation_id" form:"nation_id" binding:"required"`
	Type       int8   `json:"type" form:"type" binding:"required"`
	City       string `json:"city" form:"city" binding:"required"`
}

//login request
type ReqLogin struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

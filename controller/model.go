package controller

type GetReq struct {
	Id    string `form:"id" binding:"required"`
	Token string `form:"token"`
}

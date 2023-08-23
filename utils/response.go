package utils

import "github.com/gin-gonic/gin"

type ResponseData struct {
	Code 	int		`json:"code"`
	Msg		string	`json:"msg"`
	Data  	any		`json:"data"`
}

type ResponseObject struct {
	Handle 	*gin.Context
	Data	ResponseData
}

func CreateResponse(c *gin.Context) *ResponseObject {
	return &ResponseObject{Handle: c}
}

func (r *ResponseObject) Json(code int,msg string ,obj any){
	d 	:= &ResponseData{
		Code: code,Msg: msg,Data: obj,
	}
	r.Handle.JSON(code,d)
}

func (r *ResponseObject) Unauthorized(){
	d 	:= &ResponseData{
		Code: 401,Msg: "Unauthorized",
	}
	r.Handle.JSON(401,d)
}

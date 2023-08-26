package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
	r.Handle.JSON(http.StatusOK,d)
}

func (r *ResponseObject) Success(obj any){
	d 	:= &ResponseData{
		Code: http.StatusOK,Msg: "success",Data: obj,
	}
	r.Handle.JSON(http.StatusOK,d)
}

func (r *ResponseObject) Unauthorized(){
	d 	:= &ResponseData{
		Code: http.StatusUnauthorized,Msg: "Unauthorized",
	}
	r.Handle.JSON(http.StatusOK,d)
}

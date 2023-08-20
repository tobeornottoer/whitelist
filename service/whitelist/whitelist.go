package whitelist

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
	"whitelist/model"
	"whitelist/utils"
)

type WhiteList struct {

}

func GetWhiteListHandle() *WhiteList {
	return &WhiteList{}
}

// Get 获取email列表
func Get(c *gin.Context) {
	page,_		:= strconv.Atoi(c.DefaultQuery("page","1"))
	pageSize,_	:= strconv.Atoi(c.DefaultQuery("pageSize","10"))
	search		:= c.DefaultQuery("search","")
	searchWl	:= c.DefaultQuery("whitelist","")
	referral	:= c.DefaultQuery("referral","")
	db			:= utils.GetDB()
	var whitelist []model.WaitList
	var total int64
	handle		:= db.Model(&model.WaitList{})
	if search 	!= "" {
		handle	= handle.Where("email = ?",search)
	}
	if searchWl == "Yes" {
		handle	= handle.Where("white_list_flag = ?",true)
	}
	if searchWl == "No" {
		handle	= handle.Where("white_list_flag = ?",false)
	}
	handle		= handle.Count(&total)
	if referral	== "asc" {
		handle	= handle.Order("referral")
	}
	if referral	== "desc" {
		handle	= handle.Order("referral desc")
	}
	result		:= handle.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&whitelist)
	if result.Error != nil {
		utils.CreateResponse(c).Json(http.StatusInternalServerError,result.Error.Error(),nil)
	} else {
		utils.CreateResponse(c).Json(http.StatusOK,"success",gin.H{"list":whitelist,"count":total})
	}

}

// Approve 添加email到白名单
func Approve(c *gin.Context) {
	emailArr,ok := c.GetPostFormArray("email")
	if !ok {
		utils.CreateResponse(c).Json(http.StatusBadRequest,"param email invalid",nil)
		return
	}
	db			:= utils.GetDB()
	result		:= db.Model(model.WaitList{}).
		Where("email IN ?",emailArr).
		Updates(model.WaitList{WhiteListFlag: true})
	if result.Error != nil {
		utils.CreateResponse(c).Json(http.StatusInternalServerError,result.Error.Error(),nil)
		return
	}
	utils.CreateResponse(c).Json(http.StatusOK,"success",nil)
}

func AddEmail(c *gin.Context) {
	email 	:= c.DefaultPostForm("email","")
	if email == "" {
		utils.CreateResponse(c).Json(http.StatusBadRequest,"param email invalid",nil)
		return
	}
	var whitelist model.WaitList
	db		:= utils.GetDB()
	result	:= db.First(&whitelist,"email = ?",email)
	if errors.Is(result.Error,gorm.ErrRecordNotFound) {
		newEmail := &model.WaitList{
			Email			: email,
			Code			: "",
			Referral		: 0,
			IP				: c.ClientIP(),
			WhiteListFlag	: true,
			RegisterFlag	: false,
			Unsubscribe		: true,
		}
		insert 	:= db.Create(newEmail)
		if insert.Error != nil {
			utils.CreateResponse(c).Json(http.StatusInternalServerError,"add fail",nil)
			return
		}
		utils.CreateResponse(c).Json(http.StatusOK,"success",nil)
		return
	}
	if whitelist.WhiteListFlag == true {
		utils.CreateResponse(c).Json(http.StatusInternalServerError,"The email is already on the whitelist",nil)
		return
	}
	update 	:= db.Model(&whitelist).Update("white_list_flag",true)
	if update.Error != nil {
		utils.CreateResponse(c).Json(http.StatusInternalServerError,"add fail",nil)
		return
	}
	utils.CreateResponse(c).Json(http.StatusOK,"success",nil)
	return
}

func UploadEmailExcel(c *gin.Context) {
	file,uErr	:= c.FormFile("file")
	if uErr != nil {
		utils.CreateResponse(c).Json(http.StatusBadRequest,"file not found",nil)
		return
	}
	today		:= time.Now().Format("20060102")
	uploadDir 	:= os.Getenv("UPLOAD_DIR")
	rootPath,_	:= os.Getwd()
	fileDir		:= uploadDir + today + "/"
	_,err		:= os.Stat(fileDir)
	if err != nil {
		mkdirErr := os.MkdirAll(rootPath + fileDir,0777)
		if mkdirErr != nil {
			utils.CreateResponse(c).Json(http.StatusInternalServerError,"can not create upload dir",nil)
			return
		}
	}
	ext 		:= path.Ext(file.Filename)
	allowExt	:= map[string]bool{".xls":true,".xlsx":true,".csv":true}
	if _, ok 	:= allowExt[ext]; !ok {
		utils.CreateResponse(c).Json(http.StatusInternalServerError,"Unsupported File Format",nil)
		return
	}
	filename 	:= fileDir + time.Now().Format("20060102150405") + "_" + file.Filename
	uploadErr	:= c.SaveUploadedFile(file,rootPath + filename)
	if uploadErr != nil {
		utils.CreateResponse(c).Json(http.StatusInternalServerError,"upload fail",nil)
		return
	}
	code		:= base64.StdEncoding.EncodeToString([]byte(filename))
	utils.CreateResponse(c).Json(http.StatusOK,"success",gin.H{"file":file.Filename,"code":code})
	return
}

func ImportEmails(c *gin.Context){
	code 	:= c.DefaultPostForm("code","")
	if code == "" {
		utils.CreateResponse(c).Json(http.StatusBadRequest,"file not found",nil)
		return
	}
	fileByte,err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		utils.CreateResponse(c).Json(http.StatusBadRequest,"file not found",nil)
		return
	}
	file		:= string(fileByte)
	rootPath,_	:= os.Getwd()
	filePath	:= rootPath + file
	importErr 	:= utils.ImportEmails(filePath,c.ClientIP())
	_			= os.Remove(filePath)
	if importErr != nil {
		utils.CreateResponse(c).Json(http.StatusInternalServerError,importErr.Error(),nil)
		return
	}
	utils.CreateResponse(c).Json(http.StatusOK,"success",nil)
	return
}

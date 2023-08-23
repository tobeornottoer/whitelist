package admins

import (
	"errors"
	passwordKit "github.com/dwin/goSecretBoxPassword"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"os"
	"whitelist/model"
	"whitelist/utils"
)


const AdminPasswordKey = "whitelistAdminSalt"

func Login(c *gin.Context){
	account 	:= c.PostForm("account")
	password	:= c.PostForm("password")
	if account == "" || password == "" {
		utils.CreateResponse(c).Json(400,"account or password can not be null",nil)
		return
	}
	var admin model.Admins
	db 	:= utils.GetDB()
	result := db.Where("account = ?",account).First(&admin)
	if errors.Is(result.Error,gorm.ErrRecordNotFound) {
		utils.CreateResponse(c).Json(400,"account or password invalid",nil)
		return
	}
	err := passwordKit.Verify(password, AdminPasswordKey,admin.Password)
	if err != nil {
		utils.CreateResponse(c).Json(400,"account or password invalid",nil)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"scope"		: "whitelist_manager",
		"adminId" 	: admin.ID,
		"account"	: admin.Account,
	})

	tokenString, issueErr := token.SignedString([]byte(AdminPasswordKey))
	if issueErr != nil {
		utils.CreateResponse(c).Json(500,issueErr.Error(),nil)
		return
	}
	utils.CreateResponse(c).Json(200,"success",
		gin.H{"admin_id":admin.ID,"account":admin.Account,"token":tokenString})
	return
}

func Register(c *gin.Context){
	account 	:= os.Getenv("ADMIN_ACCOUNT")
	password 	:= os.Getenv("ADMIN_PASSWORD")
	if account == "" || password == "" {
		utils.CreateResponse(c).Json(500,"admin config not found",nil)
		return
	}
	db 		:= utils.GetDB()
	var admin model.Admins
	search 		:= db.Where("account = ?",account).First(&admin)
	if errors.Is(search.Error,gorm.ErrRecordNotFound) {
		secret,err 	:= passwordKit.Hash(
			password,
			AdminPasswordKey,
			0,
			passwordKit.ScryptParams{N: 32768, R: 16, P: 1},
			passwordKit.DefaultParams,
		)
		if err != nil {
			utils.CreateResponse(c).Json(500,err.Error(),nil)
			return
		}
		admin 	:= &model.Admins{Account: account,Password: secret}
		result 	:= db.Create(&admin)
		if result.Error != nil {
			utils.CreateResponse(c).Json(500,"init admin fail",nil)
			return
		}
	}

	utils.CreateResponse(c).Json(200,"success",nil)
	return
}

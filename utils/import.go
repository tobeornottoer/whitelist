package utils

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm/clause"
	"os"
	"regexp"
	"strconv"
	"whitelist/model"
)

// ImportEmails 从excel导入emails
func ImportEmails(filePath string,ip string) error {
	_,err		:= os.Stat(filePath)
	if err != nil {
		return errors.New("file not found")
	}
	f,openErr 		:= excelize.OpenFile(filePath)
	if openErr != nil {
		return openErr
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, exlErr := f.GetRows("Sheet1")
	if exlErr != nil {
		return exlErr
	}
	emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	var emails []model.WaitList
	for rIndex, row := range rows {
		if !emailRegexp.MatchString(row[0]) {
			return errors.New("The email in line " + strconv.Itoa(rIndex + 1) + " is incorrect")
		}
		emails = append(emails, model.WaitList{
			Email: row[0],
			Code: "",
			Referral: 0,
			IP: ip,
			WhiteListFlag: true,
			RegisterFlag: false,
			Unsubscribe: true,
		})
	}
	db 	:= GetDB()
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"white_list_flag"}),
	}).Create(&emails)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
/*
func BatchGrant(filePath string) error {
	_,err		:= os.Stat(filePath)
	if err != nil {
		return errors.New("file not found")
	}
	f,openErr 		:= excelize.OpenFile(filePath)
	if openErr != nil {
		return openErr
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, exlErr := f.GetRows("Sheet1")
	if exlErr != nil {
		return exlErr
	}
	emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	var upById	map[uint64]uint64
	var upByEmail map[string]uint64
	var emails	[]string
	for rIndex, row := range rows {
		token,err := strconv.ParseUint(row[2],10,64)
		if err != nil {
			return errors.New(fmt.Sprintf("The value in line %d is incorrect",rIndex + 1))
		}
		if row[0] != "" {
			uid,err := strconv.ParseUint(row[0],10,64)
			if err != nil {
				return errors.New(fmt.Sprintf("The value in line %d is incorrect",rIndex + 1))
			}
			upById[uid] = token
		} else if row[1] != "" {
			if !emailRegexp.MatchString(row[1]) {
				return errors.New(fmt.Sprintf("The email in line %d is incorrect",rIndex + 1))
			}
			upByEmail[row[1]] = token
			emails	= append(emails,row[1])
		}
	}
	db 	:= GetDB()
	type emailMap struct {
		ID		uint64
		Email 	string
	}
	var maps []emailMap
	if len(emails) > 0 {
		result := db.Model(&model.Account{}).Where("email in ?",emails).Select("id", "email").Scan(&maps)
		if result.Error == nil {
			for _,value := range maps {
				if t,ok := upByEmail[value.Email];ok {
					upById[value.ID] = t
				}
			}
		}
	}
	var grants []model.TokenGrant
	//var updates []model.Account
	for uid,ut := range upById {
		grants	= append(grants,model.TokenGrant{
			AccountId: uid,
			Date:time.Now().Format("2006-01-02"),
			Token:ut,
			Way:2,
		})
		//todo 这里更新有问题，给每个用户增加相对应的token
		//updates	= append(updates,model.Account{ID:uid,TokenCount: })
	}
	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&grants).Error; err != nil {
			return err
		}
		//result := db.Clauses(clause.OnConflict{
		//	Columns:   []clause.Column{{Name: "id"}},
		//	DoUpdates: clause.AssignmentColumns([]string{"token_count"}),
		//}).Create(&emails)
		//if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
		//	return err
		//}

		// 返回 nil 提交事务
		return nil
	})

}
*/
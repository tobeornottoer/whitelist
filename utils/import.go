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

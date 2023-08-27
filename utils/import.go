package utils

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"regexp"
	"strconv"
	"time"
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

// BatchGrant 批量发放token
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
	var upById	map[uint]uint64
	var upByEmail map[string]uint64
	var emails	[]string
	var uids	[]uint
	for rIndex, row := range rows {
		token,err := strconv.ParseUint(row[2],10,64)
		if err != nil {
			return errors.New(fmt.Sprintf("The token in line %d is incorrect",rIndex + 1))
		}
		if row[0] != "" {
			uid,err := strconv.ParseUint(row[0],10,64)
			if err != nil {
				return errors.New(fmt.Sprintf("The UID in line %d is incorrect",rIndex + 1))
			}
			u 			:= uint(uid)
			upById[u] 	= token
			uids		= append(uids,u)
		} else if row[1] != "" {
			if !emailRegexp.MatchString(row[1]) {
				return errors.New(fmt.Sprintf("The email in line %d is incorrect",rIndex + 1))
			}
			upByEmail[row[1]] = token
			emails	= append(emails,row[1])
		}
	}
	db 	:= GetDB()
	type acMap struct {
		ID		uint
		Email 	string
		Token	uint64
	}
	var maps []acMap
	handle	:= db.Model(&model.Account{})
	if len(uids) > 0 {
		handle	= handle.Where("id in ",uids)
	}
	if len(emails) > 0 {
		handle	= handle.Or("email in ",emails)
	}
	result	:= handle.Select("id", "email","token_count").Scan(&maps)
	if result.Error != nil {
		return errors.New("an error occurred while searching for users")
	}
	var grants []model.TokenGrant
	var updates []model.Account
	var incr uint64
	for _,value := range maps {
		if t,ok := upById[value.ID];ok {
			incr = t
		} else if t,ok := upByEmail[value.Email];ok {
			incr = t
		}
		if incr > 0 {
			updates	= append(updates,model.Account{Model:gorm.Model{ID: value.ID},TokenCount: value.Token + incr})
			grants	= append(grants,model.TokenGrant{
				AccountId: value.ID,
				Date:      time.Now().Format("2006-01-02"),
				Token:     incr,
				Way:       2,
			})
		}
	}
	if len(grants) <= 0 {
		return errors.New("some errors resulted in no data updates")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err 	:= tx.Create(&grants).Error; err != nil {
			return err
		}
		err 	:= tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"token_count"}),
		}).Create(&updates)
		if err.Error != nil {
			return err.Error
		}
		return nil
	})

}

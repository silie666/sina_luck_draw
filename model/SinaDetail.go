package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sina/drivers/mysql"
)

type SinaDetail struct {
	Id int
	HtmlStr string `gorm:"mediumtext"`
	Mid string
	Uid string
	ToUid int
	IsLike int `orm:"tinyint"`
	IsRepost int `orm:"tinyint"`
	LuckTime string
}

func (sl *SinaDetail) SinaDetailAdd(params SinaDetail)error  {
	var result *gorm.DB
	var info SinaDetail
	//and  is_ok = 1
	err := mysql.Db.Where("mid = ? and uid = ?", params.Mid,params.Uid).First(&info).Error
	if errors.Is(err,gorm.ErrRecordNotFound) {
		result = mysql.Db.Create(&params)
		fmt.Println("添加成功")
	} else {
		fmt.Println("记录已经存在")
		return nil
	}
	return result.Error
}


func (sl *SinaDetail) SinaDetailSave(params SinaDetail)error  {
	var result *gorm.DB
	var info SinaDetail
	err := mysql.Db.Where("id = ?", params.Id).First(&info).Error
	if !errors.Is(err,gorm.ErrRecordNotFound) {
		result = mysql.Db.Updates(&params)
	}
	return result.Error
}




func (sl *SinaDetail)SinaDetailList(where string)[]SinaDetail {
	var sina_luck_list []SinaDetail
	mysql.Db.Where(where).Find(&sina_luck_list)
	return sina_luck_list
}
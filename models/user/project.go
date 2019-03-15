package user

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/paylm/myweb/pkg/gmysql"
)

type Project struct {
	gorm.Model
	Id           int    `gorm:"primary_key"`
	Title        string `type:"varchar(100)"`
	Tag          string
	Content      string
	Del          int8
	ByUser       int
	PubTime      time.Time `gorm:"column:pub_time"`
	CommentCount int       `gorm:"column:comment_count"`
}

// 设置User的表名为`profiles`
func (Project) TableName() string {
	return "project"
}

func GetProjects() []Project {
	var articles []Project
	err := gmysql.DB.Find(&articles).Error
	if err != nil {
		return nil
	}
	return articles
}

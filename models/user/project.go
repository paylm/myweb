package user

import (
	"fmt"
	"time"

	"github.com/paylm/myweb/pkg/gmysql"
	"github.com/paylm/myweb/pkg/gredis"
)

type Project struct {
	Id           int    `gorm:"primary_key"`
	Title        string `type:"varchar(100)"`
	Tag          string
	Content      string
	Del          int8 `gorm:"column:del"`
	ByUser       int
	PubTime      time.Time `gorm:"column:pub_time"`
	CommentCount int       `gorm:"column:comment_count"`
}

// 设置User的表名为`profiles`
func (Project) TableName() string {
	return "project"
}

type Book struct {
	Id        int `gorm:"primary_key"`
	ProjectId int `gorm:"cloumn:project_id"`
	ByUser    int `gorm:"cloumn:by_user"`
	//PubTime   time.Time `gorm:"cloumn:create_time"`
	Stat int8 `gorm:"cloumn:stat"`
}

func (Book) TableName() string {
	return "project_book"
}

func GetProjects() []Project {
	var articles []Project
	err := gmysql.DB.Find(&articles).Error
	if err != nil {
		return nil
	}
	return articles
}

/**
* 订顶目,返回项目单号
*
 */
func BookPj() int {
	lock := gredis.SetNX("lock_book", "pj")
	if !lock {
		fmt.Printf("can't get lock to book project\n")
		return 0
	}
	var b Book
	err := gmysql.DB.Where("stat = ?", "0").First(&b).Error
	if err != nil {
		fmt.Printf("get book with err:%v\n", err)
		return 0
	}
	if b.Stat != 0 {
		fmt.Printf("all the project is not stat:0\n")
		return 0
	}
	fmt.Printf("book:%v\n", b)
	b.Stat = 1
	gmysql.DB.Save(&b)
	gredis.UnLock("book")
	return b.Id
}

/**
* 订顶目,不带锁返回项目单号
* 高并发存在重复消费情况
 */
func BookUoLockPj() int {
	var b Book
	err := gmysql.DB.Where("stat = ?", "0").First(&b).Error
	if err != nil {
		fmt.Printf("get book with err:%v\n", err)
		return 0
	}
	if b.Stat != 0 {
		fmt.Printf("all the project is not stat:0\n")
		return 0
	}
	fmt.Printf("book:%v\n", b)
	b.Stat = 1
	gmysql.DB.Save(&b)
	return b.Id
}

package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func Setup() error {
	var err error
	Db, err = gorm.Open("mysql", "paylm:123456@127.0.0.1/ppl?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Println("connection succedssed")
	}
	return err
}

func Close() {
	defer Db.Close()
}

package gmysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/paylm/myweb/pkg/setting"
)

var DB *gorm.DB

func Setup() error {
	var err error
	DB, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", setting.DatabaseSetting.User, setting.DatabaseSetting.Password, setting.DatabaseSetting.Host, setting.DatabaseSetting.Name))
	if err != nil {
		fmt.Printf("conn %s with err:%v", setting.DatabaseSetting.Type, err)
	}
	return err
}

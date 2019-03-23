package models

import (
	"fmt"
	"time"

	"github.com/paylm/myweb/pkg/gmysql"
)

type Stat struct {
	ID         int
	Atype      string `gorm:"type:varchar(12)"`
	Month      string `gorm:"type:varchar(16)"`
	CreateTime time.Time
	Count      int `gorm:"type:int"`
}

// 设置User的表名为`profiles`
func (Stat) TableName() string {
	return "active_stat"
}

func FindByType(atype string) []Stat {
	var stats []Stat

	err := gmysql.DB.Where("atype=`?`", atype).Find(&stats).Error
	if err != nil {
		fmt.Printf("FindByType atype=%s ,throw err:%v\n", atype, err)
		return nil
	}

	return stats
}

func InsertStat(s Stat) error {
	err := gmysql.DB.Create(&s).Error
	return err
}

package user

import (
	"fmt"
	"time"

	"github.com/paylm/myweb/pkg/gmysql"
)

type Stat struct {
	Id         int       `gorm:"type:int;primary_key;auto_increment" json:"id"`
	Atype      string    `gorm:"type:varchar(12)" json:"atype"`
	Month      string    `gorm:"type:varchar(16)" json:"month"`
	CreateTime time.Time `time_format:"2006-01-02" time_utc:"1" json:"createTime"`
	Count      int       `gorm:"type:int" json:"count"`
}

// 设置User的表名为`profiles`
func (Stat) TableName() string {
	return "active_stat"
}

func FindByType(atype string) []Stat {
	var stats []Stat
	err := gmysql.DB.Where("atype=?", atype).Find(&stats).Error
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

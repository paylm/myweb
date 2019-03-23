package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/paylm/myweb/pkg/gmysql"
	"github.com/paylm/myweb/pkg/gredis"
)

const (
	ONLINE_KEY = "online"
	ACTIVE_KEY = "active"
)

type User interface {
	Verlogin() error
	Reg() error
	Logout()
}

type UserInfo struct {
	tel    string
	openid string
}

type UserData struct {
	Id       int    `gorm:"primary_key" sql:"auto_increment;primary_key;unique" json:"id"`
	Img      string `form:"img" json:"img"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Email    string `form:"email" json:"email"`
	Job      string `form:"job"`
	Stat     int    `form:"stat" json:"stat"`
	//UserInfo
}

func (UserData) TableName() string {
	return "user"
}

func (u *UserData) Verlogin() (UserData, error) {

	var u2 UserData
	var err error
	if err = gmysql.DB.Where("username = ?", u.Username).First(&u2).Error; err != nil {
		fmt.Printf("Verlogin with err:%v\n", err)
		return u2, err
	}
	if &u2 == nil {
		fmt.Printf("Verlogin with err:%v\n", err)
		return u2, errors.New("帐号不存在")
	}

	if u2.Stat < 0 {
		return u2, errors.New("此帐号已被禁用")
	}

	strPwd := string(u2.Password)
	if u.Password != strPwd {
		fmt.Printf("pwd:%s,post pwd:%s\n", strPwd, u.Password)
		return u2, errors.New("密码错误")
	}
	//统计活跃
	gredis.Exec("SETBIT", ONLINE_KEY, u2.Id, 1)
	regActive(u2.Id)
	return u2, nil
}

func (u *UserData) Reg() error {

	if u.Password == "" || u.Username == "" {
		return errors.New("username or password can't be null")
	}
	err := gmysql.DB.Create(u).Error
	if err != nil {
		return err
	}
	gredis.Set(u.Username, u.Password, 3600)
	gredis.Exec("SETBIT", ONLINE_KEY, u.Id, 1)
	fmt.Printf("reg save:%v\n", u)
	return err
}

func OnlineCount() int {
	res, err := gredis.Exec("BITCOUNT", ONLINE_KEY)
	if err != nil {
		fmt.Printf("Online count whith err:%v\n", err)
		return 0
	}

	count, err1 := redis.Int(res, nil)
	if err1 != nil {
		return 0
	}

	return count
}

/**
*
* 注册用户当天登录的
**/
func regActive(uid int) {
	d := time.Now()
	today := fmt.Sprintf("%d-%d-%d", d.Year(), d.Month(), d.Day())
	gredis.Exec("SETBIT", fmt.Sprintf("%s-%s", ACTIVE_KEY, today), uid, 1)
	fmt.Printf("active %d at %s", uid, today)
}

func RecentActive() int {
	d := time.Now()
	dstkey := fmt.Sprintf("%s-7", ACTIVE_KEY)
	gredis.Exec("BITOP", "OR", dstkey, fmt.Sprintf("%s-%d-%d-%d", ACTIVE_KEY, d.Year(), d.Month(), d.Day()), fmt.Sprintf("%s-%d-%d-%d", ACTIVE_KEY, d.Year(), d.Month(), d.Day()-1), fmt.Sprintf("%s-%d-%d-%d", ACTIVE_KEY, d.Year(), d.Month(), d.Day()-2), fmt.Sprintf("%s-%d-%d-%d", ACTIVE_KEY, d.Year(), d.Month(), d.Day()-3), fmt.Sprintf("%s-%d-%d-%d", ACTIVE_KEY, d.Year(), d.Month(), d.Day()-4), fmt.Sprintf("%s-%d-%d-%d", ACTIVE_KEY, d.Year(), d.Month(), d.Day()-5), fmt.Sprintf("%s-%d-%d-%d", ACTIVE_KEY, d.Year(), d.Month(), d.Day()-6))

	res, err := gredis.Exec("BITCOUNT", dstkey)
	if err != nil {
		fmt.Printf("Online count whith err:%v\n", err)
		return 0
	}

	count, err1 := redis.Int(res, nil)
	if err1 != nil {
		return 0
	}
	return count
}

func GetAllUser(limit int) []UserData {
	var users []UserData
	err := gmysql.DB.Limit(limit).Find(&users).Error
	if err != nil {
		return nil
	}
	return users
}

func Logout(id int) {
	fmt.Printf("logout u:%d\n", id)
	gredis.Exec("SETBIT", "online", id, 0)
}

package user

import (
	"errors"
	"fmt"

	"github.com/paylm/myweb/pkg/gredis"
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
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Email    string `form:"email" json:"email"`
	//UserInfo
}

func (u *UserData) Verlogin() error {

	pwd, err := gredis.Get(u.Username)
	if err != nil {
		return err
	}
	strPwd := string(pwd)
	if u.Password != strPwd {
		fmt.Printf("pwd:%s,post pwd:%s\n", strPwd, u.Password)
		return errors.New("密码错误")
	}
	return nil
}

func (u *UserData) Reg() error {

	err := gredis.Set(u.Username, u.Password, 3600)
	fmt.Printf("reg save:%v\n", u)
	return err
}

func (u *UserData) Logout() {

}

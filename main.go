package main

import (
	"fmt"
	"net/http"

	"github.com/paylm/myweb/pkg/gmysql"
	"github.com/paylm/myweb/pkg/gredis"
	"github.com/paylm/myweb/pkg/setting"
	"github.com/paylm/myweb/routers"
)

func main() {
	setting.Setup()
	err := gredis.Setup()
	if err != nil {
		fmt.Printf("connect redis with err:%v", err)
		return
	}
	err = gmysql.Setup()
	if err != nil {
		fmt.Printf("connect mysql with err:%v", err)
		return
	}

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	fmt.Printf("listen at %s , HttpPort : %v", endPoint, setting.ServerSetting.HttpPort)
	fmt.Println(setting.RedisSetting)
	server.ListenAndServe()
	return

	//endless.DefaultReadTimeOut = readTimeout
	//	endless.DefaultWriteTimeOut = writeTimeout
	//	endless.DefaultMaxHeaderBytes = maxHeaderBytes
	//server := endless.NewServer(endPoint, routersInit)
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}

}

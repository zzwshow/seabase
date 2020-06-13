package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"seabase/extend/conf"
	"seabase/extend/redis"
	"seabase/model"
	"seabase/model/base"
	"seabase/router"
)

var (
	s string
	m bool
	h bool
	u bool
)

func usage() {
	_,err:= fmt.Fprintf(os.Stderr,`version: SeaBase/1.0.0 auther: zzw Usage: SeaBase [s runMode][-m migrate][-u rerest root password][-h help]
	Options:
		-s 项目运行模式(默认debug 生产release)
		-m 数据库初始化(初次部署安装时需要执行)
		-u 初始化admin用户密码
		-h 帮助
`)
	if err != nil {
		panic(err.Error())
	}
}

func init(){
	flag.StringVar(&s,"s","debug","项目运行模式(默认debug 生产release)")
	flag.BoolVar(&m,"m",false,"数据库初始化(初次部署安装时需要执行)")
	flag.BoolVar(&u,"u",false,"初始化admin账户密码")
	flag.BoolVar(&h,"h",false,"帮助")
	flag.Usage = usage
	log.SetPrefix("SeaBase_error: ")
	log.SetFlags(log.Ldate|log.Lmicroseconds|log.Llongfile)
}


func main() {
	flag.Parse()
	if h {
		flag.Usage()
	} else {
		if err := redis.Init(); err != nil {
			os.Exit(3)
		}
		conf.InitConf(s)
		model.Init()
		if m {
			model.DB.AutoMigrate(&base.UserModel{})
		}
		gin.SetMode(conf.ServerConf.RunMode)
		r := router.InitRouter()
		_ = r.Run(fmt.Sprintf(":%d",conf.ServerConf.Port))
	}
}

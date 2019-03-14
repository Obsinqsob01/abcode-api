package main

import (
	_ "abcode/routers"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"

	"github.com/astaxie/beego"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	orm.RegisterDriver("sqlite3", orm.DRSqlite)

	orm.RegisterDataBase("default", "sqlite3", "file:data.db")
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	err := orm.RunSyncdb("default", false, false)
	if err != nil {
		beego.Error(err)
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
	}))

	beego.Run()
}

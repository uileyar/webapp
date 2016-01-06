package controllers

import (
	"flag"

	"github.com/revel/revel"
)

func init() {
	flag.Set("logtostderr", "true")
	//flag.Set("v", "2")
	revel.OnAppStart(InitDB)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	//revel.InterceptMethod(Application.AddUser, revel.BEFORE)
	//revel.InterceptMethod(Hotels.checkUser, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}
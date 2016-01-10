package controllers

import (
	"flag"
	"regexp"

	"github.com/golang/glog"
	"github.com/revel/revel"
)

func init() {
	flag.Set("logtostderr", "true")
	//flag.Set("v", "2")
	revel.OnAppStart(InitDB)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod(Application.AddUser, revel.BEFORE)
	revel.InterceptMethod(Accounts.checkUser, revel.BEFORE)
	revel.InterceptMethod(Catelogs.checkUser, revel.BEFORE)
	revel.InterceptMethod(Bills.checkUser, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}

func CheckSqlStr(sqlStr string) bool {
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	re, err := regexp.Compile(str)
	if err != nil {
		glog.Errorf(err.Error())
		return true
	}

	return re.MatchString(sqlStr)
}

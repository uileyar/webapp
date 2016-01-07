package controllers

import (
	"database/sql"

	"github.com/go-gorp/gorp"
	"github.com/golang/glog"
	_ "github.com/mattn/go-sqlite3"
	"github.com/revel/modules/db/app"
	r "github.com/revel/revel"
	"github.com/uileyar/webapp/app/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	Dbm       *gorp.DbMap
	GLogLevel glog.Level = 3 //set glog V level
)

func InitDB() {
	glog.Infoln("InitDB in")
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.SqliteDialect{}}

	t := Dbm.AddTable(models.User{}).SetKeys(true, "UserId")
	t.ColMap("Password").Transient = true

	Dbm.AddTableWithName(models.Account{}, "jzb_accounts")
	Dbm.AddTableWithName(models.Bill{}, "jzb_bills")
	Dbm.AddTableWithName(models.Catelog{}, "jzb_catelogs")

	Dbm.TraceOn("[gorp]", r.INFO)

	err := Dbm.CreateTables()
	if err == nil {
		insertDemoData()
	}
}

func insertDemoData() {
	glog.V(GLogLevel).Infoln("insertDemoData in")

	bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("demo"), bcrypt.DefaultCost)
	demoUser := &models.User{0, "Demo User", "demo", "demo", bcryptPassword}
	if err := Dbm.Insert(demoUser); err != nil {
		panic(err)
	}

	glog.V(GLogLevel).Infoln("insertDemoData out")
}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	glog.V(GLogLevel).Infoln("Begin in")

	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	glog.V(GLogLevel).Infoln("Begin out")
	return nil
}

func (c *GorpController) Commit() r.Result {
	glog.V(GLogLevel).Infoln("Commit in")
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	glog.V(GLogLevel).Infoln("Commit out")
	return nil
}

func (c *GorpController) Rollback() r.Result {
	glog.V(GLogLevel).Infoln("Rollback in")
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	glog.V(GLogLevel).Infoln("Rollback out")
	return nil
}

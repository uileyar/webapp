package controllers

import (
	"database/sql"
	"fmt"

	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/revel/modules/db/app"
	r "github.com/revel/revel"
	"github.com/uileyar/webapp/app/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	Dbm *gorp.DbMap
)

func InitDB() {
	fmt.Println("InitDB in")
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.SqliteDialect{}}

	t := Dbm.AddTable(models.User{}).SetKeys(true, "UserId")
	t.ColMap("Password").Transient = true

	Dbm.TraceOn("[gorp]", r.INFO)

	err := Dbm.CreateTables()
	if err == nil {
		insertDemoData()
	}
}

func insertDemoData() {
	fmt.Println("insertDemoData in")

	bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("demo"), bcrypt.DefaultCost)
	demoUser := &models.User{0, "Demo User", "demo", "demo", bcryptPassword}
	if err := Dbm.Insert(demoUser); err != nil {
		panic(err)
	}

	fmt.Println("insertDemoData out")
}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	fmt.Println("Begin in")

	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	fmt.Println("Begin out")
	return nil
}

func (c *GorpController) Commit() r.Result {
	fmt.Println("Commit in")
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	fmt.Println("Commit out")
	return nil
}

func (c *GorpController) Rollback() r.Result {
	fmt.Println("Rollback in")
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	fmt.Println("Rollback out")
	return nil
}

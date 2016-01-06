package controllers

import (
	"github.com/golang/glog"
	"github.com/revel/revel"
	"github.com/uileyar/webapp/app/models"
)

type Accounts struct {
	Application
}

func (c Accounts) Index() revel.Result {
	results, err := c.Txn.Select(models.Account{},
		`select * from jzb_accounts`)
	if err != nil {
		panic(err)
	}

	var accounts []*models.Account
	for _, r := range results {
		b := r.(*models.Account)
		accounts = append(accounts, b)
		glog.Infof("%v\n", b)
	}

	return c.Render(accounts)
}

func (c Accounts) Add() revel.Result {
	return c.Render()
}

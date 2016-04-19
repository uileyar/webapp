package controllers

import (
	//"github.com/golang/glog"

	//. "github.com/aerospike/aerospike-client-go"
	"github.com/revel/revel"
	"github.com/uileyar/webapp/app/models"
	"github.com/uileyar/webapp/app/routes"
)

type Analytics struct {
	Application
}

func (c Analytics) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error(c.Message("please_login_first"))
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

func (c Analytics) JsonData() revel.Result {
	results, err := c.Txn.Select(models.Account{},
		`select * from jzb_accounts ORDER BY name DESC`)
	if err != nil {
		panic(err)
	}

	return c.RenderJson(results)
}

func (c Analytics) Bills() revel.Result {
	return c.Redirect(routes.Bills.Index())
}

func (c Analytics) Catelogs() revel.Result {
	return c.Redirect(routes.Catelogs.Index())
}

func (c Analytics) Index() revel.Result {
	return c.Render()
}

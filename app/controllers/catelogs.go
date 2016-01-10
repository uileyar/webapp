package controllers

import (
	"github.com/golang/glog"
	"github.com/revel/revel"
	"github.com/uileyar/webapp/app/models"
	"github.com/uileyar/webapp/app/routes"
)

type Catelogs struct {
	Application
}

func (c Catelogs) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error(c.Message("please_login_first"))
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

func (c Catelogs) Index() revel.Result {
	results, err := c.Txn.Select(models.Catelog{},
		`select * from jzb_catelogs`)
	if err != nil {
		panic(err)
	}

	var catelogs []*models.Catelog
	for _, r := range results {
		b := r.(*models.Catelog)
		catelogs = append(catelogs, b)
		glog.Infof("%v\n", b)
	}

	return c.Render(catelogs)
}

func (c Catelogs) Add() revel.Result {
	return c.Render()
}

func (c Catelogs) SaveCatelog() revel.Result {
	name := "test1"
	kind := "income"

	saveData := &models.Catelog{
		Name: name,
		Kind: kind,
	}

	err := c.Txn.Insert(saveData)
	if err != nil {
		panic(err)
	}
	return c.Redirect(c.Index())
}

package controllers

import (
	"github.com/golang/glog"
	"github.com/revel/revel"
	"github.com/uileyar/webapp/app/models"
	"github.com/uileyar/webapp/app/routes"
)

type Bills struct {
	Application
}

func (c Bills) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error(c.Message("please_login_first"))
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

func (c Bills) Index() revel.Result {
	results, err := c.Txn.Select(models.Bill{},
		`select bill_id,amount,title,description,date,month,catelog_id,account_id,kind,shared,version from jzb_bills`)
	if err != nil {
		panic(err)
	}

	var bills []*models.Bill
	for _, r := range results {
		b := r.(*models.Bill)
		bills = append(bills, b)
		glog.Infof("%v\n", b)
	}

	return c.Render(bills)
}

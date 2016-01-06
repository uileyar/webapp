package controllers

import (
	"github.com/golang/glog"
	"github.com/revel/revel"
	"github.com/uileyar/webapp/app/models"
)

type Catelogs struct {
	Application
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

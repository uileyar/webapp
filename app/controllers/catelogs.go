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
		//glog.Infof("%v\n", b)
	}

	return c.Render(catelogs)
}

func (c Catelogs) New() revel.Result {
	return c.Render()
}

func (c Catelogs) Save() revel.Result {
	name := c.Params.Get("name")
	kind := c.Params.Get("kind")
	glog.Infof("new name = %v, kind = %v", name, kind)

	c.Validation.Required(name).Message(c.Message("catelog.name.require"))
	c.Validation.MaxSize(name, 30).Message(c.Message("catelog.name.maxsize"))
	c.Validation.Required(kind).Message(c.Message("catelog.kind.require"))

	if CheckSqlStr(name) {
		c.Validation.Error("%s %s", name, c.Message("wrong_format"))
	}

	results, _ := c.Txn.Select(models.Catelog{},
		`select * from jzb_catelogs where name=?`, name)
	if len(results) > 0 {
		c.Validation.Error("%s %s", name, c.Message("catelog.name.exist"))
	}

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Catelogs.New())
	}

	data := &models.Catelog{
		Name: name,
		Kind: kind,
	}

	err := c.Txn.Insert(data)
	if err != nil {
		panic(err)
	}

	c.Flash.Success("%s %s %s!", c.Message("add"), name, c.Message("successed"))

	return c.Redirect(routes.Catelogs.Index())
}

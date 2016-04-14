package controllers

import (
	"github.com/golang/glog"
	"github.com/revel/revel"
	"github.com/uileyar/webapp/app/models"
	"github.com/uileyar/webapp/app/routes"
)

type Accounts struct {
	Application
}

func (c Accounts) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error(c.Message("please_login_first"))
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

func (c Accounts) JsonData() revel.Result {
	results, err := c.Txn.Select(models.Account{},
		`select * from jzb_accounts ORDER BY name DESC`)
	if err != nil {
		panic(err)
	}

	var total models.Account
	total.Name = "总计"
	for _, r := range results {
		account := r.(*models.Account)
		total.Income += account.Income
		total.Expense += account.Expense
		total.Balance += account.Balance
	}
	results = append(results, total)

	glog.Infof("account total = %v\n", total)
	return c.RenderJson(results)
}

func (c Accounts) csvData() revel.Result {

}

func (c Accounts) Index() revel.Result {
	return c.Render()
}

func (c Accounts) New() revel.Result {
	return c.Render()
}

func (c Accounts) Save() revel.Result {
	kind := "cash"
	name := c.Params.Get("name")
	glog.Infof("new name = %v", name)

	c.Validation.Required(name).Message(c.Message("account.name.require"))
	c.Validation.MaxSize(name, 30).Message(c.Message("account.name.maxsize"))
	//c.Validation.Match(name, regexp.MustCompile(`^([\u4e00-\u9fa5]{1,20}|[a-zA-Z\.\s]{1,20})$`)).Message(c.Message("wrong_format"))

	if CheckSqlStr(name) {
		c.Validation.Error("%s %s", name, c.Message("wrong_format"))
	}

	results, _ := c.Txn.Select(models.Account{},
		`select * from jzb_accounts where name=?`, name)
	if len(results) > 0 {
		c.Validation.Error("%s %s", name, c.Message("account.name.exist"))
	}

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Accounts.New())
	}

	data := &models.Account{
		Name: name,
		Kind: kind,
	}

	err := c.Txn.Insert(data)
	if err != nil {
		panic(err)
	}

	c.Flash.Success("%s %s %s!", c.Message("add"), name, c.Message("successed"))

	return c.Redirect(routes.Accounts.Index())
}

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
		glog.Infof("%v\n", b.Name)
	}

	return c.Render(accounts)
}

func (c Accounts) New() revel.Result {
	return c.Render()
}

func (c Accounts) SaveAccount() revel.Result {
	kind := "cash"
	name := c.Params.Get("name")
	glog.Infof("new name = %v", name)

	c.Validation.Required(name).Message(c.Message("require_account"))
	c.Validation.MinSize(name, 1).Message(c.Message("account_minsize"))
	c.Validation.MaxSize(name, 30).Message(c.Message("account_maxsize"))

	results, _ := c.Txn.Select(models.Account{},
		`select * from jzb_accounts where name=?`, name)
	if len(results) > 0 {
		c.Validation.Error("%s %s", name, c.Message("account_exist"))
	}

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Accounts.New())
	}

	account := &models.Account{
		Name: name,
		Kind: kind,
	}

	err := c.Txn.Insert(account)
	if err != nil {
		panic(err)
	}

	c.Flash.Success("%s %s %s!", c.Message("add"), name, c.Message("successed"))

	return c.Redirect(routes.Accounts.Index())
}

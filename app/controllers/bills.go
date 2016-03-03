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
		`select bill_id,amount,title,description,date,month,catelog_id,account_id,kind,shared,version from jzb_bills order by version DESC`)
	if err != nil {
		panic(err)
	}

	var bills []*models.Bill
	for _, r := range results {
		b := r.(*models.Bill)
		bills = append(bills, b)
		//glog.Infof("%v\n", b)
	}

	return c.Render(bills)
}

func (c Bills) New() revel.Result {
	results, err := c.Txn.Select(models.Account{},
		`select account_id,name from jzb_accounts`)
	if err != nil {
		panic(err)
	}
	var accounts []*models.Account
	for _, r := range results {
		b := r.(*models.Account)
		accounts = append(accounts, b)
	}

	results, err = c.Txn.Select(models.Catelog{},
		`select catelog_id,name from jzb_catelogs`)
	if err != nil {
		panic(err)
	}
	var catelogs []*models.Catelog
	for _, r := range results {
		b := r.(*models.Catelog)
		catelogs = append(catelogs, b)
	}

	return c.Render(accounts, catelogs)
}

func (c Bills) Save(bill models.Bill) revel.Result {

	glog.Infof("bill.Date = %v", bill.Date)
	c.Validate(bill)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Bills.New())
	}

	err := c.Txn.Insert(&bill)
	if err != nil {
		panic(err)
	}

	c.Flash.Success("%v %v %v %v%v %v!", c.Message("add"), bill.Title, c.Message(bill.Kind), "￥", bill.Amount, c.Message("successed"))

	return c.Redirect(routes.Bills.Index())
}

func (c Bills) Validate(bill models.Bill) {
	c.Validation.Check(bill.Amount,
		revel.Required{},
	).Message(c.Message("bill.amount.require"))

	c.Validation.Check(bill.Title,
		revel.Required{},
		revel.MaxSize{50 * 3},
	).Message(c.Message("bill.title.maxsize"))

	c.Validation.Check(bill.Description,
		revel.MaxSize{1000 * 3},
	).Message(c.Message("bill.description.maxsize"))

	c.Validation.Check(bill.Date,
		revel.Required{},
	).Message(c.Message("bill.date.require"))

	c.Validation.Check(bill.Catelog_id,
		revel.Required{},
	).Message(c.Message("bill.catelog.require"))

	c.Validation.Check(bill.Account_id,
		revel.Required{},
	).Message(c.Message("bill.account.require"))

	c.Validation.Check(bill.Kind,
		revel.Required{},
		revel.MaxSize{10 * 3},
	).Message(c.Message("bill.kind.require"))

}

package controllers

import (
	"fmt"
	"time"

	//"github.com/golang/glog"
	"github.com/golang/glog"
	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"
	"github.com/uileyar/webapp/app/models"
	"github.com/uileyar/webapp/app/routes"
)

type Bills struct {
	Application
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

func (c Bills) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error(c.Message("please_login_first"))
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

func (c Bills) Delete() revel.Result {
	var bill models.Bill
	uid := c.Request.PostFormValue("uid")
	if len(uid) < 10 {
		return c.NotFound(uid)
	}

	err := c.Txn.SelectOne(&bill,
		`select amount,title,description,date,month,catelog_id,account_id,kind,shared,version from jzb_bills WHERE bill_id = ?`, uid)
	if err != nil {
		panic(err)
	}

	_, err = c.Txn.Select(models.Bill{},
		`delete from jzb_bills WHERE bill_id = ?`, uid)
	if err != nil {
		panic(err)
	}

	// 立即发送电子邮件（异步）
	amount := bill.Amount
	if bill.Amount < 0 {
		amount = -bill.Amount
	}
	subject := fmt.Sprintf("%v%v %v %v￥%v", c.Message("delete"), c.Message("bill"), bill.Title, c.Message(bill.Kind), amount)
	body := fmt.Sprintf("%v%v%v %v:%v; %v:￥%v; %v:%v; %v:%v", c.connected().Name, c.Message("delete"), c.Message(bill.Kind),
		c.Message("date"), bill.Date.Format("2006-01-02"), c.Message("amount"), amount,
		c.Message("title"), bill.Title, c.Message("description"), bill.Description)
	if !revel.Config.BoolDefault("mode.dev", false) {
		jobs.Now(models.SendConfirmationEmail{
			Subject: subject,
			Body:    body,
		})
	}

	return c.RenderText("ok")
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
	var bill models.Bill

	uid := c.Request.FormValue("uid")
	if len(uid) > 10 {
		//glog.Infof("get bill id = %v\n", uid)
		err := c.Txn.SelectOne(&bill, `select bill_id,amount,title,description,date,month,catelog_id,account_id,kind,shared,version from jzb_bills WHERE bill_id = ?`, uid)
		if err != nil {
			panic(err)
		}
		if bill.Amount < 0 {
			bill.Amount = -bill.Amount
		}
		//glog.Infof("Description = %v", bill.Description)
	} else {
		bill.Date = time.Now()
	}

	accounts, err := c.Txn.Select(models.Account{},
		`select account_id,name from jzb_accounts`)
	if err != nil {
		panic(err)
	}

	catelogs, err := c.Txn.Select(models.Catelog{},
		`select catelog_id,name from jzb_catelogs`)
	if err != nil {
		panic(err)
	}

	return c.Render(accounts, catelogs, bill)
}

func (c Bills) Save(bill models.Bill) revel.Result {
	c.Validate(bill)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Bills.New())
	}

	status := "add"
	if len(bill.Bill_id) > 0 {
		status = "change"
		if _, err := c.Txn.Update(&bill); err != nil {
			panic(err)
		}
	} else {
		if err := c.Txn.Insert(&bill); err != nil {
			panic(err)
		}
	}

	// 立即发送电子邮件（异步）
	amount := bill.Amount
	if bill.Amount < 0 {
		amount = -bill.Amount
	}
	account, _ := c.Txn.SelectStr(`select name from jzb_accounts WHERE account_id = ?`, bill.Account_id)
	catelog, _ := c.Txn.SelectStr(`select name from jzb_catelogs WHERE catelog_id = ?`, bill.Catelog_id)
	subject := fmt.Sprintf("%v%v %v %v￥%v", c.Message(status), c.Message("bill"), bill.Title, c.Message(bill.Kind), amount)
	body := fmt.Sprintf("%v%v%v %v:%v; %v:￥%v; %v:%v; %v:%v; %v:%v; %v:%v", c.connected().Name, c.Message(status), c.Message(bill.Kind),
		c.Message("date"), bill.Date.Format("2006-01-02"), c.Message("amount"), amount,
		c.Message("account"), account, c.Message("catelog"), catelog,
		c.Message("title"), bill.Title, c.Message("description"), bill.Description)

	if !revel.Config.BoolDefault("mode.dev", false) {
		jobs.Now(models.SendConfirmationEmail{
			Subject: subject,
			Body:    body,
		})
	}

	c.Flash.Success("%v %v %v %v%v %v!", c.Message(status), bill.Title, c.Message(bill.Kind), "￥", bill.Amount, c.Message("successed"))

	if st := c.Params.Get("st"); st == "con" {
		return c.Redirect(routes.Bills.New())
	} else {
		return c.Redirect(routes.Bills.Index())
	}
}

type Serie struct {
	Name string
	Data []float32
}

type BillsJsonData struct {
	Categories []string
	Series     []Serie
}

type billResult struct {
	Month string  `db:"month"`
	Money float32 `db:"SUM(amount)"`
}

func (c Bills) JsonData() revel.Result {
	t := c.Params.Get("t")
	if len(t) < 1 {
		return c.NotFound("404")
	}
	glog.Infof("t=%v\n", t)

	var bjd BillsJsonData
	catelogs := map[string]int{}

	incomeResults, err := c.Txn.Select(billResult{},
		`SELECT month,SUM(amount) FROM jzb_bills where kind = "income" GROUP BY month`)
	panicOnError(err)

	expenseResults, err := c.Txn.Select(billResult{},
		`SELECT month,SUM(amount) FROM jzb_bills where kind = "expense" GROUP BY month`)
	panicOnError(err)

	balanceResults, err := c.Txn.Select(billResult{},
		`SELECT month,SUM(amount) FROM jzb_bills GROUP BY month`)
	panicOnError(err)

	for _, r := range incomeResults {
		b := r.(*billResult)
		catelogs[b.Month] = 1
	}

	for _, r := range expenseResults {
		b := r.(*billResult)
		catelogs[b.Month] = 1
	}
	glog.Infof("catelogs = %v\n", catelogs)
	glog.Infof("balanceResults = %v\n", balanceResults)
	return c.RenderJson(bjd)
}

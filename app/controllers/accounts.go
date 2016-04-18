package controllers

import (
	"bytes"
	"encoding/csv"
	"strconv"

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

	//glog.Infof("account total = %v\n", total)
	return c.RenderJson(results)
}

func (c Accounts) CsvData() revel.Result {
	kind := c.Params.Get("t")
	if len(kind) < 1 {
		return c.NotFound("404")
	}
	//glog.Infof("t=%v\n", kind)

	results, err := c.Txn.Select(models.Account{},
		`select * from jzb_accounts`)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	s := make([]string, 2)
	s[0] = c.Message("account")
	s[1] = c.Message(kind)
	w.Write(s)
	if kind != "expense" && kind != "income" {
		glog.Errorf("unknow t = %v\n", kind)
		return c.NotFound("404")
	}

	for _, r := range results {
		s := make([]string, 2)
		account := r.(*models.Account)
		s[0] = account.Name
		if kind == "expense" {
			s[1] = strconv.FormatFloat(float64(-account.Expense), 'f', 0, 32)
		} else if kind == "income" && account.Income > 0 {
			s[1] = strconv.FormatFloat(float64(account.Income), 'f', 0, 32)
		}
		w.Write(s)
	}
	w.Flush()
	//glog.Infof("buf = %v\n", buf)
	return c.RenderText(buf.String())
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

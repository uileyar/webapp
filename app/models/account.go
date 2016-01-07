package models

import (
	"fmt"
	"time"

	"github.com/go-gorp/gorp"
)

/*
 [account_id] VARCHAR(36) NOT NULL,
  [server_id] INTEGER DEFAULT (0),
  [user_id] integer DEFAULT (0),
  [name] nvarchar(10) NOT NULL,
  [kind] nvarchar(10) NOT NULL,
  [amount] float NOT NULL,
  [description] NVARCHAR(200),
  [sort] INTEGER DEFAULT (0),
  [version] TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  CONSTRAINT [sqlite_autoindex_jzb_accounts_1] PRIMARY KEY ([account_id])
*/
type Account struct {
	Account_id  string `db:", size:36, primarykey"`
	Server_id   int
	User_id     int
	Name        string `db:", size:10"`
	Kind        string `db:", size:10"`
	Amount      float32
	Description string `db:", size:200"`
	Sort        int
	Version     time.Time `db:", default:CURRENT_TIMESTAMP"`
	Income      float32   `db:"-"`
	Expense     float32   `db:"-"`
	Balance     float32   `db:"-"`
}

func (u Account) String() string {
	return fmt.Sprintf("Account(%#v)", u)
}

func (u *Account) PostGet(s gorp.SqlExecutor) error {
	var val float32
	if err := s.SelectOne(&val, `SELECT sum(amount) from jzb_bills WHERE account_id = ? and kind = 'income'`, u.Account_id); err == nil {
		u.Income = val
	}

	if err := s.SelectOne(&val, `SELECT sum(amount) from jzb_bills WHERE account_id = ? and kind = 'expense'`, u.Account_id); err == nil {
		u.Expense = val
	}

	if err := s.SelectOne(&val, `SELECT sum(amount) from jzb_bills WHERE account_id = ?`, u.Account_id); err == nil {
		u.Balance = val
	}

	return nil
}

func (u *Account) PreInsert(s gorp.SqlExecutor) error {
	u.Account_id = CreateGUID()
	var val int
	if err := s.SelectOne(&val, "select max(sort) from jzb_accounts"); err == nil {
		u.Sort = val + 1
	}

	return nil
}

func (u *Account) PreUpdate(s gorp.SqlExecutor) error {

	return nil
}

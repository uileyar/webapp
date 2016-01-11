package models

import (
	"fmt"
	"math"
	"time"

	"github.com/go-gorp/gorp"
)

/*
CREATE TABLE [jzb_catelogs] (
  [catelog_id] VARCHAR(36) NOT NULL,
  [server_id] INTEGER DEFAULT (0),
  [user_id] int DEFAULT (0),
  [name] NVARCHAR(10) NOT NULL,
  [kind] nvarchar(10) NOT NULL,
  [sort] INTEGER DEFAULT (0),
  [version] TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  CONSTRAINT [sqlite_autoindex_jzb_catelogs_1] PRIMARY KEY ([catelog_id]))
*/
type Catelog struct {
	Catelog_id string `db:", size:36, primarykey"`
	Server_id  int
	User_id    int
	Name       string `db:", size:10"`
	Kind       string `db:", size:10"`
	Sort       int
	Version    time.Time `db:", default:CURRENT_TIMESTAMP"`
	Balance    float32   `db:"-"`
	Number     int       `db:"-"`
	Percent    float32   `db:"-"`
}

func Round(f float32, n int) float32 {
	pow10_n := math.Pow10(n)
	return float32(math.Trunc((float64(f)+0.5/pow10_n)*pow10_n) / pow10_n)
}

func (u *Catelog) String() string {
	return fmt.Sprintf("Catelog(%#v)", u.Name)
}

func (u *Catelog) PostGet(s gorp.SqlExecutor) error {
	var val float32
	if err := s.SelectOne(&val, `SELECT sum(amount) from jzb_bills WHERE catelog_id = ?`, u.Catelog_id); err == nil {
		u.Balance = val
	}

	var number int
	if err := s.SelectOne(&number, `SELECT count(amount) from jzb_bills WHERE catelog_id = ?`, u.Catelog_id); err == nil {
		u.Number = number
	}

	if u.Balance == 0 {
		u.Percent = 0
	} else if err := s.SelectOne(&val, `SELECT sum(amount) from jzb_bills WHERE kind = ?`, u.Kind); err == nil {
		u.Percent = (u.Balance * 100) / (val * 100) * 100
		u.Percent = Round(u.Percent, 2)
	}
	return nil
}

func (u *Catelog) PreInsert(s gorp.SqlExecutor) error {
	u.Catelog_id = CreateGUID()
	var val int
	if err := s.SelectOne(&val, "select max(sort) from jzb_catelogs"); err == nil {
		u.Sort = val + 1
	}

	return nil
}

func (u *Catelog) PreUpdate(s gorp.SqlExecutor) error {

	return nil
}

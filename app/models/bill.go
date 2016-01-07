package models

import (
	"fmt"
	"time"
)

/*
 CREATE TABLE [jzb_bills] (
  [bill_id] VARCHAR(36) NOT NULL,
  [server_id] INTEGER DEFAULT (0),
  [user_id] INTEGER DEFAULT (0),
  [amount] FLOAT NOT NULL,
  [title] NVARCHAR(50) NOT NULL,
  [description] NVARCHAR(1000) NOT NULL,
  [date] DATE NOT NULL,
  [month] INTEGER NOT NULL,
  [catelog_id] VARCHAR(36) DEFAULT (0),
  [server_category_id] INTEGER,
  [account_id] VARCHAR(36) DEFAULT (0),
  [server_account_id] INTEGER,
  [to_account_id] VARCHAR(36) DEFAULT (0),
  [server_to_account_id] INTEGER,
  [borrower_id] VARCHAR(36),
  [server_borrower_id] INTEGER,
  [kind] NVARCHAR(10) NOT NULL,
  [shared] CHAR(10) NOT NULL DEFAULT ('false'),
  [version] TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  CONSTRAINT [sqlite_autoindex_jzb_bills_1] PRIMARY KEY ([bill_id]))
*/
type Bill struct {
	Bill_id              string `db:", size:36, primarykey"`
	Server_id            int
	User_id              int
	Amount               float32
	Title                string `db:", size:50"`
	Description          string `db:", size:1000"`
	Date                 time.Time
	Month                int
	Catelog_id           string    `db:", size:36"`
	Server_category_id   int       `db:"-"`
	Account_id           string    `db:", size:36"`
	Server_account_id    int       `db:"-"`
	To_account_id        string    `db:", size:36"`
	Server_to_account_id int       `db:"-"`
	Borrower_id          string    `db:", size:36"`
	Server_borrower_id   int       `db:"-"`
	Kind                 string    `db:", size:10"`
	Shared               string    `db:", size:10`
	Version              time.Time `db:", default:CURRENT_TIMESTAMP"`
}

func (u *Bill) String() string {
	return fmt.Sprintf("Bill(%v)", u.Title)
}

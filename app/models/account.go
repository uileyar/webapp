package models

import (
	"fmt"
	"time"
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
	Sort        int    `db:", size:200, default:0"`
	Version     time.Time
}

func (u *Account) String() string {
	return fmt.Sprintf("Account(%v)", u.Name)
}

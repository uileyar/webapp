package models

import (
	"fmt"
	"time"
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
	Sort       int    `db:",  default:0"`
	Version    time.Time
}

func (u *Catelog) String() string {
	return fmt.Sprintf("Catelog(%v)", u.Name)
}

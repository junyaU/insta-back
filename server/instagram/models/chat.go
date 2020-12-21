package models

import "time"

type Chat struct {
	Id      int64     `orm:"auto"`
	From    *User     `orm:"rel(fk)"`
	To      *User     `orm:"rel(fk)"`
	Text    string    `orm:"size(100)"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

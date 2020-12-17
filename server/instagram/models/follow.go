package models

import "time"

type Follow struct {
	Id      int64     `orm:"auto"`
	Me      *User     `orm:"rel(fk)"`
	To      *User     `orm:"rel(fk)"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

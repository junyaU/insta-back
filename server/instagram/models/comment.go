package models

import "time"

type Comment struct {
	Id      int64     `orm:"auto"`
	User    *User     `orm:"rel(fk)"`
	Post    *Post     `orm:"rel(fk)"`
	Comment string    `orm:"size(20)" form:"Comment" valid:"Required;MaxSize(20)"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

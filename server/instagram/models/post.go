package models

import (
	"time"
)

type Post struct {
	Id      int64     `orm:"auto"`
	User    *User     `orm:"rel(fk)"`
	Comment string    `orm:"size(100)" form:"Comment" valid:"Required;MaxSize(100)"`
	Image   string    `orm:"size(50)" form:"Image" valid:"Required"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

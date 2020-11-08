package models

import "time"

type Imageprofile struct {
	Id      int64     `orm:"auto"`
	User    *User     `orm:"reverse(one)"`
	Image   string    `orm:"size(50)" form:"Image" valid:"Required"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

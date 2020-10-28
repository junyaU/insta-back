package models

import (
	"time"
)

type User struct {
	Id         int64     `orm:"auto"`
	Name       string    `orm:"size(10)" form:"Name" valid:"Required"`
	Email      string    `orm:"size(64)" form:"Email" valid:"Required;Email"`
	Password   string    `orm:"size(32)" form:"Password" valid:"Required;MinSize(6)"`
	Repassword string    `orm:"-" form:"Repassword" valid:"Required"`
	Posts      []*Post   `orm:"reverse(many)"`
	Created    time.Time `orm:"auto_now_add;type(datetime)"`
	Updated    time.Time `orm:"auto_now;type(datetime)"`
}

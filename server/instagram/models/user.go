package models

import (
	"time"
)

type User struct {
	Id             int64         `orm:"auto"`
	Name           string        `orm:"size(10)" form:"Name" valid:"Required"`
	Email          string        `orm:"size(64)" form:"Email" valid:"Required;Email"`
	Password       string        `orm:"size(250)" valid:"Required;MinSize(6)"`
	Imageprofile   *Imageprofile `orm:"null;rel(one);on_delete(set_null)"`
	Repassword     string        `orm:"-" form:"Repassword" valid:"Required"`
	TotalFavorited int64         `orm:"-"`
	SessionId      string        `orm:"size(100);null"`
	Posts          []*Post       `orm:"reverse(many)"`
	FavoritePosts  []*Post       `orm:"rel(m2m);rel_table(favorite)"`
	Created        time.Time     `orm:"auto_now_add;type(datetime)"`
	Updated        time.Time     `orm:"auto_now;type(datetime)"`
}

package models

type Test struct {
	Id      int64  `orm:"auto"`
	Content string `orm:"size(128)"`
}

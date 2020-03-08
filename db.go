package main

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	host string = "192.168.1.1"
	port int32 = 3306
	user string = "root"
	password string = "123456"
	dbname string = "demo"
	charset string = "utf8"
)

type Shortlink struct {
	Id int64 `db:"id"`
	Key string `db:"key"`
	Url string `db:"url"`
	CreateTime string `db:"create_time"`
}

type Link struct {
	Key string `db:"key"`
	Url string `db:"url"`
}

var Db *sqlx.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", user, password, host, port, dbname, charset)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic("数据源配置不正确: " + err.Error())
	}
	Db = db
}

func getRecordByKey(key string) (*Shortlink, error) {
	var shortlink  *Shortlink = new(Shortlink)
	er := Db.Get(shortlink,"SELECT * FROM short_link WHERE `key` = ?",key)
	if er != nil {
		return shortlink, StatusError{500,errors.New("data get buy key faied, error: "+er.Error())}
	}
	return shortlink, nil
}

func getRecordByUrl(url string) (*Shortlink, error) {
	var shortlink *Shortlink = new(Shortlink)
	er := Db.Get(shortlink,"SELECT * FROM short_link WHERE url = ?",url)
	if er != nil {
		return shortlink, nil
	}
	return shortlink, nil
}

func getRecordLastId() (int64, error) {
	var lastId int64
	err := Db.Get(&lastId,"SELECT id FROM short_link order by id desc limit 1")
	if err != nil {
		return 0, nil
	}
	return lastId, nil
}

func addRecord(links *Link) (int64, error) {
	result, err := Db.Exec("insert into short_link(`key`,url) values(?,?)",links.Key,links.Url)
	if err != nil {
		return 0, StatusError{500,errors.New("data insert faied, error: "+err.Error())}
	}
	id, _ := result.LastInsertId()
	return id, nil
}
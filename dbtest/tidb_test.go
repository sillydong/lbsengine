package dbtest

import (
	"database/sql"
	"github.com/go-gorp/gorp"
	"github.com/go-sql-driver/mysql"
	"log"
	"testing"
	"time"
)

var client *gorp.DbMap

func init() {
	dbconfig := &mysql.Config{
		User:   "root",
		Passwd: "",
		DBName: "lbsexample",
		Net:    "tcp",
		Addr:   "127.0.0.1:4000",
	}
	db, err := sql.Open("mysql", dbconfig.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	client = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}

	client.AddTable(Tidbtest{}).SetKeys(true, "id")
}

type Tidbtest struct {
	Id   int64  `db:"id"`
	Test string `db:"test"`
}

func TestAdd(t *testing.T) {
	data := &Tidbtest{Test: "hello"}
	err := client.Insert(data)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", data)
	}
}

func TestRemove(t *testing.T) {
	data := &Tidbtest{Id: 1}
	i, err := client.Delete(data)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(i)
	}
}

func TestSelect(t *testing.T) {
	var datas []Tidbtest
	x, err := client.Select(&datas, "select * from tidbtest")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", x)
		t.Logf("%+v", datas)
	}
}

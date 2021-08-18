/**
 * @Author: jinjiaji
 * @Description: 创建模拟运行环境
 * @File:  db_init
 * @Version: 1.0.0
 * @Date: 2021/7/22 下午6:46
 */

package module2

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func mockDb() {
	os.Remove("./module1.db")
	var err error

	sqlStmt := `
	create table users (id integer not null primary key, name text);
	delete from users;
	`
	_, err = GetDB().Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	sqlStmt2 := `
	insert into  users (id,name) values(1,"zzz");
	`
	_, err = GetDB().Exec(sqlStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt2)
		return
	}
}

func mockClean() {
	os.Remove("./module1.db")
}

func GetDB() *sql.DB {
	var (
		db  *sql.DB
		err error
	)

	db, err = sql.Open("sqlite3", "./module1.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

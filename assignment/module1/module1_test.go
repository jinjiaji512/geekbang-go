/*
 * @Author: jinjiaji
 * @Description:

极客时间go训练营，模块一作业：
我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

 * @File:  module1_test
 * @Version: 1.0.0
 * @Date: 2021/7/22 下午6:36
*/
package module1

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

//Users model of table "users"
type Users struct {
	ID   int64
	Name string
}

//GetUserByID 获取用户基础信息
func GetUserByID(id int64) (user *Users, ok bool, err error) {
	var (
		id2  int
		name string
	)
	err = GetDB().QueryRow(" select id,name from users where id = ? limit 1", id).Scan(&id2, &name)
	if errors.Is(err, sql.ErrNoRows) {
		ok = false
		err = nil
		return
	}
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("id:%d", id))
		return
	}
	ok = true
	user = &Users{
		ID:   id,
		Name: name,
	}
	return
}

func TestModule1(t *testing.T) {
	mockDb()
	defer mockClean()
	//测试服务器错误

	//测试获取用户信息
	user, ok, err := GetUserByID(2)
	if err != nil {
		fmt.Println(user, ok, err)
	}
	if !ok {
		fmt.Println("入参错误")
	}
}

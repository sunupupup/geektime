//我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，
//是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"practise/极客时间/作业1/dao"
	_ "practise/极客时间/作业1/tool"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {

	http.HandleFunc("/query", query)

	http.ListenAndServe(":6789", nil)
}

func query(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("name")
	agestring := r.Form.Get("age")
	age, _ := strconv.Atoi(agestring)

	person, err := dao.QueryRow(name, age)
	if err != nil {
		log.Println("查无此人")
		log.Println(errors.Unwrap(err))
		log.Println(err)
		return
	}
	log.Println("查询成功，此人的编号是", person.ID)

}

package dao

import (
	"database/sql"
	"fmt"
	"practise/极客时间/作业1/model"
	"practise/极客时间/作业1/tool"
)

func QueryRow(name string, age int) (*model.Person, error) {

	if exit, err := tool.EngineIsExit(); exit {

		//查询
		db := tool.GetEngine()
		rows := db.QueryRow("select * from user where name=? and age=?", name, age)
		var p model.Person
		err := rows.Scan(&p.ID, &p.Name, &p.Age)
		if err != nil {
			//需要wrap吗。。。
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("error: %v:%w", "query by name and age error", err)
			}
			return nil, err
		}
		return &p, nil
	} else {
		return nil, err
	}
}

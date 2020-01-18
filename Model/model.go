package Model

import (
	"database/sql"
	"fmt"
	"strconv"
)

func PutTextDB(text []byte, db *sql.DB) string {
	row, err:=db.Exec("insert into content (userContent) values (?)", text)
	if err!=nil {
		fmt.Println(err.Error())
	}
	id, e:=row.LastInsertId()
	if e!= nil {
		fmt.Println(e.Error())
	}

	i:= strconv.FormatInt(id, 10)
	return i
}
func GetTextDB(id int,db *sql.DB) []byte{
	var text []byte
	row := db.QueryRow("SELECT userContent FROM content WHERE id=?", id)
err:=row.Scan(&text)
if err!=nil{
	fmt.Println(err.Error())
}
return text



}

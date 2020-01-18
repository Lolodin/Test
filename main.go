package main

import (
	"NewTest/Encode"
	"NewTest/Model"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	_"github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

type userData struct {
	Content string
	Key string
	Id string
}
type clientResponse struct {
	DecodeText string
	Link string
	Key string
}
type templ struct {
	Text string
}
type config struct {
	DriverName string `json:"drivername"`
	Username string `json:"username"`
	Password string `json:"password"`
	Protocol string `json:"protocol"`
	Address string `json:"address"`
	DBname string `json:"dbname"`
	ADDR string `json:"addr"`
	Mailtheme string `json:"mailtheme"`

}
func getConfig() config  {
	var c config
	f, err:=os.Open("config.json")
	if err!=nil {
		fmt.Println("Error open config.json")

	}
	defer f.Close()
	b,_:=ioutil.ReadAll(f)
	_= json.Unmarshal(b, &c)

	return c


}
var(

	PORT = ":8080"
	ADDR = "localhost"+PORT
	con = getConfig()
	MAILTHEME = con.Mailtheme
	configString = con.Username+":"+con.Password+"@"+ con.Protocol+"("+con.Address+")/"+con.DBname
	db, err = sql.Open(con.DriverName, configString)
)

func main() {
	if err!=nil {
		fmt.Println(err.Error())
	}
	// Создаем таблицу в БД
	createSсhemeDB()
	//Демон для почты
	go initDaemon()

	http.HandleFunc("/ajax", ajaxHandler)
	http.HandleFunc("/id/", textHandler)
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/", indexHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.Handle("/Js/", http.StripPrefix("/Js/", http.FileServer(http.Dir("./Js/"))))

	err:= http.ListenAndServe(PORT, nil)
	if err!= nil {
		fmt.Println(err.Error())
	}


}
//Главная страница
func indexHandler (w http.ResponseWriter, r *http.Request) {
	t, err:=template.New("index.html").Delims("<<",">>").ParseFiles("Template/index.html")
	if err!=nil{
		fmt.Println(err.Error())
	}

	err =t.Execute(w, "index")

	if err!=nil{
		fmt.Println(err.Error())
	}
}
//Кодирование текста
func dataHandler (w http.ResponseWriter, r *http.Request) {
	var data userData
	databody, err:= ioutil.ReadAll(r.Body)
	if err!=nil {
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(databody, &data)
	if err!=nil {
		fmt.Println("error parce")
	}
	fmt.Println(data.Content)
	var response clientResponse
	key, textEndcode:=Encode.EncodeAes([]byte(data.Content))
	id:=Model.PutTextDB(textEndcode, db)
	response.Key = string(key)
	response.Link= "http://"+ADDR+"/id/"+id

	res,err:= json.Marshal(response)
	if err!=nil {
		fmt.Println(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)




}
//Выдача данных по id
func textHandler (w http.ResponseWriter, r *http.Request) {
	arr:=strings.Split(r.URL.Path, "/")
	fmt.Println(arr[2])
	a, err:= strconv.Atoi(arr[2])
	if err!=nil {
		return
	}
	text:=Model.GetTextDB(a, db)
t,err:=template.New("id.html").Delims("<<",">>").ParseFiles("Template/id.html")
if err!=nil {
	fmt.Println(err.Error())
}
tp:= templ{Text:string(text)}
 e:=t.Execute(w, tp)
 if e!= nil {
 	fmt.Println(e.Error())
 }
}
//Декодирование текста
func ajaxHandler (w http.ResponseWriter, r *http.Request) {
	var data userData
body, err:= ioutil.ReadAll(r.Body)
	if err!=nil {
		fmt.Println(err.Error())
	}

err = json.Unmarshal(body, &data)
fmt.Println(data, "DATA")
id, _:= strconv.Atoi(data.Id)
	text:= Model.GetTextDB(id, db)
	if len([]byte(data.Key))>16 {
		answ:= clientResponse{DecodeText:string(text)}
		b,_:= json.Marshal(answ)
		w.Write(b)
		return
	}
decodeText:=Encode.DecodeAes([]byte(data.Key), text)
answ:= clientResponse{DecodeText:string(decodeText)}
b,_:= json.Marshal(answ)
w.Write(b)


}
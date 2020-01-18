package main

import (
	"NewTest/Encode"
	"NewTest/Model"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)


// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func initDaemon() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "testmailgolangapi@gmail.com"
	DaemonSendMail(srv, user)


	}
func SendMail(srv *gmail.Service, user string) {
	// Получаем
	Thr, err:=srv.Users.Threads.List(user).Q("subject: "+ MAILTHEME).Do()
	if err!=nil{
		fmt.Println("Error Get mail")
	}
	fmt.Println(Thr.Threads)
	for _,i:=range Thr.Threads{
		//Получаем сообщение из треда
		message, err:=srv.Users.Messages.Get(user,i.Id).Do()
		if err!= nil {
			fmt.Println("1",err.Error())
			continue
		}

			//Сообщение для отправки
		var sendmessage gmail.Message
		//Если данных нет, пропускаем цикл
		if len(message.Payload.Parts) ==0 {
			continue
		}
		txt:=message.Payload.Parts[0].Body.Data
		txt = strings.Replace(txt, "_", "/", -1)
		//кодируем в обычную строку контент из сообщения
		text, e:=base64.RawStdEncoding.DecodeString(txt)
		if e!= nil {
			fmt.Println(126,e.Error())

		}
		//Если отсутсвует заголовок отправитель, пропускаем
		if len(message.Payload.Headers)<16 {
			continue
		}
		// Шифруем сообщение пользователя
		fmt.Println(txt, "Сырой текст")
		fmt.Println(message.Payload.Parts)
		fmt.Println(text)
		key, textEndcode:=Encode.EncodeAes(text)
		//Сохраняем в DB
		id:=Model.PutTextDB(textEndcode, db)
		//Генерируем тело письма
		link:="You link: "+"http://"+ADDR+"/id/"+id +"\n"+"You secret key: "+ string(key)
		//index для получение адреса
		index:=strings.Index(message.Payload.Headers[16].Value, "<")
		//Шаблон заголовков

		temp:=[]byte("From:"+user+"\r\n"+
			"In-Reply-To:" +message.Payload.Headers[16].Value +"\r\n"+
			"To:"+message.Payload.Headers[16].Value[index:]+"\r\n"+
			"Subject: Encoding\r\n"+
			"\r\n" + link)
		//Кодируем в формат base64
		sendmessage.Raw = base64.RawStdEncoding.EncodeToString(temp)
		//Удаляем лишние символы
		sendmessage.Raw = strings.Replace(sendmessage.Raw, "/", "_", -1)
		sendmessage.Raw = strings.Replace(sendmessage.Raw, "+", "-", -1)
		sendmessage.Raw = strings.Replace(sendmessage.Raw, "=", "", -1)

		//Отправляем наше сообщение
		_, err= srv.Users.Messages.Send(user, &sendmessage).Do()

		if err!=nil{
			fmt.Println(err.Error())
		}
		//Удаляем сообщение
		srv.Users.Threads.Delete(user,i.Id).Do()














	}

}
func DaemonSendMail(srv *gmail.Service, user string) {
	for  {
		SendMail(srv , user)
		time.Sleep(10*time.Second)
	}
}




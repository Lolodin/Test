# Описание работы сервиса
[Ссылка на ТЗ](https://docs.google.com/document/d/1SdPYIov3exYCFqPcQYepGJavjVheayRSO69DRjBnU3g/edit)

Для написания сервиса использовались следующие технологии
* База данных MySQL
* [API для работы с Gmail](https://godoc.org/google.golang.org/api/gmail/v1)
* Алгоритм шифрования AES
* Vue.js

Сервис представляет из себя http-сервер и демона на gmail api для работы с почтой.
Перед началом работы необходимо задать настрокий для работы сервиса, они находятся в файле config.json:
```
{
    "drivername": "mysql",
    "username": "root",
    "password": "root",
    "protocol": "tcp",
    "address" : "localhost:1994",
    "dbname": "mydata",   
    "addr" : "localhost",
    "mailtheme" : "SHA256"  
}
```
Первые 6 пунктов относятся к базе данных. 
Строка ```addr``` определяет адрес хоста сервера, необходим для генерации URL с пользовательскими данными.
Строка ```mailtheme``` определяет ключевую тему почты, из котой мы будем парсить данные для шифрования.

Работа с почтой начинается вместе с стартом сервера в отдельной горутине. 
Каждые 10 секунд в цикле начинается проверка писем с заданной темой, после получения данных, письмо удаляется, а данные шифруются и сохраняют в БД, пользователь получает ответное письмо с ссылкой и ключем шифрования

Фронт построен на компонентах Vue и ajax запросах через fetch с json в теле запроса, ответ так же получает в json

## Запуск
Для запуска в докер контейнере достаточно будет клонировать репозиторий и запустить команду: 
```docker-compose build```
```docker-compose up -d ```


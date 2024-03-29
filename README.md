# StreamTelegram
Notifications about YouTube/Twitch streams in Telegram.  
Уведомления о YouTube/Twitch стримах в Телеграм.
### Docker
```
docker build --platform linux/amd64 -t streamtelegram:latest -t streamtelegram:1.x.x .
docker run -d -v /path/files/:/app/files/ --net=host --name stg streamtelegram
```
### Environment
обязательные - *
* LOGLVL (panic, fatal, error, warn or warning, info, debug, trace. По дефолту info)
* TOKENTG* (telegram bot api token)
* USERLIST*,ADMINLIST*,ERRORLIST* (user IDs - "id,id,id")
* LOC (локация для времени, смотреть tzdata)
* PINGPORT (порт для проверки работоспособности бота, например UptimeRobot. Пример ссылки по которой будет доступ - "http://[ip]:6975/pingLaurene", отсутствие PINGPORT - сервер для пинга не запуститься.)
* PINGON (запускать ли порт для пинга)

Пример:  
LOGLVL=INFO  
TOKENTG=19209:AAFSsiJY  
USERLIST=123456789,352536  
ADMINLIST=123456789,352536   
ERRORLIST=123456789,352536  
LOC=Europe/Moscow  
PINGPORT=6975  
PINGON=true

### Папки
```
files/          (папка и рабочее место бота)
    cfg.env     (конфиг)
    logrus.log  (файл логов)
    my.db       (база данных)
```
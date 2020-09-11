# StreamTelegram
Notifications about YouTube streams in Telegram
## Environment
required parameters - *
* CONFIG_NAME* (env file)
* LOGLVL (panic, fatal, error, warn or warning, info, debug, trace)
* NAMEDB* (Database name)
* PROXY (telegram bot, socks5://login:pass@ip:port)
* TOKEN* (telegram bot api token)
* USERLIST* (user IDs, "id,id,id")
* TOID* (chat where messages will be sent)
* ERRORTOID* (chat with logs)
* YTAPIKEY* (youtube api key)
* CHANNELID* (youtube channel id. For example UC2_vpnza621Sa0cf_xhqJ8Q)
* LOC (time zone, default - UTC)
* LANGUAGETEXT (language of text, rus or eng, default - eng)
## Telegram bot commands
* status - uptime & number of RSS check iterations
* search - search channel id by the link to the video ("/seach url" or "/seach" response to a link)
* lastrss - last RSS received
* getrss - getting an rss feed by channel id
* settings - settings
* toid - edit targets for notifications 
  ### BotFather commands
status - uptime & N iterations  
search - [URL] search channel  
lastrss - last RSS received  
getrss - [channel ID] get rss feed  
settings - settings
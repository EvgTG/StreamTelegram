# StreamTelegram
Notifications about YouTube streams in Telegram
## Environment
* CONFIG_NAME
* LOGLVL (panic, fatal, error, warn or warning, info, debug, trace)
* NAMEDB 
* TOKEN (telegram bot api token)
* USERLIST (user IDs, "id,id,id")
* TOID (chat where messages will be sent)
* ERRORTOID (chat with logs)
* YTAPIKEY (youtube api key)
* CHANNELID (youtube channel id. For example UC2_vpnza621Sa0cf_xhqJ8Q)
## Telegram bot commands
* status - uptime & number of RSS check iterations
* search - search channel id by the link to the video ("/seach url" or "/seach" response to a link)

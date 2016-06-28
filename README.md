# twitter-monitor

This service was designed for monitoring Twitter accounts. If the user will tweet a message the service will send reply message from template

sample config.json
```
{
	"TM": {
		"LogLevel" : "LogLevelTrace",
		"ConsumerKey":"--------------------",
		"ConsumerSecret" : "--------------------",
		"AccessToken" : "--------------------",
		"AccessSecret" : "--------------------",

		"Usernames" : [
			"DmitryTest1902"
			
			],
		"Message"  : "Hello, test message"          
	}
	
}
```

LogLevel: must be: "LogLevelTrace" for trace messages or "LogLevelError" for output only error messages

ConsumerKey, ConsumerSecret, AccessToken, AccessSecret: Twitter account for this service

Usernames: list of monitored Twitter accounts comma separated

Message: template of reply message

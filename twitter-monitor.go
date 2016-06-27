// twitter-monitor
package main

import (
	"net/url"
	"os"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/DmitryBugrov/log"

	"twitter-monitor/cfg"
)

var (
	Log    *log.Log
	client *anaconda.TwitterApi
	config *cfg.Cfg
)

func main() {
	err := Init("./config.json")
	if !err {
		os.Exit(1)
	}
	if Monitor() {
		Log.Print(log.LogLevelError, "Exit with errors")
	}

	Log.Print(log.LogLevelTrace, "Exit")
}

func Init(configFileName string) bool {
	//Init logging
	Log = new(log.Log)
	Log.Init(log.LogLevelError, true, true, true)

	//init config
	config = new(cfg.Cfg)
	err := config.Init(Log, configFileName)
	if err != nil {
		Log.Print(log.LogLevelError, "No configuration file loaded: ", configFileName)
		return false
	}

	//update loglevel from config file
	switch config.TM.LogLevel {
	case "LogLevelError":
		Log.LogLevel = log.LogLevelError
	case "LogLevelTrace":
		Log.LogLevel = log.LogLevelTrace
	default:
		Log.Print(log.LogLevelError, "config parametr LogLevel must be LogLevelError or LogLevelTrace")
	}

	//Init twitter client
	anaconda.SetConsumerKey(config.TM.ConsumerKey)
	anaconda.SetConsumerSecret(config.TM.ConsumerSecret)
	client = anaconda.NewTwitterApi(config.TM.AccessToken, config.TM.AccessSecret)
	return true
}

func Monitor() bool {
	Log.Print(log.LogLevelTrace, "Enter to Monitor")

	Log.Print(log.LogLevelTrace, "Add to following")
	err := AddFollowUsers(config.TM.Usernames)
	if !err {
		return err
	}

	v, err := GetUsersID(config.TM.Usernames)
	if !err {
		return err
	}

	twitterStream := client.PublicStreamFilter(v)

	for {

		item := <-twitterStream.C
		switch tweet := item.(type) {
		case anaconda.Tweet:
			Log.Print(log.LogLevelTrace, "Receiving tweet:", tweet.Text)

			go SendMessage(client, Log, config.TM.Message, tweet.InReplyToStatusID) //item.(anaconda.Tweet).User.ScreenName)

		default:
			Log.Print(log.LogLevelError, "recived unknown type")
			return true
		}

	}
	return false
}

func SendMessage(client *anaconda.TwitterApi, Log *log.Log, msg string, status int64) bool {
	Log.Print(log.LogLevelTrace, "Enter to SendMessage")

	v := url.Values{}
	if status != 0 {
		v.Set("in_reply_to_status_id", strconv.FormatInt(status, 10))
	}
	_, err := client.PostTweet(msg, v)
	if err != nil {
		Log.Print(log.LogLevelError, "Error send message:", err.Error())
		return false
	} else {
		Log.Print(log.LogLevelTrace, "Send message successfully")
		return true
	}

}

func AddFollowUsers(usernames []string) bool {
	Log.Print(log.LogLevelTrace, "Enter to AddFollowUsers")
	for _, item := range usernames {
		_, err := client.FollowUser(item)

		if err != nil {
			Log.Print(log.LogLevelError, "Error add to following:", item, err.Error())
			return false
		}
	}
	return true
}

func GetUsersID(usernames []string) (url.Values, bool) {
	Log.Print(log.LogLevelTrace, "Enter to GetUsersID")
	v := url.Values{}
	result := ""
	for _, item := range config.TM.Usernames {
		user, err := client.GetUsersShow(item, v)
		if err != nil {
			Log.Print(log.LogLevelError, "Error get userid for:", item)
			return v, false
		}
		result = result + strconv.FormatInt(user.Id, 10) + ","
	}
	if len(result) == 0 {
		Log.Print(log.LogLevelError, "Error get userid")
		return v, false
	}
	v.Set("follow", result[:len(result)-1])

	Log.Print(log.LogLevelTrace, "UsersID:", result[:len(result)-1])
	return v, true
}

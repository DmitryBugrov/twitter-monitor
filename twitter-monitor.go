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

const (
	//	consumerKey    = "AerKUoivGDOTEcr7qEz3Mln0e"
	//	consumerSecret = "dwsdF8YKrFlKd5gBcfO6nruVzQYKOLIP4lieS0tOoClewfRPVb"
	//	accessToken    = "182329995-bCwddVYl4Z9GoeU0asFEyrjUuKfdEVtAhKOji5lg"
	//	accessSecret   = "4CbxJeG7wce24ctoRW3S5jRddcILyOY9wkMMZlFYrfWpL"

	//	username = "DmitryTest1902" //twitter user
	//	message  = "Hello"          //reply message
	configFileName = "./config.json"
)

var (
	Log    *log.Log
	client *anaconda.TwitterApi
	config *cfg.Cfg
)

func main() {
	err := Init()
	if !err {
		os.Exit(1)
	}
	if Monitor() {
		Log.Print(log.LogLevelError, "Exit with errors")
	}

	Log.Print(log.LogLevelTrace, "Exit")
}

func Init() bool {
	//Init logging
	Log = new(log.Log)
	Log.Init(log.LogLevelTrace, true, true, true)

	//init config
	config = new(cfg.Cfg)
	err := config.Init(Log, configFileName)
	if err != nil {
		Log.Print(log.LogLevelError, "No configuration file loaded: ", configFileName)
		return false
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
	err := AddFollowUsers()
	if err {
		return err
	}

	v, err := GetUsersID(config.TM.Usernames)
	if err {
		return err
	}

	twitterStream := client.PublicStreamFilter(v)

	for {

		item := <-twitterStream.C
		switch tweet := item.(type) {
		case anaconda.Tweet:
			Log.Print(log.LogLevelTrace, "Receiving tweet:", tweet.Text)

			go SendMessage(client, Log, config.TM.Message, item.(anaconda.Tweet).User.ScreenName)

		default:
			Log.Print(log.LogLevelError, "recived unknown type")
			return true
		}

	}
	return false
}

func SendMessage(client *anaconda.TwitterApi, Log *log.Log, msg string, user string) {
	Log.Print(log.LogLevelTrace, "Enter to SendMessage")
	_, err := client.PostDMToScreenName(msg, user)
	if err != nil {
		Log.Print(log.LogLevelError, "Error send message:", err.Error())
	} else {
		Log.Print(log.LogLevelTrace, "Send message successfully to:", user)
	}

}

func AddFollowUsers() bool {
	Log.Print(log.LogLevelTrace, "Enter to AddFollowUsers")
	for _, item := range config.TM.Usernames {
		_, err := client.FollowUser(item)

		if err != nil {
			Log.Print(log.LogLevelError, "Error add to following:", item, err.Error())
			return true
		}
	}
	return false
}

func GetUsersID(usernames []string) (url.Values, bool) {
	Log.Print(log.LogLevelTrace, "Enter to GetUsersID")
	v := url.Values{}
	result := ""
	for _, item := range config.TM.Usernames {
		user, err := client.GetUsersShow(item, v)
		if err != nil {
			Log.Print(log.LogLevelError, "Error get userid for:", item)
			return v, true
		}
		result = result + strconv.FormatInt(user.Id, 10) + ","
	}
	if len(result) == 0 {
		Log.Print(log.LogLevelError, "Error get userid")
		return v, true
	}
	v.Set("follow", result[:len(result)-1])

	Log.Print(log.LogLevelTrace, "UsersID:", result[:len(result)-1])
	return v, false
}

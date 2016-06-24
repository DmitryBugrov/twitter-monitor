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
	Monitor()

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

func Monitor() {
	Log.Print(log.LogLevelTrace, "Enter to Monitor")

	Log.Print(log.LogLevelTrace, "Add to following")
	v := url.Values{}
	_, err := client.FollowUser(config.TM.Username)

	if err != nil {
		Log.Print(log.LogLevelError, "Error add to following:", config.TM.Username)
	}
	v = url.Values{}
	user, err := client.GetUsersShow(config.TM.Username, v)
	if err != nil {
		Log.Print(log.LogLevelError, "Error get userid for:", config.TM.Username)
	}
	v.Set("follow", strconv.FormatInt(user.Id, 10))
	twitterStream := client.PublicStreamFilter(v)

	for {

		item := <-twitterStream.C
		switch tweet := item.(type) {
		case anaconda.Tweet:
			Log.Print(log.LogLevelTrace, "Receiving tweet:", tweet.Text)

			go SendMessage(client, Log, config.TM.Message, item.(anaconda.Tweet).User.ScreenName)

		default:
			Log.Print(log.LogLevelError, "recived unknown type")
		}

	}
}

func SendMessage(client *anaconda.TwitterApi, Log *log.Log, msg string, user string) {
	_, err := client.PostDMToScreenName(msg, user)
	if err != nil {
		Log.Print(log.LogLevelError, "Error send message:", err.Error())
	} else {
		Log.Print(log.LogLevelTrace, "Send message successfully to:", user)
	}

}

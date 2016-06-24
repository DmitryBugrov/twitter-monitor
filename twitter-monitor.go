// twitter-monitor
package main

import (
	"net/url"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/DmitryBugrov/log"
)

const (
	consumerKey    = "AerKUoivGDOTEcr7qEz3Mln0e"
	consumerSecret = "dwsdF8YKrFlKd5gBcfO6nruVzQYKOLIP4lieS0tOoClewfRPVb"
	accessToken    = "182329995-bCwddVYl4Z9GoeU0asFEyrjUuKfdEVtAhKOji5lg"
	accessSecret   = "4CbxJeG7wce24ctoRW3S5jRddcILyOY9wkMMZlFYrfWpL"

	username = "DmitryTest1902" //twitter user
	message  = "Hello"          //reply message
)

var (
	Log    *log.Log
	client *anaconda.TwitterApi
)

func main() {
	Init()
	Monitor()

	Log.Print(log.LogLevelTrace, "Exit")
}

func Init() {
	//Init logging
	Log = new(log.Log)
	Log.Init(log.LogLevelTrace, true, true, true)

	//Init twitter client
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	client = anaconda.NewTwitterApi(accessToken, accessSecret)
}

func Monitor() {
	Log.Print(log.LogLevelTrace, "Enter to Monitor")

	Log.Print(log.LogLevelTrace, "Add to following")
	v := url.Values{}
	_, err := client.FollowUser(username)

	if err != nil {
		Log.Print(log.LogLevelError, "Error add to following:", username)
	}
	v = url.Values{}
	user, err := client.GetUsersShow(username, v)
	if err != nil {
		Log.Print(log.LogLevelError, "Error get userid for:", username)
	}
	v.Set("follow", strconv.FormatInt(user.Id, 10))
	twitterStream := client.PublicStreamFilter(v)

	for {

		item := <-twitterStream.C
		switch tweet := item.(type) {
		case anaconda.Tweet:
			Log.Print(log.LogLevelTrace, "Receiving tweet:", tweet.Text)

			go SendMessage(client, Log, message, item.(anaconda.Tweet).User.ScreenName)

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

// twitter-monitor
package main

import (
	"net/url"
	//	"time"

	"github.com/DmitryBugrov/log"
	//	"github.com/darkhelmet/twitterstream"
	"github.com/ChimeraCoder/anaconda"
)

const (
	consumerKey      = "AerKUoivGDOTEcr7qEz3Mln0e"
	consumerSecret   = "dwsdF8YKrFlKd5gBcfO6nruVzQYKOLIP4lieS0tOoClewfRPVb"
	accessToken      = "182329995-bCwddVYl4Z9GoeU0asFEyrjUuKfdEVtAhKOji5lg"
	accessSecret     = "4CbxJeG7wce24ctoRW3S5jRddcILyOY9wkMMZlFYrfWpL"
	userids          = "182329995"
	username         = "DmitryTest1902"
	TimeForReconnect = 5000
	TimeOutForStatus = 30000
)

var (
	Log    *log.Log
	client *anaconda.TwitterApi
)

func main() {
	Init()
	Monitor(Log)

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

func Monitor(Log *log.Log) {
	Log.Print(log.LogLevelTrace, "Enter to Monitor")
	//	client := twitterstream.NewClient(consumerKey, consumerSecret, accessToken, accessSecret)
	//	go Status(Log, client)
	Log.Print(log.LogLevelTrace, "Start follow")
	_, err := client.FollowUser(username)
	if err != nil {
		Log.Print(log.LogLevelError, "Error following:", username)
	}
	v := url.Values{}
	v.Set("follow", "746103433594347520")
	twitterStream := client.PublicStreamFilter(v)
	for {

		item := <-twitterStream.C
		switch tweet := item.(type) {
		case anaconda.Tweet:
			Log.Print(log.LogLevelTrace, "Receiving tweet:")
			Log.Print(log.LogLevelTrace, tweet.Text)

		default:
			//Log.Print(log.LogLevelError, "recived unknown type")
		}

		//		Log.Print(log.LogLevelTrace, "Receiving tweet")
		//		tweets, err := client.GetFavorites(v)

		//conn, err := client.Follow(userids)

		//		if err != nil {
		//			Log.Print(log.LogLevelError, "Error getting tweets", err.Error())
		//			time.Sleep(TimeForReconnect * time.Millisecond)
		//			continue
		//		}

		//		tweet, err := conn.Next()
		//		if err != nil {
		//			Log.Print(log.LogLevelError, "Error decode tweet", err.Error())
		//		}
		//		for i := 0; i < len(tweets); i++ {
		//			Log.Print(log.LogLevelTrace, "Received tweet", tweets[i].Text)

		//		}
	}
}

//func Status(Log *log.Log, client *twitterstream.Client) {
//	Log.Print(log.LogLevelTrace, "Enter to Status")
//	for {
//		_, err := client.Sample()
//		if err != nil {
//			Log.Print(log.LogLevelError, "Error get status", err.Error())
//		} else {
//			Log.Print(log.LogLevelTrace, "get status successfully")
//		}
//		time.Sleep(TimeOutForStatus * time.Millisecond)

//	}
//}

// twitter-monitor_test.go
package main

import (
	"testing"

	//	"twitter-monitor/cfg"

	. "github.com/smartystreets/goconvey/convey"
)

var (
//	Log    *log.Log
//	client *anaconda.TwitterApi
//	config *cfg.Cfg
)

func TestTwittermonitor(t *testing.T) {

	Convey("Init", t, func() {
		So(Init("./config_test.json"), ShouldBeTrue)

		Convey("Following", nil)
		So(AddFollowUsers(config.TM.Usernames), ShouldBeTrue)

		Convey("GetUsersID", nil)
		_, err := GetUsersID(config.TM.Usernames)
		So(err, ShouldBeTrue)

		Convey("SendMessage", nil)
		So(SendMessage(client, Log, config.TM.Message, 0), ShouldBeTrue)

	})

}

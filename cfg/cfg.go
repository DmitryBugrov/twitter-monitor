// cfg.go
package cfg

import (
	"encoding/json"
	"os"

	"github.com/DmitryBugrov/log"
)

var (
	c        *Cfg
	Log      *log.Log
	Filename string
)

type Cfg struct {
	TM struct {
		LogLevel       string
		ConsumerKey    string
		ConsumerSecret string
		AccessToken    string
		AccessSecret   string

		Usernames []string //twitter user
		Message   string   //reply message
	}
}

func (c *Cfg) Init(_Log *log.Log, _filename string) error {
	Log = _Log
	Log.Print(log.LogLevelTrace, "Enter to cfg.Init")
	Filename = _filename
	err := c.load()
	return err
}

func (c *Cfg) load() error {
	Log.Print(log.LogLevelTrace, "Enter to cfg.Load")
	file, err := os.Open(Filename)
	if err != nil {
		Log.Print(log.LogLevelError, "Configuration file cannot be loaded: ", Filename)
		return err
	}
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&c)
	if err != nil {
		Log.Print(log.LogLevelError, "Unable to decode config into struct", err.Error())
		return err
	}

	return nil
}

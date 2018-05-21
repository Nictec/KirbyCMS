package db 

import (
	"gopkg.in/mgo.v2"
)

var Session *mgo.Session
var Database *mgo.Database
var err error

func Init() {
	//load the needed config
	type conf struct {
		debug   bool
		db_host string
	}
	settings := conf{
		db_host: "localhost",
	}

	//setup the connection to mongodb
	Session, err = mgo.Dial(settings.db_host)
	if err != nil {
		panic(err)
	}
	Session.SetMode(mgo.Monotonic, true)
}
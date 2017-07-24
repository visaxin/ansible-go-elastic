package metadata

import (
	"gopkg.in/mgo.v2"
)

func MongoService() {
	session, err := mgo.Dial("server1.example.com,server2.example.com")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	//c := session.DB("test").C("people")
}

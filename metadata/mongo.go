package metadata

import (
	"github.com/go-ansible-elastic-cluster/core"
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

type MongoRegister struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func (this *MongoRegister) Do(ctx core.Context, target []byte) error {
	return this.collection.Insert(target)
}

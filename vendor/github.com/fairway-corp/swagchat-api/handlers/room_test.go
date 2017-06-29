package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/utils"
	"github.com/go-zoo/bone"
)

func BenchmarkPostRoom(b *testing.B) {
	datastoreProvider := datastore.GetProvider()
	err := datastoreProvider.Connect()
	if err != nil {
		log.Println(err.Error())
	}
	datastoreProvider.Init()
	Mux = bone.New().Prefix("/" + utils.API_VERSION)
	SetRoomMux()
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	for i := 0; i < b.N; i++ {
		in := `
				{
					"name": "dennis room"
				}
			`
		reader := strings.NewReader(in)
		http.Post(ts.URL+"/"+utils.API_VERSION+"/rooms", "application/json", reader)
	}
	err = datastoreProvider.DropDatabase()
	if err != nil {
		log.Println(err.Error())
	}
}

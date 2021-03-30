package dslookup

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"context"
	"cloud.google.com/go/datastore"
	log "github.com/sirupsen/logrus"
)

type Request struct {
	Identifier string `json:"identifier"`
}

type Elements struct {
	Ip		  string   `datastore:"ip"`
	Name      string   `datastore:"name"`
}

type Return struct {
	Ip   string `json:"ip"`
}

var request []Elements

func List(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var initial Request

	err = json.Unmarshal([]byte(body), &initial)

	if err != nil {
		log.Println(err)
	}

    address := initial.Identifier

	const collection = "iptable-store"
	const projectID = "projectid"

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
      log.Println(err)
	}

	var elems []*Elements
	
	q := datastore.NewQuery("ipt").Filter("name =", address)

	keys, err := client.GetAll(ctx, q, &elems)
 	if err != nil {
 		log.Println(err)
	}
	log.Println(keys[0].Name)

	k := keys[0].Name

	group := Return{
		Ip: k,
	}

	b, err := json.Marshal(group)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(b))
}

func Entry() {
	http.HandleFunc("/list", List)
//	http.HandleFunc("/update", Update)
	http.ListenAndServe(":8080", nil)
}

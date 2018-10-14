package main

import (
	"fmt"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
	"prepaidcard/datastore"
	"prepaidcard/models"
	"prepaidcard/server"
)

var merchants = [...]models.Merchant{
	{
		ID: "amazon",
		Name: "Amazon",
		Type: "Shopping",
		Address: "Money Trail",
	},
	{
		ID: "apple",
		Name: "Apple",
		Type: "Technology",
		Address: "Lotsa Money Trail",
	},
	{
		ID: "mcdonalds",
		Name: "Mcdonalds",
		Type: "Food & Drink",
		Address: "Less Money Trail, but still got Money",
	},
}

func duplicateHelper(err error) bool {
	switch dbErr := err.(type) {
	case *pq.Error:
		if dbErr.Code == "23505" {
			return true
		}
	}
	return false
}

func main() {
	value, ok := os.LookupEnv("DB_HOST")
	if ok == false || value == "" {
		value = "localhost"
	}
	connStr := fmt.Sprintf("user=postgres host=%s dbname=postgres sslmode=disable", value)
	ds, err := datastore.New("postgres", connStr)
	if err != nil {
		panic(err)
	}
	for _, m := range merchants{
		 _, err := ds.CreateMerchant(&m)
		 if err != nil {
		 	if duplicateHelper(err) {
		 		continue
			}
			log.Fatal(err)
		 }
	}
	apiServer := server.InitServer(ds)
	apiServer.Router.Run(":8080")
}

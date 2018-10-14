package datastore

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"prepaidcard/models"
)

func New(dbType string, dbUrl string) (models.CardStore, error) {
	switch dbType {
	case "postgres":
		db, err := sqlx.Connect("postgres", dbUrl)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{"url": dbUrl}).Fatal("bad DB URL")
		}
		ds, err := InitDB(db)
		if err != nil {
			log.WithError(err).Fatal("database failed to initialise")
		}
		return ds, nil
	}
	return nil, fmt.Errorf("invalid datastore type %s", dbType)
}

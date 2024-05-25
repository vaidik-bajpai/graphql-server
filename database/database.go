package database

import (
	"log"

	"github.com/vaidik-bajpai/hackernews/prisma/db"
)

var Db *db.PrismaClient

func OpenDB() {
	db := db.NewClient()

	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}

	Db = db
}

func CloseDB() {
	Db.Prisma.Disconnect()
}

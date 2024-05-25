package data

import (
	"context"
	"log"

	"github.com/vaidik-bajpai/hackernews/database"
	"github.com/vaidik-bajpai/hackernews/prisma/db"
)

type Link struct {
	ID      int
	Title   string
	Address string
	User    User
}

func (link Link) Save() int64 {
	createdLink, err := database.Db.Link.CreateOne(
		db.Link.Title.Set(link.Title),
		db.Link.Address.Set(link.Address),
		db.Link.User.Link(
			db.User.ID.Equals(link.User.ID),
		),
	).Exec(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	return int64(createdLink.ID)
}

func GetAllLinks() []Link {
	var links []Link
	allLinks, err := database.Db.Link.FindMany().With(
		db.Link.User.Fetch(),
	).Exec(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range allLinks {
		user := item.User()
		links = append(links, Link{
			ID:      item.ID,
			Title:   item.Title,
			Address: item.Address,
			User: User{
				ID:       user.ID,
				Username: user.Name,
				Password: user.Password,
			},
		})
	}

	return links
}

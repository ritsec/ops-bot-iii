package data

import (
	"context"

	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/ent"

	_ "github.com/mattn/go-sqlite3"
)

var (
	// Client is the ent client
	Client *ent.Client

	// Ctx is the context for the ent client
	Ctx context.Context = context.Background()

	// User is the struct references user
	User *user_s = &user_s{}

	// Signin is the struct references signin
	Signin *signin_s = &signin_s{}

	// Vote is the struct references vote
	Vote *vote_s = &vote_s{}

	// VoteResult is the struct reference vote_result
	VoteResult *vote_result_s = &vote_result_s{}
)

func init() {
	client, err := ent.Open("sqlite3", "file:data.sqlite?_loc=auto&cache=shared&_fk=1")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to open database connection")
	}
	Client = client

	if err := client.Schema.Create(Ctx); err != nil {
		logrus.WithError(err).Fatal("Failed to create schema")
	}

}

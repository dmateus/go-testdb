package testmongo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func Test_Launches_Mongo(t *testing.T) {
	db := NewMongo().
		WithTest(t).
		MustStart().
		GetDB()
	type user struct {
		FirstName string `bson:"firstName"`
		LastName  string `bson:"lastName"`
	}
	_, err := db.Collection("users").InsertOne(context.Background(), user{
		FirstName: "diogo",
		LastName:  "mateus",
	})
	assert.NoError(t, err)

	resp, err := db.Collection("users").Find(context.Background(), bson.M{})
	assert.NoError(t, err)
	var users []user
	err = resp.All(context.Background(), &users)
	assert.NoError(t, err)

	assert.Equal(t, []user{
		{
			FirstName: "diogo",
			LastName:  "mateus",
		},
	}, users)
}

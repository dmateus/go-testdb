package mongo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func Test_Launches_Mongo(t *testing.T) {
	m := NewMongo().
		WithTag("5.0")
	client, err := m.Start()
	require.NoError(t, err)
	defer m.Stop()
	db := client.Database("my-database")
	type user struct {
		FirstName string `bson:"firstName"`
		LastName  string `bson:"lastName"`
	}
	_, err = db.Collection("users").InsertOne(context.Background(), user{
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

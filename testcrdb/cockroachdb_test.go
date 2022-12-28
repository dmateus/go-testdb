package testcrdb

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed migrations
var migrationsFolder embed.FS

func Test_Launches_CockroachDB(t *testing.T) {
	db := NewCockroachDB().
		WithTag("v21.2.4").
		WithMigrations(migrationsFolder).
		WithTest(t).
		MustStart().
		GetDB()

	_, _ = db.Exec(`INSERT INTO users (id, name) VALUES (1, 'diogo');`)

	rows, err := db.Query(`SELECT * FROM users`)
	assert.NoError(t, err)
	defer rows.Close()

	type user struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}
	var users []user
	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			assert.NoError(t, err)
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		assert.NoError(t, err)

	}
	assert.Equal(t, []user{
		{
			ID:   1,
			Name: "diogo",
		},
	}, users)
}

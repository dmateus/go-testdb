package mysql

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

//go:embed migrations
var migrationsFolder embed.FS

func Test_Launches_MySQL(t *testing.T) {
	mysql := NewMySQL().
		WithTag("10.5.8").
		WithMigrations(migrationsFolder)
	db, err := mysql.Start()
	require.NoError(t, err)
	defer mysql.Stop()

	_, _ = db.Exec(`INSERT INTO users (name) VALUES ('diogo');`)

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

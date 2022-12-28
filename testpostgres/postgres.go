package testmysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dmateus/go-testdb/base"
	"github.com/golang-migrate/migrate/v4"
	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/lib/pq"
	"io/fs"
	"log"
	"net/http"
	"testing"
)

const dbName = "dbname"

type Postgres struct {
	base.Base
	db           *sql.DB
	migrationsFS fs.FS
}

func NewPostgres() *Postgres {
	m := Postgres{}
	m.DockerConfigs = &base.DockerConfigs{
		Image: "postgres",
		Tag:   "11",
		EnvVars: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			fmt.Sprintf("POSTGRES_DB=%s", dbName),
			"listen_addresses = '*'",
		},
	}
	return &m
}

// WithTag Sets the image tag. Default: 11
func (m *Postgres) WithTag(tag string) *Postgres {
	m.DockerConfigs.Tag = tag
	return m
}

func (m *Postgres) WithMigrations(migrationsFS fs.FS) *Postgres {
	m.migrationsFS = migrationsFS
	return m
}

func (m *Postgres) WithTest(t *testing.T) *Postgres {
	t.Cleanup(func() {
		m.Stop()
	})
	return m
}

func (m *Postgres) Ping() error {
	err := m.db.Ping()
	return err
}

func (m *Postgres) GetPort() string {
	return "5432"
}

func (m *Postgres) Connect(port string) error {
	conn := fmt.Sprintf("postgres://user_name:secret@localhost:%s/%s?sslmode=disable", port, dbName)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return err
	}
	m.db = db

	return nil
}

func (m *Postgres) migrateUp() error {
	source, err := httpfs.New(http.FS(m.migrationsFS), "migrations")
	if err != nil {
		return err
	}
	driver, err := postgresMigrate.WithInstance(m.db, &postgresMigrate.Config{})
	if err != nil {
		return err
	}
	mgrt, err := migrate.NewWithInstance("httpfs", source, "postgres", driver)
	if err != nil {
		return err
	}
	err = mgrt.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}
	return nil
}

func (m *Postgres) MustStart() *Postgres {
	err := base.LaunchDocker(m)
	if err != nil {
		log.Fatal(err)
	}

	if m.migrationsFS != nil {
		err = m.migrateUp()
		if err != nil {
			log.Fatal(err)
		}
	}

	return m
}

func (m *Postgres) GetDB() *sql.DB {
	return m.db
}

func (m *Postgres) ResetDB() error {
	return base.ResetSQL(m.db, dbName)
}

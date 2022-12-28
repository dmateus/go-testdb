package testcrdb

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dmateus/go-testdb/base"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migrateCockroachdb "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"io/fs"
	"log"
	"net/http"
	"testing"
)

const dbName = "defaultdb"

type CockroachDB struct {
	base.Base
	db           *sql.DB
	migrationsFS fs.FS
}

func NewCockroachDB() *CockroachDB {
	m := CockroachDB{}
	m.DockerConfigs = &base.DockerConfigs{
		Image: "cockroachdb/cockroach",
		Tag:   "v21.2.4",
		Cmd:   []string{"start-single-node", "--insecure"},
	}
	return &m
}

// WithTag Sets the image tag. Default: v21.2.4
func (m *CockroachDB) WithTag(tag string) *CockroachDB {
	m.DockerConfigs.Tag = tag
	return m
}

func (m *CockroachDB) WithMigrations(migrationsFS fs.FS) *CockroachDB {
	m.migrationsFS = migrationsFS
	return m
}

func (m *CockroachDB) WithTest(t *testing.T) *CockroachDB {
	t.Cleanup(func() {
		m.Stop()
	})
	return m
}

func (m *CockroachDB) Ping() error {
	err := m.db.Ping()
	return err
}

func (m *CockroachDB) GetPort() string {
	return "26257"
}

func (m *CockroachDB) Connect(port string) error {
	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://root@localhost:%s/%s?sslmode=disable", port, dbName))
	if err != nil {
		return err
	}
	m.db = db

	return nil
}

func (m *CockroachDB) migrateUp() error {
	source, err := httpfs.New(http.FS(m.migrationsFS), "migrations")
	if err != nil {
		return err
	}
	driver, err := migrateCockroachdb.WithInstance(m.db, &migrateCockroachdb.Config{})
	if err != nil {
		return err
	}
	mgrt, err := migrate.NewWithInstance("httpfs", source, "cockroachdb", driver)
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

func (m *CockroachDB) MustStart() *CockroachDB {
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

func (m *CockroachDB) GetDB() *sql.DB {
	return m.db
}

func (m *CockroachDB) ResetDB() error {
	return base.ResetSQL(m.db, dbName)
}

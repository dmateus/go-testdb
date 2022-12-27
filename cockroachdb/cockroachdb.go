package cockroachdb

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
	"net/http"
)

type cockroachDB struct {
	base.Base
	db           *sql.DB
	migrationsFS fs.FS
}

func NewCockroachDB() *cockroachDB {
	m := cockroachDB{}
	m.DockerConfigs = &base.DockerConfigs{
		Image: "cockroachdb/cockroach",
		Tag:   "v21.2.4",
		Cmd:   []string{"start-single-node", "--insecure"},
	}
	return &m
}

// WithTag Sets the image tag. Default: v21.2.4
func (m *cockroachDB) WithTag(tag string) *cockroachDB {
	m.DockerConfigs.Tag = tag
	return m
}

func (m *cockroachDB) WithMigrations(migrationsFS fs.FS) *cockroachDB {
	m.migrationsFS = migrationsFS
	return m
}

func (m *cockroachDB) Ping() error {
	err := m.db.Ping()
	return err
}

func (m *cockroachDB) GetPort() string {
	return "26257"
}

func (m *cockroachDB) Connect(port string) error {
	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://root@localhost:%s/defaultdb?sslmode=disable", port))
	if err != nil {
		return err
	}
	m.db = db

	return nil
}

func (m *cockroachDB) migrateUp() error {
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

func (m *cockroachDB) Start() (*sql.DB, error) {
	err := base.LaunchDocker(m)

	if m.migrationsFS != nil {
		err = m.migrateUp()
		if err != nil {
			return nil, err
		}
	}

	return m.db, err
}

package testmysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dmateus/go-testdb/base"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mySqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"io/fs"
	"log"
	"net/http"
	"testing"
)

const dbName = "defaultdb"

type MySQL struct {
	base.Base
	db           *sql.DB
	migrationsFS fs.FS
}

func NewMySQL() *MySQL {
	m := MySQL{}
	m.DockerConfigs = &base.DockerConfigs{
		Image: "mysql",
		Tag:   "5.7",
		EnvVars: []string{
			"MYSQL_ROOT_PASSWORD=secret",
		},
	}
	return &m
}

// WithTag Sets the image tag. Default: 5.7
func (m *MySQL) WithTag(tag string) *MySQL {
	m.DockerConfigs.Tag = tag
	return m
}

func (m *MySQL) WithMigrations(migrationsFS fs.FS) *MySQL {
	m.migrationsFS = migrationsFS
	return m
}

func (m *MySQL) WithTest(t *testing.T) *MySQL {
	t.Cleanup(func() {
		m.Stop()
	})
	return m
}

func (m *MySQL) Ping() error {
	err := m.db.Ping()
	return err
}

func (m *MySQL) GetPort() string {
	return "3306"
}

func (m *MySQL) Connect(port string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/%s", port, dbName))
	if err != nil {
		return err
	}
	m.db = db

	return nil
}

func (m *MySQL) migrateUp() error {
	source, err := httpfs.New(http.FS(m.migrationsFS), "migrations")
	if err != nil {
		return err
	}
	driver, err := mySqlMigrate.WithInstance(m.db, &mySqlMigrate.Config{})
	if err != nil {
		return err
	}
	mgrt, err := migrate.NewWithInstance("httpfs", source, "mysql", driver)
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

func (m *MySQL) MustStart() *MySQL {
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

func (m *MySQL) GetDB() *sql.DB {
	return m.db
}

func (m *MySQL) ResetDB() error {
	return base.ResetSQL(m.db, dbName)
}

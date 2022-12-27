package mysql

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
	"net/http"
)

type mysql struct {
	base.Base
	db           *sql.DB
	migrationsFS fs.FS
}

func NewMySQL() *mysql {
	m := mysql{}
	m.DockerConfigs = base.NewDockerConfigs("mysql", "5.7", []string{
		"MYSQL_ROOT_PASSWORD=secret",
	})
	return &m
}

// WithTag Sets the image tag. Default: 5.7
func (m *mysql) WithTag(tag string) *mysql {
	m.DockerConfigs.Tag = tag
	return m
}

func (m *mysql) WithMigrations(migrationsFS fs.FS) *mysql {
	m.migrationsFS = migrationsFS
	return m
}

func (m *mysql) Ping() error {
	err := m.db.Ping()
	return err
}

func (m *mysql) GetPort() string {
	return "3306"
}

func (m *mysql) Connect(port string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", port))
	if err != nil {
		return err
	}
	m.db = db

	return nil
}

func (m *mysql) migrateUp() error {
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

func (m *mysql) Start() (*sql.DB, error) {
	err := base.LaunchDocker(m)

	if m.migrationsFS != nil {
		err = m.migrateUp()
		if err != nil {
			return nil, err
		}
	}

	return m.db, err
}

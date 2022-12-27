package mysql

//
//import (
//	"database/sql"
//	"fmt"
//	"github.com/dmateus/go-test-db/base"
//	"github.com/dmateus/go-test-db/mongo"
//	_ "github.com/go-sql-driver/mysql"
//)
//
//type mysql struct {
//	base.base
//	db *sql.DB
//}
//
//func NewMySQL() *mongo.mongo {
//	m := mongo.mongo{}
//	m.dockerConfigs = base.NewDockerConfigs("mysql", "5.7", []string{
//		"MYSQL_ROOT_PASSWORD=secret",
//	})
//	return &m
//}
//
//// WithTag Sets the image tag. Default: 5.7
//func (m *mysql) WithTag(tag string) *mysql {
//	m.dockerConfigs.Tag = tag
//	return m
//}
//
//func (m *mysql) ping() error {
//	return m.db.Ping()
//}
//
//func (m *mysql) getPort() string {
//	return "27017"
//}
//
//func (m *mysql) connect(port string) error {
//	db, err := sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", port))
//	if err != nil {
//		return err
//	}
//	m.db = db
//	return nil
//}
//
//func (m *mysql) Start() (*sql.DB, error) {
//	err := base.LaunchDocker(m)
//	return m.db, err
//}
//
//func (m *mysql) Stop() {
//	err := base.StopDocker(m)
//	if err != nil {
//		// todo: log
//	}
//}

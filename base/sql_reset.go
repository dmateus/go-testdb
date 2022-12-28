package base

import (
	"database/sql"
	"fmt"
)

func ResetSQL(db *sql.DB, databaseName string) error {
	rows, err := db.Query(fmt.Sprintf(`SELECT 'TRUNCATE ' || input_table_name || ' CASCADE;' AS truncate_query FROM(SELECT table_schema || '.' || table_name AS input_table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_catalog = '%s');`, databaseName))
	if err != nil {
		return err
	}
	var queries []string
	for rows.Next() {
		var query string
		err = rows.Scan(&query)
		if err != nil {
			return err
		}
		queries = append(queries, query)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	err = rows.Close()
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, query := range queries {
		_, err = tx.Exec(query)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

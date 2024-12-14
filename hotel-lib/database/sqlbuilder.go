package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

type Database struct {
	db *sql.DB
}

func InitConnection(database string) (*Database, error) {
	conn := fmt.Sprintf("host=localhost port=5432 user=username password=1234 dbname=%s sslmode=disable", database)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	log.Print("connected to database")
	return &Database{
		db: db,
	}, nil
}

func (d *Database) GetAll(sqlBuilder *SqlBuilder) (*sql.Rows, error) {
	rows, err := d.db.Query(sqlBuilder.query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *Database) Insert(sqlBuilder *SqlBuilder) (uuid.UUID, error) {
	var id uuid.UUID
	err := d.db.QueryRow(sqlBuilder.query).Scan(&id)
	if err != nil {
		log.Print(err.Error())
		return uuid.Nil, err
	}
	return id, nil
}

func (d *Database) Exists(sqlBuilder *SqlBuilder) (bool, error) {
	var result int
	err := d.db.QueryRow(sqlBuilder.query).Scan(&result)
	if err != nil {
		log.Print(err.Error())
		return false, err
	}
	return result > 0, nil
}

func (d *Database) Close() {
	d.db.Close()
}

type SqlBuilder struct {
	query string
}

func NewSqlBuilder() *SqlBuilder {
	return &SqlBuilder{
		query: "",
	}
}

func (s *SqlBuilder) Select(columns []string) *SqlBuilder {
	if len(columns) == 0 {
		s.query += fmt.Sprintf("select * ")
	} else {
		s.query += fmt.Sprintf("select %s ", strings.Join(columns, ","))
	}
	return s
}

func (s *SqlBuilder) From(table string) *SqlBuilder {
	s.query += fmt.Sprintf("from %s %c ", table, table[0])
	return s
}

func (s *SqlBuilder) Join(table string, condition string) *SqlBuilder {
	s.query += fmt.Sprintf("left join %s %c on %s ", table, table[0], condition)
	return s
}

func (s *SqlBuilder) Insert(table string, values []string, keys []string) *SqlBuilder {
	s.query += fmt.Sprintf("insert into %s (%s) values (%s) ", table, strings.Join(values, ","), strings.Join(keys, ","))
	return s
}

func (s *SqlBuilder) Where(condition string) *SqlBuilder {
	s.query += fmt.Sprintf("where %s ", condition)
	return s
}

func (s *SqlBuilder) Returning(name string) *SqlBuilder {
	s.query += fmt.Sprintf("returning %s", name)
	return s
}

func (s *SqlBuilder) Clear() {
	s.query = ""
}

func (s *SqlBuilder) CustomQuery(query string) *SqlBuilder {
	s.query += query
	return s
}

package mysqluser

import "gapp/repository/mysql"

type DB struct {
	conn *mysql.MySQLDB
}

func New(conn *mysql.MySQLDB) *DB {
	return &DB{
		conn: conn,
	}
}

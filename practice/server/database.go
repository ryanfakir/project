package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"context"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

func init() {
	//TODO
	db, err := sql.Open("mysql", config.DbURL())
	err = createTable(db)
	if err != nil {
		fmt.Printf("Database failed with create table with Err: %v", err)
		panic(err)
	}
}

var (
	createTableStatements = []string{
		`CREATE DATABASE IF NOT EXISTS ip;`,
		`USE ip;`,
		`CREATE TABLE IF NOT EXISTS port_t (
		pid INT(10) NOT NULL AUTO_INCREMENT,
		portnumber INT(10) NOT NULL,
		PRIMARY KEY (pid)
	);`,
		`CREATE TABLE IF NOT EXISTS history_t (
		hid INT(10) NOT NULL AUTO_INCREMENT,
		querykey VARCHAR(64) NOT NULL,
		PRIMARY KEY (hid)
	);`,
		`CREATE TABLE IF NOT EXISTS rel_t (
		rid INT(10) NOT NULL AUTO_INCREMENT,
		hid INT(10) NOT NULL,
		pid INT(10) NOT NULL ,
		PRIMARY KEY (rid),
		FOREIGN KEY (hid) REFERENCES history_t(hid),
		FOREIGN KEY (pid) REFERENCES port_t(pid)
	)ENGINE=INNODB;`,
	}
	transCreate TransactionSQL = func(tx *Tx, val interface{}) error {
		res, ok := val.(Result)
		if !ok {
			tx.Logger.Printf("Casting with Err: %v", res)
		}
		oldRIds, err := tx.SelectRidByQuerykey(res.Hostname)
		if err != nil {
			return err
		}
		for _, rid := range oldRIds {
			err = tx.DeleteRelation(rid)
			if err != nil {
				return err
			}
		}
		hid, err := tx.SelectHistoryByQuerykey(res.Hostname)
		if err != nil {
			return err
		}
		if hid == 0 {
			hid, err = tx.CreateHistory(History{Querykey: res.Hostname})
			tx.Logger.Printf("insert history id is :%d", hid)
			if err != nil {
				return err
			}
		}
		var pIds []int64
		for _, val := range res.Ports {
			pID, err := tx.SelectPortByPortNumber(val)
			if err != nil {
				return err
			}
			if pID == 0 {
				pID, err = tx.CreatePort(Port{PortNumber: val})
				tx.Logger.Printf("insert port id is :%d", pID)
				if err != nil {
					return err
				}
			}
			pIds = append(pIds, pID)
		}
		for _, pid := range pIds {
			_, err = tx.CreateRelation(Relation{HID: hid, PID: pid})
			if err != nil {
				return err
			}
		}
		return err
	}
)

// DB wrapper for abstraction
type DB struct {
	*sql.DB
}

// Tx Wrapper for abstraction
type Tx struct {
	*sql.Tx
	Logger *log.Logger
}

// TransactionSQL for excute transaction SQL
type TransactionSQL func(tx *Tx, args interface{}) error

// InitialDB for creating instance
func InitialDB() (*DB, error) {
	db, err := sql.Open("mysql", config.DbURL()+"ip")
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, err
}

// Begin starts an returns a new transaction.
func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx, Logger: log.New(os.Stderr, "", log.LstdFlags)}, nil
}

// History DAO
type History struct {
	ID       int64
	Querykey string
}

// Port DAO
type Port struct {
	ID         int64
	PortNumber int
}

// Relation DAO
type Relation struct {
	ID  int64
	HID int64
	PID int64
}

// CreateEntry creates row for each query
func (db *DB) CreateEntry(ctx context.Context, res interface{}) error {
	return Transact(ctx, db, res, transCreate)
}

//SelectRidByQuerykey atomic select
func (tx *Tx) SelectRidByQuerykey(key string) ([]int64, error) {
	var rids []int64
	stmt, err := tx.Prepare("SELECT r.rid FROM history_t h, port_t p, rel_t r where h.hid = r.hid and r.pid = p.pid and h.querykey = ?")
	if err != nil {
		return rids, err
	}
	rows, err := stmt.Query(key)
	for rows.Next() {
		var rid int64
		if err = rows.Scan(&rid); err != nil {
			return rids, err
		}
		rids = append(rids, rid)
	}
	return rids, err
}

// SelectPortByPortNumber atomic select
func (tx *Tx) SelectPortByPortNumber(portNumber int) (int64, error) {
	stmt, err := tx.Prepare("SELECT pid FROM port_t where portnumber = ?")
	if err != nil {
		return 0, err
	}
	return readBoilerplate(stmt, portNumber)
}

// SelectHistoryByQuerykey atomic select
func (tx *Tx) SelectHistoryByQuerykey(key string) (int64, error) {
	stmt, err := tx.Prepare("SELECT hid FROM history_t where querykey = ?")
	if err != nil {
		return 0, err
	}
	return readBoilerplate(stmt, key)
}

// SelectRelationByID atomic select
func (tx *Tx) SelectRelationByID(pid, hid int64) (int64, error) {
	stmt, err := tx.Prepare("SELECT rid FROM rel_t where pid = ? and hid = ?")
	if err != nil {
		return 0, err
	}
	return readBoilerplate(stmt, pid, hid)
}

//DeleteHistory atomic delete
func (tx *Tx) DeleteHistory(hid int64) error {
	stmt, err := tx.Prepare("DELETE FROM history_t WHERE hid = ?")
	if err != nil {
		return err
	}
	return deleteBoilerplate(stmt, hid)
}

//DeleteRelation atomic delete
func (tx *Tx) DeleteRelation(rid int64) error {
	stmt, err := tx.Prepare("DELETE FROM rel_t WHERE rid = ?")
	if err != nil {
		return err
	}
	return deleteBoilerplate(stmt, rid)
}

//DeletePort atomic delete
func (tx *Tx) DeletePort(pid int64) error {
	stmt, err := tx.Prepare("DELETE FROM port_t WHERE pid = ?")
	if err != nil {
		return err
	}
	return deleteBoilerplate(stmt, pid)
}

// CreateHistory atomic insert
func (tx *Tx) CreateHistory(u History) (int64, error) {
	stmt, err := tx.Prepare("INSERT `history_t` SET `querykey` = ?")
	if err != nil {
		return 0, err
	}
	return writeBoilerplate(stmt, u.Querykey)
}

// CreatePort atomic insert
func (tx *Tx) CreatePort(p Port) (int64, error) {
	stmt, err := tx.Prepare("INSERT `port_t` SET `portnumber` = ?")
	if err != nil {
		return 0, err
	}
	return writeBoilerplate(stmt, p.PortNumber)
}

// CreateRelation atomic insert
func (tx *Tx) CreateRelation(rel Relation) (int64, error) {
	stmt, err := tx.Prepare("INSERT `rel_t` SET `hid` = ?, `pid` = ?")
	if err != nil {
		return 0, err
	}
	return writeBoilerplate(stmt, rel.HID, rel.PID)
}

func writeBoilerplate(stmt *sql.Stmt, args ...interface{}) (int64, error) {
	r, err := execSQL(stmt, args...)
	if err != nil {
		return 0, err
	}
	lastInsertID, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

func readBoilerplate(stmt *sql.Stmt, args ...interface{}) (int64, error) {
	rows, err := querySQL(stmt, args...)
	if err != nil {
		return 0, err
	}
	var id int64
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return 0, err
		}
	}
	return id, err
}

func deleteBoilerplate(stmt *sql.Stmt, args ...interface{}) error {
	r, err := execSQL(stmt, args...)
	if err != nil {
		return err
	}
	_, err = r.RowsAffected()
	if err != nil {
		return err
	}
	return err
}

func execSQL(stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	r, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func querySQL(stmt *sql.Stmt, args ...interface{}) (*sql.Rows, error) {
	r, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// createTable creates the table, and if necessary
func createTable(db *sql.DB) error {
	for _, stmt := range createTableStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// Transact as wrapper extract logic
func Transact(ctx context.Context, db *DB, res interface{}, txSQL TransactionSQL) error {
	errc := make(chan error, 1)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	go func() {
		defer close(errc)
		errc <- txSQL(tx, res)
	}()
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	case err := <-errc:
		return err
	}
}

package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	DB_DRIVER_NAME = "mysql"
	MAX_CONNECTION = 25
)

var (
	db *sql.DB
)

func InitDB() (err error) {
	dsn := config.GetDsn()

	// open
	Assert(dsn != "", "dsn must not be empty")
	db, err = sql.Open(DB_DRIVER_NAME, dsn)
	if err != nil {
		err = errors.Wrapf(err, "failed to open db: %+v", err)
		return
	}

	// ping
	Assert(db != nil, "db must not be empty")
	if err = db.Ping(); err != nil {
		err = errors.Wrapf(err, "failed to ping with db: %+v", err)
		return
	}

	// set db connection pool
	db.SetMaxOpenConns(MAX_CONNECTION)
	db.SetMaxIdleConns(MAX_CONNECTION)
	db.SetConnMaxLifetime(5 * time.Minute)
	return
}

func GenTables() (err error) {
	err = genSessionTable()
	if err != nil {
		err = errors.Wrapf(err, "failed to create tbl (%s): %+v", TblAlpsTokens, err)
		return
	}
	return
}

func genSessionTable() (err error) {
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
	token VARCHAR(1024) NOT NULL,
	createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	INDEX idxCreatedAt (createdAt DESC)
)`, TblAlpsTokens)
	_, err = db.Exec(query)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return err
}

func GetLatestAlpsToken() (token string, err error) {
	query := fmt.Sprintf(`
SELECT 
	token 
FROM %s 
ORDER BY createdAt DESC	
LIMIT 1
`, TblAlpsTokens)
	err = db.QueryRow(query).Scan(&token)
	if err != nil {
		err = errors.Wrap(err, "failed to get latest token")
		return
	}
	return
}

func GetOrderNums(orderStat OrderStat) (orderNums []string, err error) {
	orderNums = []string{}
	query := fmt.Sprintf(`
SELECT 
	o_num 
FROM %s 
WHERE o_stat = ?
`, TblDelivery)
	rows, err := db.Query(query, int(orderStat))
	if err != nil {
		err = errors.Wrapf(err, "failed to query order nums for stat: %+v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var oNum string
		if err = rows.Scan(&oNum); err != nil {
			err = errors.Wrap(err, "failed to scan o_num")
			return
		}
		orderNums = append(orderNums, oNum)
	}

	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func UpdateOrderInvoice(orderNum, invoice string) (cnt int64, err error) {
	Assert(orderNum != "", "order number must not be empty")
	Assert(invoice != "", "invoice must not be empty")

	query := fmt.Sprintf(`
UPDATE %s 
SET o_no = ?,
	o_stat = 5
WHERE o_num = ?
`, TblDelivery)
	var result sql.Result
	result, err = db.Exec(query, invoice, orderNum)
	if err != nil {
		err = errors.Wrapf(err, "failed to update invoice (orderNum: %s, invoice: %s)", orderNum, invoice)
		return
	}
	cnt, err = result.RowsAffected()
	if err != nil {
		err = errors.Wrap(err, "failed to get rows affected")
		return
	}
	return
}

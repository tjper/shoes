package main

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

// Error codes returned by failures within library.
var (
	ErrTrueToSizeInvalid = errors.New("TrueToSize needs to be <= 5 && >= 1.")
	ErrZeroShoesIds      = errors.New("Zero shoes_ids passed to SelectTrueToSizeByShoesId.")
)

// Conn establishes a connection with the shoes db, if successful,
// returns db connection.
func Db() *sql.DB {
	connStr := "dbname=shoes user=shoes host=shoesdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		LogErr.Fatalln(err)
	}
	return db
}

// DbClient serves as an access point for methods that access a Db.
type DbClient struct {
	Db interface {
		Exec(string, ...interface{}) (sql.Result, error)
		Query(string, ...interface{}) (*sql.Rows, error)
	}
}

// InsertTrueToSize inserts shoeId and truetosize into DbClient.Db, if successful,
// returns the number of sets inserted.
func (c DbClient) InsertTrueToSize(shoeId, trueToSize int) (int, error) {
	if trueToSize < 1 || trueToSize > 5 {
		LogErr.Println(ErrTrueToSizeInvalid)
		return 0, ErrTrueToSizeInvalid
	}

	query := "INSERT INTO truetosize (shoes_id, truetosize) VALUES ($1, $2);"

	res, err := c.Db.Exec(query, shoeId, trueToSize)
	if err != nil {
		LogErr.Println(err)
		return 0, err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		LogErr.Println(err)
		return 0, err
	}
	return int(aff), nil
}

// SelectTrueToSizeByShoeId retrieves truetosize set by shoeId from DbClient.Db, if successful,
// returns the set.
func (c DbClient) SelectTrueToSizeByShoeId(shoeId int) ([]int, error) {

	query := `SELECT truetosize
		  FROM truetosize 
		  WHERE shoes_id = $1`

	rows, err := c.Db.Query(query, shoeId)
	if err != nil {
		LogErr.Println(err)
		return nil, err
	}
	defer rows.Close()

	res := make([]int, 0)
	for rows.Next() {
		var tts int
		if err := rows.Scan(&tts); err != nil {
			LogErr.Println(err)
			return nil, err
		}
		res = append(res, tts)
	}
	if err := rows.Err(); err != nil {
		LogErr.Println(err)
		return nil, err
	}
	return res, nil
}

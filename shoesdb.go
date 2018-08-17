// Package shoes provides a library for interacting with the shoes Postgres database.
package shoes

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/lib/pq"
)

// Error codes returned by failures within library.
var (
	ErrZeroTrueToSizes = errors.New("Zero truetosizes passed to InsertTrueToSizes.")
	ErrZeroShoesIds    = errors.New("Zero shoes_ids passed to SelectTrueToSizeByShoesId.")
)

// Conn establishes a connection with the shoes db, if successful,
// returns db connection.
func Db() *sql.DB {
	connStr := "user=shoes dbname=shoes sslmode=verify-full"
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

// InsertTrueToSize inserts shoeId and truetosize sets into DbClient.Db, if successful,
// returns the number of sets inserted.
func (c DbClient) InsertTrueToSizes(shoeTrueToSizes map[int][]int) (int, error) {
	if len(shoeTrueToSizes) == 0 {
		LogErr.Println(ErrZeroTrueToSizes)
		return 0, ErrZeroTrueToSizes
	}

	// generate query string and arguements
	query := "INSERT INTO truetosize (shoes_id, truetosize) VALUES "
	values := make([]string, 0)
	args := make([]interface{}, 0)

	for shoe, tts := range shoeTrueToSizes {
		for _, v := range tts {
			values = append(values, "(?, ?)")
			args = append(args, shoe, v)
		}
	}

	query += strings.Join(values, ",") + ";"

	res, err := c.Db.Exec(query, args...)
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

// SelectTrueToSizeByShoesId retrieves shoeId and truetosize sets from DbClient.Db, if successful,
// returns the sets.
func (c DbClient) SelectTrueToSizeByShoesId(shoesIds []int) (map[int][]int, error) {
	if len(shoesIds) == 0 {
		return nil, ErrZeroShoesIds
	}

	query := `SELECT shoes_id, truetosize
		  FROM truetosize 
		  WHERE shoes_id IN (`
	values := make([]string, 0)
	args := make([]interface{}, 0)
	for _, id := range shoesIds {
		values = append(values, "?")
		args = append(args, id)
	}

	query += strings.Join(values, ", ") + ");"

	rows, err := c.Db.Query(query, args...)
	if err != nil {
		LogErr.Println(err)
		return nil, err
	}
	defer rows.Close()

	res := make(map[int][]int)
	for rows.Next() {
		var shoes_id, truetosize int
		if err := rows.Scan(&shoes_id, &truetosize); err != nil {
			LogErr.Println(err)
			return nil, err
		}
		res[shoes_id] = append(res[shoes_id], truetosize)
	}
	if err := rows.Err(); err != nil {
		LogErr.Println(err)
		return nil, err
	}
	return res, nil
}

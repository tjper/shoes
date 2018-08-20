package main

import (
	"errors"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
)

var ErrMock = errors.New("Mock Error")

func Test_Conn(t *testing.T) {
	db := Db()
	dbType := reflect.TypeOf(db)
	if dbType.String() != "*sql.DB" {
		t.Errorf("Expected type = *sql.DB, Actual type = %v\n", dbType.String())
	}
}

func Test_InsertTrueToSizes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		shoeId      int
		trueToSize  int
		expectedRes int
		expectedErr error
	}{
		{"trueToSize < 1.", 1, 0, 0, ErrTrueToSizeInvalid},
		{"trueToSize > 5.", 1, 6, 0, ErrTrueToSizeInvalid},
		{"trueToSize valid, shoeId valid.", 1, 3, 1, nil},
		{"trueToSize valid, shoeId DNE.", 600, 3, 0, ErrMock},
	}

	mock.ExpectExec("INSERT INTO truetosize").WithArgs(1, 3).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO truetosize").WithArgs(600, 3).WillReturnError(ErrMock)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := DbClient{Db: db}
			res, err := c.InsertTrueToSize(test.shoeId, test.trueToSize)
			if res != test.expectedRes {
				t.Errorf("Expected Result = %v, Actual Result = %v\n", test.expectedRes, res)
			}
			if err != test.expectedErr {
				t.Errorf("Expected Error = %v, Actual Error = %v\n", test.expectedErr, err)
			}
		})
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func Test_SelectTrueToSize(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		shoeId      int
		expectedRes []int
		expectedErr error
	}{
		{"Valid shoeId", 1, []int{1, 2, 3}, nil},
		{"Invalid shoeId", -1, []int{}, nil},
	}

	rows := sqlmock.NewRows([]string{"truetosize"}).AddRow(1).AddRow(2).AddRow(3)
	mock.ExpectQuery(`SELECT truetosize FROM truetosize`).WithArgs(1).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"truetosize"})
	mock.ExpectQuery(`SELECT truetosize FROM truetosize`).WithArgs(-1).WillReturnRows(rows)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := DbClient{Db: db}
			res, err := c.SelectTrueToSizeByShoeId(test.shoeId)
			if err != test.expectedErr {
				t.Errorf("Expected Error = %v, Actual Error = %v\n", test.expectedErr, err)
			}
			if !reflect.DeepEqual(test.expectedRes, res) {
				t.Errorf("Expected Res = %v, Actual Res = %v\n", test.expectedRes, res)
			}
		})
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

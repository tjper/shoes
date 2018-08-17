package shoes

import (
	"reflect"
	"testing"
	//	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func Test_Conn(t *testing.T) {
	db := Db()
	dbType := reflect.TypeOf(db)
	if dbType.String() != "*sql.DB" {
		t.Errorf("Expected type = *sql.DB, Actual type = %v\n", dbType.String())
	}
}

var (
	shoeSizes = map[int][]int{
		1: []int{1, 2, 3},
	}
	shoesSizes = map[int][]int{
		1: []int{1, 2},
		2: []int{2, 4},
	}
)

func Test_InsertTrueToSizes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	type test struct {
		name            string
		expectedErr     error
		expectedRes     int
		shoeTrueToSizes map[int][]int
	}
	tests := []test{
		test{"Zero trueToSizes.", ErrZeroTrueToSizes, 0, map[int][]int{}},
		test{"One shoe, three sizes.", nil, 3, shoeSizes},
		test{"Two shoes, two sizes.", nil, 4, shoesSizes},
	}

	mock.ExpectExec("INSERT INTO truetosize").WithArgs(1, 1, 1, 2, 1, 3).WillReturnResult(sqlmock.NewResult(3, 3))
	mock.ExpectExec("INSERT INTO truetosize").WithArgs(1, 1, 1, 2, 2, 2, 2, 4).WillReturnResult(sqlmock.NewResult(4, 4))

	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			c := DbClient{Db: db}
			res, err := c.InsertTrueToSizes(te.shoeTrueToSizes)
			if res != te.expectedRes {
				t.Errorf("Expected Result = %v, Actual Result = %v\n", te.expectedRes, res)
			}
			if err != te.expectedErr {
				t.Errorf("Expected Error = %v, Actual Error = %v\n", te.expectedErr, err)
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

var (
	oneShoe  = []int{1}
	twoShoes = []int{1, 2}
)

func Test_SelectTrueToSize(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	type test struct {
		name        string
		expectedRes map[int][]int
		expectedErr error
		shoesIds    []int
	}
	tests := []test{
		test{"Zero shoes.", map[int][]int{}, ErrZeroShoesIds, []int{}},
		test{"One shoe.", shoeSizes, nil, oneShoe},
		test{"Two shoes.", shoesSizes, nil, twoShoes},
	}

	rows := sqlmock.NewRows([]string{"shoes_id", "truetosize"}).AddRow(1, 1).AddRow(1, 2).AddRow(1, 3)
	mock.ExpectQuery(`SELECT shoes_id, truetosize FROM truetosize`).WithArgs(1).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"shoes_id", "truetosize"}).AddRow(1, 1).AddRow(1, 2).AddRow(2, 2).AddRow(2, 4)
	mock.ExpectQuery(`SELECT shoes_id, truetosize FROM truetosize`).WithArgs(1, 2).WillReturnRows(rows)

	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			c := DbClient{Db: db}
			res, err := c.SelectTrueToSizeByShoesId(te.shoesIds)
			if err != te.expectedErr {
				t.Errorf("Expected Error = %v, Actual Error = %v\n", te.expectedErr, err)
				t.SkipNow()
			}
			if !mapsEqual(res, te.expectedRes) {
				t.Errorf("Expected Res = %v, Actual Res = %v\n", te.expectedRes, res)
			}
		})
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func mapsEqual(m, s map[int][]int) bool {
	for id, sizes := range m {
		sizes2, exists := s[id]
		if !exists {
			return false
		}
		for i, size := range sizes {
			if size != sizes2[i] {
				return false
			}
		}
	}
	return true
}

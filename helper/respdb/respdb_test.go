package respdb

import (
	"reflect"
	"testing"

	"github.com/GeoNet/delta/resp"
	"github.com/hashicorp/go-memdb"
)

func TestRespDB_Schema(t *testing.T) {

	db, err := memdb.NewMemDB(NewSchema())
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		t string
		i string
		a []interface{}
		r interface{}
	}{

		{"datalogger", "id", []interface{}{"TEST"}, resp.DataloggerModel{Name: "TEST", Type: "TEST"}},
		{"sensor", "id", []interface{}{"TEST"}, resp.SensorModel{Name: "TEST", Type: "TEST"}},
		{"response", "id", []interface{}{"TEST"}, resp.Response{Name: "TEST"}},
	}

	txn := db.Txn(true)
	for _, test := range tests {
		if err := txn.Insert(test.t, test.r); err != nil {
			t.Fatal(err)
		}
	}
	txn.Commit()

	txn = db.Txn(false)
	defer txn.Abort()

	for _, test := range tests {
		v, err := txn.First(test.t, test.i, test.a...)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(v, test.r) {
			t.Errorf("db invalid check: %s [%s] got %v, expected %v", test.t, test.i, v, test.r)
		}
	}

}

func TestMetaDB_Load(t *testing.T) {
	db, err := NewRespDB("../..")
	if err != nil {
		t.Fatal(err)
	}
	txn := db.Txn(false)
	defer txn.Abort()

	var tests = []struct {
		t string
		i string
		a []interface{}
	}{
		{"datalogger", "id", []interface{}{"BASALT"}},
		{"sensor", "id", []interface{}{"L4C"}},
		{"response", "id", []interface{}{"Short Period"}},
	}

	for _, test := range tests {
		v, err := txn.First(test.t, test.i, test.a...)
		if err != nil {
			t.Fatal(err)
		}
		if v == nil {
			t.Errorf("db invalid check: %s [%s] unable to find entry: %v", test.t, test.i, test.a)
		}
	}
}

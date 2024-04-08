package postgres

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"parking/api/storage"
)

type ParkingLotTestStruct struct {
	methodName string
	desc       string
	in         interface{}
	want       interface{}
	tops       cmp.Options
	wantErr    bool
}

var (
	parkingID                    string
	createParkingLotTestData = &storage.ParkingLotResponse{
		Name:       "Parking Lot Sample",
		TotalSpace: 15,
	}

	listParkingLotTestData = []*storage.ParkingLotResponse{
		{
			Name:       "Parking Lot Sample",
			TotalSpace: 15,
			Total:      1,
		},
	}
)

func TestParkingLot(t *testing.T) {
	t.Parallel()
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("missing database connection")
	}
	st, cl := NewTestStorage(conn, filepath.Join("..", "..", "migrations"))
	t.Cleanup(cl)

	// NOTE: Can be added scenarios here to test more acceptance criteria
	tests := []ParkingLotTestStruct{
		{
			methodName: "CREATE_PARKING",
			desc:       "Success Create ParkingLot",
			wantErr:    false,
			tops: cmp.Options{
				cmpopts.IgnoreFields(storage.ParkingLotResponse{}, "ID", "CreatedAt", "UpdatedAt"),
			},
			in: &storage.InsertParkingLotRequest{
				Name:        "Parking Lot Sample",
				TotalSpaces: 15,
			},
			want: createParkingLotTestData,
		},
		{
			methodName: "GET_PARKING",
			desc:       "Success List ParkingLot",
			wantErr:    false,
			tops: cmp.Options{
				cmpopts.IgnoreFields(storage.ParkingLotResponse{}, "ID", "CreatedAt", "UpdatedAt"),
			},
			in:   &storage.ListParkingLotRequest{},
			want: listParkingLotTestData,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			switch test.methodName {
			case "CREATE_PARKING":
				CreateParkingLot(t, test, st)
			case "GET_PARKING":
				ListParkingLot(t, test, st)

			}
		})
	}
}

func CreateParkingLot(t *testing.T, test ParkingLotTestStruct, st *Storage) {
	ctx := context.Background()
	req, ok := test.in.(*storage.InsertParkingLotRequest)
	if !ok {
		t.Error("request type conversion error")
	}

	got, err := st.InsertParkingLot(ctx, req)
	if err != nil {
		t.Error("failed to create parking lot")
	}

	parkingID = got.ID
	if !test.wantErr {
		o := test.tops
		if !cmp.Equal(test.want, got, o) {
			t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
		}
	}
}

func ListParkingLot(t *testing.T, test ParkingLotTestStruct, st *Storage) {
	ctx := context.Background()
	req, ok := test.in.(*storage.ListParkingLotRequest)
	if !ok {
		t.Error("request type conversion error")
	}
	got, err := st.ListParkingLot(ctx, req)
	if err != nil && !test.wantErr {
		t.Fatal(err)
	}

	if !test.wantErr {
		o := test.tops
		if !cmp.Equal(test.want, got, o) {
			t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
		}
	}
}

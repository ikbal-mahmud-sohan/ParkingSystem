package storage

import (
	"database/sql"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// NotFound is returned when the requested resource does not exist.
	NotFound = status.Error(codes.NotFound, "not found")
	// Conflict is returned when trying to create the same resource twice.
	Conflict = status.Error(codes.AlreadyExists, "conflict")
	// Triggers when request arguments are invalid
	InvalidArgument = status.Error(codes.InvalidArgument, "invalid arguments")
)

const (
	PAGESIZE int = 10
)

// OrderBy represents an order for returning query results.
type OrderBy string

const (
	OrderByAscending  OrderBy = "ASC"
	OrderByDescending OrderBy = "DESC"
)

type SortByColumn string

const (
	SortByColumn_Date SortByColumn = "created_at"
	SortByColumn_Name SortByColumn = "name"
)

// Parking lot data models
type (
	InsertParkingLotRequest struct {
		Name        string `db:"name"`
		TotalSpaces int32  `db:"total_spaces"`
	}

	ParkingLotResponse struct {
		ID         string    `db:"id"`
		Name       string    `db:"name"`
		TotalSpace int32     `db:"total_spaces"`
		CreatedAt  time.Time `db:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"`
		Total      int32     `db:"total"`
	}
	ParkingLotDetailsRequest struct {
		ID string `db:"id"`
	}

	ParkingLotDetailsResponse struct {
		ID          string         `db:"id"`
		Name        string         `db:"name"`
		TotalSpace  int32          `db:"total_spaces"`
		VehicleNo   sql.NullString `db:"vehicle_no"`
		Entry       sql.NullTime   `db:"entry"`
		SlotID      sql.NullString         `db:"slot_id"`
		TotalParked int32          `db:"total_parked"`
		CreatedAt   time.Time      `db:"created_at"`
		UpdatedAt   time.Time      `db:"updated_at"`
	}

	ListParkingLotRequest struct {
		OrderBy      OrderBy
		SearchTerm   string
		SortByColumn SortByColumn
		Limit        int
		Offset       int
	}
)

type SlotStatus string

const (
	SlotStatusAVAILABLE        SlotStatus = "AVAILABLE"
	SlotStatusUNDERMAINTENANCE SlotStatus = "UNDERMAINTENANCE"
	SlotStatusBOOKED           SlotStatus = "BOOKED"
)

func (s SlotStatus) String() string { return string(s) }

// Parking slot data models
type (
	UpdateParkingSlotRequest struct {
		LotID      string     `db:"lot_id"`
		SlotNumber int32      `db:"slot_number"`
		Status     SlotStatus `db:"status"`
	}

	ParkingSlotResponse struct {
		ID         string     `db:"id"`
		LotID      string     `db:"lot_id"`
		SlotNumber int32      `db:"slot_number"`
		Status     SlotStatus `db:"status"`
		TotalSpace int32      `db:"total_spaces"`
		CreatedAt  time.Time  `db:"created_at"`
		UpdatedAt  time.Time  `db:"updated_at"`
	}
)

// Vehicle data models
type (
	ParkVehicleRequest struct {
		VehicleNo string    `db:"vehicle_no"`
		Entry     time.Time `db:"entry"`
		SlotID    string    `db:"slot_id"`
	}

	UnParkVehicleRequest struct {
		ID          string    `db:"id"`
		Exit        time.Time `db:"exit"`
		ParkedHours int32     `db:"parked_hours"`
		Fee         int32     `db:"fee"`
		UpdatedAt   time.Time `db:"updated_at"`
	}

	VehicleRequest struct {
		ID        string `db:"id"`
		VehicleNo string `db:"vehicle_no"`
	}

	VehicleResponse struct {
		ID          string       `db:"id"`
		VehicleNo   string       `db:"vehicle_no"`
		Entry       time.Time    `db:"entry"`
		Exit        sql.NullTime `db:"exit"`
		ParkedHours int32        `db:"parked_hours"`
		Fee         int32        `db:"fee"`
		SlotID      string       `db:"slot_id"`
		SlotNumber  int32        `db:"slot_number"`
		LotID       string       `db:"lot_id"`
		CreatedAt   time.Time    `db:"created_at"`
		UpdatedAt   time.Time    `db:"updated_at"`
	}
)

type (
	Statement struct {
		TotalEarning sql.NullInt32 `db:"total_earning"`
		TotalParked  sql.NullInt32 `db:"total_parked"`
		TotalHours   sql.NullInt32 `db:"total_hours"`
	}
)

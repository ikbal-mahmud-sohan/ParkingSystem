package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"parking/api/storage"
	"parking/utility/logging"
)

const insertParkingLot = `
	INSERT INTO parking_lot (
		name,
		total_spaces
	) VALUES (
		:name,
		:total_spaces
	) RETURNING * `

func (s *Storage) InsertParkingLot(ctx context.Context, req *storage.InsertParkingLotRequest) (*storage.ParkingLotResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.InsertParkingLot")
	if err := s.InsertParkingLotValidation(ctx, req); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, storage.InvalidArgument
	}
	pstmt, err := s.db.PrepareNamedContext(ctx, insertParkingLot)
	if err != nil {
		logging.WithError(err, log).Error("insert ParkingLot")
		return nil, err
	}
	defer pstmt.Close()

	res := &storage.ParkingLotResponse{}
	if err := pstmt.Get(res, req); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			logging.WithError(err, log).Error("failed to insert ParkingLot to duplicate row")
			return nil, storage.Conflict
		}
		logging.WithError(err, log).Error("failed to insert ParkingLot")
		return nil, err
	}

	return res, nil
}

const listParkingLot = `SELECT *, COUNT(*) OVER() as total FROM parking_lot WHERE name != ''`

func (s *Storage) ListParkingLot(ctx context.Context, req *storage.ListParkingLotRequest) ([]*storage.ParkingLotResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.ListParkingLot")
	data, statement, arguments := []*storage.ParkingLotResponse{}, listParkingLot, []interface{}{}

	if req.SearchTerm != "" {
		statement += " AND name = ?"
		arguments = append(arguments, req.SearchTerm)
	}

	orderBy := " ORDER BY %s DESC"
	if req.OrderBy == storage.OrderByAscending {
		orderBy = " ORDER BY %s ASC"
	}
	sortByColumn := storage.SortByColumn_Date
	if req.SortByColumn == storage.SortByColumn_Name {
		sortByColumn = storage.SortByColumn_Name
	}
	statement += fmt.Sprintf(orderBy, sortByColumn)
	if req.Limit > 0 {
		statement += fmt.Sprintf(" LIMIT NULLIF(%d, 0) OFFSET %d;", req.Limit, req.Offset)
	}

	completeQuery, args, err := sqlx.In(statement, arguments...)
	if err != nil {
		log.WithError(err).Error("failed to sqlx")
	}
	if err := s.db.SelectContext(ctx, &data, s.db.Rebind(completeQuery), args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		} else {
			log.WithError(err).Error("failed to get parking lot list")
			return nil, err
		}
	}
	return data, nil
}

func (s *Storage) GetParkinglotByID(ctx context.Context, lotID string) (*storage.ParkingLotResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.GetParkinglotByID")
	getParkingLot := `SELECT * FROM parking_lot WHERE id = $1`
	var parkingLot storage.ParkingLotResponse
	if err := s.db.GetContext(ctx, &parkingLot, getParkingLot, lotID); err != nil {
		logging.WithError(err, log).Error("failed to get parking lot in storage")
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}

	return &parkingLot, nil
}

const getParkingLot = `
SELECT 
	pl.id,
	pl.name,
	pl.total_spaces,
	v.vehicle_no,
	v.entry,
	v.slot_id,
	pl.created_at,
	pl.updated_at,
	COUNT(pl.*) OVER() AS total_parked
FROM parking_lot pl
	LEFT JOIN parking_slot ps ON ps.lot_id = pl.id
	LEFT JOIN vehicle v ON v.slot_id = ps.id
WHERE pl.id = :id AND ps.status = 'BOOKED' AND v.exit IS NULL
GROUP BY pl.id, v.vehicle_no, v.entry, v.slot_id
`

func (s *Storage) GetParkinglotDetailsByID(ctx context.Context, req *storage.ParkingLotDetailsRequest) ([]*storage.ParkingLotDetailsResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.GetVehicle")
	log.Trace("request received")

	getPar := []*storage.ParkingLotDetailsResponse{}
	stmt, err := s.db.PrepareNamed(getParkingLot)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Select(&getPar, req); err != nil {
		logging.WithError(err, log).Error("failed to getVehicle")
		return nil, err
	}

	return getPar, nil
}

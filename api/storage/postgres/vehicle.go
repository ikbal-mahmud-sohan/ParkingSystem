package postgres

import (
	"context"

	"github.com/lib/pq"

	"parking/api/storage"
	"parking/utility/logging"
)

const insertParkVehicle = `
	INSERT INTO vehicle (
		vehicle_no, 
		entry,
		slot_id
	) VALUES (
		:vehicle_no, 
		:entry,
		:slot_id
	) RETURNING *
`

func (s *Storage) ParkVehicle(ctx context.Context, req *storage.ParkVehicleRequest) (*storage.VehicleResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.ParkVehicle")
	if err := s.ParkVehicleValidation(ctx, req); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, storage.InvalidArgument
	}

	pstmt, err := s.db.PrepareNamedContext(ctx, insertParkVehicle)
	if err != nil {
		logging.WithError(err, log).Error("failed to prepare query")
		return nil, err
	}
	defer pstmt.Close()

	res := &storage.VehicleResponse{}
	if err := pstmt.Get(res, req); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			logging.WithError(err, log).Error("failed to park vehicle")
			return nil, storage.Conflict
		}
		logging.WithError(err, log).Error("failed to park vehicle")
		return nil, err
	}

	return res, nil
}

const updateUnParkingSlot = `
	UPDATE vehicle SET
		exit 	 	 = :exit,
		parked_hours = :parked_hours,
		fee 		 = :fee,
		updated_at 	 = :updated_at
	WHERE id = :id
	RETURNING *
`

func (s *Storage) UnParkVehicle(ctx context.Context, req *storage.UnParkVehicleRequest) (*storage.VehicleResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.UnParkVehicle")
	if err := s.UnParkVehicleValidation(ctx, req); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, storage.InvalidArgument
	}

	pstmt, err := s.db.PrepareNamedContext(ctx, updateUnParkingSlot)
	if err != nil {
		logging.WithError(err, log).Error("failed to prepare query")
		return nil, err
	}
	defer pstmt.Close()

	res := &storage.VehicleResponse{}
	if err := pstmt.Get(res, req); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			logging.WithError(err, log).Error("failed to unpark vehicle")
			return nil, storage.Conflict
		}
		logging.WithError(err, log).Error("failed to unpark vehicle")
		return nil, err
	}

	return res, nil
}

const getVehicle = `
		SELECT 
		v.*,
		ps.lot_id,
		ps.slot_number 
	FROM vehicle v LEFT JOIN parking_slot ps ON ps.id = v.slot_id
	WHERE v.id = :id AND v.vehicle_no = :vehicle_no
`

func (s *Storage) GetVehicle(ctx context.Context, req *storage.VehicleRequest) (*storage.VehicleResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.GetVehicle")
	log.Trace("request received")

	if err := s.GetVehicleValidation(ctx, req); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, storage.InvalidArgument
	}

	getVic := storage.VehicleResponse{}
	stmt, err := s.db.PrepareNamed(getVehicle)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&getVic, req); err != nil {
		logging.WithError(err, log).Error("failed to getVehicle")
		return nil, err
	}
	return &getVic, nil
}

package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"parking/api/storage"
	"parking/utility/logging"
)

const insertParkingSlot = `
	INSERT INTO parking_slot (
		lot_id, 
		slot_number, 
		status
	) VALUES (
		:lot_id, 
		:slot_number, 
		:status
	) RETURNING *
`

const updateParkingSlot = `
	UPDATE parking_slot SET 
		status = :status
	WHERE lot_id = :lot_id AND slot_number = :slot_number
	RETURNING *
`

func (s *Storage) UpdateParkingSlot(ctx context.Context, req *storage.UpdateParkingSlotRequest) (*storage.ParkingSlotResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.UpdateParkingSlot")
	if err := s.UpdateParkingSlotValidation(ctx, req); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, storage.InvalidArgument
	}
	updateQuery := insertParkingSlot
	if s.parkingSlotExist(ctx, req.LotID, req.SlotNumber) {
		updateQuery = updateParkingSlot
	}
	pstmt, err := s.db.PrepareNamedContext(ctx, updateQuery)
	if err != nil {
		logging.WithError(err, log).Error("failed to prepare query")
		return nil, err
	}
	defer pstmt.Close()

	res := &storage.ParkingSlotResponse{}
	if err := pstmt.Get(res, req); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			logging.WithError(err, log).Error("failed to update parking slot to duplicate row")
			return nil, storage.Conflict
		}
		logging.WithError(err, log).Error("failed to update parking slot")
		return nil, err
	}

	return res, nil
}
const getParkingSlots = `
	SELECT 
		ps.id,
		ps.lot_id,
		ps.slot_number,
		ps.status,
		ps.created_at,
		ps.updated_at,
		pl.total_spaces 
	FROM parking_slot ps
		LEFT JOIN parking_lot pl ON lot_id = pl.id
	WHERE lot_id = $1
	ORDER BY ps.slot_number DESC`

func (s *Storage) GetParkingSlotsByLotID(ctx context.Context, lotID string) ([]*storage.ParkingSlotResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.GetParkingSlotsByLotID")
	

	var parkingSlots []*storage.ParkingSlotResponse
	if err := s.db.SelectContext(ctx, &parkingSlots, getParkingSlots, lotID); err != nil {
		logging.WithError(err, log).Error("failed to get parking slots in storage")
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}

	return parkingSlots, nil
}

func (s *Storage) parkingSlotExist(ctx context.Context, lotID string, slotNumber int32) bool {
	var (
		parkingSlotQuery = `SELECT id FROM parking_slot WHERE lot_id = $1 AND slot_number = $2`
		parkingLot       storage.ParkingLotResponse
	)
	if err := s.db.GetContext(ctx, &parkingLot, parkingSlotQuery, lotID, slotNumber); err != nil {
		return false
	}
	return parkingLot.ID != ""
}

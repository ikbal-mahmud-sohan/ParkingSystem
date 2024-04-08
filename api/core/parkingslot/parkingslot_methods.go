package parkingslot

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	rpc "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"parking/api/storage"
	pspb "parking/gunk/api/parkingslot"
	"parking/utility/logging"
)

func (s *Svc) UpdateParkingSlot(ctx context.Context, req *pspb.UpdateParkingSlotRequest) (*pspb.UpdateParkingSlotResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "core.UpdateParkingSlot")
	log.Trace("request received")

	if req.GetStatus().String() == pspb.Status_UNKNOWN.String() {
		return nil, rpc.Error(codes.InvalidArgument, "invalid status")
	}

	lot, err := s.store.GetParkinglotByID(ctx, req.GetLotID())
	if err != nil {
		logging.WithError(err, log).Error("failed to get parking lot response")
		return nil, rpc.Error(codes.Internal, "failed to get parking lot response")
	}

	if req.GetSlotNumber() > lot.TotalSpace {
		logging.WithError(errors.New("failed to handle parking lot slot"), log)
		return nil, rpc.Error(codes.InvalidArgument, "space is bigger than slot")
	}

	res, err := s.store.UpdateParkingSlot(ctx, &storage.UpdateParkingSlotRequest{
		LotID:      req.GetLotID(),
		SlotNumber: req.GetSlotNumber(),
		Status:     storage.SlotStatus(req.GetStatus().String()),
	})
	if err != nil {
		logging.WithError(err, log).Error("failed to update parking slot")
		if err == storage.NotFound {
			return nil, rpc.Error(codes.NotFound, "no data found")
		}
		return nil, rpc.Error(codes.Internal, "failed to update parking slot")
	}

	return &pspb.UpdateParkingSlotResponse{
		Data: &pspb.ParkingSlotData{
			ParkingSlot: &pspb.ParkingSlot{
				ID:         res.ID,
				LotID:      res.LotID,
				SlotNumber: res.SlotNumber,
				Status:     pspb.Status(pspb.Status_value[res.Status.String()]),
				CreatedAt:  timestamppb.New(res.CreatedAt),
				UpdatedAt:  timestamppb.New(res.UpdatedAt),
			},
		},
	}, nil
}

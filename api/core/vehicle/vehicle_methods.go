package vehicle

import (
	"context"
	"math"
	"time"

	"google.golang.org/grpc/codes"
	rpc "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"parking/api/storage"
	pvpb "parking/gunk/api/vehicle"
	"parking/utility/logging"
)

func (s *Svc) ParkVehicle(ctx context.Context, req *pvpb.ParkVehicleRequest) (*pvpb.ParkVehicleResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "core.ParkVehicle")
	log.Trace("request received")

	slots, err := s.store.GetParkingSlotsByLotID(ctx, req.GetLotID())
	if err != nil {
		logging.WithError(err, log).Error("failed to get parking slot response")
		return nil, rpc.Error(codes.Internal, "failed to get slot")
	}

	slotNumber := int32(1)
	if len(slots) > 0 {
		totalSpace := int(slots[0].TotalSpace)
		numbers := make([]int, totalSpace)
		for i := 1; i <= totalSpace; i++ {
			numbers[i-1] = i
		}

		toDrop := []int{}
		booked := storage.SlotStatusBOOKED.String()
		undermaintenance := storage.SlotStatusUNDERMAINTENANCE.String()
		for _, v := range slots {
			if v.Status.String() == booked || v.Status.String() == undermaintenance {
				toDrop = append(toDrop, int(v.SlotNumber))
			}
		}

		for _, drop := range toDrop {
			for i, num := range numbers {
				if num == drop {
					numbers = append(numbers[:i], numbers[i+1:]...)
					break
				}
			}
		}

		minimum := math.MaxInt
		for _, num := range numbers {
			if num < minimum {
				minimum = num
			}
		}
		slotNumber = int32(minimum)
	}

	slot, err := s.store.UpdateParkingSlot(ctx, &storage.UpdateParkingSlotRequest{
		LotID:      req.GetLotID(),
		SlotNumber: slotNumber,
		Status:     storage.SlotStatusBOOKED,
	})
	if err != nil {
		logging.WithError(err, log).Error("failed to update parking slot")
		return nil, rpc.Error(codes.Internal, "failed to park")
	}

	res, err := s.store.ParkVehicle(ctx, &storage.ParkVehicleRequest{
		VehicleNo: req.GetVehicleNo(),
		Entry:     time.Now().UTC(),
		SlotID:    slot.ID,
	})
	if err != nil {
		logging.WithError(err, log).Error("failed to update parking slot")
		if err == storage.NotFound {
			return nil, rpc.Error(codes.NotFound, "no data found")
		}
		return nil, rpc.Error(codes.Internal, "failed to park")
	}
	exitTime := &timestamppb.Timestamp{}
	if res.Exit.Valid {
		exitTime = timestamppb.New(res.Exit.Time)
	}

	return &pvpb.ParkVehicleResponse{
		Data: &pvpb.VehicleData{
			Vehicle: &pvpb.Vehicle{
				ID:          res.ID,
				VehicleNo:   res.VehicleNo,
				Entry:       timestamppb.New(res.Entry),
				Exit:        exitTime,
				ParkedHours: res.ParkedHours,
				Fee:         res.Fee,
				SlotID:      res.SlotID,
				CreatedAt:   timestamppb.New(res.CreatedAt),
				UpdatedAt:   timestamppb.New(res.UpdatedAt),
			},
		},
	}, nil
}

func (s *Svc) UnParkVehicle(ctx context.Context, req *pvpb.UnParkVehicleRequest) (*pvpb.UnParkVehicleResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "core.UnParkVehicle")
	log.Trace("request received")
	var vData *storage.VehicleResponse

	vData, err := s.store.GetVehicle(ctx, &storage.VehicleRequest{
		ID:        req.GetID(),
		VehicleNo: req.GetVehicleNo(),
	})
	if err != nil {
		logging.WithError(err, log).Error("failed to get vehicle")
		return nil, rpc.Error(codes.Internal, "failed to get vehicle")
	}

	parkedHours := int32(math.Ceil(time.Now().UTC().Sub(vData.Entry).Hours()))
	res, err := s.store.UnParkVehicle(ctx, &storage.UnParkVehicleRequest{
		ID:          req.GetID(),
		Exit:        time.Now().UTC(),
		ParkedHours: parkedHours,
		Fee:         parkedHours * FEE_PER_HOUR,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		logging.WithError(err, log).Error("failed to handle vehicle response")
		return nil, rpc.Error(codes.Internal, "failed to handle vehicle response")
	}

	if _, err := s.store.UpdateParkingSlot(ctx, &storage.UpdateParkingSlotRequest{
		LotID:      vData.LotID,
		SlotNumber: vData.SlotNumber,
		Status:     storage.SlotStatusAVAILABLE,
	}); err != nil {
		logging.WithError(err, log).Error("failed to handle parking slot response")
		return nil, rpc.Error(codes.Internal, "failed to handle parking slot response")
	}

	exitTime := &timestamppb.Timestamp{}
	if res.Exit.Valid {
		exitTime = timestamppb.New(res.Exit.Time)
	}

	return &pvpb.UnParkVehicleResponse{
		Data: &pvpb.VehicleData{
			Vehicle: &pvpb.Vehicle{
				ID:          res.ID,
				VehicleNo:   res.VehicleNo,
				Entry:       timestamppb.New(res.Entry),
				Exit:        exitTime,
				ParkedHours: res.ParkedHours,
				Fee:         res.Fee,
				SlotID:      res.SlotID,
				CreatedAt:   timestamppb.New(res.CreatedAt),
				UpdatedAt:   timestamppb.New(res.UpdatedAt),
			},
		},
	}, nil
}

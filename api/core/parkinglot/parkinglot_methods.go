package parkinglot

import (
	"context"
	"math"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"parking/api/storage"
	plpb "parking/gunk/api/parkinglot"
	"parking/utility/logging"
)

func (s *Svc) InsertParkingLot(ctx context.Context, req *plpb.InsertParkingLotRequest) (*plpb.InsertParkingLotResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "core.InsertParkingLot")
	log.Trace("request received")

	res, err := s.store.InsertParkingLot(ctx, &storage.InsertParkingLotRequest{
		Name:        req.GetName(),
		TotalSpaces: req.GetTotalSpace(),
	})
	if err != nil {
		logging.WithError(err, log).Error("failed to insert ParkingLot")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "no data found")
		}
		return nil, status.Error(codes.Internal, "failed to insert ParkingLot")
	}

	return &plpb.InsertParkingLotResponse{
		Data: &plpb.ParkingLotData{
			ParkingLot: &plpb.ParkingLot{
				ID:         res.ID,
				Name:       res.Name,
				TotalSpace: res.TotalSpace,
				CreatedAt:  timestamppb.New(res.CreatedAt),
				UpdatedAt:  timestamppb.New(res.UpdatedAt),
			},
		},
	}, nil
}

func (s *Svc) ListParkingLot(ctx context.Context, req *plpb.ListParkingLotRequest) (*plpb.ListParkingLotResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "core.ListParkingLot")
	log.Trace("request received")

	res, err := s.store.ListParkingLot(ctx, &storage.ListParkingLotRequest{
		OrderBy:      storage.OrderBy(req.OrderBy.String()),
		SearchTerm:   req.GetSearchTerm(),
		SortByColumn: storage.SortByColumn(req.GetSortByColumn()),
		Limit:        int(req.GetLimit()),
		Offset:       int(req.GetOffset()),
	})
	if err != nil || res == nil {
		logging.WithError(err, log).Error("failed to list parking lo")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "no data found")
		}
		return nil, status.Error(codes.Internal, "failed to list parking lo")
	}
	if len(res) == 0 {
		return &plpb.ListParkingLotResponse{
			Data: &plpb.ListParkingLotData{
				ParkingLot: nil,
			},
			Total: 0,
		}, nil
	}
	plRes, err := parkingLotListResponse(res)
	if err != nil {
		logging.WithError(err, log).Error("failed to list parking lot response")
		return nil, status.Error(codes.Internal, "failed to list parking lot response")
	}

	var totalParkingLot int32 = 0
	if len(res) > 0 {
		totalParkingLot = res[0].Total
	}

	return &plpb.ListParkingLotResponse{
		Data: &plpb.ListParkingLotData{
			ParkingLot: plRes,
		},
		Total: totalParkingLot,
	}, nil
}

func (s *Svc) GetParkingLot(ctx context.Context, req *plpb.GetParkingLotRequest) (*plpb.GetParkingLotResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "core.GetParkingLot")
	log.Trace("request received")

	res, err := s.store.GetParkinglotDetailsByID(ctx, &storage.ParkingLotDetailsRequest{ID: req.GetID()})
	if err != nil || res == nil {
		logging.WithError(err, log).Error("failed to get parking lot")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "no data found")
		}
		return nil, status.Error(codes.Internal, "failed to get parking lot")
	}
	parkedVehicles := make(map[string][]*plpb.ParkedVehicles, 0)
	for _, v := range res {
		parkedHours := int32(math.Ceil(time.Now().Sub(v.Entry.Time).Hours()))
		parkedVehicles[v.ID] = append(parkedVehicles[v.ID], &plpb.ParkedVehicles{
			VehicleNo:   v.VehicleNo.String,
			EntryTime:   timestamppb.New(v.Entry.Time),
			ParkedHours: parkedHours,
			Fee:         parkedHours * FEE_PER_HOUR,
			SlotID:      v.SlotID.String,
		})
	}
	var parkingLot *plpb.GetParkingLot
	if len(res) > 0 {
		parkingLot = &plpb.GetParkingLot{
			ID:             res[0].ID,
			Name:           res[0].Name,
			TotalSpace:     res[0].TotalSpace,
			TotalParked:    res[0].TotalParked,
			ParkedVehicles: parkedVehicles[res[0].ID],
			CreatedAt:      timestamppb.New(res[0].CreatedAt),
			UpdatedAt:      timestamppb.New(res[0].UpdatedAt),
		}
	}
	return &plpb.GetParkingLotResponse{
		Data: &plpb.GetParkingLotData{
			ParkingLot: parkingLot,
		},
	}, nil
}

func parkingLotListResponse(res []*storage.ParkingLotResponse) ([]*plpb.ParkingLot, error) {
	if res == nil {
		return nil, status.Error(codes.NotFound, "no data found")
	}
	parkinglot := make([]*plpb.ParkingLot, 0, len(res))

	for _, v := range res {
		parkinglot = append(parkinglot, &plpb.ParkingLot{
			ID:         v.ID,
			Name:       v.Name,
			TotalSpace: v.TotalSpace,
			CreatedAt:  timestamppb.New(v.CreatedAt),
			UpdatedAt:  timestamppb.New(v.UpdatedAt),
		})
	}
	if len(parkinglot) == 0 {
		return nil, nil
	}

	return parkinglot, nil
}

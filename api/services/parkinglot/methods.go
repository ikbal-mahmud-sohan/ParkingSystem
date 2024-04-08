package parkinglot

import (
	"context"

	plpb "parking/gunk/api/parkinglot"
	uErr "parking/utility/error/error"
	"parking/utility/logging"
)

func (s *Svc) InsertParkingLot(ctx context.Context, req *plpb.InsertParkingLotRequest) (*plpb.InsertParkingLotResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.InsertParkingLot")
	log.Trace("request received")

	if err := s.InsertParkingLotValidation(ctx, req); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, uErr.HandleServiceErr(err)
	}

	res, err := s.core.InsertParkingLot(ctx, req)
	if err != nil {
		errM := "failed to Create ParkingLot"
		log.WithError(err).Error(errM)
		return nil, uErr.HandleServiceErr(err)
	}

	return res, nil
}

func (s *Svc) ListParkingLot(ctx context.Context, req *plpb.ListParkingLotRequest) (*plpb.ListParkingLotResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "Service.ListParkingLot")
	log.Trace("request received")

	res, err := s.core.ListParkingLot(ctx, req)
	if err != nil {
		errM := "failed to get ParkingLot list"
		log.WithError(err).Error(errM)
		return nil, uErr.HandleServiceErr(err)
	}

	return res, nil
}

func (s *Svc) GetParkingLot(ctx context.Context, req *plpb.GetParkingLotRequest) (*plpb.GetParkingLotResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "Service.GetParkingLot")
	log.Trace("request received")
	if err := s.GetParkingLotValidation(ctx, req); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, uErr.HandleServiceErr(err)
	}
	res, err := s.core.GetParkingLot(ctx, req)
	if err != nil {
		errM := "failed to get parking lot details"
		log.WithError(err).Error(errM)
		return nil, uErr.HandleServiceErr(err)
	}

	return res, nil
}

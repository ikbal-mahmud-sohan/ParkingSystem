package parkingstatement

import (
	"context"

	"google.golang.org/grpc/codes"
	rpc "google.golang.org/grpc/status"

	pspb "parking/gunk/api/parkingstatement"
	"parking/utility/logging"
)

func (s *Svc) GetParkingStatement(ctx context.Context, req *pspb.GetParkingStatementRequest) (*pspb.GetParkingStatementResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "core.GetParkingStatement")
	log.Trace("request received")

	res, err := s.store.GetParkingStatement(ctx, req.StartDate.AsTime(), req.GetEndDate().AsTime())
	if err != nil {
		logging.WithError(err, log).Error("failed to get parking statement response")
		return nil, rpc.Error(codes.Internal, "failed to get parking statement response")
	}

	return &pspb.GetParkingStatementResponse{
		Data: &pspb.StatementData{
			Statement: &pspb.Statement{
				TotalEarning: res.TotalEarning.Int32,
				TotalParked:  res.TotalParked.Int32,
				TotalHours:   res.TotalHours.Int32,
			},
		},
	}, nil
}

package parkinglot

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	plpb "parking/gunk/api/parkinglot"
)

type ParkingLotCore interface {
	InsertParkingLot(ctx context.Context, req *plpb.InsertParkingLotRequest) (*plpb.InsertParkingLotResponse, error)
	ListParkingLot(ctx context.Context, req *plpb.ListParkingLotRequest) (*plpb.ListParkingLotResponse, error)
	GetParkingLot(ctx context.Context, req *plpb.GetParkingLotRequest) (*plpb.GetParkingLotResponse, error)
}

type Svc struct {
	plpb.UnimplementedParkingLotServiceServer
	core ParkingLotCore
}

func New(core ParkingLotCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterSvc RegisterService with grpc server.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	plpb.RegisterParkingLotServiceServer(srv, s)
	return nil
}

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return plpb.RegisterParkingLotServiceHandlerFromEndpoint(ctx, mux, address, options)
}

package vehicle

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	pvpb "parking/gunk/api/vehicle"
)

type VehicleCore interface {
	ParkVehicle(ctx context.Context, req *pvpb.ParkVehicleRequest) (*pvpb.ParkVehicleResponse, error)
	UnParkVehicle(ctx context.Context, req *pvpb.UnParkVehicleRequest) (*pvpb.UnParkVehicleResponse, error)
}

type Svc struct {
	pvpb.UnimplementedVehicleServiceServer
	core VehicleCore
}

func New(core VehicleCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterSvc RegisterService with grpc server.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	pvpb.RegisterVehicleServiceServer(srv, s)
	return nil
}

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return pvpb.RegisterVehicleServiceHandlerFromEndpoint(ctx, mux, address, options)
}

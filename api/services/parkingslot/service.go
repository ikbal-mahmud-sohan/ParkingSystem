package parkingslot

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	pspb "parking/gunk/api/parkingslot"
)

type ParkingSlotCore interface {
	UpdateParkingSlot(ctx context.Context, req *pspb.UpdateParkingSlotRequest) (*pspb.UpdateParkingSlotResponse, error)
}

type Svc struct {
	pspb.UnimplementedParkingSlotServiceServer
	core ParkingSlotCore
}

func New(core ParkingSlotCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterSvc RegisterService with grpc server.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	pspb.RegisterParkingSlotServiceServer(srv, s)
	return nil
}

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return pspb.RegisterParkingSlotServiceHandlerFromEndpoint(ctx, mux, address, options)
}

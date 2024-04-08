package parkingstatement

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	pspb "parking/gunk/api/parkingstatement"
)

type StatementCore interface {
	GetParkingStatement(ctx context.Context, req *pspb.GetParkingStatementRequest) (*pspb.GetParkingStatementResponse, error)
}

type Svc struct {
	pspb.UnimplementedStatementServiceServer
	core StatementCore
}

func New(core StatementCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterSvc RegisterService with grpc server.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	pspb.RegisterStatementServiceServer(srv, s)
	return nil
}

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return pspb.RegisterStatementServiceHandlerFromEndpoint(ctx, mux, address, options)
}

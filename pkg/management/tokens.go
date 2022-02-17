package management

import (
	"context"

	core "github.com/rancher/opni-monitoring/pkg/core"
	"github.com/rancher/opni-monitoring/pkg/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *Server) CreateBootstrapToken(
	ctx context.Context,
	req *CreateBootstrapTokenRequest,
) (*core.BootstrapToken, error) {
	if err := validation.Validate(req); err != nil {
		return nil, err
	}
	token, err := m.tokenStore.CreateToken(ctx, req.Ttl.AsDuration(), req.GetLabels())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return token, nil
}

func (m *Server) RevokeBootstrapToken(
	ctx context.Context,
	ref *core.Reference,
) (*emptypb.Empty, error) {
	if err := validation.Validate(ref); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, grpcError(m.tokenStore.DeleteToken(ctx, ref))
}

func (m *Server) ListBootstrapTokens(
	ctx context.Context,
	_ *emptypb.Empty,
) (*core.BootstrapTokenList, error) {
	tokens, err := m.tokenStore.ListTokens(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	tokenList := &core.BootstrapTokenList{}
	tokenList.Items = append(tokenList.Items, tokens...)
	return tokenList, nil
}

func (m *Server) GetBootstrapToken(
	ctx context.Context,
	ref *core.Reference,
) (*core.BootstrapToken, error) {
	if err := validation.Validate(ref); err != nil {
		return nil, err
	}
	token, err := m.tokenStore.GetToken(ctx, ref)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return token, nil
}

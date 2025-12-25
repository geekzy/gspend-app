package client

import (
	"context"
	"time"

	authv1 "github.com/imam/gspend-app/apps/financial-service/pkg/proto/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthGRPCClient struct {
	client authv1.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthGRPCClient(authServiceAddr string) (*AuthGRPCClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, authServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	client := authv1.NewAuthServiceClient(conn)

	return &AuthGRPCClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *AuthGRPCClient) ValidateToken(ctx context.Context, token string) (*authv1.ValidateTokenResponse, error) {
	return c.client.ValidateToken(ctx, &authv1.ValidateTokenRequest{
		Token: token,
	})
}

func (c *AuthGRPCClient) GetUserProfile(ctx context.Context, userID string) (*authv1.GetUserProfileResponse, error) {
	return c.client.GetUserProfile(ctx, &authv1.GetUserProfileRequest{
		UserId: userID,
	})
}

func (c *AuthGRPCClient) CheckUserExists(ctx context.Context, userID string) (*authv1.CheckUserExistsResponse, error) {
	return c.client.CheckUserExists(ctx, &authv1.CheckUserExistsRequest{
		UserId: userID,
	})
}

func (c *AuthGRPCClient) Close() error {
	return c.conn.Close()
}

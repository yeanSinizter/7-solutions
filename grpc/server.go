package grpc

import (
	userpb "7-solutions/proto"
	"7-solutions/usecase"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCServer struct {
	userpb.UnimplementedUserServiceServer
	Usecase usecase.UserUsecase
}

func NewUserGRPCServer(uc usecase.UserUsecase) *UserGRPCServer {
	return &UserGRPCServer{Usecase: uc}
}

func (s *UserGRPCServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	err := s.Usecase.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{
		User: &userpb.User{
			Id:        uuid.New().String(),
			Name:      req.Name,
			Email:     req.Email,
			CreatedAt: time.Now().Format(time.RFC3339),
		},
	}, nil
}

// func (s *UserGRPCServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
// 	user, err := s.Usecase.GetUser(ctx, req.Id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &userpb.GetUserResponse{
// 		User: &userpb.User{
// 			Id:        user.ID.Hex(),
// 			Name:      user.Name,
// 			Email:     user.Email,
// 			CreatedAt: user.CreatedAt.Format(time.RFC3339),
// 		},
// 	}, nil
// }

func (s *UserGRPCServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	// ดึง user_id จาก context (อาจจะต้องส่งมาจาก interceptor)
	userIDFromToken, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	// เช็คสิทธิ์ ว่าขอข้อมูล user ตัวเองเท่านั้น (ถ้าต้องการ)
	if userIDFromToken != req.Id {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	user, err := s.Usecase.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserResponse{
		User: &userpb.User{
			Id:        user.ID.Hex(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return "", errors.New("userID not found in context")
	}
	return userID, nil
}

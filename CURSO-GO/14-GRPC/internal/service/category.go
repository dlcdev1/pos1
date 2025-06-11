package service

import (
	"context"
	"github.com/dlcdev1/pos/14-gRPC/internal/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
import "github.com/dlcdev1/pos/14-gRPC/internal/pb"

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	categoy, err := s.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	categoryResponse := &pb.Category{
		Id:          categoy.ID,
		Name:        categoy.Name,
		Description: categoy.Description,
	}

	return categoryResponse, nil
}

package user

import (
	"acution_dlcdev/configuration/logger"
	"acution_dlcdev/internal/entity/user_entity"
	"acution_dlcdev/internal/internal_error"
	"context"
)

func (ur *UserRepository) CreateUser(ctx context.Context, userEntity user_entity.User) (user_entity.User, *internal_error.InternalError) {
	userEntityMongo := UserEntityMongo{
		Id:   userEntity.Id,
		Name: userEntity.Name,
	}

	_, err := ur.Collection.InsertOne(ctx, userEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert user", err)
		return userEntity, internal_error.NewInternalServerError("Error trying to insert user")
	}

	return userEntity, nil
}

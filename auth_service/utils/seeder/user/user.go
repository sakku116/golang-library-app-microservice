package seeder_util

import (
	"auth_service/config"
	"auth_service/domain/model"
	author_pb "auth_service/interface/grpc/genproto/author"
	"auth_service/repository"
	bcrypt_util "auth_service/utils/bcrypt"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/op/go-logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = logging.MustGetLogger("main")

func SeedUser(userRepo repository.IUserRepo, authorGrpcServiceClient author_pb.AuthorServiceClient) error {
	users := []model.User{}

	if config.Envs.INITIAL_ADMIN_USERNAME != "" && config.Envs.INITIAL_ADMIN_PASSWORD != "" {
		hashedPassword, _ := bcrypt_util.Hash(config.Envs.INITIAL_ADMIN_PASSWORD)
		users = append(users, model.User{
			UUID:     uuid.New(),
			Username: config.Envs.INITIAL_ADMIN_USERNAME,
			Password: hashedPassword,
			Email:    fmt.Sprint(config.Envs.INITIAL_ADMIN_USERNAME, "@gmail.com"),
			Role:     "admin",
		})
	} else {
		logger.Warningf("initial admin username and password not set")
	}

	if config.Envs.INITIAL_USER_USERNAME != "" && config.Envs.INITIAL_USER_PASSWORD != "" {
		hashedPassword, _ := bcrypt_util.Hash(config.Envs.INITIAL_USER_PASSWORD)
		users = append(users, model.User{
			UUID:     uuid.New(),
			Username: config.Envs.INITIAL_USER_USERNAME,
			Password: hashedPassword,
			Email:    fmt.Sprint(config.Envs.INITIAL_USER_USERNAME, "@gmail.com"),
			Role:     "user",
		})
	} else {
		logger.Warningf("initial user username and password not set")
	}

	for _, user := range users {
		logger.Infof("seeding user: %s", user.Username)

		existing, _ := userRepo.GetByUsername(user.Username)
		if existing != nil {
			logger.Warningf("user already exists: %s", user.Username)
		}

		err := userRepo.Create(&user)
		if err != nil {
			logger.Warningf("failed to seed user: %s", user.Username)
		}

		logger.Infof("user seeded: %s", user.Username)

		// create author through author service
		userUUID := user.UUID.String()
		if existing != nil {
			logger.Debugf("user exist")
			userUUID = existing.UUID.String()
		}
		logger.Debugf("userUUID: %s", userUUID)
		createAuthorResp, err := authorGrpcServiceClient.CreateAuthor(
			context.Background(),
			&author_pb.CreateAuthorReq{
				UserUuid:  userUUID,
				FirstName: user.Username,
			},
		)

		grpcCode := status.Code(err)

		if grpcCode != codes.OK || err != nil {
			logger.Warningf("failed to seed author: %s; error: %s", user.Username, err.Error())
			continue
		}

		if createAuthorResp == nil {
			logger.Warningf("failed to seed author: %s; response is nil", user.Username)
			continue
		}

		logger.Infof("author seeded: %s %s", createAuthorResp.FirstName, createAuthorResp.LastName)

	}

	return nil
}

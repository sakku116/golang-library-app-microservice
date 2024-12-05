package handler

import (
	"book_service/domain/dto"
	author_pb "book_service/interface/grpc/genproto/author"
	ucase "book_service/usecase"
	error_utils "book_service/utils/error"
	"context"

	"github.com/op/go-logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthorServiceHandler struct {
	author_pb.UnimplementedAuthorServiceServer
	authorUcase ucase.IAuthorUcase
}

var logger = logging.MustGetLogger("main")

func NewAuthorServiceHandler(authorUcase ucase.IAuthorUcase) *AuthorServiceHandler {
	handler := &AuthorServiceHandler{authorUcase: authorUcase}
	// logger.Debugf("ucase: %v", authorUcase)
	return handler
}

func (r *AuthorServiceHandler) CreateAuthor(
	ctx context.Context,
	in *author_pb.CreateAuthorReq,
) (*author_pb.CreateAuthorResp, error) {
	logger.Debugf("incoming request: %v", in)
	// payload validation
	payloadDto := dto.CreateNewAuthorReq{
		LastName: in.LastName,
	}
	if in.UserUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "user uuid is required")
	} else {
		payloadDto.UserUUID = &in.UserUuid
	}

	if in.FirstName == "" {
		return nil, status.Error(codes.InvalidArgument, "first name is required")
	} else {
		payloadDto.FirstName = in.FirstName
	}

	if in.BirthDate == "" {
		payloadDto.BirthDate = nil
	} else {
		payloadDto.BirthDate = &in.BirthDate
	}

	if in.Bio == "" {
		payloadDto.Bio = nil
	} else {
		payloadDto.Bio = &in.Bio
	}

	logger.Debugf("calling create new author")
	raw, err := r.authorUcase.CreateNewAuthor(ctx, payloadDto)
	logger.Debugf("create new author done")
	if err != nil {
		logger.Errorf("error creating author: %v", err)
		customErr, ok := err.(*error_utils.CustomErr)
		if ok {
			return nil, status.Errorf(customErr.GrpcCode, customErr.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &author_pb.CreateAuthorResp{
		Uuid:      raw.UUID.String(),
		UserUuid:  raw.UserUUID.String(),
		FirstName: raw.FirstName,
		LastName:  raw.LastName,
	}

	if raw.BirthDate != nil {
		resp.BirthDate = *raw.BirthDate
	} else {
		resp.BirthDate = ""
	}

	if raw.Bio != nil {
		resp.Bio = *raw.Bio
	} else {
		resp.Bio = ""
	}

	return resp, nil
}
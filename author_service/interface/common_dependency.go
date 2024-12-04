package interface_pkg

import (
	"author_service/interface/grpc/genproto/auth"
	"author_service/repository"
	ucase "author_service/usecase"
)

type CommonDependency struct {
	AuthorUcase ucase.IAuthorUcase
	AuthorRepo  repository.IAuthorRepo
	AuthGrpc    auth.AuthServiceClient
}

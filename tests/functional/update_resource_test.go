package functional_test

import (
	"github.com/stretchr/testify/suite"
)

type UpdateResource struct {
	suite.Suite
}

func NewUpdateResource() *UpdateResource {
	return &UpdateResource{}
}

func (s *UpdateResource) SetupSuite() {}

func (s *UpdateResource) TearDownSuite() {}

func (s *UpdateResource) Test_Success() {}

func (s *UpdateResource) Test_InvalidArgument() {}

func (s *UpdateResource) Test_NotFound() {}

func (s *UpdateResource) Test_AlreadyExists() {}

func (s *UpdateResource) Test_PermissionDenied() {}

func (s *UpdateResource) Test_Unauthenticated() {}

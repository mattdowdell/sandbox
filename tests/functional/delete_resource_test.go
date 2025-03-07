package functional_test

import (
	"github.com/stretchr/testify/suite"
)

type DeleteResource struct {
	suite.Suite
}

func NewDeleteResource() *DeleteResource {
	return &DeleteResource{}
}

func (s *DeleteResource) SetupSuite() {}

func (s *DeleteResource) TearDownSuite() {}

func (s *DeleteResource) Test_Success() {}

func (s *DeleteResource) Test_InvalidArgument() {}

func (s *DeleteResource) Test_NotFound() {}

func (s *DeleteResource) Test_PermissionDenied() {}

func (s *DeleteResource) Test_Unauthenticated() {}

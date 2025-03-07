package functional_test

import (
	"github.com/stretchr/testify/suite"
)

type GetResource struct {
	suite.Suite
}

func NewGetResource() *GetResource {
	return &GetResource{}
}

func (s *GetResource) SetupSuite() {}

func (s *GetResource) TearDownSuite() {}

func (s *GetResource) Test_Success() {}

func (s *GetResource) Test_InvalidArgument() {}

func (s *GetResource) Test_NotFound() {}

func (s *GetResource) Test_PermissionDenied() {}

func (s *GetResource) Test_Unauthenticated() {}

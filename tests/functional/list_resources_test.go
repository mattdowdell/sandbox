package functional_test

import (
	"github.com/stretchr/testify/suite"
)

type ListResources struct {
	suite.Suite
}

func NewListResources() *ListResources {
	return &ListResources{}
}

func (s *ListResources) SetupSuite() {}

func (s *ListResources) TearDownSuite() {}

func (s *ListResources) Test_Success() {}

func (s *ListResources) Test_InvalidArgument() {}

func (s *ListResources) Test_PermissionDenied() {}

func (s *ListResources) Test_Unauthenticated() {}

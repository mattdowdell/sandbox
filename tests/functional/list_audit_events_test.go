package functional_test

import (
	"github.com/stretchr/testify/suite"
)

type ListAuditEvents struct {
	suite.Suite
}

func NewListAuditEvents() *ListAuditEvents {
	return &ListAuditEvents{}
}

func (s *ListAuditEvents) SetupSuite() {}

func (s *ListAuditEvents) TearDownSuite() {}

func (s *ListAuditEvents) Test_Success() {}

func (s *ListAuditEvents) Test_InvalidArgument() {}

func (s *ListAuditEvents) Test_PermissionDenied() {}

func (s *ListAuditEvents) Test_Unauthenticated() {}

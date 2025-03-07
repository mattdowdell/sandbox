package functional_test

import (
	"github.com/stretchr/testify/suite"
)

type WatchAuditEvents struct {
	suite.Suite
}

func NewWatchAuditEvents() *WatchAuditEvents {
	return &WatchAuditEvents{}
}

func (s *WatchAuditEvents) SetupSuite() {}

func (s *WatchAuditEvents) TearDownSuite() {}

func (s *WatchAuditEvents) Test_Success() {}

func (s *WatchAuditEvents) Test_InvalidArgument() {}

func (s *WatchAuditEvents) Test_PermissionDenied() {}

func (s *WatchAuditEvents) Test_Unauthenticated() {}

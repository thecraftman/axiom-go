//go:build integration
// +build integration

package axiom_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// VersionTestSuite tests all methods of the Axiom Version API against a live
// deployment.
type VersionTestSuite struct {
	IntegrationTestSuite
}

func TestVersionTestSuite(t *testing.T) {
	suite.Run(t, new(VersionTestSuite))
}

func (s *VersionTestSuite) Test() {
	version, err := s.client.Version.Get(s.ctx)
	s.Require().NoError(err)
	s.Require().NotEmpty(version)

	s.Contains(version, "v1.")
}

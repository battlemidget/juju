// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package ec2_test

import (
	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/cloud"
	"github.com/juju/juju/environs"
	envtesting "github.com/juju/juju/environs/testing"
)

type credentialsSuite struct {
	testing.IsolationSuite
	provider environs.EnvironProvider
}

var _ = gc.Suite(&credentialsSuite{})

func (s *credentialsSuite) SetUpTest(c *gc.C) {
	s.IsolationSuite.SetUpTest(c)

	var err error
	s.provider, err = environs.Provider("ec2")
	c.Assert(err, jc.ErrorIsNil)
}

func (s *credentialsSuite) TestCredentialSchemas(c *gc.C) {
	envtesting.AssertProviderAuthTypes(c, s.provider, "access-key")
}

func (s *credentialsSuite) TestAccessKeyCredentialsValid(c *gc.C) {
	envtesting.AssertProviderCredentialsValid(c, s.provider, "access-key", map[string]string{
		"access-key": "key",
		"secret-key": "secret",
	})
}

func (s *credentialsSuite) TestAccessKeyHiddenAttributes(c *gc.C) {
	envtesting.AssertProviderCredentialsAttributesHidden(c, s.provider, "access-key", "secret-key")
}

func (s *credentialsSuite) TestDetectCredentialsNotFound(c *gc.C) {
	// No environment variables set, so no credentials should be found.
	credentials, err := s.provider.DetectCredentials()
	c.Assert(err, jc.Satisfies, errors.IsNotFound)
	c.Assert(credentials, gc.HasLen, 0)
}

func (s *credentialsSuite) TestDetectCredentialsEnvironmentVariables(c *gc.C) {
	s.PatchEnvironment("AWS_ACCESS_KEY_ID", "key-id")
	s.PatchEnvironment("AWS_SECRET_ACCESS_KEY", "secret-access-key")

	credentials, err := s.provider.DetectCredentials()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(credentials, gc.HasLen, 1)
	c.Assert(credentials[0], jc.DeepEquals, environs.LabeledCredential{
		Credential: cloud.NewCredential(
			cloud.AccessKeyAuthType, map[string]string{
				"access-key": "key-id",
				"secret-key": "secret-access-key",
			},
		),
	})
}

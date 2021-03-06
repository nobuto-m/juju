// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state_test

import (
	gitjujutesting "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/juju/names.v2"

	"github.com/juju/juju/constraints"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/mongo/mongotest"
	"github.com/juju/juju/state"
	statetesting "github.com/juju/juju/state/testing"
	"github.com/juju/juju/testing"
)

type InitializeSuite struct {
	gitjujutesting.MgoSuite
	testing.BaseSuite
	State *state.State
}

var _ = gc.Suite(&InitializeSuite{})

func (s *InitializeSuite) SetUpSuite(c *gc.C) {
	s.BaseSuite.SetUpSuite(c)
	s.MgoSuite.SetUpSuite(c)
}

func (s *InitializeSuite) TearDownSuite(c *gc.C) {
	s.MgoSuite.TearDownSuite(c)
	s.BaseSuite.TearDownSuite(c)
}

func (s *InitializeSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)
	s.MgoSuite.SetUpTest(c)
}

func (s *InitializeSuite) openState(c *gc.C, modelTag names.ModelTag) {
	st, err := state.Open(
		modelTag,
		statetesting.NewMongoInfo(),
		mongotest.DialOpts(),
		state.Policy(nil),
	)
	c.Assert(err, jc.ErrorIsNil)
	s.State = st
}

func (s *InitializeSuite) TearDownTest(c *gc.C) {
	if s.State != nil {
		err := s.State.Close()
		c.Check(err, jc.ErrorIsNil)
	}
	s.MgoSuite.TearDownTest(c)
	s.BaseSuite.TearDownTest(c)
}

func (s *InitializeSuite) TestInitialize(c *gc.C) {
	cfg := testing.ModelConfig(c)
	uuid := cfg.UUID()
	initial := cfg.AllAttrs()
	owner := names.NewLocalUserTag("initialize-admin")
	st, err := state.Initialize(owner, statetesting.NewMongoInfo(), "dummy", "some-region", nil, cfg, mongotest.DialOpts(), nil)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(st, gc.NotNil)
	modelTag := st.ModelTag()
	c.Assert(modelTag.Id(), gc.Equals, uuid)
	err = st.Close()
	c.Assert(err, jc.ErrorIsNil)

	s.openState(c, modelTag)

	cfg, err = s.State.ModelConfig()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(cfg.AllAttrs(), gc.DeepEquals, initial)
	// Check that the model has been created.
	model, err := s.State.Model()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(model.Tag(), gc.Equals, modelTag)
	c.Assert(model.CloudRegion(), gc.Equals, "some-region")
	// Check that the owner has been created.
	c.Assert(model.Owner(), gc.Equals, owner)
	// Check that the owner can be retrieved by the tag.
	entity, err := s.State.FindEntity(model.Owner())
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(entity.Tag(), gc.Equals, owner)
	// Check that the owner has an ModelUser created for the bootstrapped model.
	modelUser, err := s.State.ModelUser(model.Owner())
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(modelUser.UserTag(), gc.Equals, owner)
	c.Assert(modelUser.ModelTag(), gc.Equals, model.Tag())

	// Check that the model can be found through the tag.
	entity, err = s.State.FindEntity(modelTag)
	c.Assert(err, jc.ErrorIsNil)
	cons, err := s.State.ModelConstraints()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(&cons, jc.Satisfies, constraints.IsEmpty)

	addrs, err := s.State.APIHostPorts()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(addrs, gc.HasLen, 0)

	info, err := s.State.ControllerInfo()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(info, jc.DeepEquals, &state.ControllerInfo{ModelTag: modelTag, CloudName: "dummy"})
}

func (s *InitializeSuite) TestInitializeWithCloudConfig(c *gc.C) {
	cfg := testing.ModelConfig(c)
	uuid := cfg.UUID()
	initial := cfg.AllAttrs()
	cloudAttrs := map[string]interface{}{
		"default-series": initial["default-series"],
	}
	owner := names.NewLocalUserTag("initialize-admin")
	st, err := state.Initialize(owner, statetesting.NewMongoInfo(), "dummy", "some-region", cloudAttrs, cfg, mongotest.DialOpts(), nil)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(st, gc.NotNil)
	modelTag := st.ModelTag()
	c.Assert(modelTag.Id(), gc.Equals, uuid)
	err = st.Close()
	c.Assert(err, jc.ErrorIsNil)

	s.openState(c, modelTag)

	cloudCfg, err := s.State.CloudConfig()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(cloudCfg, jc.DeepEquals, cloudAttrs)

	cfg, err = s.State.ModelConfig()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(cfg.AllAttrs(), gc.DeepEquals, initial)
}

func (s *InitializeSuite) TestDoubleInitializeConfig(c *gc.C) {
	cfg := testing.ModelConfig(c)
	owner := names.NewLocalUserTag("initialize-admin")

	mgoInfo := statetesting.NewMongoInfo()
	dialOpts := mongotest.DialOpts()
	st, err := state.Initialize(owner, mgoInfo, "dummy", "some-region", nil, cfg, dialOpts, state.Policy(nil))
	c.Assert(err, jc.ErrorIsNil)
	err = st.Close()
	c.Check(err, jc.ErrorIsNil)

	st, err = state.Initialize(owner, mgoInfo, "dummy", "some-region", nil, cfg, dialOpts, state.Policy(nil))
	c.Check(err, gc.ErrorMatches, "already initialized")
	if !c.Check(st, gc.IsNil) {
		err = st.Close()
		c.Check(err, jc.ErrorIsNil)
	}
}

func (s *InitializeSuite) TestModelConfigWithAdminSecret(c *gc.C) {
	// admin-secret blocks Initialize.
	good := testing.ModelConfig(c)
	badUpdateAttrs := map[string]interface{}{"admin-secret": "foo"}
	bad, err := good.Apply(badUpdateAttrs)
	owner := names.NewLocalUserTag("initialize-admin")

	_, err = state.Initialize(owner, statetesting.NewMongoInfo(), "dummy", "some-region", nil, bad, mongotest.DialOpts(), state.Policy(nil))
	c.Assert(err, gc.ErrorMatches, "admin-secret should never be written to the state")

	// admin-secret blocks UpdateModelConfig.
	st := statetesting.Initialize(c, owner, good, nil)
	st.Close()

	s.openState(c, st.ModelTag())
	err = s.State.UpdateModelConfig(badUpdateAttrs, nil, nil)
	c.Assert(err, gc.ErrorMatches, "admin-secret should never be written to the state")

	// ModelConfig remains inviolate.
	cfg, err := s.State.ModelConfig()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(cfg.AllAttrs(), jc.DeepEquals, good.AllAttrs())
}

func (s *InitializeSuite) TestModelConfigWithoutAgentVersion(c *gc.C) {
	// admin-secret blocks Initialize.
	good := testing.ModelConfig(c)
	attrs := good.AllAttrs()
	delete(attrs, "agent-version")
	bad, err := config.New(config.NoDefaults, attrs)
	c.Assert(err, jc.ErrorIsNil)
	owner := names.NewLocalUserTag("initialize-admin")

	_, err = state.Initialize(owner, statetesting.NewMongoInfo(), "dummy", "some-region", nil, bad, mongotest.DialOpts(), state.Policy(nil))
	c.Assert(err, gc.ErrorMatches, "agent-version must always be set in state")

	st := statetesting.Initialize(c, owner, good, nil)
	// yay side effects
	st.Close()

	s.openState(c, st.ModelTag())
	err = s.State.UpdateModelConfig(map[string]interface{}{}, []string{"agent-version"}, nil)
	c.Assert(err, gc.ErrorMatches, "agent-version must always be set in state")

	// ModelConfig remains inviolate.
	cfg, err := s.State.ModelConfig()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(cfg.AllAttrs(), jc.DeepEquals, good.AllAttrs())
}

func (s *InitializeSuite) TestCloudConfigWithForbiddenValues(c *gc.C) {
	badAttrNames := []string{
		config.AdminSecretKey,
		config.AgentVersionKey,
	}
	for _, attr := range config.ControllerOnlyConfigAttributes {
		badAttrNames = append(badAttrNames, attr)
	}
	good := testing.ModelConfig(c)

	for _, badAttrName := range badAttrNames {
		badCloudAttrs := map[string]interface{}{badAttrName: "foo"}
		owner := names.NewLocalUserTag("initialize-admin")

		_, err := state.Initialize(owner, statetesting.NewMongoInfo(), "dummy", "some-region", badCloudAttrs, good, mongotest.DialOpts(), state.Policy(nil))
		c.Assert(err, gc.ErrorMatches, "cloud config cannot contain .*")
	}
}

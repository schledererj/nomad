package nomad

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/nomad/command/agent/consul"
	"github.com/hashicorp/nomad/helper"
	"github.com/hashicorp/nomad/helper/testlog"
	"github.com/hashicorp/nomad/helper/uuid"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/stretchr/testify/require"
)

var _ ConsulACLsAPI = (*consulACLsAPI)(nil)
var _ ConsulACLsAPI = (*mockConsulACLsAPI)(nil)

type mockConsulACLsAPI struct {
	revokeRequests []string
}

func (m *mockConsulACLsAPI) CreateToken(
	_ context.Context,
	_ ServiceIdentityIndex) (*structs.SIToken, error) {
	panic("not used yet")
}

func (m *mockConsulACLsAPI) RevokeTokens(
	_ context.Context,
	accessors []*structs.SITokenAccessor) error {
	for _, accessor := range accessors {
		m.revokeRequests = append(m.revokeRequests, accessor.AccessorID)
	}
	return nil
}

func (m *mockConsulACLsAPI) ListTokens() ([]string, error) {
	panic("not used yet")
}

func TestConsulACLsAPI_CreateToken(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	try := func(expErr error) {
		logger := testlog.HCLogger(t)
		aclAPI := consul.NewMockACLsAPI(logger)
		aclAPI.SetError(expErr)

		c, err := NewConsulACLsAPI(aclAPI, logger)
		r.NoError(err)

		ctx := context.Background()
		sii := ServiceIdentityIndex{
			AllocID:   uuid.Generate(),
			ClusterID: uuid.Generate(),
			TaskName:  "my-task1",
		}

		token, err := c.CreateToken(ctx, sii)

		if expErr != nil {
			r.Equal(expErr, err)
			r.Nil(token)
		} else {
			r.NoError(err)
			r.Equal("my-task1", token.TaskName)
			r.True(helper.IsUUID(token.AccessorID))
			r.True(helper.IsUUID(token.SecretID))
		}
	}

	t.Run("create token success", func(t *testing.T) {
		try(nil)
	})

	t.Run("create token error", func(t *testing.T) {
		try(errors.New("consul broke"))
	})
}

func TestConsulACLsAPI_RevokeTokens(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	setup := func(exp error) (context.Context, ConsulACLsAPI, *structs.SIToken) {
		logger := testlog.HCLogger(t)
		aclAPI := consul.NewMockACLsAPI(logger)

		c, err := NewConsulACLsAPI(aclAPI, logger)
		r.NoError(err)

		ctx := context.Background()
		generated, err := c.CreateToken(ctx, ServiceIdentityIndex{TaskName: "task1"})
		r.NoError(err)

		// set the mock error after calling CreateToken for setting up
		aclAPI.SetError(exp)

		return context.Background(), c, generated
	}

	accessors := func(ids ...string) (result []*structs.SITokenAccessor) {
		for _, id := range ids {
			result = append(result, &structs.SITokenAccessor{AccessorID: id})
		}
		return
	}

	t.Run("revoke token success", func(t *testing.T) {
		ctx, c, token := setup(nil)
		err := c.RevokeTokens(ctx, accessors(token.AccessorID))
		r.NoError(err)
	})

	t.Run("revoke token non-existent", func(t *testing.T) {
		ctx, c, _ := setup(nil)
		err := c.RevokeTokens(ctx, accessors(uuid.Generate()))
		r.EqualError(err, "token does not exist")
	})

	t.Run("revoke token error", func(t *testing.T) {
		exp := errors.New("consul broke")
		ctx, c, token := setup(exp)
		err := c.RevokeTokens(ctx, accessors(token.AccessorID))
		r.EqualError(err, exp.Error())
	})
}

func TestConsulACLsAPI_ListTokens(t *testing.T) {
	t.Parallel()
	// todo
}

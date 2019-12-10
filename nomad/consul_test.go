package nomad

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/nomad/command/agent/consul"
	"github.com/hashicorp/nomad/helper"
	"github.com/hashicorp/nomad/helper/testlog"
	"github.com/hashicorp/nomad/helper/uuid"
	"github.com/stretchr/testify/require"
)

var _ ConsulACLsAPI = (*consulACLsAPI)(nil)

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
		} else {
			r.NoError(err)
			r.True(helper.IsUUID(token))
		}
	}

	t.Run("create token success", func(t *testing.T) {
		try(nil)
	})

	t.Run("create token error", func(t *testing.T) {
		try(errors.New("consul busted"))
	})
}

func TestConsulACLsAPI_RevokeTokens(t *testing.T) {
	t.Parallel()
	// todo
}

func TestConsulACLsAPI_ListTokens(t *testing.T) {
	t.Parallel()
	// todo
}

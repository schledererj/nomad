package nomad

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/armon/go-metrics"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/command/agent/consul"
	"github.com/hashicorp/nomad/nomad/structs"

	"golang.org/x/time/rate"
)

const (
	// siTokenDescriptionFmt is the format for the .Description field of
	// service identity tokens generated on behalf of Nomad.
	siTokenDescriptionFmt = "_nomad_si [%s] [%s] [%s]"

	// siTokenRequestRateLimit is the maximum number of requests per second Nomad
	// will make against Consul for requesting SI tokens.
	siTokenRequestRateLimit rate.Limit = 500

	// siTokenMaxParallelRevokes is the maximum number of parallel SI token
	// revocation requests Nomad will make against Consul.
	siTokenMaxParallelRevokes = 64

	// todo: more revocation things
)

type ServiceIdentityIndex struct {
	ClusterID string
	AllocID   string
	TaskName  string
}

func (sii ServiceIdentityIndex) Validate() error {
	switch {
	case sii.ClusterID == "":
		return errors.New("cluster id not set")
	case sii.AllocID == "":
		return errors.New("alloc id not set")
	case sii.TaskName == "":
		return errors.New("task name not set")
	default:
		return nil
	}
}

func (sii ServiceIdentityIndex) Description() string {
	return fmt.Sprintf(siTokenDescriptionFmt, sii.ClusterID, sii.AllocID, sii.TaskName)
}

// ConsulACLsAPI is the consul/api.ACL API used by Nomad Server.
type ConsulACLsAPI interface {
	CreateToken(context.Context, ServiceIdentityIndex) (*structs.SIToken, error)
	RevokeTokens(context.Context, []*structs.SITokenAccessor) error
	ListTokens() ([]string, error) // used for reconciliation
}

type consulACLsAPI struct {
	// aclClient is the API subset of the real consul client we need for
	// managing Service Identity tokens.
	aclClient consul.ACLsAPI

	// limiter is used to rate limit requests to consul
	limiter *rate.Limiter

	// logger is used to log messages
	logger hclog.Logger
}

func NewConsulACLsAPI(aclClient consul.ACLsAPI, logger hclog.Logger) (ConsulACLsAPI, error) {
	c := &consulACLsAPI{
		aclClient: aclClient,
		logger:    logger.Named("consul_acl"),
		limiter:   rate.NewLimiter(requestRateLimit, int(requestRateLimit)),
	}
	return c, nil
}

func (c *consulACLsAPI) CreateToken(ctx context.Context, sii ServiceIdentityIndex) (*structs.SIToken, error) {
	defer metrics.MeasureSince([]string{"nomad", "consul", "create_token"}, time.Now())

	// sanity check the metadata for the token we want
	if err := sii.Validate(); err != nil {
		return nil, err
	}

	// todo: use ctx

	// todo: rate limiting

	partial := &api.ACLToken{
		Description:       sii.Description(),
		ServiceIdentities: []*api.ACLServiceIdentity{{ServiceName: sii.TaskName}},
	}

	token, _, err := c.aclClient.TokenCreate(partial, nil)
	if err != nil {
		return nil, err
	}

	return &structs.SIToken{
		TaskName:   sii.TaskName,
		AccessorID: token.AccessorID,
		SecretID:   token.SecretID,
	}, nil
}

func (c *consulACLsAPI) RevokeTokens(ctx context.Context, accessors []*structs.SITokenAccessor) error {
	defer metrics.MeasureSince([]string{"nomad", "consul", "revoke_tokens"}, time.Now())

	// todo: use ctx

	// todo: rate limiting

	for _, accessor := range accessors {
		if err := c.revokeToken(ctx, accessor); err != nil {
			// todo: accumulate errors and IDs that are going to need another attempt
			return err
		}
	}

	return nil
}

func (c *consulACLsAPI) revokeToken(_ context.Context, accessor *structs.SITokenAccessor) error {
	c.logger.Trace("revoke SI token", "task", accessor.TaskName, "alloc_id", accessor.AllocID, "node_id", accessor.NodeID)
	_, err := c.aclClient.TokenDelete(accessor.AccessorID, nil)
	return err
}

func (c *consulACLsAPI) ListTokens() ([]string, error) {
	defer metrics.MeasureSince([]string{"nomad", "consul", "list_tokens"}, time.Now())

	return nil, errors.New("not yet implemented")
}
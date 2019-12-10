package consul

import (
	"errors"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/helper/uuid"
)

// MockACLsAPI is a mock of consul.ACLsAPI
type MockACLsAPI struct {
	logger hclog.Logger

	lock  sync.Mutex
	state struct {
		index  uint64
		error  error
		tokens map[string]*api.ACLToken
	}
}

func NewMockACLsAPI(l hclog.Logger) *MockACLsAPI {
	return &MockACLsAPI{
		logger: l.Named("mock_consul"),
		state: struct {
			index  uint64
			error  error
			tokens map[string]*api.ACLToken
		}{tokens: make(map[string]*api.ACLToken)},
	}
}

// SetError is a helper method for configuring an error that will be returned
// on future calls to mocked methods.
func (m *MockACLsAPI) SetError(err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.state.error = err
}

// TokenCreate is a mock of ACLsAPI.TokenCreate
func (m *MockACLsAPI) TokenCreate(token *api.ACLToken, opts *api.WriteOptions) (*api.ACLToken, *api.WriteMeta, error) {
	index, created, meta, err := m.tokenCreate(token, opts)

	services := func(token *api.ACLToken) []string {
		if token == nil {
			return nil
		}
		var names []string
		for _, id := range token.ServiceIdentities {
			names = append(names, id.ServiceName)
		}
		return names
	}(created)

	description := func(token *api.ACLToken) string {
		if token == nil {
			return "<nil>"
		}
		return token.Description
	}(created)

	accessor := func(token *api.ACLToken) string {
		if token == nil {
			return "<nil>"
		}
		return token.AccessorID
	}

	secret := func(token *api.ACLToken) string {
		if token == nil {
			return "<nil>"
		}
		return token.SecretID
	}

	m.logger.Trace("TokenCreate()", "description", description, "service_identities", services, "accessor", accessor, "secret", secret, "index", index, "error", err)
	return created, meta, err
}

func (m *MockACLsAPI) tokenCreate(token *api.ACLToken, _ *api.WriteOptions) (uint64, *api.ACLToken, *api.WriteMeta, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.state.index++

	if m.state.error != nil {
		return m.state.index, nil, nil, m.state.error
	}

	secret := &api.ACLToken{
		CreateIndex:       m.state.index,
		ModifyIndex:       m.state.index,
		AccessorID:        uuid.Generate(),
		SecretID:          uuid.Generate(),
		Description:       token.Description,
		ServiceIdentities: token.ServiceIdentities,
		CreateTime:        time.Now(),
	}

	m.state.tokens[secret.AccessorID] = secret

	w := &api.WriteMeta{
		RequestTime: 1 * time.Millisecond,
	}

	return m.state.index, secret, w, nil
}

// TokenDelete is a mock of ACLsAPI.TokenDelete
func (m *MockACLsAPI) TokenDelete(accessorID string, opts *api.WriteOptions) (*api.WriteMeta, error) {
	meta, err := m.tokenDelete(accessorID, opts)
	m.logger.Trace("TokenDelete()", "accessor", accessorID, "error", err)
	return meta, err
}

func (m *MockACLsAPI) tokenDelete(tokenID string, _ *api.WriteOptions) (*api.WriteMeta, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.state.index++

	if m.state.error != nil {
		return nil, m.state.error
	}

	if _, exists := m.state.tokens[tokenID]; !exists {
		return nil, errors.New("token does not exist")
	}

	delete(m.state.tokens, tokenID)

	m.logger.Trace("TokenDelete()")

	return nil, nil
}

// TokenList is a mock of ACLsAPI.TokenList
func (m *MockACLsAPI) TokenList(_ *api.QueryOptions) ([]*api.ACLTokenListEntry, *api.QueryMeta, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	//todo(shoenig): will need this for background token reconciliation
	// coming in another issue

	return nil, nil, nil
}

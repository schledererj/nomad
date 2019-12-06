package nomad

type TokenCreateOptions struct {
	Description string
}

// ConsulACLsAPI is the consul/api.ACL API used by Nomad Server.
type ConsulACLsAPI interface {
	TokenCreate(options TokenCreateOptions) (string, error) // what does vault iface do
	TokenDelete() error                                     // what does vault iface do
	TokenList() ([]string, error)
}

type consulClient struct {
	// todo: impl this !
}

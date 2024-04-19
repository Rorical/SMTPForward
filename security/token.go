package security

import "context"

type TokenCredential struct {
	Token string
}

func NewTokenCredential(token string) TokenCredential {
	return TokenCredential{Token: token}
}

func (c TokenCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": c.Token,
	}, nil
}

func (c TokenCredential) RequireTransportSecurity() bool {
	return true
}

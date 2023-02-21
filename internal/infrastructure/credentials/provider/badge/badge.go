package badge

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// BadgeCredentialsProvider return auth method from badge
type BadgeCredentialsProvider struct {
	methods []repository.AuthMethodConstructor
}

// NewBadgeCredentialsProvider return new instance of BadgeCredentialsProvider
func NewBadgeCredentialsProvider(methods ...repository.AuthMethodConstructor) *BadgeCredentialsProvider {
	return &BadgeCredentialsProvider{
		methods: methods,
	}
}

// Get return user password auth for docker registry
func (a *BadgeCredentialsProvider) Get(badge *credentials.Badge) (repository.AuthMethodReader, error) {
	var err error
	var method repository.AuthMethodReader

	if badge == nil {
		return nil, nil
	}

	errContext := "(credentials::provider::BadgeCredentialsProvider::Get)"

	for _, m := range a.methods {
		method, err = m.AuthMethodConstructor(badge)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}

		if method != nil {
			return method, nil
		}
	}

	return nil, nil
}

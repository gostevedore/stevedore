package factory

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/basic"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/keyfile"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/sshagent"
)

// CredentialsFactory is a factory for auth providers
type CredentialsFactory struct {
	store                repository.CredentialsStorer
	credentialsProviders []repository.CredentialsProviderer
}

// NewCredentialsFactory creates a new auth provider factory
func NewCredentialsFactory(store repository.CredentialsStorer, auth ...repository.CredentialsProviderer) *CredentialsFactory {

	factory := &CredentialsFactory{
		store: store,
	}

	factory.credentialsProviders = append([]repository.CredentialsProviderer{}, auth...)

	return factory
}

func (f *CredentialsFactory) Get(id string) (repository.AuthMethodReader, error) {

	var err error
	var badge *credentials.Badge
	var method repository.AuthMethodReader
	errContext := "(factory::CredentialsFactory::Get)"

	badge, err = f.store.Get(id)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	for _, provider := range f.credentialsProviders {
		method, err = provider.Get(badge)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}

		if method != nil {
			return method, nil
		}
	}

	return nil, nil
}

func (f *CredentialsFactory) GetBasicAuthMethod(m repository.AuthMethodReader) (*basic.BasicAuthMethod, error) {
	var correct bool
	var method *basic.BasicAuthMethod
	errContext := "(factory::CredentialsFactory::GetBasicAuthMethod)"

	if m == nil {
		return nil, errors.New(errContext, "Basic auth method could not be created because the given method is nil")
	}

	method, correct = m.(*basic.BasicAuthMethod)
	if !correct {
		return nil, nil
	}

	return method, nil
}

func (f *CredentialsFactory) GetKeyFileAuthMethod(m repository.AuthMethodReader) (*keyfile.KeyFileAuthMethod, error) {
	var correct bool
	var method *keyfile.KeyFileAuthMethod
	errContext := "(factory::CredentialsFactory::GetKeyFileAuthMethod)"

	if m == nil {
		return nil, errors.New(errContext, "Key file auth method could not be created because the given method is nil")
	}

	method, correct = m.(*keyfile.KeyFileAuthMethod)
	if !correct {
		return nil, nil
	}

	return method, nil
}

func (f *CredentialsFactory) GetSSHAgentAuthMethod(m repository.AuthMethodReader) (*sshagent.SSHAgentAuthMethod, error) {
	var correct bool
	var method *sshagent.SSHAgentAuthMethod
	errContext := "(factory::CredentialsFactory::GetSSHAgentAuthMethod)"

	if m == nil {
		return nil, errors.New(errContext, "Key file auth method could not be created because the given method is nil")
	}

	method, correct = m.(*sshagent.SSHAgentAuthMethod)
	if !correct {
		return nil, nil
	}

	return method, nil
}

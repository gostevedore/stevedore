package factory

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

type PromoteFactory map[string]repository.Promoter

func NewPromoteFactory() PromoteFactory {
	return make(PromoteFactory)
}

func (f PromoteFactory) Get(id string) (repository.Promoter, error) {
	errContext := "(PromoteFactory::GetPromoter)"

	promoter, exist := f[id]
	if !exist {
		return nil, errors.New(errContext, fmt.Sprintf("Promoter '%s' has not been registered", id))
	}

	return promoter, nil
}

func (f PromoteFactory) Register(id string, promoter repository.Promoter) error {

	errContext := "(PromoteFactory::Register)"

	_, exist := f[id]
	if exist {
		return errors.New(errContext, fmt.Sprintf("Factory '%s' already registered", id))
	}

	f[id] = promoter

	return nil
}

package operation

type FilterOperationFactory struct{}

func NewFilterOperationFactory() *FilterOperationFactory {
	return &FilterOperationFactory{}
}

func (f *FilterOperationFactory) FilterOperation() *FilterOperation {
	return NewFilterOperation()
}

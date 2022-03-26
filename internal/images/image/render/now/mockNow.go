package now

// MockNow is a mock implementation of the Now interface
type MockNow struct{}

// NewMockNow returns a new instance of the MockNow
func NewMockNow() *MockNow {
	return &MockNow{}
}

// NowFunc returns a function that return the same string passed as argument
func (t *MockNow) NowFunc() func(string) string {
	return func(layout string) string {
		return layout
	}
}

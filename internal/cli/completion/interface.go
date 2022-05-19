package completion

// Consoler interface to show messages through console
type Consoler interface {
	Write(data []byte) (int, error)
}

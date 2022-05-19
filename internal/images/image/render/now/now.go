package now

import "time"

// Now returns the current local time formated
type Now struct{}

// NewNow returns a new instance of the Now
func NewNow() *Now {
	return &Now{}
}

// NowFunc returns a function that return the current local time formated
func (n *Now) NowFunc() func(string) string {
	return func(layout string) string {
		return time.Now().Format(layout)
	}
}

package render

// Nower is the interface to the timer that generates formated types
type Nower interface {
	NowFunc() func(layout string) string
}

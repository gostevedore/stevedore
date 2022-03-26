package render

// ImageSerializer is the interface for the image serializer
// type ImageSerializer interface {
// 	YAMLMarshal() ([]byte, error)
// 	YAMLUnmarshal([]byte) error
// }

// Nower is the interface to the timer that generates formated types
type Nower interface {
	NowFunc() func(layout string) string
}

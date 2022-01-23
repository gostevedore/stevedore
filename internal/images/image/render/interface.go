package render

// ImageSerializer is the interface for the image serializer
type ImageSerializer interface {
	YAMLMarshal() ([]byte, error)
	YAMLUnmarshal([]byte) error
}

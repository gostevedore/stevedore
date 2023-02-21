package configuration

type EncryptionKeyGenerator interface {
	GenerateEncryptionKey() (string, error)
}

package client

const (
	RotateEncryptionKeyOutputType         = "rotateEncryptionKeyOutput"
	RotateEncryptionKeyOutputFieldBackup  = "backup"
	RotateEncryptionKeyOutputFieldMessage = "message"
)

type RotateEncryptionKeyOutput struct {
	Backup  *EtcdBackup `json:"backup,omitempty" yaml:"backup,omitempty"`
	Message string      `json:"message,omitempty" yaml:"message,omitempty"`
}

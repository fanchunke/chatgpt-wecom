package errno

import "errors"

var (
	ErrInvalidSignature              = errors.New("Invalid signature")
	ErrInvalidEncryptMsgSize         = errors.New("Invalid encrypt_msg size")
	ErrInvalidPKCS7UnpaddingTextSize = errors.New("Invalid PKCS7Unpadding text size")
	ErrInvalidPlainPayloadSize       = errors.New("Invalid plain payload size")
	ErrInvalidJson                   = errors.New("Invalid Json")
)

package helpers

import (
	"crypto/aes"
	"fmt"
)

type CryptoHelper struct {
	key []byte
}

func NewHelper(key []byte) *CryptoHelper {
	return &CryptoHelper{
		key: key,
	}
}

func (h *CryptoHelper) Encode(s string) []byte {
	aesblock, err := aes.NewCipher(h.key)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	dst := make([]byte, aes.BlockSize) // шифруем
	aesblock.Encrypt(dst, []byte(s))

	return dst
}

func (h *CryptoHelper) Decode(s string) []byte {
	aesblock, err := aes.NewCipher(h.key)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	src2 := make([]byte, aes.BlockSize) // расшифровываем
	aesblock.Decrypt(src2, []byte(s))
	return src2
}

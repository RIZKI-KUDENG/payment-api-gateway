package credential

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)



func Encrypt(plaintext []byte, key []byte) ([]byte, error){
	if len(key) != 32 {
		return nil, errors.New("master key must be 32bytes for AES-256")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(ciphertext []byte, key []byte)([]byte, error){
	if len(key) != 32{
		return nil, errors.New("master key must be 32bytes for AES-256")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize{
		return nil, errors.New("cipher text too short")
	}

	nonce, actualCipherText := ciphertext[:nonceSize], ciphertext[nonceSize:]

	return gcm.Open(nil, nonce, actualCipherText, nil)

}

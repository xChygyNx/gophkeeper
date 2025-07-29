package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func Encrypt(plaintext string, secretKey []byte) (string, error) {
	aes, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", fmt.Errorf("error in create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", fmt.Errorf("error in get Galois Counter Mode: %w", err)
	}

	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", fmt.Errorf("error in create random byte sequence: %w", err)
	}

	// ciphertext here is actually nonce+ciphertext
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return string(ciphertext), nil
}

func Decrypt(ciphertext string, secretKey []byte) (string, error) {
	aesBlock, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", fmt.Errorf("error in create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", fmt.Errorf("error in get Galois Counter Mode: %w", err)
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", fmt.Errorf("error in decrypt data: %w", err)
	}

	return string(plaintext), nil
}

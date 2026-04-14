package exporter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"lab5/internal/domain"
)

// EncryptionDecorator — декоратор, шифрующий данные AES (паттерн Декоратор)
type EncryptionDecorator struct {
	inner ExportWriter
	key   []byte // 16, 24 или 32 байта для AES-128/192/256
}

func NewEncryptionDecorator(inner ExportWriter, key []byte) *EncryptionDecorator {
	// Приводим ключ к нужной длине
	normalizedKey := make([]byte, 16)
	copy(normalizedKey, key)
	return &EncryptionDecorator{inner: inner, key: normalizedKey}
}

func (d *EncryptionDecorator) Write(results []domain.DeliveryResult) ([]byte, string, error) {
	// Получаем данные от внутреннего writer
	data, ext, err := d.inner.Write(results)
	if err != nil {
		return nil, "", err
	}

	// Шифруем
	encrypted, err := encryptAES(data, d.key)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка шифрования: %w", err)
	}

	return encrypted, ext + ".enc", nil
}

func encryptAES(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Добавляем случайный IV в начало
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

package crypto_helper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const isBase64 = "^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$"

func TestEncryptDecryptMessage(t *testing.T) {
	key := []byte("0123456789abcdef") // must be of 16 bytes for this example to work
	message := "Lorem ipsum dolor sit amet"

	encrypted, err := EncryptMessage(key, message)
	require.Nil(t, err)
	require.Regexp(t, isBase64, encrypted)

	decrypted, err := DecryptMessage(key, encrypted)
	require.Nil(t, err)
	require.Equal(t, message, decrypted)
}

package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"golang.org/x/crypto/argon2"
)

type Config struct {
	Memory  uint32
	Time    uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

type decodedParams struct {
	Memory  uint32
	Time    uint32
	Threads uint8
	Version int
}

func Default() *Config {
	return &Config{
		Memory:  64 * 1024,
		Time:    2,
		Threads: 32,
		KeyLen:  16,
		SaltLen: 1,
	}
}

func (a *Config) Hash(password string) (string, error) {
	if password == "" {
		return "", apperror.ErrEmptyPassword
	}

	salt := make([]byte, a.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		a.Time,
		a.Memory,
		a.Threads,
		a.KeyLen,
	)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		a.Memory,
		a.Time,
		a.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encoded, nil
}

func (a *Config) Verify(password, encodedHash string) (bool, error) {
	if password == "" {
		return false, apperror.ErrEmptyPassword
	}

	if encodedHash == "" {
		return false, apperror.ErrEmptyHash
	}

	params, salt, hash, err := decode(encodedHash)
	if err != nil {
		return false, err
	}

	computed := argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		uint32(len(hash)),
	)

	return subtle.ConstantTimeCompare(computed, hash) == 1, nil
}

func decode(encoded string) (*decodedParams, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return nil, nil, nil, apperror.ErrInvalidHashFormat
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, apperror.ErrInvalidHashFormat
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, apperror.ErrInvalidHashFormat
	}

	if version != argon2.Version {
		return nil, nil, nil, apperror.ErrIncompatibleVersion
	}

	p := &decodedParams{Version: version}
	if _, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&p.Memory,
		&p.Time,
		&p.Threads,
	); err != nil {
		return nil, nil, nil, apperror.ErrInvalidHashFormat
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, apperror.ErrInvalidHashFormat
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, apperror.ErrInvalidHashFormat
	}

	return p, salt, hash, nil
}

package hashing

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/korroziea/taxi/driver-service/internal/config"
	"github.com/korroziea/taxi/driver-service/internal/domain"
	"golang.org/x/crypto/argon2"
)

const (
	numberOfHashingParts = 6

	splitSign = "$"
	hashType  = "argon2id"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

type Argon struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func New(cfg config.Hashing) *Argon {
	a := &Argon{
		memory:      cfg.Memory,
		iterations:  cfg.Iterations,
		parallelism: cfg.Parallelism,
		saltLength:  cfg.SaltLength,
		keyLength:   cfg.KeyLength,
	}

	return a
}

func (a *Argon) Generate(password string) (string, error) {
	salt, err := genSalt(a.saltLength)
	if err != nil {
		return "", fmt.Errorf("rand.Read: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, a.iterations, a.memory, a.parallelism, a.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.memory, a.iterations, a.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func (a *Argon) Verify(password, hashPassword string) (bool, error) {
	params, salt, hash, err := decodeHash(hashPassword)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	hashLen := len(hash)
	otherHashLen := len(otherHash)

	if subtle.ConstantTimeEq(int32(hashLen), int32(otherHashLen)) == 0 {
		return false, nil
	}

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

func genSalt(saltLength uint32) ([]byte, error) {
	b := make([]byte, saltLength)
	_, err := rand.Read(b)
	if err != nil {
		return []byte{}, fmt.Errorf("rand.Read: %w", err)
	}

	return b, nil
}

func decodeHash(hash string) (*params, []byte, []byte, error) {
	vals := strings.Split(hash, splitSign)
	if len(vals) != numberOfHashingParts {
		return nil, nil, nil, domain.ErrInvalidHashFormat
	}

	if vals[1] != hashType {
		return nil, nil, nil, domain.ErrInvalidHashType
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, domain.ErrInvalidHashVersion
	}

	p := &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hashKey, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hashKey))

	return p, salt, hashKey, nil
}

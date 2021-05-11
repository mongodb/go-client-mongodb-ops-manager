// Copyright 2020 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atmcfg

import (
	"crypto/hmac"
	"crypto/md5" //nolint:gosec // used as part of the sha1 standard
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/xdg-go/stringprep"
	"go.mongodb.org/ops-manager/opsmngr"
)

const (
	// using the default MongoDB values for the number of iterations depending on mechanism
	scramSha1Iterations     = 10000
	scramSha256Iterations   = 15000
	clientKeyInput          = "Client Key" // specified in RFC 5802
	serverKeyInput          = "Server Key" // specified in RFC 5802
	rfc5802MandatedSaltSize = 4
)

type hashingFunc func() hash.Hash

func computeScramCredentials(f hashingFunc, iterationCount int, base64EncodedSalt, password string) (*opsmngr.ScramShaCreds, error) {
	// password should be encrypted in the case of SCRAM-SHA-1 and unencrypted in the case of SCRAM-SHA-256
	storedKey, serverKey, err := generateB64EncodedSecrets(f, password, base64EncodedSalt, iterationCount)
	if err != nil {
		return nil, fmt.Errorf("error generating SCRAM-SHA keys: %s", err)
	}

	return &opsmngr.ScramShaCreds{IterationCount: iterationCount, Salt: base64EncodedSalt, StoredKey: storedKey, ServerKey: serverKey}, nil
}

func generateB64EncodedSecrets(f hashingFunc, password, b64EncodedSalt string, iterationCount int) (storedKey, serverKey string, err error) {
	salt, err := base64.StdEncoding.DecodeString(b64EncodedSalt)
	if err != nil {
		return "", "", fmt.Errorf("error decoding salt: %s", err)
	}

	unencodedStoredKey, unencodedServerKey, err := generateSecrets(f, password, salt, iterationCount)
	if err != nil {
		return "", "", fmt.Errorf("error generating secrets: %s", err)
	}

	storedKey = base64.StdEncoding.EncodeToString(unencodedStoredKey)
	serverKey = base64.StdEncoding.EncodeToString(unencodedServerKey)
	return storedKey, serverKey, nil
}

func generateSecrets(f hashingFunc, password string, salt []byte, iterationCount int) (storedKey, serverKey []byte, err error) {
	saltedPassword, err := generateSaltedPassword(f, password, salt, iterationCount)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating salted password: %s", err)
	}

	clientKey, err := generateClientOrServerKey(f, saltedPassword, clientKeyInput)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating client key: %s", err)
	}

	storedKey, err = generateStoredKey(f, clientKey)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating stored key: %s", err)
	}

	serverKey, err = generateClientOrServerKey(f, saltedPassword, serverKeyInput)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating server key: %s", err)
	}

	return storedKey, serverKey, err
}

func generateSaltedPassword(hashConstructor func() hash.Hash, password string, salt []byte, iterationCount int) ([]byte, error) {
	preparedPassword, err := stringprep.SASLprep.Prepare(password)
	if err != nil {
		return nil, fmt.Errorf("error SASLprep'ing password: %s", err)
	}

	result, err := hmacIteration(hashConstructor, []byte(preparedPassword), salt, iterationCount)
	if err != nil {
		return nil, fmt.Errorf("error running hmacIteration: %s", err)
	}
	return result, nil
}

func generateClientOrServerKey(f hashingFunc, saltedPassword []byte, input string) ([]byte, error) {
	hmacHash := hmac.New(f, saltedPassword)
	if _, err := hmacHash.Write([]byte(input)); err != nil {
		return nil, fmt.Errorf("error running hmacHash: %s", err)
	}

	return hmacHash.Sum(nil), nil
}

func generateStoredKey(f hashingFunc, clientKey []byte) ([]byte, error) {
	h := f()
	if _, err := h.Write(clientKey); err != nil {
		return nil, fmt.Errorf("error hashing: %s", err)
	}
	return h.Sum(nil), nil
}

func hmacIteration(f hashingFunc, input, salt []byte, iterationCount int) ([]byte, error) {
	hashSize := f().Size()

	// incorrect salt size will pass validation, but the credentials will be invalid. i.e. it will not
	// be possible to auth with the password provided to create the credentials.
	if len(salt) != hashSize-rfc5802MandatedSaltSize {
		return nil, fmt.Errorf("salt should have a size of %v bytes, but instead has a size of %v bytes", hashSize-rfc5802MandatedSaltSize, len(salt))
	}

	startKey := append(salt, 0, 0, 0, 1) //nolint:gocritic // this is assigment is correct
	result := make([]byte, hashSize)

	hmacHash := hmac.New(f, input)
	if _, err := hmacHash.Write(startKey); err != nil {
		return nil, fmt.Errorf("error running hmacHash: %s", err)
	}

	intermediateDigest := hmacHash.Sum(nil)

	for i := 0; i < len(intermediateDigest); i++ {
		result[i] = intermediateDigest[i]
	}

	for i := 1; i < iterationCount; i++ {
		hmacHash.Reset()
		if _, err := hmacHash.Write(intermediateDigest); err != nil {
			return nil, fmt.Errorf("error running hmacHash: %s", err)
		}

		intermediateDigest = hmacHash.Sum(nil)

		for i := 0; i < len(intermediateDigest); i++ {
			result[i] ^= intermediateDigest[i]
		}
	}

	return result, nil
}

// generateSalt will create a salt for use with newScramShaCreds based on the given hashConstructor.
// sha1.New should be used for MONGODB-CR/SCRAM-SHA-1 and sha256.New should be used for SCRAM-SHA-256
func generateSalt(hashConstructor func() hash.Hash) ([]byte, error) {
	saltSize := hashConstructor().Size() - rfc5802MandatedSaltSize
	salt, err := generateRandomBase64String(saltSize)
	if err != nil {
		return nil, err
	}
	return []byte(salt), nil
}

func md5Hex(s string) (string, error) {
	h := md5.New() //nolint:gosec // used as part of the sha1 standard

	if _, err := h.Write([]byte(s)); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

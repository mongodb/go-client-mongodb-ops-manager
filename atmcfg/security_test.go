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
	"crypto/sha1" // #nosec G101 // #nosec G505 // used as part of the sha1 standard
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateB64EncodedSecrets(t *testing.T) {
	// these were taken from MongoDB.  passwordHash is from authSchema
	// 3. iterationCount, salt, storedKey, and serverKey are from
	// authSchema 5 (after upgrading from authSchema 3)
	testCases := map[string]map[string]string{
		"test 1": {
			"passwordHash": "caeec61ba3b15b15b188d29e876514e8",
			"salt":         "S3cuk2Rnu/MlbewzxrmmVA==",
			"storedKey":    "sYBa3XlSPKNrgjzhOuEuRlJY4dQ=",
			"serverKey":    "zuAxRSQb3gZkbaB1IGlusK4jy1M=",
		},
		"test 2": {
			"passwordHash": "4d9625b297999b3ca786d4a9622d04f1",
			"salt":         "kW9KbCQiCOll5Ljd44cjkQ==",
			"storedKey":    "VJ8fFVHkPltibvT//mG/OWw44Hc=",
			"serverKey":    "ceDRsgj9HezpZ4/vkZX8GZNNN50=",
		},
		"test 3": {
			"passwordHash": "fd0a78e418dcef39f8c768222810b894",
			"salt":         "hhX6xsoID6FeWjXncuNgAg==",
			"storedKey":    "TxgaZJ4cIn+S9EfTcc9IOEG7RGc=",
			"serverKey":    "d6/qjwBs0qkPKfUAjSh5eemsySE=",
		},
	}
	for testType, values := range testCases {
		passwordHash := values["passwordHash"]
		salt := values["salt"]
		storedKey := values["storedKey"]
		serverKey := values["serverKey"]
		t.Run(testType, func(t *testing.T) {
			computedStoredKey, computedServerKey, err := generateB64EncodedSecrets(sha1.New, passwordHash, salt, 10)
			a := assert.New(t)
			if a.NoError(err) {
				a.Equal(storedKey, computedStoredKey)
				a.Equal(serverKey, computedServerKey)
			}
		})
	}
}

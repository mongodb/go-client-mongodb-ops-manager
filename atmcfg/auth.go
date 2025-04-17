// Copyright 2022 MongoDB Inc
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
	"crypto/sha1" // #nosec G101 // #nosec G505 // mongodb scram-sha-1 supports this tho is not recommended
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"go.mongodb.org/ops-manager/opsmngr"
	"go.mongodb.org/ops-manager/search"
)

const (
	automationAgentName            = "mms-automation"
	atmAgentWindowsKeyFilePath     = "%SystemDrive%\\MMSAutomation\\versions\\keyfile"
	atmAgentKeyFilePathInContainer = "/var/lib/mongodb-mms-automation/keyfile"
	keyLength                      = 500
	mongoCR                        = "MONGODB-CR"
	scramSha256                    = "SCRAM-SHA-256"
)

// EnableMechanism allows you to enable a given set of authentication mechanisms to an opsmngr.AutomationConfig.
// This method currently only supports MONGODB-CR, and SCRAM-SHA-256.
func EnableMechanism(out *opsmngr.AutomationConfig, m []string) error {
	out.Auth.Disabled = false
	for _, v := range m {
		if v != mongoCR && v != scramSha256 {
			return fmt.Errorf("unsupported mechanism %s", v)
		}
		if v == scramSha256 && out.Auth.AutoAuthMechanism == "" {
			out.Auth.AutoAuthMechanism = v
		}
		if !stringInSlice(out.Auth.DeploymentAuthMechanisms, v) {
			out.Auth.DeploymentAuthMechanisms = append(out.Auth.DeploymentAuthMechanisms, v)
		}
		if !stringInSlice(out.Auth.AutoAuthMechanisms, v) {
			out.Auth.AutoAuthMechanisms = append(out.Auth.AutoAuthMechanisms, v)
		}
	}

	if out.Auth.AutoUser == "" && out.Auth.AutoPwd == "" {
		if err := setAutoUser(out); err != nil {
			return err
		}
	}

	if out.Auth.Key == "" {
		var err error
		if out.Auth.Key, err = generateRandomBase64String(keyLength); err != nil {
			return err
		}
	}
	if out.Auth.Keyfile == "" {
		out.Auth.Keyfile = atmAgentKeyFilePathInContainer
	}
	if out.Auth.KeyfileWindows == "" {
		out.Auth.KeyfileWindows = atmAgentWindowsKeyFilePath
	}

	return nil
}

func setAutoUser(out *opsmngr.AutomationConfig) error {
	var err error
	out.Auth.AutoUser = automationAgentName
	if out.Auth.AutoPwd, err = generateRandomASCIIString(keyLength); err != nil {
		return err
	}

	return nil
}

// AddUser adds a opsmngr.MongoDBUser to the opsmngr.AutomationConfi.
func AddUser(out *opsmngr.AutomationConfig, u *opsmngr.MongoDBUser) {
	out.Auth.UsersWanted = append(out.Auth.UsersWanted, u)
}

// RemoveUser removes a MongoDBUser from the authentication config.
func RemoveUser(out *opsmngr.AutomationConfig, username, database string) error {
	pos, found := search.MongoDBUsers(out.Auth.UsersWanted, func(p *opsmngr.MongoDBUser) bool {
		return p.Username == username && p.Database == database
	})
	if !found {
		return fmt.Errorf("user '%s' not found for '%s'", username, database)
	}
	out.Auth.UsersWanted = append(out.Auth.UsersWanted[:pos], out.Auth.UsersWanted[pos+1:]...)
	return nil
}

// ConfigureScramCredentials creates both SCRAM-SHA-1 and SCRAM-SHA-256 credentials.
// Use this method to guarantee that password can be updated later.
func ConfigureScramCredentials(user *opsmngr.MongoDBUser, password string) error {
	scram256Creds, err := newScramSha256Creds(user, password)
	if err != nil {
		return err
	}

	scram1Creds, err := newScramSha1Creds(user, password)
	if err != nil {
		return err
	}
	user.ScramSha256Creds = scram256Creds
	user.ScramSha1Creds = scram1Creds
	return nil
}

func newScramSha1Creds(user *opsmngr.MongoDBUser, password string) (*opsmngr.ScramShaCreds, error) {
	scram1Salt, err := generateSalt(sha1.New)
	if err != nil {
		return nil, fmt.Errorf("error generating scramSha1 salt: %w", err)
	}
	scram1Creds, err := newScramShaCreds(scram1Salt, user.Username, password, mongoCR)
	if err != nil {
		return nil, fmt.Errorf("error generating scramSha1Creds: %w", err)
	}
	return scram1Creds, nil
}

func newScramSha256Creds(user *opsmngr.MongoDBUser, password string) (*opsmngr.ScramShaCreds, error) {
	scram256Salt, err := generateSalt(sha256.New)
	if err != nil {
		return nil, fmt.Errorf("error generating scramSha256 salt: %w", err)
	}
	scram256Creds, err := newScramShaCreds(scram256Salt, user.Username, password, scramSha256)
	if err != nil {
		return nil, fmt.Errorf("error generating scramSha256 creds: %w", err)
	}
	return scram256Creds, nil
}

// ErrUnsupportedMechanism means the provided mechanism wasn't valid.
var ErrUnsupportedMechanism = errors.New("unrecognized SCRAM-SHA format")

// newScramShaCreds takes a plain text password and a specified mechanism name and generates
// the ScramShaCreds which will be embedded into a MongoDBUser.
func newScramShaCreds(salt []byte, username, password, mechanism string) (*opsmngr.ScramShaCreds, error) {
	if mechanism != scramSha256 && mechanism != mongoCR {
		return nil, fmt.Errorf("%w %s", ErrUnsupportedMechanism, mechanism)
	}
	var hashConstructor hashingFunc
	iterations := 0

	switch mechanism {
	case scramSha256:
		hashConstructor = sha256.New
		iterations = scramSha256Iterations
	case mongoCR:
		hashConstructor = sha1.New
		iterations = scramSha1Iterations

		// MONGODB-CR/SCRAM-SHA-1 requires the hash of the password being passed computeScramCredentials
		// instead of the plain text password.
		var err error
		password, err = md5Hex(username + ":mongo:" + password)
		if err != nil {
			return nil, err
		}
	}

	base64EncodedSalt := base64.StdEncoding.EncodeToString(salt)
	return computeScramCredentials(hashConstructor, iterations, base64EncodedSalt, password)
}

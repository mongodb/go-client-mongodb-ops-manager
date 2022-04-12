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
	"testing"

	"go.mongodb.org/ops-manager/opsmngr"
)

func TestAddUser(t *testing.T) {
	config := automationConfigWithoutMongoDBUsers()
	u := mongoDBUsers()
	AddUser(config, u)
	if len(config.Auth.UsersWanted) != 1 {
		t.Error("User not added\n")
	}
}

func TestRemoveUser(t *testing.T) {
	config := automationConfigWithMongoDBUsers()
	t.Run("user exists", func(t *testing.T) {
		u := mongoDBUsers()
		err := RemoveUser(config, u.Username, u.Database)
		if err != nil {
			t.Fatalf("RemoveUser unexpecter err: %#v\n", err)
		}
		if len(config.Auth.UsersWanted) != 0 {
			t.Error("User not removed\n")
		}
	})
	t.Run("user does not exists", func(t *testing.T) {
		err := RemoveUser(config, "random", "random")
		if err == nil {
			t.Fatal("RemoveUser should return an error\n")
		}
	})
}

func TestEnableMechanism(t *testing.T) {
	config := automationConfigWithoutMongoDBUsers()
	t.Run("enable invalid", func(t *testing.T) {
		if e := EnableMechanism(config, []string{"invalid"}); e == nil {
			t.Fatalf("EnableMechanism() expected an error but got none\n")
		}
	})
	t.Run("enable SCRAM-SHA-256", func(t *testing.T) {
		if e := EnableMechanism(config, []string{"SCRAM-SHA-256"}); e != nil {
			t.Fatalf("EnableMechanism() unexpected error: %v\n", e)
		}

		if config.Auth.Disabled {
			t.Error("config.Auth.Disabled is true\n")
		}

		if config.Auth.AutoAuthMechanisms[0] != "SCRAM-SHA-256" {
			t.Error("AutoAuthMechanisms not set\n")
		}

		if config.Auth.AutoUser == "" || config.Auth.AutoPwd == "" {
			t.Error("config.Auth.Auto* not set\n")
		}

		if config.Auth.Key == "" || config.Auth.KeyfileWindows == "" || config.Auth.Keyfile == "" {
			t.Error("config.Auth.Key* not set\n")
		}

		if len(config.Auth.UsersWanted) != 0 {
			t.Errorf("expected 0 user got: %d\n", len(config.Auth.UsersWanted))
		}
	})
}

func TestConfigureScramCredentials(t *testing.T) {
	u := &opsmngr.MongoDBUser{
		Username: "test",
	}
	if err := ConfigureScramCredentials(u, "password"); err != nil {
		t.Fatalf("ConfigureScramCredentials() unexpected error: %v\n", err)
	}
	if u.ScramSha1Creds == nil {
		t.Fatalf("ConfigureScramCredentials() unexpected error: %v\n", u.ScramSha1Creds)
	}
	if u.ScramSha256Creds == nil {
		t.Fatalf("ConfigureScramCredentials() unexpected error: %v\n", u.ScramSha256Creds)
	}
}

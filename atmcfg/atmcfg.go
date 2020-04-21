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
	"errors"
	"fmt"

	"github.com/mongodb/go-client-mongodb-ops-manager/opsmngr"
	"github.com/mongodb/go-client-mongodb-ops-manager/search"
)

func setDisabledByClusterName(out *opsmngr.AutomationConfig, name string, disabled bool) {
	// This value may not be present and is mandatory
	if out.Auth.DeploymentAuthMechanisms == nil {
		out.Auth.DeploymentAuthMechanisms = make([]string, 0)
	}
	for _, rs := range out.ReplicaSets {
		if rs.ID == name {
			for _, m := range rs.Members {
				for k, p := range out.Processes {
					if p.Name == m.Host {
						out.Processes[k].Disabled = disabled
					}
				}
			}
			break
		}
	}
}

// Shutdown disables all processes of the given cluster name
func Shutdown(out *opsmngr.AutomationConfig, name string) {
	setDisabledByClusterName(out, name, true)
}

// Startup enables all processes of the given cluster name
func Startup(out *opsmngr.AutomationConfig, name string) {
	setDisabledByClusterName(out, name, false)
}

// AddUser adds a MongoDBUser to the authentication config
func AddUser(out *opsmngr.AutomationConfig, u *opsmngr.MongoDBUser) {
	out.Auth.Users = append(out.Auth.Users, u)
}

// AddIndexConfig adds an IndexConfig to the authentication config
func AddIndexConfig(out *opsmngr.AutomationConfig, newIndex *opsmngr.IndexConfigs) error {
	if out == nil {
		return errors.New("the Automation Config has not been initialized")
	}
	_, exists := search.MongoDBIndexes(out.IndexConfigs, func(index *opsmngr.IndexConfigs) bool {
		if index.RSName == newIndex.RSName && index.CollectionName == newIndex.CollectionName && index.DBName == newIndex.DBName && len(index.Key) == len(newIndex.Key) {
			// if keys are the equal the two indexes are considered to be the same
			for i := 0; i < len(index.Key); i++ {
				if index.Key[i][0] != newIndex.Key[i][0] || index.Key[i][1] != newIndex.Key[i][1] {
					return false
				}
			}

			return true
		}

		return false
	})
	if !exists {
		out.IndexConfigs = append(out.IndexConfigs, newIndex)
	}

	return nil
}

// RemoveUser removes a MongoDBUser from the authentication config
func RemoveUser(out *opsmngr.AutomationConfig, username string, database string) error {
	pos, found := search.MongoDBUsers(out.Auth.Users, func(p *opsmngr.MongoDBUser) bool {
		return p.Username == username && p.Database == database
	})
	if !found {
		return fmt.Errorf("user '%s' not found for '%s'", username, database)
	}
	out.Auth.Users = append(out.Auth.Users[:pos], out.Auth.Users[pos+1:]...)
	return nil
}

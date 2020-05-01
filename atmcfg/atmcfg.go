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

	"go.mongodb.org/ops-manager/opsmngr"
	"go.mongodb.org/ops-manager/search"
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
func AddIndexConfig(out *opsmngr.AutomationConfig, newIndex *opsmngr.IndexConfig) error {
	if out == nil {
		return errors.New("the Automation Config has not been initialized")
	}
	_, exists := search.MongoDBIndexes(out.IndexConfigs, compareIndexConfig(newIndex))

	if exists {
		return errors.New("index already exists")
	}
	out.IndexConfigs = append(out.IndexConfigs, newIndex)

	return nil
}

// compareIndexConfig returns a function that compares two indexConfig struts
func compareIndexConfig(newIndex *opsmngr.IndexConfig) func(index *opsmngr.IndexConfig) bool {
	return func(index *opsmngr.IndexConfig) bool {
		if newIndex.RSName == index.RSName && newIndex.CollectionName == index.CollectionName && newIndex.DBName == index.DBName && len(newIndex.Key) == len(index.Key) {
			// if keys are equal the two indexes are considered to be the same
			for i := 0; i < len(newIndex.Key); i++ {
				if newIndex.Key[i][0] != index.Key[i][0] || newIndex.Key[i][1] != index.Key[i][1] {
					return false
				}
			}

			return true
		}
		return false
	}

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

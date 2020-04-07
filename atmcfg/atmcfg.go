package atmcfg

import (
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

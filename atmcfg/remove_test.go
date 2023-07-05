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

import "testing"

func TestRemoveByClusterName(t *testing.T) {
	t.Run("replica set", func(t *testing.T) {
		config := automationConfigWithOneReplicaSet(clusterName, false)

		RemoveByClusterName(config, clusterName)
		if len(config.Processes) != 0 {
			t.Errorf("Got = %#v, want = 0", len(config.Processes))
		}
	})
	t.Run("sharded cluster", func(t *testing.T) {
		config := automationConfigWithThreeShardsCluster(clusterName, false)

		RemoveByClusterName(config, clusterName)
		if len(config.Processes) != 0 {
			t.Errorf("Got = %#v, want = 0", len(config.Processes))
		}
	})
}

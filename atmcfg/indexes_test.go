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

func TestAddIndexConfig(t *testing.T) {
	newIndex := &opsmngr.IndexConfig{
		DBName:         "test",
		CollectionName: "test",
		RSName:         "myReplicaSet",
		Key: [][]string{
			{
				"test1", "test",
			},
		},
		Options:   nil,
		Collation: nil,
	}
	t.Run("AutomationConfig not initialized", func(t *testing.T) {
		err := AddIndexConfig(nil, newIndex)
		if err == nil {
			t.Error("AddIndexConfig should return an error")
		}
	})

	t.Run("empty IndexConfig", func(t *testing.T) {
		a := &opsmngr.AutomationConfig{}
		err := AddIndexConfig(a, newIndex)
		if err != nil {
			t.Fatalf("AddIndexConfig unexpected error: %v", err)
		}
		if len(a.IndexConfigs) != 1 {
			t.Error("indexConfig has not been added to the AutomationConfig")
		}
	})

	t.Run("add an index with different keys", func(t *testing.T) {
		config := automationConfigWithIndexConfig()
		err := AddIndexConfig(config, newIndex)
		if err != nil {
			t.Fatalf("AutomationConfig() returned an unexpected error: %v", err)
		}
		if len(config.IndexConfigs) != 2 {
			t.Error("indexConfig has not been added to the AutomationConfig")
		}
	})

	t.Run("add an index with different rsName", func(t *testing.T) {
		newIndex := &opsmngr.IndexConfig{
			DBName:         "test",
			CollectionName: "test",
			RSName:         "myReplicaSet_1",
			Key: [][]string{
				{
					"test1", "test",
				},
			},
			Options:   nil,
			Collation: nil,
		}
		config := automationConfigWithIndexConfig()
		err := AddIndexConfig(config, newIndex)
		if err != nil {
			t.Fatalf("AutomationConfig() returned an unexpected error: %v", err)
		}
		if len(config.IndexConfigs) != 2 {
			t.Error("indexConfig has not been added to the AutomationConfig")
		}
	})

	t.Run("trying to add an index that is already in the AutomationConfig", func(t *testing.T) {
		config := automationConfigWithIndexConfig()
		index := &opsmngr.IndexConfig{
			DBName:         "test",
			CollectionName: "test",
			RSName:         "myReplicaSet",
			Key: [][]string{
				{
					"test", "test",
				},
			},
			Options:   nil,
			Collation: nil,
		}
		err := AddIndexConfig(config, index)

		if err == nil {
			t.Fatalf("AddIndexConfig should return an error")
		}
	})
}

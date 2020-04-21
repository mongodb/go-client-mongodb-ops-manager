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

package search

import (
	"github.com/mongodb/go-client-mongodb-ops-manager/opsmngr"
)

// Processes return the smallest index i
// in [0, n) at which f(i) is true, assuming that on the range [0, n),
// f(i) == true implies f(i+1) == true.
// returns the first true index. If there is no such index, Processes returns n and false
func Processes(a []*opsmngr.Process, f func(*opsmngr.Process) bool) (int, bool) {
	for i, p := range a {
		if f(p) {
			return i, true
		}
	}
	return len(a), false
}

// Members return the smallest index i
// in [0, n) at which f(i) is true, assuming that on the range [0, n),
// f(i) == true implies f(i+1) == true.
// returns the first true index. If there is no such index, Members returns n and false
func Members(a []opsmngr.Member, f func(opsmngr.Member) bool) (int, bool) {
	for i, m := range a {
		if f(m) {
			return i, true
		}
	}
	return len(a), false
}

// ReplicaSets return the smallest index i
// in [0, n) at which f(i) is true, assuming that on the range [0, n),
// f(i) == true implies f(i+1) == true.
// returns the first true index. If there is no such index, ReplicaSets returns n and false
func ReplicaSets(a []*opsmngr.ReplicaSet, f func(*opsmngr.ReplicaSet) bool) (int, bool) {
	for i, m := range a {
		if f(m) {
			return i, true
		}
	}
	return len(a), false
}

// MongoDBUser return the smallest index i
// in [0, n) at which f(i) is true, assuming that on the range [0, n),
// f(i) == true implies f(i+1) == true.
// returns the first true index. If there is no such index, MongoDBUsers returns n and false
func MongoDBUsers(a []*opsmngr.MongoDBUser, f func(*opsmngr.MongoDBUser) bool) (int, bool) {
	for i, m := range a {
		if f(m) {
			return i, true
		}
	}
	return len(a), false
}

// MongoDBIndexes return the smallest index i
// in [0, n) at which f(i) is true, assuming that on the range [0, n),
// f(i) == true implies f(i+1) == true.
// returns the first true index. If there is no such index, MongoDBIndexes returns n and false
func MongoDBIndexes(a []*opsmngr.IndexConfig, f func(configs *opsmngr.IndexConfig) bool) (int, bool) {
	for i, m := range a {
		if f(m) {
			return i, true
		}
	}
	return len(a), false
}

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

package opsmngr

import "encoding/json"

// JSONArray is the type alias for []string
type JSONArrayStr []string

// MarshalJSON encodes JSONArray object into bytes. If object is nil it provides empty array
func (a JSONArrayStr) MarshalJSON() ([]byte, error) {
	tmp := make([]string, 0)
	if a != nil {
		tmp = a
	}
	return json.Marshal(&tmp)
}

// JSONArrayOfArray is the type alias for [][]string
type JSONArrayOfArrayStr [][]string

// MarshalJSON encodes JSONArrayOfArray object into bytes. If object is nil it provides empty array
func (aa JSONArrayOfArrayStr) MarshalJSON() ([]byte, error) {
	tmp := make([][]string, 0)
	if aa != nil {
		tmp = aa
	}
	return json.Marshal(&tmp)
}

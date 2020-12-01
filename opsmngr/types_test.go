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

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONArrayStr_Empty(t *testing.T) {
	type Test struct {
		A JSONArrayStr
	}
	originStr := `{}`
	expectedStr := `{"A":[]}`

	data := Test{}
	err := json.Unmarshal([]byte(originStr), &data)
	if err != nil {
		t.Fatalf("Unmarshal for '%s' returned error: %v", originStr, err)
	}

	raw, err := json.Marshal(&data)
	if err != nil {
		t.Fatalf("Marshal for '%v' returned error: %v", data, err)
	}

	assert.Equal(t, expectedStr, string(raw))
}

func TestJSONArrayStr_Correct(t *testing.T) {
	type Test struct {
		A JSONArrayStr
	}
	originStr := `{"A":["1","2"]}`

	data := Test{}
	err := json.Unmarshal([]byte(originStr), &data)
	if err != nil {
		t.Fatalf("Unmarshal for '%s' returned error: %v", originStr, err)
	}

	raw, err := json.Marshal(&data)
	if err != nil {
		t.Fatalf("Marshal for '%v' returned error: %v", data, err)
	}

	assert.Equal(t, originStr, string(raw))
}

func TestJSONArrayOfArrayStr_Empty(t *testing.T) {
	type Test struct {
		A JSONArrayOfArrayStr
	}
	originStr := `{}`
	expectedStr := `{"A":[]}`

	data := Test{}
	err := json.Unmarshal([]byte(originStr), &data)
	if err != nil {
		t.Fatalf("Unmarshal for '%s' returned error: %v", originStr, err)
	}

	raw, err := json.Marshal(&data)
	if err != nil {
		t.Fatalf("Marshal for '%v' returned error: %v", data, err)
	}

	assert.Equal(t, expectedStr, string(raw))
}

func TestJSONArrayOfArrayStr_Correct(t *testing.T) {
	type Test struct {
		A JSONArrayOfArrayStr
	}
	originStr := `{"A":[["0","1"],["2","3"]]}`

	data := Test{}
	err := json.Unmarshal([]byte(originStr), &data)
	if err != nil {
		t.Fatalf("Unmarshal for '%s' returned error: %v", originStr, err)
	}

	raw, err := json.Marshal(&data)
	if err != nil {
		t.Fatalf("Marshal for '%v' returned error: %v", data, err)
	}

	assert.Equal(t, originStr, string(raw))
}

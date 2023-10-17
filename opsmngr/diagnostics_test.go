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
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func TestDiagnostics_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	path := fmt.Sprintf("/api/public/v1.0/groups/%s/diagnostics", projectID)
	const expected = "test"
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, expected)
	})

	buf := new(bytes.Buffer)
	_, err := client.Diagnostics.Get(ctx, projectID, nil, buf)
	if err != nil {
		t.Fatalf("Diagnostics.Get returned error: %v", err)
	}

	if buf.String() != expected {
		t.Fatalf("Diagnostics.Get returned error: %v", err)
	}
}

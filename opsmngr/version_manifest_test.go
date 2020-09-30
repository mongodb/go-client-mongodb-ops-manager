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
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
)

func TestVersionManifest_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()
	const version = "4.4.json"

	mux.HandleFunc(fmt.Sprintf("/static/version_manifest/%s", version), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `

		{
		  "updated": 1599868800044,
		  "versions": [
			{
			  "builds": [
				{
				  "architecture": "amd64",
				  "gitVersion": "1c1c76aeca21c5983dc178920f5052c298db616c",
				  "platform": "linux",
				  "url": "/linux/mongodb-linux-x86_64-2.6.0.tgz"
				},
				{
				  "architecture": "amd64",
				  "gitVersion": "1c1c76aeca21c5983dc178920f5052c298db616c",
				  "platform": "osx",
				  "url": "/osx/mongodb-osx-x86_64-2.6.0.tgz"
				}
			  ],
			  "name": "2.6.0"
			}]
		}`)
	})

	newManifest, _, err := client.VersionManifest.Get(ctx, version)
	if err != nil {
		t.Fatalf("VersionManifest.Get returned error: %v", err)
	}

	expected := &VersionManifest{
		Updated: 1599868800044,
		Versions: []*Version{
			{
				Name: "2.6.0",
				Builds: []*Build{
					{
						Architecture: "amd64",
						GitVersion:   "1c1c76aeca21c5983dc178920f5052c298db616c",
						Platform:     "linux",
						URL:          "/linux/mongodb-linux-x86_64-2.6.0.tgz",
					},
					{
						Architecture: "amd64",
						GitVersion:   "1c1c76aeca21c5983dc178920f5052c298db616c",
						Platform:     "osx",
						URL:          "/osx/mongodb-osx-x86_64-2.6.0.tgz",
					},
				},
			},
		},
	}

	if diff := deep.Equal(newManifest, expected); diff != nil {
		t.Error(diff)
	}
}

func TestVersionManifest_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/versionManifest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		_, _ = fmt.Fprint(w, `

		{
		  "updated": 1599868800044,
		  "versions": [
			{
			  "builds": [
				{
				  "architecture": "amd64",
				  "gitVersion": "1c1c76aeca21c5983dc178920f5052c298db616c",
				  "platform": "linux",
				  "url": "/linux/mongodb-linux-x86_64-2.6.0.tgz"
				},
				{
				  "architecture": "amd64",
				  "gitVersion": "1c1c76aeca21c5983dc178920f5052c298db616c",
				  "platform": "osx",
				  "url": "/osx/mongodb-osx-x86_64-2.6.0.tgz"
				}
			  ],
			  "name": "2.6.0"
			}]
		}`)
	})

	newManifest := &VersionManifest{
		Updated: 1599868800044,
		Versions: []*Version{
			{
				Name: "2.6.0",
				Builds: []*Build{
					{
						Architecture: "amd64",
						GitVersion:   "1c1c76aeca21c5983dc178920f5052c298db616c",
						Platform:     "linux",
						URL:          "/linux/mongodb-linux-x86_64-2.6.0.tgz",
					},
					{
						Architecture: "amd64",
						GitVersion:   "1c1c76aeca21c5983dc178920f5052c298db616c",
						Platform:     "osx",
						URL:          "/osx/mongodb-osx-x86_64-2.6.0.tgz",
					},
				},
			},
		},
	}

	newManifest, _, err := client.VersionManifest.Update(ctx, newManifest)
	if err != nil {
		t.Fatalf("VersionManifest.Update returned error: %v", err)
	}

	expected := &VersionManifest{
		Updated: 1599868800044,
		Versions: []*Version{
			{
				Name: "2.6.0",
				Builds: []*Build{
					{
						Architecture: "amd64",
						GitVersion:   "1c1c76aeca21c5983dc178920f5052c298db616c",
						Platform:     "linux",
						URL:          "/linux/mongodb-linux-x86_64-2.6.0.tgz",
					},
					{
						Architecture: "amd64",
						GitVersion:   "1c1c76aeca21c5983dc178920f5052c298db616c",
						Platform:     "osx",
						URL:          "/osx/mongodb-osx-x86_64-2.6.0.tgz",
					},
				},
			},
		},
	}

	if diff := deep.Equal(newManifest, expected); diff != nil {
		t.Error(diff)
	}
}

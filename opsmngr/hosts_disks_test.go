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
	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestHostDisksService_Get(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/groups/12345678/hosts/1/disks/xvdb", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			 "links":[
				{
				   "href":"https://local/api/public/v1.0/groups/12345678/hosts/1/disks/xvdb",
				   "rel":"self"
				}
			 ],
			 "partitionName":"xvdb"
		}`)
	})

	disks, _, err := client.HostDisks.Get(ctx, "12345678", "1", "xvdb")
	if err != nil {
		t.Fatalf("HostDisks.Get returned error: %v", err)
	}

	expected := &atlas.ProcessDisk{
		Links: []*atlas.Link{
			{
				Rel:  "self",
				Href: "https://local/api/public/v1.0/groups/12345678/hosts/1/disks/xvdb",
			},
		},
		PartitionName: "xvdb",
	}

	if diff := deep.Equal(disks, expected); diff != nil {
		t.Error(diff)
	}
}

func TestHostDisksService_List(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/groups/12345678/hosts/1/disks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
		   "links":[
			  {
				 "href":"https://local/api/public/v1.0/groups/12345678/hosts/1/disks?pageNum=1&itemsPerPage=100",
				 "rel":"self"
			  }
		   ],
		   "results":[
			  {
				 "links":[
					{
					   "href":"https://local/api/public/v1.0/groups/12345678/hosts/1/disks/xvdb",
					   "rel":"self"
					}
				 ],
				 "partitionName":"xvdb"
			  }
		   ],
		   "totalCount":1
		}`)
	})

	disks, _, err := client.HostDisks.List(ctx, "12345678", "1", nil)
	if err != nil {
		t.Fatalf("HostDisks.List returned error: %v", err)
	}

	expected := &atlas.ProcessDisksResponse{
		Links: []*atlas.Link{
			{
				Rel:  "self",
				Href: "https://local/api/public/v1.0/groups/12345678/hosts/1/disks?pageNum=1&itemsPerPage=100",
			},
		},
		Results: []*atlas.ProcessDisk{
			{
				Links: []*atlas.Link{
					{
						Rel:  "self",
						Href: "https://local/api/public/v1.0/groups/12345678/hosts/1/disks/xvdb",
					},
				},
				PartitionName: "xvdb",
			},
		},
		TotalCount: 1,
	}

	if diff := deep.Equal(disks, expected); diff != nil {
		t.Error(diff)
	}
}

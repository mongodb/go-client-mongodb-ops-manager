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

func TestDeployments_GetPartition(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()
	mux.HandleFunc("/api/public/v1.0/groups/12345678/hosts/1/disks/xvdb", func(w http.ResponseWriter, r *http.Request) {
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

	disks, _, err := client.Deployments.GetPartition(ctx, "12345678", "1", "xvdb")
	if err != nil {
		t.Fatalf("Deployments.GetPartition returned error: %v", err)
	}

	expected := &ProcessDisk{
		Links: []*Link{
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

func TestDeployments_ListPartitions(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()
	mux.HandleFunc("/api/public/v1.0/groups/12345678/hosts/1/disks", func(w http.ResponseWriter, r *http.Request) {
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

	disks, _, err := client.Deployments.ListPartitions(ctx, "12345678", "1", nil)
	if err != nil {
		t.Fatalf("Deployments.ListPartitions returned error: %v", err)
	}

	expected := &ProcessDisksResponse{
		Links: []*Link{
			{
				Rel:  "self",
				Href: "https://local/api/public/v1.0/groups/12345678/hosts/1/disks?pageNum=1&itemsPerPage=100",
			},
		},
		Results: []*ProcessDisk{
			{
				Links: []*Link{
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

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
	"context"
	"fmt"
	"net/http"

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	versionManifestBasePath   = "versionManifest"
	versionManifestStaticPath = "static/version_manifest/%s"
)

// VersionManifestService is an interface for using the Version Manifest
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/version-manifest/
type VersionManifestService interface {
	Get(context.Context, string) (*VersionManifest, *atlas.Response, error)
	Update(context.Context, *VersionManifest) (*VersionManifest, *atlas.Response, error)
}

// VersionManifestServiceOp provides an implementation of the VersionManifestService interface
type VersionManifestServiceOp service

var _ VersionManifestService = &VersionManifestServiceOp{}

type VersionManifest struct {
	Updated  int64      `json:"updated,omitempty"`
	Versions []*Version `json:"versions,omitempty"`
}

type Version struct {
	Name   string   `json:"name,omitempty"`
	Builds []*Build `json:"builds,omitempty"`
}

type Build struct {
	Architecture       string    `json:"architecture,omitempty"`
	GitVersion         string    `json:"gitVersion,omitempty"`
	Platform           string    `json:"platform,omitempty"`
	URL                string    `json:"url,omitempty"`
	MaxOSVersion       string    `json:"maxOsVersion,omitempty"`
	MinOSVersion       string    `json:"minOsVersion,omitempty"`
	Win2008plus        bool      `json:"win2008plus,omitempty"`
	WinVCRedistDLL     string    `json:"winVCRedistDll,omitempty"`     //nolint:tagliatelle // correct from API
	WinVCRedistOptions []*string `json:"winVCRedistOptions,omitempty"` //nolint:tagliatelle // correct from API
	WinVCRedistURL     string    `json:"winVCRedistUrl,omitempty"`     //nolint:tagliatelle // correct from API
	WinVCRedistVersion string    `json:"winVCRedistVersion,omitempty"` //nolint:tagliatelle // correct from API
	Flavor             string    `json:"flavor,omitempty"`
}

// Get retrieves the version manifest
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/version-manifest/view-version-manifest/
func (s *VersionManifestServiceOp) Get(ctx context.Context, version string) (*VersionManifest, *atlas.Response, error) {
	if version == "" {
		return nil, nil, atlas.NewArgError("version", "must be set")
	}

	path := fmt.Sprintf(versionManifestStaticPath, version)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(VersionManifest)
	resp, err := s.Client.Do(ctx, req, root)
	return root, resp, err
}

// Update updates the version manifest
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/version-manifest/update-version-manifest/
func (s *VersionManifestServiceOp) Update(ctx context.Context, versionManifest *VersionManifest) (*VersionManifest, *atlas.Response, error) {
	if versionManifest == nil {
		return nil, nil, atlas.NewArgError("versionManifest", "must be set")
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPut, versionManifestBasePath, versionManifest)
	if err != nil {
		return nil, nil, err
	}

	root := new(VersionManifest)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

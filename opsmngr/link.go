// Copyright 2021 MongoDB Inc
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

import "net/url"

// Link is the link to sub-resources and/or related resources.
type Link struct {
	Rel  string `json:"rel,omitempty"`
	Href string `json:"href,omitempty"`
}

//nolint:all // this is a temp method and will be called later
func (l *Link) getHrefURL() (*url.URL, error) {
	return url.Parse(l.Href)
}

//nolint:all // this is a temp method and will be called later
func (l *Link) getHrefQueryParam(param string) (string, error) {
	hrefURL, err := l.getHrefURL()
	if err != nil {
		return "", err
	}
	return hrefURL.Query().Get(param), nil
}
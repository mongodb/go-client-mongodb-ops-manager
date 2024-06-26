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

package auth

import (
	"net/http"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`  //nolint:tagliatelle // used as in the API
	RefreshToken string `json:"refresh_token"` //nolint:tagliatelle // used as in the API
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`   //nolint:tagliatelle // used as in the API
	TokenType    string `json:"token_type"` //nolint:tagliatelle // used as in the API
	ExpiresIn    int    `json:"expires_in"` //nolint:tagliatelle // used as in the API
	Expiry       time.Time
}

func (t *Token) SetAuthHeader(r *http.Request) {
	r.Header.Set("Authorization", "Bearer "+t.AccessToken)
}

const expiryDelta = 10 * time.Second

func (t *Token) expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Round(0).Add(-expiryDelta).Before(time.Now())
}

func (t *Token) Valid() bool {
	return t != nil && t.AccessToken != "" && !t.expired()
}

// Copyright © 2018 ocicert authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dist

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/kinvolk/ocicert/pkg/auth"
	"github.com/kinvolk/ocicert/pkg/image"
)

var (
	homeDir    string
	regAuthCtx auth.RegAuthContext
)

func init() {
	homeDir = os.Getenv("HOME")
}

func TestCheckAPIVersion(t *testing.T) {
	reqPath := ""

	indexServer := image.GetIndexName(image.DefaultIndexURLAuth)

	regAuthCtx = auth.NewRegAuthContext()
	regAuthCtx.Scope.RemoteName = reqPath
	regAuthCtx.Scope.Actions = "pull"

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	_, res, err := regAuthCtx.SendRequestWithToken(inputURL, "GET")
	if err != nil {
		t.Fatalf("failed to send request with token to %s: %v", inputURL, err)
	}

	switch res.StatusCode {
	case http.StatusCreated:
		t.Fatalf("got an unexpected reply: 201 Created")
	case http.StatusUnauthorized:
		t.Fatalf("got an unexpected reply: 401 Unauthorized")
	case http.StatusNotFound:
		t.Fatalf("got an unexpected reply: 404 Not Found")
	case http.StatusOK:
		break
	default:
		t.Fatalf("statusCode = %v, request URL = %v", res.StatusCode, inputURL)
	}
}

func TestPullManifest(t *testing.T) {
	imageName := "busybox"
	refName := "latest"

	indexServer := image.GetIndexName(image.DefaultIndexURLAuth)

	remoteName := filepath.Join(image.DefaultRepoPrefix, imageName)
	reqPath := filepath.Join(remoteName, "manifests", refName)

	regAuthCtx = auth.NewRegAuthContext()
	regAuthCtx.Scope.RemoteName = remoteName
	regAuthCtx.Scope.Actions = "pull"

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	_, res, err := regAuthCtx.SendRequestWithToken(inputURL, "GET")
	if err != nil {
		t.Fatalf("failed to send request with token to %s: %v", inputURL, err)
	}

	switch res.StatusCode {
	case http.StatusOK:
		return
	case http.StatusCreated:
		fallthrough
	case http.StatusUnauthorized:
		fallthrough
	case http.StatusNotFound:
		t.Fatalf("got an unexpected reply: %v", res.StatusCode)
	default:
		t.Fatalf("statusCode = %v, request URL = %v", res.StatusCode, inputURL)
	}
}

func TestPushManifest(t *testing.T) {
	imageName := "busybox"
	refName := "latest"

	indexServer := image.GetIndexName(image.DefaultIndexURLAuth)

	remoteName := filepath.Join(image.DefaultRepoPrefix, imageName)
	reqPath := filepath.Join(remoteName, "manifests", refName)

	regAuthCtx = auth.NewRegAuthContext()
	regAuthCtx.Scope.RemoteName = remoteName
	regAuthCtx.Scope.Actions = "push"

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	_, res, err := regAuthCtx.SendRequestWithToken(inputURL, "PUT")
	if err != nil {
		t.Fatalf("failed to send request with token to %s: %v", inputURL, err)
	}

	switch res.StatusCode {
	case http.StatusOK:
		return
	case http.StatusCreated:
		fallthrough
	case http.StatusUnauthorized:
		fallthrough
	case http.StatusNotFound:
		t.Fatalf("got an unexpected reply: %v", res.StatusCode)
	default:
		t.Fatalf("statusCode = %v, request URL = %v", res.StatusCode, inputURL)
	}
}

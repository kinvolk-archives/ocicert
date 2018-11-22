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
	distp "github.com/kinvolk/ocicert/pkg/distp"
	"github.com/kinvolk/ocicert/pkg/image"
)

var (
	homeDir    string
	regAuthCtx auth.RegAuthContext

	testImageName string = "busybox"
	testRefName   string = "latest"
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

	_, res, err := regAuthCtx.SendRequestWithToken(inputURL, "GET", nil)
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

	if vers := res.Header.Get(distp.DistAPIVersionKey); vers != distp.DistAPIVersionValue {
		t.Fatalf("got an unexpected API version %v", vers)
	}
}

func TestPullManifest(t *testing.T) {
	indexServer := image.GetIndexName(image.DefaultIndexURLAuth)

	remoteName := filepath.Join(image.DefaultRepoPrefix, testImageName)
	reqPath := filepath.Join(remoteName, "manifests", testRefName)

	regAuthCtx = auth.NewRegAuthContext()
	regAuthCtx.Scope.RemoteName = remoteName
	regAuthCtx.Scope.Actions = "pull"

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	_, res, err := regAuthCtx.SendRequestWithToken(inputURL, "GET", nil)
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
	indexServer := image.GetIndexName(image.DefaultIndexURLAuth)

	remoteName := filepath.Join(image.DefaultRepoPrefix, testImageName)
	reqPath := filepath.Join(remoteName, "manifests", testRefName)

	regAuthCtx = auth.NewRegAuthContext()
	regAuthCtx.Scope.RemoteName = remoteName
	regAuthCtx.Scope.Actions = "push"

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	_, res, err := regAuthCtx.SendRequestWithToken(inputURL, "PUT", nil)
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

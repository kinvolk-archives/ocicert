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
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	regURL        string
)

func init() {
	homeDir = os.Getenv("HOME")

	regURL = regAuthCtx.RegURL
}

func TestCheckAPIVersion(t *testing.T) {
	reqPath := ""

	regAuthCtx := auth.NewRegAuthContext()
	regAuthCtx.Scope.RemoteName = reqPath
	regAuthCtx.Scope.Actions = "pull"

	indexServer := auth.GetIndexServer(regURL)

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	res, err := regAuthCtx.GetResponse(inputURL, "GET", nil, []int{http.StatusOK})
	if err != nil {
		t.Fatalf("got an unexpected reply: %v", err)
	}

	if vers := res.Header.Get(distp.DistAPIVersionKey); vers != distp.DistAPIVersionValue {
		t.Fatalf("got an unexpected API version %v", vers)
	}
}

func getDigestFromManifest(regURL, testImageName, testRefName string) (string, error) {
	indexServer := auth.GetIndexServer(regURL)

	remoteName := filepath.Join(auth.DefaultRepoPrefix, testImageName)
	reqPath := filepath.Join(remoteName, "manifests", testRefName)

	regAuthCtx := auth.NewRegAuthContext()
	regAuthCtx.Scope.RemoteName = remoteName
	regAuthCtx.Scope.Actions = "pull"

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		return "", fmt.Errorf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	res, err := regAuthCtx.GetResponse(inputURL, "HEAD", nil, []int{http.StatusOK})
	if err != nil {
		return "", fmt.Errorf("got an unexpected reply: %v", err)
	}

	return res.Header.Get(distp.ContentDigest), nil
}

func testUploadLayer(t *testing.T) {
	indexServer := auth.GetIndexServer(regURL)

	remoteName := filepath.Join(auth.DefaultRepoPrefix, testImageName)
	reqPath := filepath.Join(remoteName, "blobs/uploads", testRefName)

	regAuthCtx = auth.NewRegAuthContext()
	regAuthCtx.Scope.RemoteName = remoteName
	regAuthCtx.Scope.Actions = "push"

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	// 1. POST
	// Send a POST request without any body specified.
	postURL := "https://" + indexServer + "/v2/" + reqPath

	res, err := regAuthCtx.GetResponse(postURL, "POST", nil, []int{http.StatusOK, http.StatusAccepted})
	if err != nil {
		t.Fatalf("got an unexpected reply: %v", err)
	}

	uuid := res.Header.Get(distp.UploadUuidKey)

	// 2. PATCH
	// Generate a 100-byte blob of a randomly generated string.
	// Send a PATCH request with the blob.
	blob := image.GenRandomBlob(100)

	if _, err := regAuthCtx.GetResponse(postURL, "PATCH", strings.NewReader(blob),
		[]int{http.StatusOK, http.StatusAccepted}); err != nil {
		t.Fatalf("got an unexpected reply: %v", err)
	}

	// 3. PUT
	// Generate a blob's digest, generated as a sha256 checksum of the blob.
	// Send a PUT request with a "digest=..." option appended to its URL.
	digest := image.GetHash(blob)
	putURL := "https://" + indexServer + "/v2/" + reqPath + "/" + uuid + "?digest=" + digest

	if _, err := regAuthCtx.GetResponse(putURL, "PUT", strings.NewReader(blob),
		[]int{http.StatusCreated}); err != nil {
		t.Fatalf("got an unexpected reply: %v", err)
	}
}

func testPushManifest(t *testing.T) {
	indexServer := auth.GetIndexServer(regURL)

	remoteName := filepath.Join(auth.DefaultRepoPrefix, testImageName)
	reqPath := filepath.Join(remoteName, "manifests", testRefName)

	regAuthCtx = auth.NewRegAuthContext()
	regAuthCtx.Scope.RemoteName = remoteName
	regAuthCtx.Scope.Actions = "push"

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	if _, err := regAuthCtx.GetResponse(inputURL, "PUT", nil, []int{http.StatusCreated}); err != nil {
		t.Fatalf("got an unexpected reply: %v", err)
	}
}

func testPullManifest(t *testing.T) {
	indexServer := auth.GetIndexServer(regURL)

	remoteName := filepath.Join(auth.DefaultRepoPrefix, testImageName)
	reqPath := filepath.Join(remoteName, "manifests", testRefName)

	regAuthCtx := auth.NewRegAuthContext()
	regAuthCtx.Scope.RemoteName = remoteName
	regAuthCtx.Scope.Actions = "pull"

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	if _, err := regAuthCtx.GetResponse(inputURL, "GET", nil, []int{http.StatusOK}); err != nil {
		t.Fatalf("got an unexpected reply: %v", err)
	}
}

func testPullLayer(t *testing.T) {
	indexServer := auth.GetIndexServer(regURL)

	remoteName := filepath.Join(auth.DefaultRepoPrefix, testImageName)

	testDigest, err := getDigestFromManifest(regURL, testImageName, testRefName)
	if err != nil {
		t.Fatalf("failed to get digest from %s: %v", indexServer, err)
	}

	reqPath := filepath.Join(remoteName, "blobs", testDigest)

	regAuthCtx := auth.NewRegAuthContext()
	regAuthCtx.Scope.RemoteName = remoteName
	regAuthCtx.Scope.Actions = "pull"

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	if _, err := regAuthCtx.GetResponse(inputURL, "GET", nil, []int{http.StatusOK}); err != nil {
		t.Fatalf("got an unexpected reply: %v", err)
	}
}

func testDeleteLayer(t *testing.T) {
	indexServer := auth.GetIndexServer(regURL)

	remoteName := filepath.Join(auth.DefaultRepoPrefix, testImageName)

	testDigest, err := getDigestFromManifest(regURL, testImageName, testRefName)
	if err != nil {
		t.Fatalf("failed to get digest from %s: %v", indexServer, err)
	}

	reqPath := filepath.Join(remoteName, "blobs", testDigest)

	regAuthCtx = auth.NewRegAuthContext()
	regAuthCtx.Scope.RemoteName = remoteName
	regAuthCtx.Scope.Actions = "push"

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	if _, err := regAuthCtx.GetResponse(inputURL, "DELETE", nil, []int{http.StatusAccepted}); err != nil {
		t.Fatalf("got an unexpected reply: %v", err)
	}
}

func testDeleteManifest(t *testing.T) {
	indexServer := auth.GetIndexServer(regURL)

	remoteName := filepath.Join(auth.DefaultRepoPrefix, testImageName)
	reqPath := filepath.Join(remoteName, "manifests", testRefName)

	regAuthCtx = auth.NewRegAuthContext()
	regAuthCtx.Scope.RemoteName = remoteName
	regAuthCtx.Scope.Actions = "push"

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	if _, err := regAuthCtx.GetResponse(inputURL, "DELETE", nil, []int{http.StatusAccepted}); err != nil {
		t.Fatalf("got an unexpected reply: %v", err)
	}
}

func TestPushPullLayer(t *testing.T) {
	testUploadLayer(t)
	testPushManifest(t)
	testPullManifest(t)
	testPullLayer(t)
	testDeleteLayer(t)
	testDeleteManifest(t)
}

func TestListRepos(t *testing.T) {
	reqPath := "catalog"

	regAuthCtx.Scope.RemoteName = reqPath
	regAuthCtx.Scope.Actions = "pull"

	indexServer := auth.GetIndexServer(regURL)

	// NOTE: it will fail when testing against docker.io, as '/v2/_catalog' endpoint
	// will not be supported.
	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	_, err := regAuthCtx.GetResponse(inputURL, "GET", nil, []int{http.StatusOK})
	if err != nil {
		t.Fatalf("got an unexpected reply: %v", err)
	}
}

func TestListTags(t *testing.T) {
	indexServer := auth.GetIndexServer(regURL)

	regAuthCtx := auth.NewRegAuthContext()
	remoteName := filepath.Join(auth.DefaultRepoPrefix, testImageName)
	reqPath := filepath.Join(remoteName, "tags/list")

	regAuthCtx.Scope.RemoteName = remoteName
	regAuthCtx.Scope.Actions = "pull"

	if err := regAuthCtx.PrepareAuth(indexServer); err != nil {
		t.Fatalf("failed to prepare auth to %s for %s: %v", indexServer, reqPath, err)
	}

	inputURL := "https://" + indexServer + "/v2/" + reqPath

	if _, err := regAuthCtx.GetResponse(inputURL, "GET", nil, []int{http.StatusOK}); err != nil {
		t.Fatalf("got an unexpected reply: %v", err)
	}
}

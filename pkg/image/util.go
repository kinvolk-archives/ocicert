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

package image

import (
	"strings"
)

const (
	DefaultIndexURLPlain = "registry-1.docker.io"
	DefaultIndexURLAuth  = "index.docker.io"
	DefaultRepoPrefix    = "library/"
)

// GetIndexName returns the index server from a registry URL.
func GetIndexName(regURL string) string {
	index, _ := SplitReposName(regURL)
	return index
}

// SplitReposName breaks a repo name into an index name and remote name.
func SplitReposName(name string) (indexName, remoteName string) {
	i := strings.IndexRune(name, '/')
	if i == -1 || (!strings.ContainsAny(name[:i], ".:") && name[:i] != "localhost") {
		indexName, remoteName = DefaultIndexURLPlain, name
	} else {
		indexName, remoteName = name[:i], name[i+1:]
	}
	if indexName == DefaultIndexURLPlain && !strings.ContainsRune(remoteName, '/') {
		remoteName = DefaultRepoPrefix + remoteName
	}
	return
}

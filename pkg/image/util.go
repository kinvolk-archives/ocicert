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
	"crypto/sha256"
	"math/rand"
	"strconv"
)

func GenRandomBlob(blobLen int) string {
	blob := ""
	for i := 0; i < blobLen; i++ {
		blob += strconv.Itoa(rand.Intn(9))
	}
	return blob
}

func GetHash(inputStr string) string {
	h := sha256.New()
	h.Write([]byte(inputStr))
	return string(h.Sum(nil))
}

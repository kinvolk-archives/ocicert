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
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type testRegistry struct {
	cmd *exec.Cmd
	url string
	dir string
}

var (
	privateRegURL string = "127.0.0.1:5000"

	testReg *testRegistry
)

func init() {
	localreg := os.Getenv("OCICERT_LOCALREG")

	if localreg != "" {
		regEnv := os.Getenv("OCICERT_REGISTRY")
		if regEnv != "" {
			privateRegURL = regEnv
		}

		testReg, _ = setupRegistry(privateRegURL)
	}
}

func newTestRegistry(url string) (*testRegistry, error) {
	tmpDir, err := ioutil.TempDir("", "ocicert-test-")
	if err != nil {
		return nil, err
	}

	template := `version: 0.1
loglevel: debug
storage:
    filesystem:
        rootdirectory: %s
    delete:
        enabled: true
http:
    addr: %s
`

	confPath := filepath.Join(tmpDir, "config.yaml")

	config, err := os.Create(confPath)
	if err != nil {
		return nil, err
	}
	defer config.Close()

	if _, err := fmt.Fprintf(config, template, tmpDir, url); err != nil {
		return nil, err
	}

	regBinary := ""
	if regBinary, err = exec.LookPath("registry"); err != nil {
		return nil, fmt.Errorf("cannot find binary registry")
	}

	// NOTE: it works only for registry from docker/distribution.
	cmd := exec.Command(regBinary, "serve", confPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &testRegistry{
		cmd: cmd,
		url: url,
		dir: tmpDir,
	}, nil
}

func setupRegistry(url string) (*testRegistry, error) {
	reg, err := newTestRegistry(url)
	if err != nil {
		return nil, err
	}

	var errPing error

	timeout := 5 * time.Second
	alarm := time.After(timeout)
	ticker := time.Tick(200 * time.Millisecond)
	for {
		select {
		case <-alarm:
			return nil, fmt.Errorf("timeout waiting on server being available within %v", timeout)
		case <-ticker:
			if errPing = reg.Ping(); err == nil {
				return reg, nil
			}
		}
	}

	return nil, errPing
}

func (t *testRegistry) Ping() error {
	resp, err := http.Get(fmt.Sprintf("http://%s/v2/", t.url))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("ping returned with an unexpected status code %d", resp.StatusCode)
	}
	return nil
}

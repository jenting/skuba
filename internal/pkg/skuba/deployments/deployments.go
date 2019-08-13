/*
 * Copyright (c) 2019 SUSE LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package deployments

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/klog"
)

var (
	stateMap = map[string]Runner{}
)

type Runner func(t *Target, data interface{}) error

type Actionable interface {
	Do(command string, args ...string) (string, string, error)
	SilentDo(command string, args ...string) (string, string, error)
	UploadFile(sourcePath, targetPath string) error
	UploadFileContents(targetPath, contents string) error
	DownloadFileContents(sourcePath string) (string, error)
}

type TargetCache struct {
	OsRelease map[string]string
}

type Target struct {
	Actionable
	Hostname string
	Nodename string
	Cache    TargetCache
}

func (t *Target) Apply(data interface{}, states ...string) error {
	for _, stateName := range states {
		if stateName == "" {
			continue
		}

		klog.V(1).Infof("=== applying state %s ===", stateName)
		if state, stateExists := stateMap[stateName]; stateExists {
			if err := state(t, data); err != nil {
				return errors.Wrapf(err, "failed to apply state %s", stateName)
			}
			klog.V(1).Infof("=== state %s applied successfully ===", stateName)
		} else {
			return fmt.Errorf("state does not exist: %s", stateName)
		}
	}

	return nil
}

func (t *Target) UploadFile(sourcePath, targetPath string) error {
	return t.Actionable.UploadFile(sourcePath, targetPath)
}

func (t *Target) UploadFileContents(targetPath, contents string) error {
	return t.Actionable.UploadFileContents(targetPath, contents)
}

func (t *Target) DownloadFileContents(sourcePath string) (string, error) {
	return t.Actionable.DownloadFileContents(sourcePath)
}

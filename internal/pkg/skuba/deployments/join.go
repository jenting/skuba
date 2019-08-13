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
	"strings"

	"k8s.io/klog"

	"github.com/SUSE/skuba/pkg/skuba"
)

func MustGetRoleFromString(s string) (role skuba.Role) {
	switch strings.ToLower(s) {
	case "master":
		role = skuba.MasterRole
	case "worker":
		role = skuba.WorkerRole
	default:
		klog.Fatalf("[join] invalid role provided: %q, 'master' or 'worker' are the only accepted roles", s)
	}
	return
}

type JoinConfiguration struct {
	Role             skuba.Role
	KubeadmExtraArgs map[string]string
}

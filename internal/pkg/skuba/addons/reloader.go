/*
 * Copyright (c) 2020 SUSE LLC.
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

package addons

import (
	"k8s.io/kubernetes/cmd/kubeadm/app/images"

	"github.com/SUSE/skuba/internal/pkg/skuba/kubernetes"
	"github.com/SUSE/skuba/internal/pkg/skuba/skuba"
	skubaconstants "github.com/SUSE/skuba/pkg/skuba"
)

func init() {
	registerAddon(kubernetes.Reloader, renderReloaderTemplate, nil, reloaderCallbacks{}, normalPriority, []getImageCallback{GetReloaderImage})
}

func GetReloaderImage(imageTag string) string {
	return images.GetGenericImage(skubaconstants.ImageRepository, "reloader", imageTag)
}

func (renderContext renderContext) ReloaderImage() string {
	return GetReloaderImage(kubernetes.AddonVersionForClusterVersion(kubernetes.Reloader, renderContext.config.ClusterVersion).Version)
}

func renderReloaderTemplate(addonConfiguration AddonConfiguration) string {
	return reloaderManifest
}

type reloaderCallbacks struct{}

func (reloaderCallbacks) beforeApply(addonConfiguration AddonConfiguration, skubaConfiguration *skuba.SkubaConfiguration) error {
	return nil
}

func (reloaderCallbacks) afterApply(addonConfiguration AddonConfiguration, skubaConfiguration *skuba.SkubaConfiguration) error {
	return nil
}

const (
	// copy from https://raw.githubusercontent.com/stakater/Reloader/v0.0.58/deployments/kubernetes/reloader.yaml
	reloaderManifest = `---
---
# Source: reloader/templates/clusterrole.yaml

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  labels:
    app: reloader-reloader
    chart: "reloader-v0.0.58"
    release: "reloader"
    heritage: "Tiller"
  name: reloader-reloader-role
rules:
  - apiGroups:
      - ""
    resources:
      - secrets
      - configmaps
    verbs:
      - list
      - get
      - watch
  - apiGroups:
      - "apps"
    resources:
      - deployments
      - daemonsets
      - statefulsets
    verbs:
      - list
      - get
      - update
      - patch
  - apiGroups:
      - "extensions"
    resources:
      - deployments
      - daemonsets
    verbs:
      - list
      - get
      - update
      - patch

---
# Source: reloader/templates/clusterrolebinding.yaml

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  labels:
    app: reloader-reloader
    chart: "reloader-v0.0.58"
    release: "reloader"
    heritage: "Tiller"
  name: reloader-reloader-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: reloader-reloader-role
subjects:
  - kind: ServiceAccount
    name: reloader-reloader
    namespace: kube-system

---
# Source: reloader/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: reloader-reloader
    chart: "reloader-v0.0.58"
    release: "reloader"
    heritage: "Tiller"
    group: com.stakater.platform
    provider: stakater
    version: v0.0.58
  name: reloader-reloader
  namespace: kube-system
spec:
  replicas: 1
  revisionHistoryLimit: 3
  selector:
    matchLabels:
      app: reloader-reloader
      release: "reloader"
  template:
    metadata:
      labels:
        app: reloader-reloader
        chart: "reloader-v0.0.58"
        release: "reloader"
        heritage: "Tiller"
        group: com.stakater.platform
        provider: stakater
        version: v0.0.58
      annotations:
        {{.AnnotatedVersion}}		
    spec:
      containers:
      - env:
        image: {{.ReloaderImage}}
        imagePullPolicy: IfNotPresent
        name: reloader-reloader
        args:
      serviceAccountName: reloader-reloader

---
# Source: reloader/templates/role.yaml


---
# Source: reloader/templates/rolebinding.yaml


---
# Source: reloader/templates/service.yaml

---
# Source: reloader/templates/serviceaccount.yaml

apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app: reloader-reloader
    chart: "reloader-v0.0.58"
    release: "reloader"
    heritage: "Tiller"
  name: reloader-reloader
  namespace: kube-system
`
)

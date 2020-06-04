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
	registerAddon(kubernetes.CertManager, renderCertManagerTemplate, nil, certManagerCallbacks{}, highPriority, []getImageCallback{GetCertManagerCAInjectorImage, GetCertManagerControllerImage, GetCertManagerWebhookImage})
}

func GetCertManagerCAInjectorImage(imageTag string) string {
	return images.GetGenericImage(skubaconstants.ImageRepository, "cert-manager-cainjector", imageTag)
}

func GetCertManagerControllerImage(imageTag string) string {
	return images.GetGenericImage(skubaconstants.ImageRepository, "cert-manager-controller", imageTag)
}

func GetCertManagerWebhookImage(imageTag string) string {
	return images.GetGenericImage(skubaconstants.ImageRepository, "cert-manager-webhook", imageTag)
}

func (renderContext renderContext) CertManagerCAInjectorImage() string {
	return GetCertManagerCAInjectorImage(kubernetes.AddonVersionForClusterVersion(kubernetes.CertManager, renderContext.config.ClusterVersion).Version)
}

func (renderContext renderContext) CertManagerControllerImage() string {
	return GetCertManagerControllerImage(kubernetes.AddonVersionForClusterVersion(kubernetes.CertManager, renderContext.config.ClusterVersion).Version)
}

func (renderContext renderContext) CertManagerWebhookImage() string {
	return GetCertManagerWebhookImage(kubernetes.AddonVersionForClusterVersion(kubernetes.CertManager, renderContext.config.ClusterVersion).Version)
}

func renderCertManagerTemplate(addonConfiguration AddonConfiguration) string {
	return certManagerManifest
}

type certManagerCallbacks struct{}

func (certManagerCallbacks) beforeApply(addonConfiguration AddonConfiguration, skubaConfiguration *skuba.SkubaConfiguration) error {
	return nil
}

func (certManagerCallbacks) afterApply(addonConfiguration AddonConfiguration, skubaConfiguration *skuba.SkubaConfiguration) error {
	return nil
}

const (
	// generated from helm repo https://charts.jetstack.io name jetstack/cert-manager with chart version v0.15.1
	// helm template \
	//   cert-manager \
	//   --name cert-manager \
	//   --namespace kube-system \
	//   --set installCRDs=true
	certManagerManifest = `---
# Source: cert-manager/templates/cainjector-serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cert-manager-cainjector
  namespace: kube-system
  labels:
    app: cainjector
    app.kubernetes.io/name: cainjector
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "cainjector"
    helm.sh/chart: cert-manager-v0.15.1
---
# Source: cert-manager/templates/serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cert-manager
  namespace: kube-system
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1

---
# Source: cert-manager/templates/webhook-serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cert-manager-webhook
  namespace: kube-system
  labels:
    app: webhook
    app.kubernetes.io/name: webhook
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "webhook"
    helm.sh/chart: cert-manager-v0.15.1
---
# Source: cert-manager/templates/crds.legacy.yaml

apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    cert-manager.io/inject-ca-from-secret: 'kube-system/cert-manager-webhook-ca'
  labels:
    app: 'cert-manager'
    app.kubernetes.io/instance: 'cert-manager'
    app.kubernetes.io/managed-by: 'Tiller'
    app.kubernetes.io/name: 'cert-manager'
    helm.sh/chart: 'cert-manager-v0.15.1'
  name: certificaterequests.cert-manager.io
spec:
  additionalPrinterColumns:
  - JSONPath: .status.conditions[?(@.type=="Ready")].status
    name: Ready
    type: string
  - JSONPath: .spec.issuerRef.name
    name: Issuer
    priority: 1
    type: string
  - JSONPath: .status.conditions[?(@.type=="Ready")].message
    name: Status
    priority: 1
    type: string
  - JSONPath: .metadata.creationTimestamp
    description: CreationTimestamp is a timestamp representing the server time when
      this object was created. It is not guaranteed to be set in happens-before order
      across separate operations. Clients may not set this value. It is represented
      in RFC3339 form and is in UTC.
    name: Age
    type: date
  group: cert-manager.io
  names:
    kind: CertificateRequest
    listKind: CertificateRequestList
    plural: certificaterequests
    shortNames:
    - cr
    - crs
    singular: certificaterequest
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: CertificateRequest is a type to represent a Certificate Signing
        Request
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: CertificateRequestSpec defines the desired state of CertificateRequest
          properties:
            csr:
              description: Byte slice containing the PEM encoded CertificateSigningRequest
              format: byte
              type: string
            duration:
              description: Requested certificate default Duration
              type: string
            isCA:
              description: IsCA will mark the resulting certificate as valid for signing.
                This implies that the 'cert sign' usage is set
              type: boolean
            issuerRef:
              description: IssuerRef is a reference to the issuer for this CertificateRequest.  If
                the 'kind' field is not set, or set to 'Issuer', an Issuer resource
                with the given name in the same namespace as the CertificateRequest
                will be used.  If the 'kind' field is set to 'ClusterIssuer', a ClusterIssuer
                with the provided name will be used. The 'name' field in this stanza
                is required at all times. The group field refers to the API group
                of the issuer which defaults to 'cert-manager.io' if empty.
              properties:
                group:
                  type: string
                kind:
                  type: string
                name:
                  type: string
              required:
              - name
              type: object
            usages:
              description: Usages is the set of x509 actions that are enabled for
                a given key. Defaults are ('digital signature', 'key encipherment')
                if empty
              items:
                description: 'KeyUsage specifies valid usage contexts for keys. See:
                  https://tools.ietf.org/html/rfc5280#section-4.2.1.3      https://tools.ietf.org/html/rfc5280#section-4.2.1.12
                  Valid KeyUsage values are as follows: "signing", "digital signature",
                  "content commitment", "key encipherment", "key agreement", "data
                  encipherment", "cert sign", "crl sign", "encipher only", "decipher
                  only", "any", "server auth", "client auth", "code signing", "email
                  protection", "s/mime", "ipsec end system", "ipsec tunnel", "ipsec
                  user", "timestamping", "ocsp signing", "microsoft sgc", "netscape
                  sgc"'
                enum:
                - signing
                - digital signature
                - content commitment
                - key encipherment
                - key agreement
                - data encipherment
                - cert sign
                - crl sign
                - encipher only
                - decipher only
                - any
                - server auth
                - client auth
                - code signing
                - email protection
                - s/mime
                - ipsec end system
                - ipsec tunnel
                - ipsec user
                - timestamping
                - ocsp signing
                - microsoft sgc
                - netscape sgc
                type: string
              type: array
          required:
          - csr
          - issuerRef
          type: object
        status:
          description: CertificateStatus defines the observed state of CertificateRequest
            and resulting signed certificate.
          properties:
            ca:
              description: Byte slice containing the PEM encoded certificate authority
                of the signed certificate.
              format: byte
              type: string
            certificate:
              description: Byte slice containing a PEM encoded signed certificate
                resulting from the given certificate signing request.
              format: byte
              type: string
            conditions:
              items:
                description: CertificateRequestCondition contains condition information
                  for a CertificateRequest.
                properties:
                  lastTransitionTime:
                    description: LastTransitionTime is the timestamp corresponding
                      to the last status change of this condition.
                    format: date-time
                    type: string
                  message:
                    description: Message is a human readable description of the details
                      of the last transition, complementing reason.
                    type: string
                  reason:
                    description: Reason is a brief machine readable explanation for
                      the condition's last transition.
                    type: string
                  status:
                    description: Status of the condition, one of ('True', 'False',
                      'Unknown').
                    enum:
                    - "True"
                    - "False"
                    - Unknown
                    type: string
                  type:
                    description: Type of the condition, currently ('Ready', 'InvalidRequest').
                    type: string
                required:
                - status
                - type
                type: object
              type: array
            failureTime:
              description: FailureTime stores the time that this CertificateRequest
                failed. This is used to influence garbage collection and back-off.
              format: date-time
              type: string
          type: object
  versions:
  - name: v1alpha2
    served: true
    storage: true
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    cert-manager.io/inject-ca-from-secret: 'kube-system/cert-manager-webhook-ca'
  labels:
    app: 'cert-manager'
    app.kubernetes.io/instance: 'cert-manager'
    app.kubernetes.io/managed-by: 'Tiller'
    app.kubernetes.io/name: 'cert-manager'
    helm.sh/chart: 'cert-manager-v0.15.1'
  name: certificates.cert-manager.io
spec:
  additionalPrinterColumns:
  - JSONPath: .status.conditions[?(@.type=="Ready")].status
    name: Ready
    type: string
  - JSONPath: .spec.secretName
    name: Secret
    type: string
  - JSONPath: .spec.issuerRef.name
    name: Issuer
    priority: 1
    type: string
  - JSONPath: .status.conditions[?(@.type=="Ready")].message
    name: Status
    priority: 1
    type: string
  - JSONPath: .metadata.creationTimestamp
    description: CreationTimestamp is a timestamp representing the server time when
      this object was created. It is not guaranteed to be set in happens-before order
      across separate operations. Clients may not set this value. It is represented
      in RFC3339 form and is in UTC.
    name: Age
    type: date
  group: cert-manager.io
  names:
    kind: Certificate
    listKind: CertificateList
    plural: certificates
    shortNames:
    - cert
    - certs
    singular: certificate
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Certificate is a type to represent a Certificate from ACME
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: CertificateSpec defines the desired state of Certificate. A
            valid Certificate requires at least one of a CommonName, DNSName, or URISAN
            to be valid.
          properties:
            commonName:
              description: 'CommonName is a common name to be used on the Certificate.
                The CommonName should have a length of 64 characters or fewer to avoid
                generating invalid CSRs. This value is ignored by TLS clients when
                any subject alt name is set. This is x509 behaviour: https://tools.ietf.org/html/rfc6125#section-6.4.4'
              type: string
            dnsNames:
              description: DNSNames is a list of subject alt names to be used on the
                Certificate.
              items:
                type: string
              type: array
            duration:
              description: Certificate default Duration
              type: string
            emailSANs:
              description: EmailSANs is a list of Email Subject Alternative Names
                to be set on this Certificate.
              items:
                type: string
              type: array
            ipAddresses:
              description: IPAddresses is a list of IP addresses to be used on the
                Certificate
              items:
                type: string
              type: array
            isCA:
              description: IsCA will mark this Certificate as valid for signing. This
                implies that the 'cert sign' usage is set
              type: boolean
            issuerRef:
              description: IssuerRef is a reference to the issuer for this certificate.
                If the 'kind' field is not set, or set to 'Issuer', an Issuer resource
                with the given name in the same namespace as the Certificate will
                be used. If the 'kind' field is set to 'ClusterIssuer', a ClusterIssuer
                with the provided name will be used. The 'name' field in this stanza
                is required at all times.
              properties:
                group:
                  type: string
                kind:
                  type: string
                name:
                  type: string
              required:
              - name
              type: object
            keyAlgorithm:
              description: KeyAlgorithm is the private key algorithm of the corresponding
                private key for this certificate. If provided, allowed values are
                either "rsa" or "ecdsa" If KeyAlgorithm is specified and KeySize is
                not provided, key size of 256 will be used for "ecdsa" key algorithm
                and key size of 2048 will be used for "rsa" key algorithm.
              enum:
              - rsa
              - ecdsa
              type: string
            keyEncoding:
              description: KeyEncoding is the private key cryptography standards (PKCS)
                for this certificate's private key to be encoded in. If provided,
                allowed values are "pkcs1" and "pkcs8" standing for PKCS#1 and PKCS#8,
                respectively. If KeyEncoding is not specified, then PKCS#1 will be
                used by default.
              enum:
              - pkcs1
              - pkcs8
              type: string
            keySize:
              description: KeySize is the key bit size of the corresponding private
                key for this certificate. If provided, value must be between 2048
                and 8192 inclusive when KeyAlgorithm is empty or is set to "rsa",
                and value must be one of (256, 384, 521) when KeyAlgorithm is set
                to "ecdsa".
              maximum: 8192
              minimum: 0
              type: integer
            keystores:
              description: Keystores configures additional keystore output formats
                stored in the secretName Secret resource.
              properties:
                jks:
                  description: JKS configures options for storing a JKS keystore in
                    the spec.secretName Secret resource.
                  properties:
                    create:
                      description: Create enables JKS keystore creation for the Certificate.
                        If true, a file named keystore.jks will be created in the
                        target Secret resource, encrypted using the password stored
                        in passwordSecretRef. The keystore file will only be updated
                        upon re-issuance.
                      type: boolean
                    passwordSecretRef:
                      description: PasswordSecretRef is a reference to a key in a
                        Secret resource containing the password used to encrypt the
                        JKS keystore.
                      properties:
                        key:
                          description: The key of the secret to select from. Must
                            be a valid secret key.
                          type: string
                        name:
                          description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                            TODO: Add other useful fields. apiVersion, kind, uid?'
                          type: string
                      required:
                      - name
                      type: object
                  required:
                  - create
                  - passwordSecretRef
                  type: object
                pkcs12:
                  description: PKCS12 configures options for storing a PKCS12 keystore
                    in the spec.secretName Secret resource.
                  properties:
                    create:
                      description: Create enables PKCS12 keystore creation for the
                        Certificate. If true, a file named keystore.p12 will be
                        created in the target Secret resource, encrypted using the
                        password stored in passwordSecretRef. The keystore file
                        will only be updated upon re-issuance.
                      type: boolean
                    passwordSecretRef:
                      description: PasswordSecretRef is a reference to a key in a
                        Secret resource containing the password used to encrypt the
                        PKCS12 keystore.
                      properties:
                        key:
                          description: The key of the secret to select from. Must
                            be a valid secret key.
                          type: string
                        name:
                          description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                            TODO: Add other useful fields. apiVersion, kind, uid?'
                          type: string
                      required:
                      - name
                      type: object
                  required:
                  - create
                  - passwordSecretRef
                  type: object
              type: object
            organization:
              description: Organization is the organization to be used on the Certificate
              items:
                type: string
              type: array
            privateKey:
              description: Options to control private keys used for the Certificate.
              properties:
                rotationPolicy:
                  description: RotationPolicy controls how private keys should be
                    regenerated when a re-issuance is being processed. If set to Never,
                    a private key will only be generated if one does not already exist
                    in the target spec.secretName. If one does exists but it does
                    not have the correct algorithm or size, a warning will be raised
                    to await user intervention. If set to Always, a private key matching
                    the specified requirements will be generated whenever a re-issuance
                    occurs. Default is 'Never' for backward compatibility.
                  type: string
              type: object
            renewBefore:
              description: Certificate renew before expiration duration
              type: string
            secretName:
              description: SecretName is the name of the secret resource to store
                this secret in
              type: string
            subject:
              description: Full X509 name specification (https://golang.org/pkg/crypto/x509/pkix/#Name).
              properties:
                countries:
                  description: Countries to be used on the Certificate.
                  items:
                    type: string
                  type: array
                localities:
                  description: Cities to be used on the Certificate.
                  items:
                    type: string
                  type: array
                organizationalUnits:
                  description: Organizational Units to be used on the Certificate.
                  items:
                    type: string
                  type: array
                postalCodes:
                  description: Postal codes to be used on the Certificate.
                  items:
                    type: string
                  type: array
                provinces:
                  description: State/Provinces to be used on the Certificate.
                  items:
                    type: string
                  type: array
                serialNumber:
                  description: Serial number to be used on the Certificate.
                  type: string
                streetAddresses:
                  description: Street addresses to be used on the Certificate.
                  items:
                    type: string
                  type: array
              type: object
            uriSANs:
              description: URISANs is a list of URI Subject Alternative Names to be
                set on this Certificate.
              items:
                type: string
              type: array
            usages:
              description: Usages is the set of x509 actions that are enabled for
                a given key. Defaults are ('digital signature', 'key encipherment')
                if empty
              items:
                description: 'KeyUsage specifies valid usage contexts for keys. See:
                  https://tools.ietf.org/html/rfc5280#section-4.2.1.3      https://tools.ietf.org/html/rfc5280#section-4.2.1.12
                  Valid KeyUsage values are as follows: "signing", "digital signature",
                  "content commitment", "key encipherment", "key agreement", "data
                  encipherment", "cert sign", "crl sign", "encipher only", "decipher
                  only", "any", "server auth", "client auth", "code signing", "email
                  protection", "s/mime", "ipsec end system", "ipsec tunnel", "ipsec
                  user", "timestamping", "ocsp signing", "microsoft sgc", "netscape
                  sgc"'
                enum:
                - signing
                - digital signature
                - content commitment
                - key encipherment
                - key agreement
                - data encipherment
                - cert sign
                - crl sign
                - encipher only
                - decipher only
                - any
                - server auth
                - client auth
                - code signing
                - email protection
                - s/mime
                - ipsec end system
                - ipsec tunnel
                - ipsec user
                - timestamping
                - ocsp signing
                - microsoft sgc
                - netscape sgc
                type: string
              type: array
          required:
          - issuerRef
          - secretName
          type: object
        status:
          description: CertificateStatus defines the observed state of Certificate
          properties:
            conditions:
              items:
                description: CertificateCondition contains condition information for
                  an Certificate.
                properties:
                  lastTransitionTime:
                    description: LastTransitionTime is the timestamp corresponding
                      to the last status change of this condition.
                    format: date-time
                    type: string
                  message:
                    description: Message is a human readable description of the details
                      of the last transition, complementing reason.
                    type: string
                  reason:
                    description: Reason is a brief machine readable explanation for
                      the condition's last transition.
                    type: string
                  status:
                    description: Status of the condition, one of ('True', 'False',
                      'Unknown').
                    enum:
                    - "True"
                    - "False"
                    - Unknown
                    type: string
                  type:
                    description: Type of the condition, currently ('Ready').
                    type: string
                required:
                - status
                - type
                type: object
              type: array
            lastFailureTime:
              format: date-time
              type: string
            nextPrivateKeySecretName:
              description: The name of the Secret resource containing the private
                key to be used for the next certificate iteration. The keymanager
                controller will automatically set this field if the Issuing condition
                is set to True. It will automatically unset this field when the
                Issuing condition is not set or False.
              type: string
            notAfter:
              description: The expiration time of the certificate stored in the secret
                named by this resource in spec.secretName.
              format: date-time
              type: string
            revision:
              description: "The current 'revision' of the certificate as issued. \n
                When a CertificateRequest resource is created, it will have the cert-manager.io/certificate-revision
                set to one greater than the current value of this field. \n Upon issuance,
                this field will be set to the value of the annotation on the CertificateRequest
                resource used to issue the certificate. \n Persisting the value on
                the CertificateRequest resource allows the certificates controller
                to know whether a request is part of an old issuance or if it is part
                of the ongoing revision's issuance by checking if the revision value
                in the annotation is greater than this field."
              type: integer
          type: object
  versions:
  - name: v1alpha2
    served: true
    storage: true
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    cert-manager.io/inject-ca-from-secret: 'kube-system/cert-manager-webhook-ca'
  labels:
    app: 'cert-manager'
    app.kubernetes.io/instance: 'cert-manager'
    app.kubernetes.io/managed-by: 'Tiller'
    app.kubernetes.io/name: 'cert-manager'
    helm.sh/chart: 'cert-manager-v0.15.1'
  name: challenges.acme.cert-manager.io
spec:
  additionalPrinterColumns:
  - JSONPath: .status.state
    name: State
    type: string
  - JSONPath: .spec.dnsName
    name: Domain
    type: string
  - JSONPath: .status.reason
    name: Reason
    priority: 1
    type: string
  - JSONPath: .metadata.creationTimestamp
    description: CreationTimestamp is a timestamp representing the server time when
      this object was created. It is not guaranteed to be set in happens-before order
      across separate operations. Clients may not set this value. It is represented
      in RFC3339 form and is in UTC.
    name: Age
    type: date
  group: acme.cert-manager.io
  names:
    kind: Challenge
    listKind: ChallengeList
    plural: challenges
    singular: challenge
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Challenge is a type to represent a Challenge request with an ACME
        server
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            authzURL:
              description: AuthzURL is the URL to the ACME Authorization resource
                that this challenge is a part of.
              type: string
            dnsName:
              description: DNSName is the identifier that this challenge is for, e.g.
                example.com. If the requested DNSName is a 'wildcard', this field
                MUST be set to the non-wildcard domain, e.g. for *.example.com,
                it must be example.com.
              type: string
            issuerRef:
              description: IssuerRef references a properly configured ACME-type Issuer
                which should be used to create this Challenge. If the Issuer does
                not exist, processing will be retried. If the Issuer is not an 'ACME'
                Issuer, an error will be returned and the Challenge will be marked
                as failed.
              properties:
                group:
                  type: string
                kind:
                  type: string
                name:
                  type: string
              required:
              - name
              type: object
            key:
              description: 'Key is the ACME challenge key for this challenge For HTTP01
                challenges, this is the value that must be responded with to complete
                the HTTP01 challenge in the format: <private key JWK thumbprint>.<key
                from acme server for challenge>. For DNS01 challenges, this is the
                base64 encoded SHA256 sum of the <private key JWK thumbprint>.<key
                from acme server for challenge> text that must be set as the TXT
                record content.'
              type: string
            solver:
              description: Solver contains the domain solving configuration that should
                be used to solve this challenge resource.
              properties:
                dns01:
                  properties:
                    acmedns:
                      description: ACMEIssuerDNS01ProviderAcmeDNS is a structure containing
                        the configuration for ACME-DNS servers
                      properties:
                        accountSecretRef:
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                        host:
                          type: string
                      required:
                      - accountSecretRef
                      - host
                      type: object
                    akamai:
                      description: ACMEIssuerDNS01ProviderAkamai is a structure containing
                        the DNS configuration for Akamai DNS—Zone Record Management
                        API
                      properties:
                        accessTokenSecretRef:
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                        clientSecretSecretRef:
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                        clientTokenSecretRef:
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                        serviceConsumerDomain:
                          type: string
                      required:
                      - accessTokenSecretRef
                      - clientSecretSecretRef
                      - clientTokenSecretRef
                      - serviceConsumerDomain
                      type: object
                    azuredns:
                      description: ACMEIssuerDNS01ProviderAzureDNS is a structure
                        containing the configuration for Azure DNS
                      properties:
                        clientID:
                          description: if both this and ClientSecret are left unset
                            MSI will be used
                          type: string
                        clientSecretSecretRef:
                          description: if both this and ClientID are left unset MSI
                            will be used
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                        environment:
                          enum:
                          - AzurePublicCloud
                          - AzureChinaCloud
                          - AzureGermanCloud
                          - AzureUSGovernmentCloud
                          type: string
                        hostedZoneName:
                          type: string
                        resourceGroupName:
                          type: string
                        subscriptionID:
                          type: string
                        tenantID:
                          description: when specifying ClientID and ClientSecret then
                            this field is also needed
                          type: string
                      required:
                      - resourceGroupName
                      - subscriptionID
                      type: object
                    clouddns:
                      description: ACMEIssuerDNS01ProviderCloudDNS is a structure
                        containing the DNS configuration for Google Cloud DNS
                      properties:
                        project:
                          type: string
                        serviceAccountSecretRef:
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                      required:
                      - project
                      type: object
                    cloudflare:
                      description: ACMEIssuerDNS01ProviderCloudflare is a structure
                        containing the DNS configuration for Cloudflare
                      properties:
                        apiKeySecretRef:
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                        apiTokenSecretRef:
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                        email:
                          type: string
                      required:
                      - email
                      type: object
                    cnameStrategy:
                      description: CNAMEStrategy configures how the DNS01 provider
                        should handle CNAME records when found in DNS zones.
                      enum:
                      - None
                      - Follow
                      type: string
                    digitalocean:
                      description: ACMEIssuerDNS01ProviderDigitalOcean is a structure
                        containing the DNS configuration for DigitalOcean Domains
                      properties:
                        tokenSecretRef:
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                      required:
                      - tokenSecretRef
                      type: object
                    rfc2136:
                      description: ACMEIssuerDNS01ProviderRFC2136 is a structure containing
                        the configuration for RFC2136 DNS
                      properties:
                        nameserver:
                          description: The IP address or hostname of an authoritative
                            DNS server supporting RFC2136 in the form host:port. If
                            the host is an IPv6 address it must be enclosed in square
                            brackets (e.g [2001:db8::1]) ; port is optional. This
                            field is required.
                          type: string
                        tsigAlgorithm:
                          description: 'The TSIG Algorithm configured in the DNS supporting
                            RFC2136. Used only when tsigSecretSecretRef and tsigKeyName
                            are defined. Supported values are (case-insensitive):
                            HMACMD5 (default), HMACSHA1, HMACSHA256 or
                            HMACSHA512.'
                          type: string
                        tsigKeyName:
                          description: The TSIG Key name configured in the DNS. If
                            tsigSecretSecretRef is defined, this field is required.
                          type: string
                        tsigSecretSecretRef:
                          description: The name of the secret containing the TSIG
                            value. If tsigKeyName is defined, this field is required.
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                      required:
                      - nameserver
                      type: object
                    route53:
                      description: ACMEIssuerDNS01ProviderRoute53 is a structure containing
                        the Route 53 configuration for AWS
                      properties:
                        accessKeyID:
                          description: 'The AccessKeyID is used for authentication.
                            If not set we fall-back to using env vars, shared credentials
                            file or AWS Instance metadata see: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials'
                          type: string
                        hostedZoneID:
                          description: If set, the provider will manage only this
                            zone in Route53 and will not do an lookup using the route53:ListHostedZonesByName
                            api call.
                          type: string
                        region:
                          description: Always set the region when using AccessKeyID
                            and SecretAccessKey
                          type: string
                        role:
                          description: Role is a Role ARN which the Route53 provider
                            will assume using either the explicit credentials AccessKeyID/SecretAccessKey
                            or the inferred credentials from environment variables,
                            shared credentials file or AWS Instance metadata
                          type: string
                        secretAccessKeySecretRef:
                          description: The SecretAccessKey is used for authentication.
                            If not set we fall-back to using env vars, shared credentials
                            file or AWS Instance metadata https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                      required:
                      - region
                      type: object
                    webhook:
                      description: ACMEIssuerDNS01ProviderWebhook specifies configuration
                        for a webhook DNS01 provider, including where to POST ChallengePayload
                        resources.
                      properties:
                        config:
                          description: Additional configuration that should be passed
                            to the webhook apiserver when challenges are processed.
                            This can contain arbitrary JSON data. Secret values should
                            not be specified in this stanza. If secret values are
                            needed (e.g. credentials for a DNS service), you should
                            use a SecretKeySelector to reference a Secret resource.
                            For details on the schema of this field, consult the webhook
                            provider implementation's documentation.
                        groupName:
                          description: The API group name that should be used when
                            POSTing ChallengePayload resources to the webhook apiserver.
                            This should be the same as the GroupName specified in
                            the webhook provider implementation.
                          type: string
                        solverName:
                          description: The name of the solver to use, as defined in
                            the webhook provider implementation. This will typically
                            be the name of the provider, e.g. 'cloudflare'.
                          type: string
                      required:
                      - groupName
                      - solverName
                      type: object
                  type: object
                http01:
                  description: ACMEChallengeSolverHTTP01 contains configuration detailing
                    how to solve HTTP01 challenges within a Kubernetes cluster. Typically
                    this is accomplished through creating 'routes' of some description
                    that configure ingress controllers to direct traffic to 'solver
                    pods', which are responsible for responding to the ACME server's
                    HTTP requests.
                  properties:
                    ingress:
                      description: The ingress based HTTP01 challenge solver will
                        solve challenges by creating or modifying Ingress resources
                        in order to route requests for '/.well-known/acme-challenge/XYZ'
                        to 'challenge solver' pods that are provisioned by cert-manager
                        for each Challenge to be completed.
                      properties:
                        class:
                          description: The ingress class to use when creating Ingress
                            resources to solve ACME challenges that use this challenge
                            solver. Only one of 'class' or 'name' may be specified.
                          type: string
                        ingressTemplate:
                          description: Optional ingress template used to configure
                            the ACME challenge solver ingress used for HTTP01 challenges
                          properties:
                            metadata:
                              description: ObjectMeta overrides for the ingress used
                                to solve HTTP01 challenges. Only the 'labels' and
                                'annotations' fields may be set. If labels or annotations
                                overlap with in-built values, the values here will
                                override the in-built values.
                              properties:
                                annotations:
                                  additionalProperties:
                                    type: string
                                  description: Annotations that should be added to
                                    the created ACME HTTP01 solver ingress.
                                  type: object
                                labels:
                                  additionalProperties:
                                    type: string
                                  description: Labels that should be added to the
                                    created ACME HTTP01 solver ingress.
                                  type: object
                              type: object
                          type: object
                        name:
                          description: The name of the ingress resource that should
                            have ACME challenge solving routes inserted into it in
                            order to solve HTTP01 challenges. This is typically used
                            in conjunction with ingress controllers like ingress-gce,
                            which maintains a 1:1 mapping between external IPs and
                            ingress resources.
                          type: string
                        podTemplate:
                          description: Optional pod template used to configure the
                            ACME challenge solver pods used for HTTP01 challenges
                          properties:
                            metadata:
                              description: ObjectMeta overrides for the pod used to
                                solve HTTP01 challenges. Only the 'labels' and 'annotations'
                                fields may be set. If labels or annotations overlap
                                with in-built values, the values here will override
                                the in-built values.
                              properties:
                                annotations:
                                  additionalProperties:
                                    type: string
                                  description: Annotations that should be added to
                                    the create ACME HTTP01 solver pods.
                                  type: object
                                labels:
                                  additionalProperties:
                                    type: string
                                  description: Labels that should be added to the
                                    created ACME HTTP01 solver pods.
                                  type: object
                              type: object
                            spec:
                              description: PodSpec defines overrides for the HTTP01
                                challenge solver pod. Only the 'nodeSelector', 'affinity'
                                and 'tolerations' fields are supported currently.
                                All other fields will be ignored.
                              properties:
                                affinity:
                                  description: If specified, the pod's scheduling
                                    constraints
                                  properties:
                                    nodeAffinity:
                                      description: Describes node affinity scheduling
                                        rules for the pod.
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          description: The scheduler will prefer to
                                            schedule pods to nodes that satisfy the
                                            affinity expressions specified by this
                                            field, but it may choose a node that violates
                                            one or more of the expressions. The node
                                            that is most preferred is the one with
                                            the greatest sum of weights, i.e. for
                                            each node that meets all of the scheduling
                                            requirements (resource request, requiredDuringScheduling
                                            affinity expressions, etc.), compute a
                                            sum by iterating through the elements
                                            of this field and adding "weight" to the
                                            sum if the node matches the corresponding
                                            matchExpressions; the node(s) with the
                                            highest sum are the most preferred.
                                          items:
                                            description: An empty preferred scheduling
                                              term matches all objects with implicit
                                              weight 0 (i.e. it's a no-op). A null
                                              preferred scheduling term matches no
                                              objects (i.e. is also a no-op).
                                            properties:
                                              preference:
                                                description: A node selector term,
                                                  associated with the corresponding
                                                  weight.
                                                properties:
                                                  matchExpressions:
                                                    description: A list of node selector
                                                      requirements by node's labels.
                                                    items:
                                                      description: A node selector
                                                        requirement is a selector
                                                        that contains values, a key,
                                                        and an operator that relates
                                                        the key and values.
                                                      properties:
                                                        key:
                                                          description: The label key
                                                            that the selector applies
                                                            to.
                                                          type: string
                                                        operator:
                                                          description: Represents
                                                            a key's relationship to
                                                            a set of values. Valid
                                                            operators are In, NotIn,
                                                            Exists, DoesNotExist.
                                                            Gt, and Lt.
                                                          type: string
                                                        values:
                                                          description: An array of
                                                            string values. If the
                                                            operator is In or NotIn,
                                                            the values array must
                                                            be non-empty. If the operator
                                                            is Exists or DoesNotExist,
                                                            the values array must
                                                            be empty. If the operator
                                                            is Gt or Lt, the values
                                                            array must have a single
                                                            element, which will be
                                                            interpreted as an integer.
                                                            This array is replaced
                                                            during a strategic merge
                                                            patch.
                                                          items:
                                                            type: string
                                                          type: array
                                                      required:
                                                      - key
                                                      - operator
                                                      type: object
                                                    type: array
                                                  matchFields:
                                                    description: A list of node selector
                                                      requirements by node's fields.
                                                    items:
                                                      description: A node selector
                                                        requirement is a selector
                                                        that contains values, a key,
                                                        and an operator that relates
                                                        the key and values.
                                                      properties:
                                                        key:
                                                          description: The label key
                                                            that the selector applies
                                                            to.
                                                          type: string
                                                        operator:
                                                          description: Represents
                                                            a key's relationship to
                                                            a set of values. Valid
                                                            operators are In, NotIn,
                                                            Exists, DoesNotExist.
                                                            Gt, and Lt.
                                                          type: string
                                                        values:
                                                          description: An array of
                                                            string values. If the
                                                            operator is In or NotIn,
                                                            the values array must
                                                            be non-empty. If the operator
                                                            is Exists or DoesNotExist,
                                                            the values array must
                                                            be empty. If the operator
                                                            is Gt or Lt, the values
                                                            array must have a single
                                                            element, which will be
                                                            interpreted as an integer.
                                                            This array is replaced
                                                            during a strategic merge
                                                            patch.
                                                          items:
                                                            type: string
                                                          type: array
                                                      required:
                                                      - key
                                                      - operator
                                                      type: object
                                                    type: array
                                                type: object
                                              weight:
                                                description: Weight associated with
                                                  matching the corresponding nodeSelectorTerm,
                                                  in the range 1-100.
                                                format: int32
                                                type: integer
                                            required:
                                            - preference
                                            - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          description: If the affinity requirements
                                            specified by this field are not met at
                                            scheduling time, the pod will not be scheduled
                                            onto the node. If the affinity requirements
                                            specified by this field cease to be met
                                            at some point during pod execution (e.g.
                                            due to an update), the system may or may
                                            not try to eventually evict the pod from
                                            its node.
                                          properties:
                                            nodeSelectorTerms:
                                              description: Required. A list of node
                                                selector terms. The terms are ORed.
                                              items:
                                                description: A null or empty node
                                                  selector term matches no objects.
                                                  The requirements of them are ANDed.
                                                  The TopologySelectorTerm type implements
                                                  a subset of the NodeSelectorTerm.
                                                properties:
                                                  matchExpressions:
                                                    description: A list of node selector
                                                      requirements by node's labels.
                                                    items:
                                                      description: A node selector
                                                        requirement is a selector
                                                        that contains values, a key,
                                                        and an operator that relates
                                                        the key and values.
                                                      properties:
                                                        key:
                                                          description: The label key
                                                            that the selector applies
                                                            to.
                                                          type: string
                                                        operator:
                                                          description: Represents
                                                            a key's relationship to
                                                            a set of values. Valid
                                                            operators are In, NotIn,
                                                            Exists, DoesNotExist.
                                                            Gt, and Lt.
                                                          type: string
                                                        values:
                                                          description: An array of
                                                            string values. If the
                                                            operator is In or NotIn,
                                                            the values array must
                                                            be non-empty. If the operator
                                                            is Exists or DoesNotExist,
                                                            the values array must
                                                            be empty. If the operator
                                                            is Gt or Lt, the values
                                                            array must have a single
                                                            element, which will be
                                                            interpreted as an integer.
                                                            This array is replaced
                                                            during a strategic merge
                                                            patch.
                                                          items:
                                                            type: string
                                                          type: array
                                                      required:
                                                      - key
                                                      - operator
                                                      type: object
                                                    type: array
                                                  matchFields:
                                                    description: A list of node selector
                                                      requirements by node's fields.
                                                    items:
                                                      description: A node selector
                                                        requirement is a selector
                                                        that contains values, a key,
                                                        and an operator that relates
                                                        the key and values.
                                                      properties:
                                                        key:
                                                          description: The label key
                                                            that the selector applies
                                                            to.
                                                          type: string
                                                        operator:
                                                          description: Represents
                                                            a key's relationship to
                                                            a set of values. Valid
                                                            operators are In, NotIn,
                                                            Exists, DoesNotExist.
                                                            Gt, and Lt.
                                                          type: string
                                                        values:
                                                          description: An array of
                                                            string values. If the
                                                            operator is In or NotIn,
                                                            the values array must
                                                            be non-empty. If the operator
                                                            is Exists or DoesNotExist,
                                                            the values array must
                                                            be empty. If the operator
                                                            is Gt or Lt, the values
                                                            array must have a single
                                                            element, which will be
                                                            interpreted as an integer.
                                                            This array is replaced
                                                            during a strategic merge
                                                            patch.
                                                          items:
                                                            type: string
                                                          type: array
                                                      required:
                                                      - key
                                                      - operator
                                                      type: object
                                                    type: array
                                                type: object
                                              type: array
                                          required:
                                          - nodeSelectorTerms
                                          type: object
                                      type: object
                                    podAffinity:
                                      description: Describes pod affinity scheduling
                                        rules (e.g. co-locate this pod in the same
                                        node, zone, etc. as some other pod(s)).
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          description: The scheduler will prefer to
                                            schedule pods to nodes that satisfy the
                                            affinity expressions specified by this
                                            field, but it may choose a node that violates
                                            one or more of the expressions. The node
                                            that is most preferred is the one with
                                            the greatest sum of weights, i.e. for
                                            each node that meets all of the scheduling
                                            requirements (resource request, requiredDuringScheduling
                                            affinity expressions, etc.), compute a
                                            sum by iterating through the elements
                                            of this field and adding "weight" to the
                                            sum if the node has pods which matches
                                            the corresponding podAffinityTerm; the
                                            node(s) with the highest sum are the most
                                            preferred.
                                          items:
                                            description: The weights of all of the
                                              matched WeightedPodAffinityTerm fields
                                              are added per-node to find the most
                                              preferred node(s)
                                            properties:
                                              podAffinityTerm:
                                                description: Required. A pod affinity
                                                  term, associated with the corresponding
                                                  weight.
                                                properties:
                                                  labelSelector:
                                                    description: A label query over
                                                      a set of resources, in this
                                                      case pods.
                                                    properties:
                                                      matchExpressions:
                                                        description: matchExpressions
                                                          is a list of label selector
                                                          requirements. The requirements
                                                          are ANDed.
                                                        items:
                                                          description: A label selector
                                                            requirement is a selector
                                                            that contains values,
                                                            a key, and an operator
                                                            that relates the key and
                                                            values.
                                                          properties:
                                                            key:
                                                              description: key is
                                                                the label key that
                                                                the selector applies
                                                                to.
                                                              type: string
                                                            operator:
                                                              description: operator
                                                                represents a key's
                                                                relationship to a
                                                                set of values. Valid
                                                                operators are In,
                                                                NotIn, Exists and
                                                                DoesNotExist.
                                                              type: string
                                                            values:
                                                              description: values
                                                                is an array of string
                                                                values. If the operator
                                                                is In or NotIn, the
                                                                values array must
                                                                be non-empty. If the
                                                                operator is Exists
                                                                or DoesNotExist, the
                                                                values array must
                                                                be empty. This array
                                                                is replaced during
                                                                a strategic merge
                                                                patch.
                                                              items:
                                                                type: string
                                                              type: array
                                                          required:
                                                          - key
                                                          - operator
                                                          type: object
                                                        type: array
                                                      matchLabels:
                                                        additionalProperties:
                                                          type: string
                                                        description: matchLabels is
                                                          a map of {key,value} pairs.
                                                          A single {key,value} in
                                                          the matchLabels map is equivalent
                                                          to an element of matchExpressions,
                                                          whose key field is "key",
                                                          the operator is "In", and
                                                          the values array contains
                                                          only "value". The requirements
                                                          are ANDed.
                                                        type: object
                                                    type: object
                                                  namespaces:
                                                    description: namespaces specifies
                                                      which namespaces the labelSelector
                                                      applies to (matches against);
                                                      null or empty list means "this
                                                      pod's namespace"
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    description: This pod should be
                                                      co-located (affinity) or not
                                                      co-located (anti-affinity) with
                                                      the pods matching the labelSelector
                                                      in the specified namespaces,
                                                      where co-located is defined
                                                      as running on a node whose value
                                                      of the label with key topologyKey
                                                      matches that of any node on
                                                      which any of the selected pods
                                                      is running. Empty topologyKey
                                                      is not allowed.
                                                    type: string
                                                required:
                                                - topologyKey
                                                type: object
                                              weight:
                                                description: weight associated with
                                                  matching the corresponding podAffinityTerm,
                                                  in the range 1-100.
                                                format: int32
                                                type: integer
                                            required:
                                            - podAffinityTerm
                                            - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          description: If the affinity requirements
                                            specified by this field are not met at
                                            scheduling time, the pod will not be scheduled
                                            onto the node. If the affinity requirements
                                            specified by this field cease to be met
                                            at some point during pod execution (e.g.
                                            due to a pod label update), the system
                                            may or may not try to eventually evict
                                            the pod from its node. When there are
                                            multiple elements, the lists of nodes
                                            corresponding to each podAffinityTerm
                                            are intersected, i.e. all terms must be
                                            satisfied.
                                          items:
                                            description: Defines a set of pods (namely
                                              those matching the labelSelector relative
                                              to the given namespace(s)) that this
                                              pod should be co-located (affinity)
                                              or not co-located (anti-affinity) with,
                                              where co-located is defined as running
                                              on a node whose value of the label with
                                              key <topologyKey> matches that of any
                                              node on which a pod of the set of pods
                                              is running
                                            properties:
                                              labelSelector:
                                                description: A label query over a
                                                  set of resources, in this case pods.
                                                properties:
                                                  matchExpressions:
                                                    description: matchExpressions
                                                      is a list of label selector
                                                      requirements. The requirements
                                                      are ANDed.
                                                    items:
                                                      description: A label selector
                                                        requirement is a selector
                                                        that contains values, a key,
                                                        and an operator that relates
                                                        the key and values.
                                                      properties:
                                                        key:
                                                          description: key is the
                                                            label key that the selector
                                                            applies to.
                                                          type: string
                                                        operator:
                                                          description: operator represents
                                                            a key's relationship to
                                                            a set of values. Valid
                                                            operators are In, NotIn,
                                                            Exists and DoesNotExist.
                                                          type: string
                                                        values:
                                                          description: values is an
                                                            array of string values.
                                                            If the operator is In
                                                            or NotIn, the values array
                                                            must be non-empty. If
                                                            the operator is Exists
                                                            or DoesNotExist, the values
                                                            array must be empty. This
                                                            array is replaced during
                                                            a strategic merge patch.
                                                          items:
                                                            type: string
                                                          type: array
                                                      required:
                                                      - key
                                                      - operator
                                                      type: object
                                                    type: array
                                                  matchLabels:
                                                    additionalProperties:
                                                      type: string
                                                    description: matchLabels is a
                                                      map of {key,value} pairs. A
                                                      single {key,value} in the matchLabels
                                                      map is equivalent to an element
                                                      of matchExpressions, whose key
                                                      field is "key", the operator
                                                      is "In", and the values array
                                                      contains only "value". The requirements
                                                      are ANDed.
                                                    type: object
                                                type: object
                                              namespaces:
                                                description: namespaces specifies
                                                  which namespaces the labelSelector
                                                  applies to (matches against); null
                                                  or empty list means "this pod's
                                                  namespace"
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                description: This pod should be co-located
                                                  (affinity) or not co-located (anti-affinity)
                                                  with the pods matching the labelSelector
                                                  in the specified namespaces, where
                                                  co-located is defined as running
                                                  on a node whose value of the label
                                                  with key topologyKey matches that
                                                  of any node on which any of the
                                                  selected pods is running. Empty
                                                  topologyKey is not allowed.
                                                type: string
                                            required:
                                            - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                    podAntiAffinity:
                                      description: Describes pod anti-affinity scheduling
                                        rules (e.g. avoid putting this pod in the
                                        same node, zone, etc. as some other pod(s)).
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          description: The scheduler will prefer to
                                            schedule pods to nodes that satisfy the
                                            anti-affinity expressions specified by
                                            this field, but it may choose a node that
                                            violates one or more of the expressions.
                                            The node that is most preferred is the
                                            one with the greatest sum of weights,
                                            i.e. for each node that meets all of the
                                            scheduling requirements (resource request,
                                            requiredDuringScheduling anti-affinity
                                            expressions, etc.), compute a sum by iterating
                                            through the elements of this field and
                                            adding "weight" to the sum if the node
                                            has pods which matches the corresponding
                                            podAffinityTerm; the node(s) with the
                                            highest sum are the most preferred.
                                          items:
                                            description: The weights of all of the
                                              matched WeightedPodAffinityTerm fields
                                              are added per-node to find the most
                                              preferred node(s)
                                            properties:
                                              podAffinityTerm:
                                                description: Required. A pod affinity
                                                  term, associated with the corresponding
                                                  weight.
                                                properties:
                                                  labelSelector:
                                                    description: A label query over
                                                      a set of resources, in this
                                                      case pods.
                                                    properties:
                                                      matchExpressions:
                                                        description: matchExpressions
                                                          is a list of label selector
                                                          requirements. The requirements
                                                          are ANDed.
                                                        items:
                                                          description: A label selector
                                                            requirement is a selector
                                                            that contains values,
                                                            a key, and an operator
                                                            that relates the key and
                                                            values.
                                                          properties:
                                                            key:
                                                              description: key is
                                                                the label key that
                                                                the selector applies
                                                                to.
                                                              type: string
                                                            operator:
                                                              description: operator
                                                                represents a key's
                                                                relationship to a
                                                                set of values. Valid
                                                                operators are In,
                                                                NotIn, Exists and
                                                                DoesNotExist.
                                                              type: string
                                                            values:
                                                              description: values
                                                                is an array of string
                                                                values. If the operator
                                                                is In or NotIn, the
                                                                values array must
                                                                be non-empty. If the
                                                                operator is Exists
                                                                or DoesNotExist, the
                                                                values array must
                                                                be empty. This array
                                                                is replaced during
                                                                a strategic merge
                                                                patch.
                                                              items:
                                                                type: string
                                                              type: array
                                                          required:
                                                          - key
                                                          - operator
                                                          type: object
                                                        type: array
                                                      matchLabels:
                                                        additionalProperties:
                                                          type: string
                                                        description: matchLabels is
                                                          a map of {key,value} pairs.
                                                          A single {key,value} in
                                                          the matchLabels map is equivalent
                                                          to an element of matchExpressions,
                                                          whose key field is "key",
                                                          the operator is "In", and
                                                          the values array contains
                                                          only "value". The requirements
                                                          are ANDed.
                                                        type: object
                                                    type: object
                                                  namespaces:
                                                    description: namespaces specifies
                                                      which namespaces the labelSelector
                                                      applies to (matches against);
                                                      null or empty list means "this
                                                      pod's namespace"
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    description: This pod should be
                                                      co-located (affinity) or not
                                                      co-located (anti-affinity) with
                                                      the pods matching the labelSelector
                                                      in the specified namespaces,
                                                      where co-located is defined
                                                      as running on a node whose value
                                                      of the label with key topologyKey
                                                      matches that of any node on
                                                      which any of the selected pods
                                                      is running. Empty topologyKey
                                                      is not allowed.
                                                    type: string
                                                required:
                                                - topologyKey
                                                type: object
                                              weight:
                                                description: weight associated with
                                                  matching the corresponding podAffinityTerm,
                                                  in the range 1-100.
                                                format: int32
                                                type: integer
                                            required:
                                            - podAffinityTerm
                                            - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          description: If the anti-affinity requirements
                                            specified by this field are not met at
                                            scheduling time, the pod will not be scheduled
                                            onto the node. If the anti-affinity requirements
                                            specified by this field cease to be met
                                            at some point during pod execution (e.g.
                                            due to a pod label update), the system
                                            may or may not try to eventually evict
                                            the pod from its node. When there are
                                            multiple elements, the lists of nodes
                                            corresponding to each podAffinityTerm
                                            are intersected, i.e. all terms must be
                                            satisfied.
                                          items:
                                            description: Defines a set of pods (namely
                                              those matching the labelSelector relative
                                              to the given namespace(s)) that this
                                              pod should be co-located (affinity)
                                              or not co-located (anti-affinity) with,
                                              where co-located is defined as running
                                              on a node whose value of the label with
                                              key <topologyKey> matches that of any
                                              node on which a pod of the set of pods
                                              is running
                                            properties:
                                              labelSelector:
                                                description: A label query over a
                                                  set of resources, in this case pods.
                                                properties:
                                                  matchExpressions:
                                                    description: matchExpressions
                                                      is a list of label selector
                                                      requirements. The requirements
                                                      are ANDed.
                                                    items:
                                                      description: A label selector
                                                        requirement is a selector
                                                        that contains values, a key,
                                                        and an operator that relates
                                                        the key and values.
                                                      properties:
                                                        key:
                                                          description: key is the
                                                            label key that the selector
                                                            applies to.
                                                          type: string
                                                        operator:
                                                          description: operator represents
                                                            a key's relationship to
                                                            a set of values. Valid
                                                            operators are In, NotIn,
                                                            Exists and DoesNotExist.
                                                          type: string
                                                        values:
                                                          description: values is an
                                                            array of string values.
                                                            If the operator is In
                                                            or NotIn, the values array
                                                            must be non-empty. If
                                                            the operator is Exists
                                                            or DoesNotExist, the values
                                                            array must be empty. This
                                                            array is replaced during
                                                            a strategic merge patch.
                                                          items:
                                                            type: string
                                                          type: array
                                                      required:
                                                      - key
                                                      - operator
                                                      type: object
                                                    type: array
                                                  matchLabels:
                                                    additionalProperties:
                                                      type: string
                                                    description: matchLabels is a
                                                      map of {key,value} pairs. A
                                                      single {key,value} in the matchLabels
                                                      map is equivalent to an element
                                                      of matchExpressions, whose key
                                                      field is "key", the operator
                                                      is "In", and the values array
                                                      contains only "value". The requirements
                                                      are ANDed.
                                                    type: object
                                                type: object
                                              namespaces:
                                                description: namespaces specifies
                                                  which namespaces the labelSelector
                                                  applies to (matches against); null
                                                  or empty list means "this pod's
                                                  namespace"
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                description: This pod should be co-located
                                                  (affinity) or not co-located (anti-affinity)
                                                  with the pods matching the labelSelector
                                                  in the specified namespaces, where
                                                  co-located is defined as running
                                                  on a node whose value of the label
                                                  with key topologyKey matches that
                                                  of any node on which any of the
                                                  selected pods is running. Empty
                                                  topologyKey is not allowed.
                                                type: string
                                            required:
                                            - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                  type: object
                                nodeSelector:
                                  additionalProperties:
                                    type: string
                                  description: 'NodeSelector is a selector which must
                                    be true for the pod to fit on a node. Selector
                                    which must match a node''s labels for the pod
                                    to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/'
                                  type: object
                                tolerations:
                                  description: If specified, the pod's tolerations.
                                  items:
                                    description: The pod this Toleration is attached
                                      to tolerates any taint that matches the triple
                                      <key,value,effect> using the matching operator
                                      <operator>.
                                    properties:
                                      effect:
                                        description: Effect indicates the taint effect
                                          to match. Empty means match all taint effects.
                                          When specified, allowed values are NoSchedule,
                                          PreferNoSchedule and NoExecute.
                                        type: string
                                      key:
                                        description: Key is the taint key that the
                                          toleration applies to. Empty means match
                                          all taint keys. If the key is empty, operator
                                          must be Exists; this combination means to
                                          match all values and all keys.
                                        type: string
                                      operator:
                                        description: Operator represents a key's relationship
                                          to the value. Valid operators are Exists
                                          and Equal. Defaults to Equal. Exists is
                                          equivalent to wildcard for value, so that
                                          a pod can tolerate all taints of a particular
                                          category.
                                        type: string
                                      tolerationSeconds:
                                        description: TolerationSeconds represents
                                          the period of time the toleration (which
                                          must be of effect NoExecute, otherwise this
                                          field is ignored) tolerates the taint. By
                                          default, it is not set, which means tolerate
                                          the taint forever (do not evict). Zero and
                                          negative values will be treated as 0 (evict
                                          immediately) by the system.
                                        format: int64
                                        type: integer
                                      value:
                                        description: Value is the taint value the
                                          toleration matches to. If the operator is
                                          Exists, the value should be empty, otherwise
                                          just a regular string.
                                        type: string
                                    type: object
                                  type: array
                              type: object
                          type: object
                        serviceType:
                          description: Optional service type for Kubernetes solver
                            service
                          type: string
                      type: object
                  type: object
                selector:
                  description: Selector selects a set of DNSNames on the Certificate
                    resource that should be solved using this challenge solver.
                  properties:
                    dnsNames:
                      description: List of DNSNames that this solver will be used
                        to solve. If specified and a match is found, a dnsNames selector
                        will take precedence over a dnsZones selector. If multiple
                        solvers match with the same dnsNames value, the solver with
                        the most matching labels in matchLabels will be selected.
                        If neither has more matches, the solver defined earlier in
                        the list will be selected.
                      items:
                        type: string
                      type: array
                    dnsZones:
                      description: List of DNSZones that this solver will be used
                        to solve. The most specific DNS zone match specified here
                        will take precedence over other DNS zone matches, so a solver
                        specifying sys.example.com will be selected over one specifying
                        example.com for the domain www.sys.example.com. If multiple
                        solvers match with the same dnsZones value, the solver with
                        the most matching labels in matchLabels will be selected.
                        If neither has more matches, the solver defined earlier in
                        the list will be selected.
                      items:
                        type: string
                      type: array
                    matchLabels:
                      additionalProperties:
                        type: string
                      description: A label selector that is used to refine the set
                        of certificate's that this challenge solver will apply to.
                      type: object
                  type: object
              type: object
            token:
              description: Token is the ACME challenge token for this challenge. This
                is the raw value returned from the ACME server.
              type: string
            type:
              description: Type is the type of ACME challenge this resource represents,
                e.g. "dns01" or "http01".
              type: string
            url:
              description: URL is the URL of the ACME Challenge resource for this
                challenge. This can be used to lookup details about the status of
                this challenge.
              type: string
            wildcard:
              description: Wildcard will be true if this challenge is for a wildcard
                identifier, for example '*.example.com'.
              type: boolean
          required:
          - authzURL
          - dnsName
          - issuerRef
          - key
          - solver
          - token
          - type
          - url
          type: object
        status:
          properties:
            presented:
              description: Presented will be set to true if the challenge values for
                this challenge are currently 'presented'. This *does not* imply the
                self check is passing. Only that the values have been 'submitted'
                for the appropriate challenge mechanism (i.e. the DNS01 TXT record
                has been presented, or the HTTP01 configuration has been configured).
              type: boolean
            processing:
              description: Processing is used to denote whether this challenge should
                be processed or not. This field will only be set to true by the 'scheduling'
                component. It will only be set to false by the 'challenges' controller,
                after the challenge has reached a final state or timed out. If this
                field is set to false, the challenge controller will not take any
                more action.
              type: boolean
            reason:
              description: Reason contains human readable information on why the Challenge
                is in the current state.
              type: string
            state:
              description: State contains the current 'state' of the challenge. If
                not set, the state of the challenge is unknown.
              enum:
              - valid
              - ready
              - pending
              - processing
              - invalid
              - expired
              - errored
              type: string
          type: object
      required:
      - metadata
  versions:
  - name: v1alpha2
    served: true
    storage: true
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    cert-manager.io/inject-ca-from-secret: 'kube-system/cert-manager-webhook-ca'
  labels:
    app: 'cert-manager'
    app.kubernetes.io/instance: 'cert-manager'
    app.kubernetes.io/managed-by: 'Tiller'
    app.kubernetes.io/name: 'cert-manager'
    helm.sh/chart: 'cert-manager-v0.15.1'
  name: clusterissuers.cert-manager.io
spec:
  additionalPrinterColumns:
  - JSONPath: .status.conditions[?(@.type=="Ready")].status
    name: Ready
    type: string
  - JSONPath: .status.conditions[?(@.type=="Ready")].message
    name: Status
    priority: 1
    type: string
  - JSONPath: .metadata.creationTimestamp
    description: CreationTimestamp is a timestamp representing the server time when
      this object was created. It is not guaranteed to be set in happens-before order
      across separate operations. Clients may not set this value. It is represented
      in RFC3339 form and is in UTC.
    name: Age
    type: date
  group: cert-manager.io
  names:
    kind: ClusterIssuer
    listKind: ClusterIssuerList
    plural: clusterissuers
    singular: clusterissuer
  scope: Cluster
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: IssuerSpec is the specification of an Issuer. This includes
            any configuration required for the issuer.
          properties:
            acme:
              description: ACMEIssuer contains the specification for an ACME issuer
              properties:
                email:
                  description: Email is the email for this account
                  type: string
                externalAccountBinding:
                  description: ExternalAccountBinding is a reference to a CA external
                    account of the ACME server.
                  properties:
                    keyAlgorithm:
                      description: keyAlgorithm is the MAC key algorithm that the
                        key is used for. Valid values are "HS256", "HS384" and "HS512".
                      enum:
                      - HS256
                      - HS384
                      - HS512
                      type: string
                    keyID:
                      description: keyID is the ID of the CA key that the External
                        Account is bound to.
                      type: string
                    keySecretRef:
                      description: keySecretRef is a Secret Key Selector referencing
                        a data item in a Kubernetes Secret which holds the symmetric
                        MAC key of the External Account Binding. The key is the
                        index string that is paired with the key data in the Secret
                        and should not be confused with the key data itself, or indeed
                        with the External Account Binding keyID above. The secret
                        key stored in the Secret **must** be un-padded, base64 URL
                        encoded data.
                      properties:
                        key:
                          description: The key of the secret to select from. Must
                            be a valid secret key.
                          type: string
                        name:
                          description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                            TODO: Add other useful fields. apiVersion, kind, uid?'
                          type: string
                      required:
                      - name
                      type: object
                  required:
                  - keyAlgorithm
                  - keyID
                  - keySecretRef
                  type: object
                privateKeySecretRef:
                  description: PrivateKey is the name of a secret containing the private
                    key for this user account.
                  properties:
                    key:
                      description: The key of the secret to select from. Must be a
                        valid secret key.
                      type: string
                    name:
                      description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                        TODO: Add other useful fields. apiVersion, kind, uid?'
                      type: string
                  required:
                  - name
                  type: object
                server:
                  description: Server is the ACME server URL
                  type: string
                skipTLSVerify:
                  description: If true, skip verifying the ACME server TLS certificate
                  type: boolean
                solvers:
                  description: Solvers is a list of challenge solvers that will be
                    used to solve ACME challenges for the matching domains.
                  items:
                    properties:
                      dns01:
                        properties:
                          acmedns:
                            description: ACMEIssuerDNS01ProviderAcmeDNS is a structure
                              containing the configuration for ACME-DNS servers
                            properties:
                              accountSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              host:
                                type: string
                            required:
                            - accountSecretRef
                            - host
                            type: object
                          akamai:
                            description: ACMEIssuerDNS01ProviderAkamai is a structure
                              containing the DNS configuration for Akamai DNS—Zone
                              Record Management API
                            properties:
                              accessTokenSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              clientSecretSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              clientTokenSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              serviceConsumerDomain:
                                type: string
                            required:
                            - accessTokenSecretRef
                            - clientSecretSecretRef
                            - clientTokenSecretRef
                            - serviceConsumerDomain
                            type: object
                          azuredns:
                            description: ACMEIssuerDNS01ProviderAzureDNS is a structure
                              containing the configuration for Azure DNS
                            properties:
                              clientID:
                                description: if both this and ClientSecret are left
                                  unset MSI will be used
                                type: string
                              clientSecretSecretRef:
                                description: if both this and ClientID are left unset
                                  MSI will be used
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              environment:
                                enum:
                                - AzurePublicCloud
                                - AzureChinaCloud
                                - AzureGermanCloud
                                - AzureUSGovernmentCloud
                                type: string
                              hostedZoneName:
                                type: string
                              resourceGroupName:
                                type: string
                              subscriptionID:
                                type: string
                              tenantID:
                                description: when specifying ClientID and ClientSecret
                                  then this field is also needed
                                type: string
                            required:
                            - resourceGroupName
                            - subscriptionID
                            type: object
                          clouddns:
                            description: ACMEIssuerDNS01ProviderCloudDNS is a structure
                              containing the DNS configuration for Google Cloud DNS
                            properties:
                              project:
                                type: string
                              serviceAccountSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                            required:
                            - project
                            type: object
                          cloudflare:
                            description: ACMEIssuerDNS01ProviderCloudflare is a structure
                              containing the DNS configuration for Cloudflare
                            properties:
                              apiKeySecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              apiTokenSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              email:
                                type: string
                            required:
                            - email
                            type: object
                          cnameStrategy:
                            description: CNAMEStrategy configures how the DNS01 provider
                              should handle CNAME records when found in DNS zones.
                            enum:
                            - None
                            - Follow
                            type: string
                          digitalocean:
                            description: ACMEIssuerDNS01ProviderDigitalOcean is a
                              structure containing the DNS configuration for DigitalOcean
                              Domains
                            properties:
                              tokenSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                            required:
                            - tokenSecretRef
                            type: object
                          rfc2136:
                            description: ACMEIssuerDNS01ProviderRFC2136 is a structure
                              containing the configuration for RFC2136 DNS
                            properties:
                              nameserver:
                                description: The IP address or hostname of an authoritative
                                  DNS server supporting RFC2136 in the form host:port.
                                  If the host is an IPv6 address it must be enclosed
                                  in square brackets (e.g [2001:db8::1]) ; port is
                                  optional. This field is required.
                                type: string
                              tsigAlgorithm:
                                description: 'The TSIG Algorithm configured in the
                                  DNS supporting RFC2136. Used only when tsigSecretSecretRef
                                  and tsigKeyName are defined. Supported values
                                  are (case-insensitive): HMACMD5 (default), HMACSHA1,
                                  HMACSHA256 or HMACSHA512.'
                                type: string
                              tsigKeyName:
                                description: The TSIG Key name configured in the DNS.
                                  If tsigSecretSecretRef is defined, this field
                                  is required.
                                type: string
                              tsigSecretSecretRef:
                                description: The name of the secret containing the
                                  TSIG value. If tsigKeyName is defined, this
                                  field is required.
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                            required:
                            - nameserver
                            type: object
                          route53:
                            description: ACMEIssuerDNS01ProviderRoute53 is a structure
                              containing the Route 53 configuration for AWS
                            properties:
                              accessKeyID:
                                description: 'The AccessKeyID is used for authentication.
                                  If not set we fall-back to using env vars, shared
                                  credentials file or AWS Instance metadata see: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials'
                                type: string
                              hostedZoneID:
                                description: If set, the provider will manage only
                                  this zone in Route53 and will not do an lookup using
                                  the route53:ListHostedZonesByName api call.
                                type: string
                              region:
                                description: Always set the region when using AccessKeyID
                                  and SecretAccessKey
                                type: string
                              role:
                                description: Role is a Role ARN which the Route53
                                  provider will assume using either the explicit credentials
                                  AccessKeyID/SecretAccessKey or the inferred credentials
                                  from environment variables, shared credentials file
                                  or AWS Instance metadata
                                type: string
                              secretAccessKeySecretRef:
                                description: The SecretAccessKey is used for authentication.
                                  If not set we fall-back to using env vars, shared
                                  credentials file or AWS Instance metadata https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                            required:
                            - region
                            type: object
                          webhook:
                            description: ACMEIssuerDNS01ProviderWebhook specifies
                              configuration for a webhook DNS01 provider, including
                              where to POST ChallengePayload resources.
                            properties:
                              config:
                                description: Additional configuration that should
                                  be passed to the webhook apiserver when challenges
                                  are processed. This can contain arbitrary JSON data.
                                  Secret values should not be specified in this stanza.
                                  If secret values are needed (e.g. credentials for
                                  a DNS service), you should use a SecretKeySelector
                                  to reference a Secret resource. For details on the
                                  schema of this field, consult the webhook provider
                                  implementation's documentation.
                              groupName:
                                description: The API group name that should be used
                                  when POSTing ChallengePayload resources to the webhook
                                  apiserver. This should be the same as the GroupName
                                  specified in the webhook provider implementation.
                                type: string
                              solverName:
                                description: The name of the solver to use, as defined
                                  in the webhook provider implementation. This will
                                  typically be the name of the provider, e.g. 'cloudflare'.
                                type: string
                            required:
                            - groupName
                            - solverName
                            type: object
                        type: object
                      http01:
                        description: ACMEChallengeSolverHTTP01 contains configuration
                          detailing how to solve HTTP01 challenges within a Kubernetes
                          cluster. Typically this is accomplished through creating
                          'routes' of some description that configure ingress controllers
                          to direct traffic to 'solver pods', which are responsible
                          for responding to the ACME server's HTTP requests.
                        properties:
                          ingress:
                            description: The ingress based HTTP01 challenge solver
                              will solve challenges by creating or modifying Ingress
                              resources in order to route requests for '/.well-known/acme-challenge/XYZ'
                              to 'challenge solver' pods that are provisioned by cert-manager
                              for each Challenge to be completed.
                            properties:
                              class:
                                description: The ingress class to use when creating
                                  Ingress resources to solve ACME challenges that
                                  use this challenge solver. Only one of 'class' or
                                  'name' may be specified.
                                type: string
                              ingressTemplate:
                                description: Optional ingress template used to configure
                                  the ACME challenge solver ingress used for HTTP01
                                  challenges
                                properties:
                                  metadata:
                                    description: ObjectMeta overrides for the ingress
                                      used to solve HTTP01 challenges. Only the 'labels'
                                      and 'annotations' fields may be set. If labels
                                      or annotations overlap with in-built values,
                                      the values here will override the in-built values.
                                    properties:
                                      annotations:
                                        additionalProperties:
                                          type: string
                                        description: Annotations that should be added
                                          to the created ACME HTTP01 solver ingress.
                                        type: object
                                      labels:
                                        additionalProperties:
                                          type: string
                                        description: Labels that should be added to
                                          the created ACME HTTP01 solver ingress.
                                        type: object
                                    type: object
                                type: object
                              name:
                                description: The name of the ingress resource that
                                  should have ACME challenge solving routes inserted
                                  into it in order to solve HTTP01 challenges. This
                                  is typically used in conjunction with ingress controllers
                                  like ingress-gce, which maintains a 1:1 mapping
                                  between external IPs and ingress resources.
                                type: string
                              podTemplate:
                                description: Optional pod template used to configure
                                  the ACME challenge solver pods used for HTTP01 challenges
                                properties:
                                  metadata:
                                    description: ObjectMeta overrides for the pod
                                      used to solve HTTP01 challenges. Only the 'labels'
                                      and 'annotations' fields may be set. If labels
                                      or annotations overlap with in-built values,
                                      the values here will override the in-built values.
                                    properties:
                                      annotations:
                                        additionalProperties:
                                          type: string
                                        description: Annotations that should be added
                                          to the create ACME HTTP01 solver pods.
                                        type: object
                                      labels:
                                        additionalProperties:
                                          type: string
                                        description: Labels that should be added to
                                          the created ACME HTTP01 solver pods.
                                        type: object
                                    type: object
                                  spec:
                                    description: PodSpec defines overrides for the
                                      HTTP01 challenge solver pod. Only the 'nodeSelector',
                                      'affinity' and 'tolerations' fields are supported
                                      currently. All other fields will be ignored.
                                    properties:
                                      affinity:
                                        description: If specified, the pod's scheduling
                                          constraints
                                        properties:
                                          nodeAffinity:
                                            description: Describes node affinity scheduling
                                              rules for the pod.
                                            properties:
                                              preferredDuringSchedulingIgnoredDuringExecution:
                                                description: The scheduler will prefer
                                                  to schedule pods to nodes that satisfy
                                                  the affinity expressions specified
                                                  by this field, but it may choose
                                                  a node that violates one or more
                                                  of the expressions. The node that
                                                  is most preferred is the one with
                                                  the greatest sum of weights, i.e.
                                                  for each node that meets all of
                                                  the scheduling requirements (resource
                                                  request, requiredDuringScheduling
                                                  affinity expressions, etc.), compute
                                                  a sum by iterating through the elements
                                                  of this field and adding "weight"
                                                  to the sum if the node matches the
                                                  corresponding matchExpressions;
                                                  the node(s) with the highest sum
                                                  are the most preferred.
                                                items:
                                                  description: An empty preferred
                                                    scheduling term matches all objects
                                                    with implicit weight 0 (i.e. it's
                                                    a no-op). A null preferred scheduling
                                                    term matches no objects (i.e.
                                                    is also a no-op).
                                                  properties:
                                                    preference:
                                                      description: A node selector
                                                        term, associated with the
                                                        corresponding weight.
                                                      properties:
                                                        matchExpressions:
                                                          description: A list of node
                                                            selector requirements
                                                            by node's labels.
                                                          items:
                                                            description: A node selector
                                                              requirement is a selector
                                                              that contains values,
                                                              a key, and an operator
                                                              that relates the key
                                                              and values.
                                                            properties:
                                                              key:
                                                                description: The label
                                                                  key that the selector
                                                                  applies to.
                                                                type: string
                                                              operator:
                                                                description: Represents
                                                                  a key's relationship
                                                                  to a set of values.
                                                                  Valid operators
                                                                  are In, NotIn, Exists,
                                                                  DoesNotExist. Gt,
                                                                  and Lt.
                                                                type: string
                                                              values:
                                                                description: An array
                                                                  of string values.
                                                                  If the operator
                                                                  is In or NotIn,
                                                                  the values array
                                                                  must be non-empty.
                                                                  If the operator
                                                                  is Exists or DoesNotExist,
                                                                  the values array
                                                                  must be empty. If
                                                                  the operator is
                                                                  Gt or Lt, the values
                                                                  array must have
                                                                  a single element,
                                                                  which will be interpreted
                                                                  as an integer. This
                                                                  array is replaced
                                                                  during a strategic
                                                                  merge patch.
                                                                items:
                                                                  type: string
                                                                type: array
                                                            required:
                                                            - key
                                                            - operator
                                                            type: object
                                                          type: array
                                                        matchFields:
                                                          description: A list of node
                                                            selector requirements
                                                            by node's fields.
                                                          items:
                                                            description: A node selector
                                                              requirement is a selector
                                                              that contains values,
                                                              a key, and an operator
                                                              that relates the key
                                                              and values.
                                                            properties:
                                                              key:
                                                                description: The label
                                                                  key that the selector
                                                                  applies to.
                                                                type: string
                                                              operator:
                                                                description: Represents
                                                                  a key's relationship
                                                                  to a set of values.
                                                                  Valid operators
                                                                  are In, NotIn, Exists,
                                                                  DoesNotExist. Gt,
                                                                  and Lt.
                                                                type: string
                                                              values:
                                                                description: An array
                                                                  of string values.
                                                                  If the operator
                                                                  is In or NotIn,
                                                                  the values array
                                                                  must be non-empty.
                                                                  If the operator
                                                                  is Exists or DoesNotExist,
                                                                  the values array
                                                                  must be empty. If
                                                                  the operator is
                                                                  Gt or Lt, the values
                                                                  array must have
                                                                  a single element,
                                                                  which will be interpreted
                                                                  as an integer. This
                                                                  array is replaced
                                                                  during a strategic
                                                                  merge patch.
                                                                items:
                                                                  type: string
                                                                type: array
                                                            required:
                                                            - key
                                                            - operator
                                                            type: object
                                                          type: array
                                                      type: object
                                                    weight:
                                                      description: Weight associated
                                                        with matching the corresponding
                                                        nodeSelectorTerm, in the range
                                                        1-100.
                                                      format: int32
                                                      type: integer
                                                  required:
                                                  - preference
                                                  - weight
                                                  type: object
                                                type: array
                                              requiredDuringSchedulingIgnoredDuringExecution:
                                                description: If the affinity requirements
                                                  specified by this field are not
                                                  met at scheduling time, the pod
                                                  will not be scheduled onto the node.
                                                  If the affinity requirements specified
                                                  by this field cease to be met at
                                                  some point during pod execution
                                                  (e.g. due to an update), the system
                                                  may or may not try to eventually
                                                  evict the pod from its node.
                                                properties:
                                                  nodeSelectorTerms:
                                                    description: Required. A list
                                                      of node selector terms. The
                                                      terms are ORed.
                                                    items:
                                                      description: A null or empty
                                                        node selector term matches
                                                        no objects. The requirements
                                                        of them are ANDed. The TopologySelectorTerm
                                                        type implements a subset of
                                                        the NodeSelectorTerm.
                                                      properties:
                                                        matchExpressions:
                                                          description: A list of node
                                                            selector requirements
                                                            by node's labels.
                                                          items:
                                                            description: A node selector
                                                              requirement is a selector
                                                              that contains values,
                                                              a key, and an operator
                                                              that relates the key
                                                              and values.
                                                            properties:
                                                              key:
                                                                description: The label
                                                                  key that the selector
                                                                  applies to.
                                                                type: string
                                                              operator:
                                                                description: Represents
                                                                  a key's relationship
                                                                  to a set of values.
                                                                  Valid operators
                                                                  are In, NotIn, Exists,
                                                                  DoesNotExist. Gt,
                                                                  and Lt.
                                                                type: string
                                                              values:
                                                                description: An array
                                                                  of string values.
                                                                  If the operator
                                                                  is In or NotIn,
                                                                  the values array
                                                                  must be non-empty.
                                                                  If the operator
                                                                  is Exists or DoesNotExist,
                                                                  the values array
                                                                  must be empty. If
                                                                  the operator is
                                                                  Gt or Lt, the values
                                                                  array must have
                                                                  a single element,
                                                                  which will be interpreted
                                                                  as an integer. This
                                                                  array is replaced
                                                                  during a strategic
                                                                  merge patch.
                                                                items:
                                                                  type: string
                                                                type: array
                                                            required:
                                                            - key
                                                            - operator
                                                            type: object
                                                          type: array
                                                        matchFields:
                                                          description: A list of node
                                                            selector requirements
                                                            by node's fields.
                                                          items:
                                                            description: A node selector
                                                              requirement is a selector
                                                              that contains values,
                                                              a key, and an operator
                                                              that relates the key
                                                              and values.
                                                            properties:
                                                              key:
                                                                description: The label
                                                                  key that the selector
                                                                  applies to.
                                                                type: string
                                                              operator:
                                                                description: Represents
                                                                  a key's relationship
                                                                  to a set of values.
                                                                  Valid operators
                                                                  are In, NotIn, Exists,
                                                                  DoesNotExist. Gt,
                                                                  and Lt.
                                                                type: string
                                                              values:
                                                                description: An array
                                                                  of string values.
                                                                  If the operator
                                                                  is In or NotIn,
                                                                  the values array
                                                                  must be non-empty.
                                                                  If the operator
                                                                  is Exists or DoesNotExist,
                                                                  the values array
                                                                  must be empty. If
                                                                  the operator is
                                                                  Gt or Lt, the values
                                                                  array must have
                                                                  a single element,
                                                                  which will be interpreted
                                                                  as an integer. This
                                                                  array is replaced
                                                                  during a strategic
                                                                  merge patch.
                                                                items:
                                                                  type: string
                                                                type: array
                                                            required:
                                                            - key
                                                            - operator
                                                            type: object
                                                          type: array
                                                      type: object
                                                    type: array
                                                required:
                                                - nodeSelectorTerms
                                                type: object
                                            type: object
                                          podAffinity:
                                            description: Describes pod affinity scheduling
                                              rules (e.g. co-locate this pod in the
                                              same node, zone, etc. as some other
                                              pod(s)).
                                            properties:
                                              preferredDuringSchedulingIgnoredDuringExecution:
                                                description: The scheduler will prefer
                                                  to schedule pods to nodes that satisfy
                                                  the affinity expressions specified
                                                  by this field, but it may choose
                                                  a node that violates one or more
                                                  of the expressions. The node that
                                                  is most preferred is the one with
                                                  the greatest sum of weights, i.e.
                                                  for each node that meets all of
                                                  the scheduling requirements (resource
                                                  request, requiredDuringScheduling
                                                  affinity expressions, etc.), compute
                                                  a sum by iterating through the elements
                                                  of this field and adding "weight"
                                                  to the sum if the node has pods
                                                  which matches the corresponding
                                                  podAffinityTerm; the node(s) with
                                                  the highest sum are the most preferred.
                                                items:
                                                  description: The weights of all
                                                    of the matched WeightedPodAffinityTerm
                                                    fields are added per-node to find
                                                    the most preferred node(s)
                                                  properties:
                                                    podAffinityTerm:
                                                      description: Required. A pod
                                                        affinity term, associated
                                                        with the corresponding weight.
                                                      properties:
                                                        labelSelector:
                                                          description: A label query
                                                            over a set of resources,
                                                            in this case pods.
                                                          properties:
                                                            matchExpressions:
                                                              description: matchExpressions
                                                                is a list of label
                                                                selector requirements.
                                                                The requirements are
                                                                ANDed.
                                                              items:
                                                                description: A label
                                                                  selector requirement
                                                                  is a selector that
                                                                  contains values,
                                                                  a key, and an operator
                                                                  that relates the
                                                                  key and values.
                                                                properties:
                                                                  key:
                                                                    description: key
                                                                      is the label
                                                                      key that the
                                                                      selector applies
                                                                      to.
                                                                    type: string
                                                                  operator:
                                                                    description: operator
                                                                      represents a
                                                                      key's relationship
                                                                      to a set of
                                                                      values. Valid
                                                                      operators are
                                                                      In, NotIn, Exists
                                                                      and DoesNotExist.
                                                                    type: string
                                                                  values:
                                                                    description: values
                                                                      is an array
                                                                      of string values.
                                                                      If the operator
                                                                      is In or NotIn,
                                                                      the values array
                                                                      must be non-empty.
                                                                      If the operator
                                                                      is Exists or
                                                                      DoesNotExist,
                                                                      the values array
                                                                      must be empty.
                                                                      This array is
                                                                      replaced during
                                                                      a strategic
                                                                      merge patch.
                                                                    items:
                                                                      type: string
                                                                    type: array
                                                                required:
                                                                - key
                                                                - operator
                                                                type: object
                                                              type: array
                                                            matchLabels:
                                                              additionalProperties:
                                                                type: string
                                                              description: matchLabels
                                                                is a map of {key,value}
                                                                pairs. A single {key,value}
                                                                in the matchLabels
                                                                map is equivalent
                                                                to an element of matchExpressions,
                                                                whose key field is
                                                                "key", the operator
                                                                is "In", and the values
                                                                array contains only
                                                                "value". The requirements
                                                                are ANDed.
                                                              type: object
                                                          type: object
                                                        namespaces:
                                                          description: namespaces
                                                            specifies which namespaces
                                                            the labelSelector applies
                                                            to (matches against);
                                                            null or empty list means
                                                            "this pod's namespace"
                                                          items:
                                                            type: string
                                                          type: array
                                                        topologyKey:
                                                          description: This pod should
                                                            be co-located (affinity)
                                                            or not co-located (anti-affinity)
                                                            with the pods matching
                                                            the labelSelector in the
                                                            specified namespaces,
                                                            where co-located is defined
                                                            as running on a node whose
                                                            value of the label with
                                                            key topologyKey matches
                                                            that of any node on which
                                                            any of the selected pods
                                                            is running. Empty topologyKey
                                                            is not allowed.
                                                          type: string
                                                      required:
                                                      - topologyKey
                                                      type: object
                                                    weight:
                                                      description: weight associated
                                                        with matching the corresponding
                                                        podAffinityTerm, in the range
                                                        1-100.
                                                      format: int32
                                                      type: integer
                                                  required:
                                                  - podAffinityTerm
                                                  - weight
                                                  type: object
                                                type: array
                                              requiredDuringSchedulingIgnoredDuringExecution:
                                                description: If the affinity requirements
                                                  specified by this field are not
                                                  met at scheduling time, the pod
                                                  will not be scheduled onto the node.
                                                  If the affinity requirements specified
                                                  by this field cease to be met at
                                                  some point during pod execution
                                                  (e.g. due to a pod label update),
                                                  the system may or may not try to
                                                  eventually evict the pod from its
                                                  node. When there are multiple elements,
                                                  the lists of nodes corresponding
                                                  to each podAffinityTerm are intersected,
                                                  i.e. all terms must be satisfied.
                                                items:
                                                  description: Defines a set of pods
                                                    (namely those matching the labelSelector
                                                    relative to the given namespace(s))
                                                    that this pod should be co-located
                                                    (affinity) or not co-located (anti-affinity)
                                                    with, where co-located is defined
                                                    as running on a node whose value
                                                    of the label with key <topologyKey>
                                                    matches that of any node on which
                                                    a pod of the set of pods is running
                                                  properties:
                                                    labelSelector:
                                                      description: A label query over
                                                        a set of resources, in this
                                                        case pods.
                                                      properties:
                                                        matchExpressions:
                                                          description: matchExpressions
                                                            is a list of label selector
                                                            requirements. The requirements
                                                            are ANDed.
                                                          items:
                                                            description: A label selector
                                                              requirement is a selector
                                                              that contains values,
                                                              a key, and an operator
                                                              that relates the key
                                                              and values.
                                                            properties:
                                                              key:
                                                                description: key is
                                                                  the label key that
                                                                  the selector applies
                                                                  to.
                                                                type: string
                                                              operator:
                                                                description: operator
                                                                  represents a key's
                                                                  relationship to
                                                                  a set of values.
                                                                  Valid operators
                                                                  are In, NotIn, Exists
                                                                  and DoesNotExist.
                                                                type: string
                                                              values:
                                                                description: values
                                                                  is an array of string
                                                                  values. If the operator
                                                                  is In or NotIn,
                                                                  the values array
                                                                  must be non-empty.
                                                                  If the operator
                                                                  is Exists or DoesNotExist,
                                                                  the values array
                                                                  must be empty. This
                                                                  array is replaced
                                                                  during a strategic
                                                                  merge patch.
                                                                items:
                                                                  type: string
                                                                type: array
                                                            required:
                                                            - key
                                                            - operator
                                                            type: object
                                                          type: array
                                                        matchLabels:
                                                          additionalProperties:
                                                            type: string
                                                          description: matchLabels
                                                            is a map of {key,value}
                                                            pairs. A single {key,value}
                                                            in the matchLabels map
                                                            is equivalent to an element
                                                            of matchExpressions, whose
                                                            key field is "key", the
                                                            operator is "In", and
                                                            the values array contains
                                                            only "value". The requirements
                                                            are ANDed.
                                                          type: object
                                                      type: object
                                                    namespaces:
                                                      description: namespaces specifies
                                                        which namespaces the labelSelector
                                                        applies to (matches against);
                                                        null or empty list means "this
                                                        pod's namespace"
                                                      items:
                                                        type: string
                                                      type: array
                                                    topologyKey:
                                                      description: This pod should
                                                        be co-located (affinity) or
                                                        not co-located (anti-affinity)
                                                        with the pods matching the
                                                        labelSelector in the specified
                                                        namespaces, where co-located
                                                        is defined as running on a
                                                        node whose value of the label
                                                        with key topologyKey matches
                                                        that of any node on which
                                                        any of the selected pods is
                                                        running. Empty topologyKey
                                                        is not allowed.
                                                      type: string
                                                  required:
                                                  - topologyKey
                                                  type: object
                                                type: array
                                            type: object
                                          podAntiAffinity:
                                            description: Describes pod anti-affinity
                                              scheduling rules (e.g. avoid putting
                                              this pod in the same node, zone, etc.
                                              as some other pod(s)).
                                            properties:
                                              preferredDuringSchedulingIgnoredDuringExecution:
                                                description: The scheduler will prefer
                                                  to schedule pods to nodes that satisfy
                                                  the anti-affinity expressions specified
                                                  by this field, but it may choose
                                                  a node that violates one or more
                                                  of the expressions. The node that
                                                  is most preferred is the one with
                                                  the greatest sum of weights, i.e.
                                                  for each node that meets all of
                                                  the scheduling requirements (resource
                                                  request, requiredDuringScheduling
                                                  anti-affinity expressions, etc.),
                                                  compute a sum by iterating through
                                                  the elements of this field and adding
                                                  "weight" to the sum if the node
                                                  has pods which matches the corresponding
                                                  podAffinityTerm; the node(s) with
                                                  the highest sum are the most preferred.
                                                items:
                                                  description: The weights of all
                                                    of the matched WeightedPodAffinityTerm
                                                    fields are added per-node to find
                                                    the most preferred node(s)
                                                  properties:
                                                    podAffinityTerm:
                                                      description: Required. A pod
                                                        affinity term, associated
                                                        with the corresponding weight.
                                                      properties:
                                                        labelSelector:
                                                          description: A label query
                                                            over a set of resources,
                                                            in this case pods.
                                                          properties:
                                                            matchExpressions:
                                                              description: matchExpressions
                                                                is a list of label
                                                                selector requirements.
                                                                The requirements are
                                                                ANDed.
                                                              items:
                                                                description: A label
                                                                  selector requirement
                                                                  is a selector that
                                                                  contains values,
                                                                  a key, and an operator
                                                                  that relates the
                                                                  key and values.
                                                                properties:
                                                                  key:
                                                                    description: key
                                                                      is the label
                                                                      key that the
                                                                      selector applies
                                                                      to.
                                                                    type: string
                                                                  operator:
                                                                    description: operator
                                                                      represents a
                                                                      key's relationship
                                                                      to a set of
                                                                      values. Valid
                                                                      operators are
                                                                      In, NotIn, Exists
                                                                      and DoesNotExist.
                                                                    type: string
                                                                  values:
                                                                    description: values
                                                                      is an array
                                                                      of string values.
                                                                      If the operator
                                                                      is In or NotIn,
                                                                      the values array
                                                                      must be non-empty.
                                                                      If the operator
                                                                      is Exists or
                                                                      DoesNotExist,
                                                                      the values array
                                                                      must be empty.
                                                                      This array is
                                                                      replaced during
                                                                      a strategic
                                                                      merge patch.
                                                                    items:
                                                                      type: string
                                                                    type: array
                                                                required:
                                                                - key
                                                                - operator
                                                                type: object
                                                              type: array
                                                            matchLabels:
                                                              additionalProperties:
                                                                type: string
                                                              description: matchLabels
                                                                is a map of {key,value}
                                                                pairs. A single {key,value}
                                                                in the matchLabels
                                                                map is equivalent
                                                                to an element of matchExpressions,
                                                                whose key field is
                                                                "key", the operator
                                                                is "In", and the values
                                                                array contains only
                                                                "value". The requirements
                                                                are ANDed.
                                                              type: object
                                                          type: object
                                                        namespaces:
                                                          description: namespaces
                                                            specifies which namespaces
                                                            the labelSelector applies
                                                            to (matches against);
                                                            null or empty list means
                                                            "this pod's namespace"
                                                          items:
                                                            type: string
                                                          type: array
                                                        topologyKey:
                                                          description: This pod should
                                                            be co-located (affinity)
                                                            or not co-located (anti-affinity)
                                                            with the pods matching
                                                            the labelSelector in the
                                                            specified namespaces,
                                                            where co-located is defined
                                                            as running on a node whose
                                                            value of the label with
                                                            key topologyKey matches
                                                            that of any node on which
                                                            any of the selected pods
                                                            is running. Empty topologyKey
                                                            is not allowed.
                                                          type: string
                                                      required:
                                                      - topologyKey
                                                      type: object
                                                    weight:
                                                      description: weight associated
                                                        with matching the corresponding
                                                        podAffinityTerm, in the range
                                                        1-100.
                                                      format: int32
                                                      type: integer
                                                  required:
                                                  - podAffinityTerm
                                                  - weight
                                                  type: object
                                                type: array
                                              requiredDuringSchedulingIgnoredDuringExecution:
                                                description: If the anti-affinity
                                                  requirements specified by this field
                                                  are not met at scheduling time,
                                                  the pod will not be scheduled onto
                                                  the node. If the anti-affinity requirements
                                                  specified by this field cease to
                                                  be met at some point during pod
                                                  execution (e.g. due to a pod label
                                                  update), the system may or may not
                                                  try to eventually evict the pod
                                                  from its node. When there are multiple
                                                  elements, the lists of nodes corresponding
                                                  to each podAffinityTerm are intersected,
                                                  i.e. all terms must be satisfied.
                                                items:
                                                  description: Defines a set of pods
                                                    (namely those matching the labelSelector
                                                    relative to the given namespace(s))
                                                    that this pod should be co-located
                                                    (affinity) or not co-located (anti-affinity)
                                                    with, where co-located is defined
                                                    as running on a node whose value
                                                    of the label with key <topologyKey>
                                                    matches that of any node on which
                                                    a pod of the set of pods is running
                                                  properties:
                                                    labelSelector:
                                                      description: A label query over
                                                        a set of resources, in this
                                                        case pods.
                                                      properties:
                                                        matchExpressions:
                                                          description: matchExpressions
                                                            is a list of label selector
                                                            requirements. The requirements
                                                            are ANDed.
                                                          items:
                                                            description: A label selector
                                                              requirement is a selector
                                                              that contains values,
                                                              a key, and an operator
                                                              that relates the key
                                                              and values.
                                                            properties:
                                                              key:
                                                                description: key is
                                                                  the label key that
                                                                  the selector applies
                                                                  to.
                                                                type: string
                                                              operator:
                                                                description: operator
                                                                  represents a key's
                                                                  relationship to
                                                                  a set of values.
                                                                  Valid operators
                                                                  are In, NotIn, Exists
                                                                  and DoesNotExist.
                                                                type: string
                                                              values:
                                                                description: values
                                                                  is an array of string
                                                                  values. If the operator
                                                                  is In or NotIn,
                                                                  the values array
                                                                  must be non-empty.
                                                                  If the operator
                                                                  is Exists or DoesNotExist,
                                                                  the values array
                                                                  must be empty. This
                                                                  array is replaced
                                                                  during a strategic
                                                                  merge patch.
                                                                items:
                                                                  type: string
                                                                type: array
                                                            required:
                                                            - key
                                                            - operator
                                                            type: object
                                                          type: array
                                                        matchLabels:
                                                          additionalProperties:
                                                            type: string
                                                          description: matchLabels
                                                            is a map of {key,value}
                                                            pairs. A single {key,value}
                                                            in the matchLabels map
                                                            is equivalent to an element
                                                            of matchExpressions, whose
                                                            key field is "key", the
                                                            operator is "In", and
                                                            the values array contains
                                                            only "value". The requirements
                                                            are ANDed.
                                                          type: object
                                                      type: object
                                                    namespaces:
                                                      description: namespaces specifies
                                                        which namespaces the labelSelector
                                                        applies to (matches against);
                                                        null or empty list means "this
                                                        pod's namespace"
                                                      items:
                                                        type: string
                                                      type: array
                                                    topologyKey:
                                                      description: This pod should
                                                        be co-located (affinity) or
                                                        not co-located (anti-affinity)
                                                        with the pods matching the
                                                        labelSelector in the specified
                                                        namespaces, where co-located
                                                        is defined as running on a
                                                        node whose value of the label
                                                        with key topologyKey matches
                                                        that of any node on which
                                                        any of the selected pods is
                                                        running. Empty topologyKey
                                                        is not allowed.
                                                      type: string
                                                  required:
                                                  - topologyKey
                                                  type: object
                                                type: array
                                            type: object
                                        type: object
                                      nodeSelector:
                                        additionalProperties:
                                          type: string
                                        description: 'NodeSelector is a selector which
                                          must be true for the pod to fit on a node.
                                          Selector which must match a node''s labels
                                          for the pod to be scheduled on that node.
                                          More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/'
                                        type: object
                                      tolerations:
                                        description: If specified, the pod's tolerations.
                                        items:
                                          description: The pod this Toleration is
                                            attached to tolerates any taint that matches
                                            the triple <key,value,effect> using the
                                            matching operator <operator>.
                                          properties:
                                            effect:
                                              description: Effect indicates the taint
                                                effect to match. Empty means match
                                                all taint effects. When specified,
                                                allowed values are NoSchedule, PreferNoSchedule
                                                and NoExecute.
                                              type: string
                                            key:
                                              description: Key is the taint key that
                                                the toleration applies to. Empty means
                                                match all taint keys. If the key is
                                                empty, operator must be Exists; this
                                                combination means to match all values
                                                and all keys.
                                              type: string
                                            operator:
                                              description: Operator represents a key's
                                                relationship to the value. Valid operators
                                                are Exists and Equal. Defaults to
                                                Equal. Exists is equivalent to wildcard
                                                for value, so that a pod can tolerate
                                                all taints of a particular category.
                                              type: string
                                            tolerationSeconds:
                                              description: TolerationSeconds represents
                                                the period of time the toleration
                                                (which must be of effect NoExecute,
                                                otherwise this field is ignored) tolerates
                                                the taint. By default, it is not set,
                                                which means tolerate the taint forever
                                                (do not evict). Zero and negative
                                                values will be treated as 0 (evict
                                                immediately) by the system.
                                              format: int64
                                              type: integer
                                            value:
                                              description: Value is the taint value
                                                the toleration matches to. If the
                                                operator is Exists, the value should
                                                be empty, otherwise just a regular
                                                string.
                                              type: string
                                          type: object
                                        type: array
                                    type: object
                                type: object
                              serviceType:
                                description: Optional service type for Kubernetes
                                  solver service
                                type: string
                            type: object
                        type: object
                      selector:
                        description: Selector selects a set of DNSNames on the Certificate
                          resource that should be solved using this challenge solver.
                        properties:
                          dnsNames:
                            description: List of DNSNames that this solver will be
                              used to solve. If specified and a match is found, a
                              dnsNames selector will take precedence over a dnsZones
                              selector. If multiple solvers match with the same dnsNames
                              value, the solver with the most matching labels in matchLabels
                              will be selected. If neither has more matches, the solver
                              defined earlier in the list will be selected.
                            items:
                              type: string
                            type: array
                          dnsZones:
                            description: List of DNSZones that this solver will be
                              used to solve. The most specific DNS zone match specified
                              here will take precedence over other DNS zone matches,
                              so a solver specifying sys.example.com will be selected
                              over one specifying example.com for the domain www.sys.example.com.
                              If multiple solvers match with the same dnsZones value,
                              the solver with the most matching labels in matchLabels
                              will be selected. If neither has more matches, the solver
                              defined earlier in the list will be selected.
                            items:
                              type: string
                            type: array
                          matchLabels:
                            additionalProperties:
                              type: string
                            description: A label selector that is used to refine the
                              set of certificate's that this challenge solver will
                              apply to.
                            type: object
                        type: object
                    type: object
                  type: array
              required:
              - privateKeySecretRef
              - server
              type: object
            ca:
              properties:
                crlDistributionPoints:
                  description: The CRL distribution points is an X.509 v3 certificate
                    extension which identifies the location of the CRL from which
                    the revocation of this certificate can be checked. If not set
                    certificate will be issued without CDP. Values are strings.
                  items:
                    type: string
                  type: array
                secretName:
                  description: SecretName is the name of the secret used to sign Certificates
                    issued by this Issuer.
                  type: string
              required:
              - secretName
              type: object
            selfSigned:
              properties:
                crlDistributionPoints:
                  description: The CRL distribution points is an X.509 v3 certificate
                    extension which identifies the location of the CRL from which
                    the revocation of this certificate can be checked. If not set
                    certificate will be issued without CDP. Values are strings.
                  items:
                    type: string
                  type: array
              type: object
            vault:
              properties:
                auth:
                  description: Vault authentication
                  properties:
                    appRole:
                      description: This Secret contains a AppRole and Secret
                      properties:
                        path:
                          description: Where the authentication path is mounted in
                            Vault.
                          type: string
                        roleId:
                          type: string
                        secretRef:
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                      required:
                      - path
                      - roleId
                      - secretRef
                      type: object
                    kubernetes:
                      description: This contains a Role and Secret with a ServiceAccount
                        token to authenticate with vault.
                      properties:
                        mountPath:
                          description: The Vault mountPath here is the mount path
                            to use when authenticating with Vault. For example, setting
                            a value to /v1/auth/foo, will use the path /v1/auth/foo/login
                            to authenticate with Vault. If unspecified, the default
                            value "/v1/auth/kubernetes" will be used.
                          type: string
                        role:
                          description: A required field containing the Vault Role
                            to assume. A Role binds a Kubernetes ServiceAccount with
                            a set of Vault policies.
                          type: string
                        secretRef:
                          description: The required Secret field containing a Kubernetes
                            ServiceAccount JWT used for authenticating with Vault.
                            Use of 'ambient credentials' is not supported.
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                      required:
                      - role
                      - secretRef
                      type: object
                    tokenSecretRef:
                      description: This Secret contains the Vault token key
                      properties:
                        key:
                          description: The key of the secret to select from. Must
                            be a valid secret key.
                          type: string
                        name:
                          description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                            TODO: Add other useful fields. apiVersion, kind, uid?'
                          type: string
                      required:
                      - name
                      type: object
                  type: object
                caBundle:
                  description: Base64 encoded CA bundle to validate Vault server certificate.
                    Only used if the Server URL is using HTTPS protocol. This parameter
                    is ignored for plain HTTP protocol connection. If not set the
                    system root certificates are used to validate the TLS connection.
                  format: byte
                  type: string
                path:
                  description: Vault URL path to the certificate role
                  type: string
                server:
                  description: Server is the vault connection address
                  type: string
              required:
              - auth
              - path
              - server
              type: object
            venafi:
              description: VenafiIssuer describes issuer configuration details for
                Venafi Cloud.
              properties:
                cloud:
                  description: Cloud specifies the Venafi cloud configuration settings.
                    Only one of TPP or Cloud may be specified.
                  properties:
                    apiTokenSecretRef:
                      description: APITokenSecretRef is a secret key selector for
                        the Venafi Cloud API token.
                      properties:
                        key:
                          description: The key of the secret to select from. Must
                            be a valid secret key.
                          type: string
                        name:
                          description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                            TODO: Add other useful fields. apiVersion, kind, uid?'
                          type: string
                      required:
                      - name
                      type: object
                    url:
                      description: URL is the base URL for Venafi Cloud
                      type: string
                  required:
                  - apiTokenSecretRef
                  type: object
                tpp:
                  description: TPP specifies Trust Protection Platform configuration
                    settings. Only one of TPP or Cloud may be specified.
                  properties:
                    caBundle:
                      description: CABundle is a PEM encoded TLS certificate to use
                        to verify connections to the TPP instance. If specified, system
                        roots will not be used and the issuing CA for the TPP instance
                        must be verifiable using the provided root. If not specified,
                        the connection will be verified using the cert-manager system
                        root certificates.
                      format: byte
                      type: string
                    credentialsRef:
                      description: CredentialsRef is a reference to a Secret containing
                        the username and password for the TPP server. The secret must
                        contain two keys, 'username' and 'password'.
                      properties:
                        name:
                          description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                            TODO: Add other useful fields. apiVersion, kind, uid?'
                          type: string
                      required:
                      - name
                      type: object
                    url:
                      description: URL is the base URL for the Venafi TPP instance
                      type: string
                  required:
                  - credentialsRef
                  - url
                  type: object
                zone:
                  description: Zone is the Venafi Policy Zone to use for this issuer.
                    All requests made to the Venafi platform will be restricted by
                    the named zone policy. This field is required.
                  type: string
              required:
              - zone
              type: object
          type: object
        status:
          description: IssuerStatus contains status information about an Issuer
          properties:
            acme:
              properties:
                lastRegisteredEmail:
                  description: LastRegisteredEmail is the email associated with the
                    latest registered ACME account, in order to track changes made
                    to registered account associated with the  Issuer
                  type: string
                uri:
                  description: URI is the unique account identifier, which can also
                    be used to retrieve account details from the CA
                  type: string
              type: object
            conditions:
              items:
                description: IssuerCondition contains condition information for an
                  Issuer.
                properties:
                  lastTransitionTime:
                    description: LastTransitionTime is the timestamp corresponding
                      to the last status change of this condition.
                    format: date-time
                    type: string
                  message:
                    description: Message is a human readable description of the details
                      of the last transition, complementing reason.
                    type: string
                  reason:
                    description: Reason is a brief machine readable explanation for
                      the condition's last transition.
                    type: string
                  status:
                    description: Status of the condition, one of ('True', 'False',
                      'Unknown').
                    enum:
                    - "True"
                    - "False"
                    - Unknown
                    type: string
                  type:
                    description: Type of the condition, currently ('Ready').
                    type: string
                required:
                - status
                - type
                type: object
              type: array
          type: object
  versions:
  - name: v1alpha2
    served: true
    storage: true
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    cert-manager.io/inject-ca-from-secret: 'kube-system/cert-manager-webhook-ca'
  labels:
    app: 'cert-manager'
    app.kubernetes.io/instance: 'cert-manager'
    app.kubernetes.io/managed-by: 'Tiller'
    app.kubernetes.io/name: 'cert-manager'
    helm.sh/chart: 'cert-manager-v0.15.1'
  name: issuers.cert-manager.io
spec:
  additionalPrinterColumns:
  - JSONPath: .status.conditions[?(@.type=="Ready")].status
    name: Ready
    type: string
  - JSONPath: .status.conditions[?(@.type=="Ready")].message
    name: Status
    priority: 1
    type: string
  - JSONPath: .metadata.creationTimestamp
    description: CreationTimestamp is a timestamp representing the server time when
      this object was created. It is not guaranteed to be set in happens-before order
      across separate operations. Clients may not set this value. It is represented
      in RFC3339 form and is in UTC.
    name: Age
    type: date
  group: cert-manager.io
  names:
    kind: Issuer
    listKind: IssuerList
    plural: issuers
    singular: issuer
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: IssuerSpec is the specification of an Issuer. This includes
            any configuration required for the issuer.
          properties:
            acme:
              description: ACMEIssuer contains the specification for an ACME issuer
              properties:
                email:
                  description: Email is the email for this account
                  type: string
                externalAccountBinding:
                  description: ExternalAccountBinding is a reference to a CA external
                    account of the ACME server.
                  properties:
                    keyAlgorithm:
                      description: keyAlgorithm is the MAC key algorithm that the
                        key is used for. Valid values are "HS256", "HS384" and "HS512".
                      enum:
                      - HS256
                      - HS384
                      - HS512
                      type: string
                    keyID:
                      description: keyID is the ID of the CA key that the External
                        Account is bound to.
                      type: string
                    keySecretRef:
                      description: keySecretRef is a Secret Key Selector referencing
                        a data item in a Kubernetes Secret which holds the symmetric
                        MAC key of the External Account Binding. The key is the
                        index string that is paired with the key data in the Secret
                        and should not be confused with the key data itself, or indeed
                        with the External Account Binding keyID above. The secret
                        key stored in the Secret **must** be un-padded, base64 URL
                        encoded data.
                      properties:
                        key:
                          description: The key of the secret to select from. Must
                            be a valid secret key.
                          type: string
                        name:
                          description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                            TODO: Add other useful fields. apiVersion, kind, uid?'
                          type: string
                      required:
                      - name
                      type: object
                  required:
                  - keyAlgorithm
                  - keyID
                  - keySecretRef
                  type: object
                privateKeySecretRef:
                  description: PrivateKey is the name of a secret containing the private
                    key for this user account.
                  properties:
                    key:
                      description: The key of the secret to select from. Must be a
                        valid secret key.
                      type: string
                    name:
                      description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                        TODO: Add other useful fields. apiVersion, kind, uid?'
                      type: string
                  required:
                  - name
                  type: object
                server:
                  description: Server is the ACME server URL
                  type: string
                skipTLSVerify:
                  description: If true, skip verifying the ACME server TLS certificate
                  type: boolean
                solvers:
                  description: Solvers is a list of challenge solvers that will be
                    used to solve ACME challenges for the matching domains.
                  items:
                    properties:
                      dns01:
                        properties:
                          acmedns:
                            description: ACMEIssuerDNS01ProviderAcmeDNS is a structure
                              containing the configuration for ACME-DNS servers
                            properties:
                              accountSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              host:
                                type: string
                            required:
                            - accountSecretRef
                            - host
                            type: object
                          akamai:
                            description: ACMEIssuerDNS01ProviderAkamai is a structure
                              containing the DNS configuration for Akamai DNS—Zone
                              Record Management API
                            properties:
                              accessTokenSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              clientSecretSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              clientTokenSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              serviceConsumerDomain:
                                type: string
                            required:
                            - accessTokenSecretRef
                            - clientSecretSecretRef
                            - clientTokenSecretRef
                            - serviceConsumerDomain
                            type: object
                          azuredns:
                            description: ACMEIssuerDNS01ProviderAzureDNS is a structure
                              containing the configuration for Azure DNS
                            properties:
                              clientID:
                                description: if both this and ClientSecret are left
                                  unset MSI will be used
                                type: string
                              clientSecretSecretRef:
                                description: if both this and ClientID are left unset
                                  MSI will be used
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              environment:
                                enum:
                                - AzurePublicCloud
                                - AzureChinaCloud
                                - AzureGermanCloud
                                - AzureUSGovernmentCloud
                                type: string
                              hostedZoneName:
                                type: string
                              resourceGroupName:
                                type: string
                              subscriptionID:
                                type: string
                              tenantID:
                                description: when specifying ClientID and ClientSecret
                                  then this field is also needed
                                type: string
                            required:
                            - resourceGroupName
                            - subscriptionID
                            type: object
                          clouddns:
                            description: ACMEIssuerDNS01ProviderCloudDNS is a structure
                              containing the DNS configuration for Google Cloud DNS
                            properties:
                              project:
                                type: string
                              serviceAccountSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                            required:
                            - project
                            type: object
                          cloudflare:
                            description: ACMEIssuerDNS01ProviderCloudflare is a structure
                              containing the DNS configuration for Cloudflare
                            properties:
                              apiKeySecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              apiTokenSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                              email:
                                type: string
                            required:
                            - email
                            type: object
                          cnameStrategy:
                            description: CNAMEStrategy configures how the DNS01 provider
                              should handle CNAME records when found in DNS zones.
                            enum:
                            - None
                            - Follow
                            type: string
                          digitalocean:
                            description: ACMEIssuerDNS01ProviderDigitalOcean is a
                              structure containing the DNS configuration for DigitalOcean
                              Domains
                            properties:
                              tokenSecretRef:
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                            required:
                            - tokenSecretRef
                            type: object
                          rfc2136:
                            description: ACMEIssuerDNS01ProviderRFC2136 is a structure
                              containing the configuration for RFC2136 DNS
                            properties:
                              nameserver:
                                description: The IP address or hostname of an authoritative
                                  DNS server supporting RFC2136 in the form host:port.
                                  If the host is an IPv6 address it must be enclosed
                                  in square brackets (e.g [2001:db8::1]) ; port is
                                  optional. This field is required.
                                type: string
                              tsigAlgorithm:
                                description: 'The TSIG Algorithm configured in the
                                  DNS supporting RFC2136. Used only when tsigSecretSecretRef
                                  and tsigKeyName are defined. Supported values
                                  are (case-insensitive): HMACMD5 (default), HMACSHA1,
                                  HMACSHA256 or HMACSHA512.'
                                type: string
                              tsigKeyName:
                                description: The TSIG Key name configured in the DNS.
                                  If tsigSecretSecretRef is defined, this field
                                  is required.
                                type: string
                              tsigSecretSecretRef:
                                description: The name of the secret containing the
                                  TSIG value. If tsigKeyName is defined, this
                                  field is required.
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                            required:
                            - nameserver
                            type: object
                          route53:
                            description: ACMEIssuerDNS01ProviderRoute53 is a structure
                              containing the Route 53 configuration for AWS
                            properties:
                              accessKeyID:
                                description: 'The AccessKeyID is used for authentication.
                                  If not set we fall-back to using env vars, shared
                                  credentials file or AWS Instance metadata see: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials'
                                type: string
                              hostedZoneID:
                                description: If set, the provider will manage only
                                  this zone in Route53 and will not do an lookup using
                                  the route53:ListHostedZonesByName api call.
                                type: string
                              region:
                                description: Always set the region when using AccessKeyID
                                  and SecretAccessKey
                                type: string
                              role:
                                description: Role is a Role ARN which the Route53
                                  provider will assume using either the explicit credentials
                                  AccessKeyID/SecretAccessKey or the inferred credentials
                                  from environment variables, shared credentials file
                                  or AWS Instance metadata
                                type: string
                              secretAccessKeySecretRef:
                                description: The SecretAccessKey is used for authentication.
                                  If not set we fall-back to using env vars, shared
                                  credentials file or AWS Instance metadata https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials
                                properties:
                                  key:
                                    description: The key of the secret to select from.
                                      Must be a valid secret key.
                                    type: string
                                  name:
                                    description: 'Name of the referent. More info:
                                      https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                      TODO: Add other useful fields. apiVersion, kind,
                                      uid?'
                                    type: string
                                required:
                                - name
                                type: object
                            required:
                            - region
                            type: object
                          webhook:
                            description: ACMEIssuerDNS01ProviderWebhook specifies
                              configuration for a webhook DNS01 provider, including
                              where to POST ChallengePayload resources.
                            properties:
                              config:
                                description: Additional configuration that should
                                  be passed to the webhook apiserver when challenges
                                  are processed. This can contain arbitrary JSON data.
                                  Secret values should not be specified in this stanza.
                                  If secret values are needed (e.g. credentials for
                                  a DNS service), you should use a SecretKeySelector
                                  to reference a Secret resource. For details on the
                                  schema of this field, consult the webhook provider
                                  implementation's documentation.
                              groupName:
                                description: The API group name that should be used
                                  when POSTing ChallengePayload resources to the webhook
                                  apiserver. This should be the same as the GroupName
                                  specified in the webhook provider implementation.
                                type: string
                              solverName:
                                description: The name of the solver to use, as defined
                                  in the webhook provider implementation. This will
                                  typically be the name of the provider, e.g. 'cloudflare'.
                                type: string
                            required:
                            - groupName
                            - solverName
                            type: object
                        type: object
                      http01:
                        description: ACMEChallengeSolverHTTP01 contains configuration
                          detailing how to solve HTTP01 challenges within a Kubernetes
                          cluster. Typically this is accomplished through creating
                          'routes' of some description that configure ingress controllers
                          to direct traffic to 'solver pods', which are responsible
                          for responding to the ACME server's HTTP requests.
                        properties:
                          ingress:
                            description: The ingress based HTTP01 challenge solver
                              will solve challenges by creating or modifying Ingress
                              resources in order to route requests for '/.well-known/acme-challenge/XYZ'
                              to 'challenge solver' pods that are provisioned by cert-manager
                              for each Challenge to be completed.
                            properties:
                              class:
                                description: The ingress class to use when creating
                                  Ingress resources to solve ACME challenges that
                                  use this challenge solver. Only one of 'class' or
                                  'name' may be specified.
                                type: string
                              ingressTemplate:
                                description: Optional ingress template used to configure
                                  the ACME challenge solver ingress used for HTTP01
                                  challenges
                                properties:
                                  metadata:
                                    description: ObjectMeta overrides for the ingress
                                      used to solve HTTP01 challenges. Only the 'labels'
                                      and 'annotations' fields may be set. If labels
                                      or annotations overlap with in-built values,
                                      the values here will override the in-built values.
                                    properties:
                                      annotations:
                                        additionalProperties:
                                          type: string
                                        description: Annotations that should be added
                                          to the created ACME HTTP01 solver ingress.
                                        type: object
                                      labels:
                                        additionalProperties:
                                          type: string
                                        description: Labels that should be added to
                                          the created ACME HTTP01 solver ingress.
                                        type: object
                                    type: object
                                type: object
                              name:
                                description: The name of the ingress resource that
                                  should have ACME challenge solving routes inserted
                                  into it in order to solve HTTP01 challenges. This
                                  is typically used in conjunction with ingress controllers
                                  like ingress-gce, which maintains a 1:1 mapping
                                  between external IPs and ingress resources.
                                type: string
                              podTemplate:
                                description: Optional pod template used to configure
                                  the ACME challenge solver pods used for HTTP01 challenges
                                properties:
                                  metadata:
                                    description: ObjectMeta overrides for the pod
                                      used to solve HTTP01 challenges. Only the 'labels'
                                      and 'annotations' fields may be set. If labels
                                      or annotations overlap with in-built values,
                                      the values here will override the in-built values.
                                    properties:
                                      annotations:
                                        additionalProperties:
                                          type: string
                                        description: Annotations that should be added
                                          to the create ACME HTTP01 solver pods.
                                        type: object
                                      labels:
                                        additionalProperties:
                                          type: string
                                        description: Labels that should be added to
                                          the created ACME HTTP01 solver pods.
                                        type: object
                                    type: object
                                  spec:
                                    description: PodSpec defines overrides for the
                                      HTTP01 challenge solver pod. Only the 'nodeSelector',
                                      'affinity' and 'tolerations' fields are supported
                                      currently. All other fields will be ignored.
                                    properties:
                                      affinity:
                                        description: If specified, the pod's scheduling
                                          constraints
                                        properties:
                                          nodeAffinity:
                                            description: Describes node affinity scheduling
                                              rules for the pod.
                                            properties:
                                              preferredDuringSchedulingIgnoredDuringExecution:
                                                description: The scheduler will prefer
                                                  to schedule pods to nodes that satisfy
                                                  the affinity expressions specified
                                                  by this field, but it may choose
                                                  a node that violates one or more
                                                  of the expressions. The node that
                                                  is most preferred is the one with
                                                  the greatest sum of weights, i.e.
                                                  for each node that meets all of
                                                  the scheduling requirements (resource
                                                  request, requiredDuringScheduling
                                                  affinity expressions, etc.), compute
                                                  a sum by iterating through the elements
                                                  of this field and adding "weight"
                                                  to the sum if the node matches the
                                                  corresponding matchExpressions;
                                                  the node(s) with the highest sum
                                                  are the most preferred.
                                                items:
                                                  description: An empty preferred
                                                    scheduling term matches all objects
                                                    with implicit weight 0 (i.e. it's
                                                    a no-op). A null preferred scheduling
                                                    term matches no objects (i.e.
                                                    is also a no-op).
                                                  properties:
                                                    preference:
                                                      description: A node selector
                                                        term, associated with the
                                                        corresponding weight.
                                                      properties:
                                                        matchExpressions:
                                                          description: A list of node
                                                            selector requirements
                                                            by node's labels.
                                                          items:
                                                            description: A node selector
                                                              requirement is a selector
                                                              that contains values,
                                                              a key, and an operator
                                                              that relates the key
                                                              and values.
                                                            properties:
                                                              key:
                                                                description: The label
                                                                  key that the selector
                                                                  applies to.
                                                                type: string
                                                              operator:
                                                                description: Represents
                                                                  a key's relationship
                                                                  to a set of values.
                                                                  Valid operators
                                                                  are In, NotIn, Exists,
                                                                  DoesNotExist. Gt,
                                                                  and Lt.
                                                                type: string
                                                              values:
                                                                description: An array
                                                                  of string values.
                                                                  If the operator
                                                                  is In or NotIn,
                                                                  the values array
                                                                  must be non-empty.
                                                                  If the operator
                                                                  is Exists or DoesNotExist,
                                                                  the values array
                                                                  must be empty. If
                                                                  the operator is
                                                                  Gt or Lt, the values
                                                                  array must have
                                                                  a single element,
                                                                  which will be interpreted
                                                                  as an integer. This
                                                                  array is replaced
                                                                  during a strategic
                                                                  merge patch.
                                                                items:
                                                                  type: string
                                                                type: array
                                                            required:
                                                            - key
                                                            - operator
                                                            type: object
                                                          type: array
                                                        matchFields:
                                                          description: A list of node
                                                            selector requirements
                                                            by node's fields.
                                                          items:
                                                            description: A node selector
                                                              requirement is a selector
                                                              that contains values,
                                                              a key, and an operator
                                                              that relates the key
                                                              and values.
                                                            properties:
                                                              key:
                                                                description: The label
                                                                  key that the selector
                                                                  applies to.
                                                                type: string
                                                              operator:
                                                                description: Represents
                                                                  a key's relationship
                                                                  to a set of values.
                                                                  Valid operators
                                                                  are In, NotIn, Exists,
                                                                  DoesNotExist. Gt,
                                                                  and Lt.
                                                                type: string
                                                              values:
                                                                description: An array
                                                                  of string values.
                                                                  If the operator
                                                                  is In or NotIn,
                                                                  the values array
                                                                  must be non-empty.
                                                                  If the operator
                                                                  is Exists or DoesNotExist,
                                                                  the values array
                                                                  must be empty. If
                                                                  the operator is
                                                                  Gt or Lt, the values
                                                                  array must have
                                                                  a single element,
                                                                  which will be interpreted
                                                                  as an integer. This
                                                                  array is replaced
                                                                  during a strategic
                                                                  merge patch.
                                                                items:
                                                                  type: string
                                                                type: array
                                                            required:
                                                            - key
                                                            - operator
                                                            type: object
                                                          type: array
                                                      type: object
                                                    weight:
                                                      description: Weight associated
                                                        with matching the corresponding
                                                        nodeSelectorTerm, in the range
                                                        1-100.
                                                      format: int32
                                                      type: integer
                                                  required:
                                                  - preference
                                                  - weight
                                                  type: object
                                                type: array
                                              requiredDuringSchedulingIgnoredDuringExecution:
                                                description: If the affinity requirements
                                                  specified by this field are not
                                                  met at scheduling time, the pod
                                                  will not be scheduled onto the node.
                                                  If the affinity requirements specified
                                                  by this field cease to be met at
                                                  some point during pod execution
                                                  (e.g. due to an update), the system
                                                  may or may not try to eventually
                                                  evict the pod from its node.
                                                properties:
                                                  nodeSelectorTerms:
                                                    description: Required. A list
                                                      of node selector terms. The
                                                      terms are ORed.
                                                    items:
                                                      description: A null or empty
                                                        node selector term matches
                                                        no objects. The requirements
                                                        of them are ANDed. The TopologySelectorTerm
                                                        type implements a subset of
                                                        the NodeSelectorTerm.
                                                      properties:
                                                        matchExpressions:
                                                          description: A list of node
                                                            selector requirements
                                                            by node's labels.
                                                          items:
                                                            description: A node selector
                                                              requirement is a selector
                                                              that contains values,
                                                              a key, and an operator
                                                              that relates the key
                                                              and values.
                                                            properties:
                                                              key:
                                                                description: The label
                                                                  key that the selector
                                                                  applies to.
                                                                type: string
                                                              operator:
                                                                description: Represents
                                                                  a key's relationship
                                                                  to a set of values.
                                                                  Valid operators
                                                                  are In, NotIn, Exists,
                                                                  DoesNotExist. Gt,
                                                                  and Lt.
                                                                type: string
                                                              values:
                                                                description: An array
                                                                  of string values.
                                                                  If the operator
                                                                  is In or NotIn,
                                                                  the values array
                                                                  must be non-empty.
                                                                  If the operator
                                                                  is Exists or DoesNotExist,
                                                                  the values array
                                                                  must be empty. If
                                                                  the operator is
                                                                  Gt or Lt, the values
                                                                  array must have
                                                                  a single element,
                                                                  which will be interpreted
                                                                  as an integer. This
                                                                  array is replaced
                                                                  during a strategic
                                                                  merge patch.
                                                                items:
                                                                  type: string
                                                                type: array
                                                            required:
                                                            - key
                                                            - operator
                                                            type: object
                                                          type: array
                                                        matchFields:
                                                          description: A list of node
                                                            selector requirements
                                                            by node's fields.
                                                          items:
                                                            description: A node selector
                                                              requirement is a selector
                                                              that contains values,
                                                              a key, and an operator
                                                              that relates the key
                                                              and values.
                                                            properties:
                                                              key:
                                                                description: The label
                                                                  key that the selector
                                                                  applies to.
                                                                type: string
                                                              operator:
                                                                description: Represents
                                                                  a key's relationship
                                                                  to a set of values.
                                                                  Valid operators
                                                                  are In, NotIn, Exists,
                                                                  DoesNotExist. Gt,
                                                                  and Lt.
                                                                type: string
                                                              values:
                                                                description: An array
                                                                  of string values.
                                                                  If the operator
                                                                  is In or NotIn,
                                                                  the values array
                                                                  must be non-empty.
                                                                  If the operator
                                                                  is Exists or DoesNotExist,
                                                                  the values array
                                                                  must be empty. If
                                                                  the operator is
                                                                  Gt or Lt, the values
                                                                  array must have
                                                                  a single element,
                                                                  which will be interpreted
                                                                  as an integer. This
                                                                  array is replaced
                                                                  during a strategic
                                                                  merge patch.
                                                                items:
                                                                  type: string
                                                                type: array
                                                            required:
                                                            - key
                                                            - operator
                                                            type: object
                                                          type: array
                                                      type: object
                                                    type: array
                                                required:
                                                - nodeSelectorTerms
                                                type: object
                                            type: object
                                          podAffinity:
                                            description: Describes pod affinity scheduling
                                              rules (e.g. co-locate this pod in the
                                              same node, zone, etc. as some other
                                              pod(s)).
                                            properties:
                                              preferredDuringSchedulingIgnoredDuringExecution:
                                                description: The scheduler will prefer
                                                  to schedule pods to nodes that satisfy
                                                  the affinity expressions specified
                                                  by this field, but it may choose
                                                  a node that violates one or more
                                                  of the expressions. The node that
                                                  is most preferred is the one with
                                                  the greatest sum of weights, i.e.
                                                  for each node that meets all of
                                                  the scheduling requirements (resource
                                                  request, requiredDuringScheduling
                                                  affinity expressions, etc.), compute
                                                  a sum by iterating through the elements
                                                  of this field and adding "weight"
                                                  to the sum if the node has pods
                                                  which matches the corresponding
                                                  podAffinityTerm; the node(s) with
                                                  the highest sum are the most preferred.
                                                items:
                                                  description: The weights of all
                                                    of the matched WeightedPodAffinityTerm
                                                    fields are added per-node to find
                                                    the most preferred node(s)
                                                  properties:
                                                    podAffinityTerm:
                                                      description: Required. A pod
                                                        affinity term, associated
                                                        with the corresponding weight.
                                                      properties:
                                                        labelSelector:
                                                          description: A label query
                                                            over a set of resources,
                                                            in this case pods.
                                                          properties:
                                                            matchExpressions:
                                                              description: matchExpressions
                                                                is a list of label
                                                                selector requirements.
                                                                The requirements are
                                                                ANDed.
                                                              items:
                                                                description: A label
                                                                  selector requirement
                                                                  is a selector that
                                                                  contains values,
                                                                  a key, and an operator
                                                                  that relates the
                                                                  key and values.
                                                                properties:
                                                                  key:
                                                                    description: key
                                                                      is the label
                                                                      key that the
                                                                      selector applies
                                                                      to.
                                                                    type: string
                                                                  operator:
                                                                    description: operator
                                                                      represents a
                                                                      key's relationship
                                                                      to a set of
                                                                      values. Valid
                                                                      operators are
                                                                      In, NotIn, Exists
                                                                      and DoesNotExist.
                                                                    type: string
                                                                  values:
                                                                    description: values
                                                                      is an array
                                                                      of string values.
                                                                      If the operator
                                                                      is In or NotIn,
                                                                      the values array
                                                                      must be non-empty.
                                                                      If the operator
                                                                      is Exists or
                                                                      DoesNotExist,
                                                                      the values array
                                                                      must be empty.
                                                                      This array is
                                                                      replaced during
                                                                      a strategic
                                                                      merge patch.
                                                                    items:
                                                                      type: string
                                                                    type: array
                                                                required:
                                                                - key
                                                                - operator
                                                                type: object
                                                              type: array
                                                            matchLabels:
                                                              additionalProperties:
                                                                type: string
                                                              description: matchLabels
                                                                is a map of {key,value}
                                                                pairs. A single {key,value}
                                                                in the matchLabels
                                                                map is equivalent
                                                                to an element of matchExpressions,
                                                                whose key field is
                                                                "key", the operator
                                                                is "In", and the values
                                                                array contains only
                                                                "value". The requirements
                                                                are ANDed.
                                                              type: object
                                                          type: object
                                                        namespaces:
                                                          description: namespaces
                                                            specifies which namespaces
                                                            the labelSelector applies
                                                            to (matches against);
                                                            null or empty list means
                                                            "this pod's namespace"
                                                          items:
                                                            type: string
                                                          type: array
                                                        topologyKey:
                                                          description: This pod should
                                                            be co-located (affinity)
                                                            or not co-located (anti-affinity)
                                                            with the pods matching
                                                            the labelSelector in the
                                                            specified namespaces,
                                                            where co-located is defined
                                                            as running on a node whose
                                                            value of the label with
                                                            key topologyKey matches
                                                            that of any node on which
                                                            any of the selected pods
                                                            is running. Empty topologyKey
                                                            is not allowed.
                                                          type: string
                                                      required:
                                                      - topologyKey
                                                      type: object
                                                    weight:
                                                      description: weight associated
                                                        with matching the corresponding
                                                        podAffinityTerm, in the range
                                                        1-100.
                                                      format: int32
                                                      type: integer
                                                  required:
                                                  - podAffinityTerm
                                                  - weight
                                                  type: object
                                                type: array
                                              requiredDuringSchedulingIgnoredDuringExecution:
                                                description: If the affinity requirements
                                                  specified by this field are not
                                                  met at scheduling time, the pod
                                                  will not be scheduled onto the node.
                                                  If the affinity requirements specified
                                                  by this field cease to be met at
                                                  some point during pod execution
                                                  (e.g. due to a pod label update),
                                                  the system may or may not try to
                                                  eventually evict the pod from its
                                                  node. When there are multiple elements,
                                                  the lists of nodes corresponding
                                                  to each podAffinityTerm are intersected,
                                                  i.e. all terms must be satisfied.
                                                items:
                                                  description: Defines a set of pods
                                                    (namely those matching the labelSelector
                                                    relative to the given namespace(s))
                                                    that this pod should be co-located
                                                    (affinity) or not co-located (anti-affinity)
                                                    with, where co-located is defined
                                                    as running on a node whose value
                                                    of the label with key <topologyKey>
                                                    matches that of any node on which
                                                    a pod of the set of pods is running
                                                  properties:
                                                    labelSelector:
                                                      description: A label query over
                                                        a set of resources, in this
                                                        case pods.
                                                      properties:
                                                        matchExpressions:
                                                          description: matchExpressions
                                                            is a list of label selector
                                                            requirements. The requirements
                                                            are ANDed.
                                                          items:
                                                            description: A label selector
                                                              requirement is a selector
                                                              that contains values,
                                                              a key, and an operator
                                                              that relates the key
                                                              and values.
                                                            properties:
                                                              key:
                                                                description: key is
                                                                  the label key that
                                                                  the selector applies
                                                                  to.
                                                                type: string
                                                              operator:
                                                                description: operator
                                                                  represents a key's
                                                                  relationship to
                                                                  a set of values.
                                                                  Valid operators
                                                                  are In, NotIn, Exists
                                                                  and DoesNotExist.
                                                                type: string
                                                              values:
                                                                description: values
                                                                  is an array of string
                                                                  values. If the operator
                                                                  is In or NotIn,
                                                                  the values array
                                                                  must be non-empty.
                                                                  If the operator
                                                                  is Exists or DoesNotExist,
                                                                  the values array
                                                                  must be empty. This
                                                                  array is replaced
                                                                  during a strategic
                                                                  merge patch.
                                                                items:
                                                                  type: string
                                                                type: array
                                                            required:
                                                            - key
                                                            - operator
                                                            type: object
                                                          type: array
                                                        matchLabels:
                                                          additionalProperties:
                                                            type: string
                                                          description: matchLabels
                                                            is a map of {key,value}
                                                            pairs. A single {key,value}
                                                            in the matchLabels map
                                                            is equivalent to an element
                                                            of matchExpressions, whose
                                                            key field is "key", the
                                                            operator is "In", and
                                                            the values array contains
                                                            only "value". The requirements
                                                            are ANDed.
                                                          type: object
                                                      type: object
                                                    namespaces:
                                                      description: namespaces specifies
                                                        which namespaces the labelSelector
                                                        applies to (matches against);
                                                        null or empty list means "this
                                                        pod's namespace"
                                                      items:
                                                        type: string
                                                      type: array
                                                    topologyKey:
                                                      description: This pod should
                                                        be co-located (affinity) or
                                                        not co-located (anti-affinity)
                                                        with the pods matching the
                                                        labelSelector in the specified
                                                        namespaces, where co-located
                                                        is defined as running on a
                                                        node whose value of the label
                                                        with key topologyKey matches
                                                        that of any node on which
                                                        any of the selected pods is
                                                        running. Empty topologyKey
                                                        is not allowed.
                                                      type: string
                                                  required:
                                                  - topologyKey
                                                  type: object
                                                type: array
                                            type: object
                                          podAntiAffinity:
                                            description: Describes pod anti-affinity
                                              scheduling rules (e.g. avoid putting
                                              this pod in the same node, zone, etc.
                                              as some other pod(s)).
                                            properties:
                                              preferredDuringSchedulingIgnoredDuringExecution:
                                                description: The scheduler will prefer
                                                  to schedule pods to nodes that satisfy
                                                  the anti-affinity expressions specified
                                                  by this field, but it may choose
                                                  a node that violates one or more
                                                  of the expressions. The node that
                                                  is most preferred is the one with
                                                  the greatest sum of weights, i.e.
                                                  for each node that meets all of
                                                  the scheduling requirements (resource
                                                  request, requiredDuringScheduling
                                                  anti-affinity expressions, etc.),
                                                  compute a sum by iterating through
                                                  the elements of this field and adding
                                                  "weight" to the sum if the node
                                                  has pods which matches the corresponding
                                                  podAffinityTerm; the node(s) with
                                                  the highest sum are the most preferred.
                                                items:
                                                  description: The weights of all
                                                    of the matched WeightedPodAffinityTerm
                                                    fields are added per-node to find
                                                    the most preferred node(s)
                                                  properties:
                                                    podAffinityTerm:
                                                      description: Required. A pod
                                                        affinity term, associated
                                                        with the corresponding weight.
                                                      properties:
                                                        labelSelector:
                                                          description: A label query
                                                            over a set of resources,
                                                            in this case pods.
                                                          properties:
                                                            matchExpressions:
                                                              description: matchExpressions
                                                                is a list of label
                                                                selector requirements.
                                                                The requirements are
                                                                ANDed.
                                                              items:
                                                                description: A label
                                                                  selector requirement
                                                                  is a selector that
                                                                  contains values,
                                                                  a key, and an operator
                                                                  that relates the
                                                                  key and values.
                                                                properties:
                                                                  key:
                                                                    description: key
                                                                      is the label
                                                                      key that the
                                                                      selector applies
                                                                      to.
                                                                    type: string
                                                                  operator:
                                                                    description: operator
                                                                      represents a
                                                                      key's relationship
                                                                      to a set of
                                                                      values. Valid
                                                                      operators are
                                                                      In, NotIn, Exists
                                                                      and DoesNotExist.
                                                                    type: string
                                                                  values:
                                                                    description: values
                                                                      is an array
                                                                      of string values.
                                                                      If the operator
                                                                      is In or NotIn,
                                                                      the values array
                                                                      must be non-empty.
                                                                      If the operator
                                                                      is Exists or
                                                                      DoesNotExist,
                                                                      the values array
                                                                      must be empty.
                                                                      This array is
                                                                      replaced during
                                                                      a strategic
                                                                      merge patch.
                                                                    items:
                                                                      type: string
                                                                    type: array
                                                                required:
                                                                - key
                                                                - operator
                                                                type: object
                                                              type: array
                                                            matchLabels:
                                                              additionalProperties:
                                                                type: string
                                                              description: matchLabels
                                                                is a map of {key,value}
                                                                pairs. A single {key,value}
                                                                in the matchLabels
                                                                map is equivalent
                                                                to an element of matchExpressions,
                                                                whose key field is
                                                                "key", the operator
                                                                is "In", and the values
                                                                array contains only
                                                                "value". The requirements
                                                                are ANDed.
                                                              type: object
                                                          type: object
                                                        namespaces:
                                                          description: namespaces
                                                            specifies which namespaces
                                                            the labelSelector applies
                                                            to (matches against);
                                                            null or empty list means
                                                            "this pod's namespace"
                                                          items:
                                                            type: string
                                                          type: array
                                                        topologyKey:
                                                          description: This pod should
                                                            be co-located (affinity)
                                                            or not co-located (anti-affinity)
                                                            with the pods matching
                                                            the labelSelector in the
                                                            specified namespaces,
                                                            where co-located is defined
                                                            as running on a node whose
                                                            value of the label with
                                                            key topologyKey matches
                                                            that of any node on which
                                                            any of the selected pods
                                                            is running. Empty topologyKey
                                                            is not allowed.
                                                          type: string
                                                      required:
                                                      - topologyKey
                                                      type: object
                                                    weight:
                                                      description: weight associated
                                                        with matching the corresponding
                                                        podAffinityTerm, in the range
                                                        1-100.
                                                      format: int32
                                                      type: integer
                                                  required:
                                                  - podAffinityTerm
                                                  - weight
                                                  type: object
                                                type: array
                                              requiredDuringSchedulingIgnoredDuringExecution:
                                                description: If the anti-affinity
                                                  requirements specified by this field
                                                  are not met at scheduling time,
                                                  the pod will not be scheduled onto
                                                  the node. If the anti-affinity requirements
                                                  specified by this field cease to
                                                  be met at some point during pod
                                                  execution (e.g. due to a pod label
                                                  update), the system may or may not
                                                  try to eventually evict the pod
                                                  from its node. When there are multiple
                                                  elements, the lists of nodes corresponding
                                                  to each podAffinityTerm are intersected,
                                                  i.e. all terms must be satisfied.
                                                items:
                                                  description: Defines a set of pods
                                                    (namely those matching the labelSelector
                                                    relative to the given namespace(s))
                                                    that this pod should be co-located
                                                    (affinity) or not co-located (anti-affinity)
                                                    with, where co-located is defined
                                                    as running on a node whose value
                                                    of the label with key <topologyKey>
                                                    matches that of any node on which
                                                    a pod of the set of pods is running
                                                  properties:
                                                    labelSelector:
                                                      description: A label query over
                                                        a set of resources, in this
                                                        case pods.
                                                      properties:
                                                        matchExpressions:
                                                          description: matchExpressions
                                                            is a list of label selector
                                                            requirements. The requirements
                                                            are ANDed.
                                                          items:
                                                            description: A label selector
                                                              requirement is a selector
                                                              that contains values,
                                                              a key, and an operator
                                                              that relates the key
                                                              and values.
                                                            properties:
                                                              key:
                                                                description: key is
                                                                  the label key that
                                                                  the selector applies
                                                                  to.
                                                                type: string
                                                              operator:
                                                                description: operator
                                                                  represents a key's
                                                                  relationship to
                                                                  a set of values.
                                                                  Valid operators
                                                                  are In, NotIn, Exists
                                                                  and DoesNotExist.
                                                                type: string
                                                              values:
                                                                description: values
                                                                  is an array of string
                                                                  values. If the operator
                                                                  is In or NotIn,
                                                                  the values array
                                                                  must be non-empty.
                                                                  If the operator
                                                                  is Exists or DoesNotExist,
                                                                  the values array
                                                                  must be empty. This
                                                                  array is replaced
                                                                  during a strategic
                                                                  merge patch.
                                                                items:
                                                                  type: string
                                                                type: array
                                                            required:
                                                            - key
                                                            - operator
                                                            type: object
                                                          type: array
                                                        matchLabels:
                                                          additionalProperties:
                                                            type: string
                                                          description: matchLabels
                                                            is a map of {key,value}
                                                            pairs. A single {key,value}
                                                            in the matchLabels map
                                                            is equivalent to an element
                                                            of matchExpressions, whose
                                                            key field is "key", the
                                                            operator is "In", and
                                                            the values array contains
                                                            only "value". The requirements
                                                            are ANDed.
                                                          type: object
                                                      type: object
                                                    namespaces:
                                                      description: namespaces specifies
                                                        which namespaces the labelSelector
                                                        applies to (matches against);
                                                        null or empty list means "this
                                                        pod's namespace"
                                                      items:
                                                        type: string
                                                      type: array
                                                    topologyKey:
                                                      description: This pod should
                                                        be co-located (affinity) or
                                                        not co-located (anti-affinity)
                                                        with the pods matching the
                                                        labelSelector in the specified
                                                        namespaces, where co-located
                                                        is defined as running on a
                                                        node whose value of the label
                                                        with key topologyKey matches
                                                        that of any node on which
                                                        any of the selected pods is
                                                        running. Empty topologyKey
                                                        is not allowed.
                                                      type: string
                                                  required:
                                                  - topologyKey
                                                  type: object
                                                type: array
                                            type: object
                                        type: object
                                      nodeSelector:
                                        additionalProperties:
                                          type: string
                                        description: 'NodeSelector is a selector which
                                          must be true for the pod to fit on a node.
                                          Selector which must match a node''s labels
                                          for the pod to be scheduled on that node.
                                          More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/'
                                        type: object
                                      tolerations:
                                        description: If specified, the pod's tolerations.
                                        items:
                                          description: The pod this Toleration is
                                            attached to tolerates any taint that matches
                                            the triple <key,value,effect> using the
                                            matching operator <operator>.
                                          properties:
                                            effect:
                                              description: Effect indicates the taint
                                                effect to match. Empty means match
                                                all taint effects. When specified,
                                                allowed values are NoSchedule, PreferNoSchedule
                                                and NoExecute.
                                              type: string
                                            key:
                                              description: Key is the taint key that
                                                the toleration applies to. Empty means
                                                match all taint keys. If the key is
                                                empty, operator must be Exists; this
                                                combination means to match all values
                                                and all keys.
                                              type: string
                                            operator:
                                              description: Operator represents a key's
                                                relationship to the value. Valid operators
                                                are Exists and Equal. Defaults to
                                                Equal. Exists is equivalent to wildcard
                                                for value, so that a pod can tolerate
                                                all taints of a particular category.
                                              type: string
                                            tolerationSeconds:
                                              description: TolerationSeconds represents
                                                the period of time the toleration
                                                (which must be of effect NoExecute,
                                                otherwise this field is ignored) tolerates
                                                the taint. By default, it is not set,
                                                which means tolerate the taint forever
                                                (do not evict). Zero and negative
                                                values will be treated as 0 (evict
                                                immediately) by the system.
                                              format: int64
                                              type: integer
                                            value:
                                              description: Value is the taint value
                                                the toleration matches to. If the
                                                operator is Exists, the value should
                                                be empty, otherwise just a regular
                                                string.
                                              type: string
                                          type: object
                                        type: array
                                    type: object
                                type: object
                              serviceType:
                                description: Optional service type for Kubernetes
                                  solver service
                                type: string
                            type: object
                        type: object
                      selector:
                        description: Selector selects a set of DNSNames on the Certificate
                          resource that should be solved using this challenge solver.
                        properties:
                          dnsNames:
                            description: List of DNSNames that this solver will be
                              used to solve. If specified and a match is found, a
                              dnsNames selector will take precedence over a dnsZones
                              selector. If multiple solvers match with the same dnsNames
                              value, the solver with the most matching labels in matchLabels
                              will be selected. If neither has more matches, the solver
                              defined earlier in the list will be selected.
                            items:
                              type: string
                            type: array
                          dnsZones:
                            description: List of DNSZones that this solver will be
                              used to solve. The most specific DNS zone match specified
                              here will take precedence over other DNS zone matches,
                              so a solver specifying sys.example.com will be selected
                              over one specifying example.com for the domain www.sys.example.com.
                              If multiple solvers match with the same dnsZones value,
                              the solver with the most matching labels in matchLabels
                              will be selected. If neither has more matches, the solver
                              defined earlier in the list will be selected.
                            items:
                              type: string
                            type: array
                          matchLabels:
                            additionalProperties:
                              type: string
                            description: A label selector that is used to refine the
                              set of certificate's that this challenge solver will
                              apply to.
                            type: object
                        type: object
                    type: object
                  type: array
              required:
              - privateKeySecretRef
              - server
              type: object
            ca:
              properties:
                crlDistributionPoints:
                  description: The CRL distribution points is an X.509 v3 certificate
                    extension which identifies the location of the CRL from which
                    the revocation of this certificate can be checked. If not set
                    certificate will be issued without CDP. Values are strings.
                  items:
                    type: string
                  type: array
                secretName:
                  description: SecretName is the name of the secret used to sign Certificates
                    issued by this Issuer.
                  type: string
              required:
              - secretName
              type: object
            selfSigned:
              properties:
                crlDistributionPoints:
                  description: The CRL distribution points is an X.509 v3 certificate
                    extension which identifies the location of the CRL from which
                    the revocation of this certificate can be checked. If not set
                    certificate will be issued without CDP. Values are strings.
                  items:
                    type: string
                  type: array
              type: object
            vault:
              properties:
                auth:
                  description: Vault authentication
                  properties:
                    appRole:
                      description: This Secret contains a AppRole and Secret
                      properties:
                        path:
                          description: Where the authentication path is mounted in
                            Vault.
                          type: string
                        roleId:
                          type: string
                        secretRef:
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                      required:
                      - path
                      - roleId
                      - secretRef
                      type: object
                    kubernetes:
                      description: This contains a Role and Secret with a ServiceAccount
                        token to authenticate with vault.
                      properties:
                        mountPath:
                          description: The Vault mountPath here is the mount path
                            to use when authenticating with Vault. For example, setting
                            a value to /v1/auth/foo, will use the path /v1/auth/foo/login
                            to authenticate with Vault. If unspecified, the default
                            value "/v1/auth/kubernetes" will be used.
                          type: string
                        role:
                          description: A required field containing the Vault Role
                            to assume. A Role binds a Kubernetes ServiceAccount with
                            a set of Vault policies.
                          type: string
                        secretRef:
                          description: The required Secret field containing a Kubernetes
                            ServiceAccount JWT used for authenticating with Vault.
                            Use of 'ambient credentials' is not supported.
                          properties:
                            key:
                              description: The key of the secret to select from. Must
                                be a valid secret key.
                              type: string
                            name:
                              description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                                TODO: Add other useful fields. apiVersion, kind, uid?'
                              type: string
                          required:
                          - name
                          type: object
                      required:
                      - role
                      - secretRef
                      type: object
                    tokenSecretRef:
                      description: This Secret contains the Vault token key
                      properties:
                        key:
                          description: The key of the secret to select from. Must
                            be a valid secret key.
                          type: string
                        name:
                          description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                            TODO: Add other useful fields. apiVersion, kind, uid?'
                          type: string
                      required:
                      - name
                      type: object
                  type: object
                caBundle:
                  description: Base64 encoded CA bundle to validate Vault server certificate.
                    Only used if the Server URL is using HTTPS protocol. This parameter
                    is ignored for plain HTTP protocol connection. If not set the
                    system root certificates are used to validate the TLS connection.
                  format: byte
                  type: string
                path:
                  description: Vault URL path to the certificate role
                  type: string
                server:
                  description: Server is the vault connection address
                  type: string
              required:
              - auth
              - path
              - server
              type: object
            venafi:
              description: VenafiIssuer describes issuer configuration details for
                Venafi Cloud.
              properties:
                cloud:
                  description: Cloud specifies the Venafi cloud configuration settings.
                    Only one of TPP or Cloud may be specified.
                  properties:
                    apiTokenSecretRef:
                      description: APITokenSecretRef is a secret key selector for
                        the Venafi Cloud API token.
                      properties:
                        key:
                          description: The key of the secret to select from. Must
                            be a valid secret key.
                          type: string
                        name:
                          description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                            TODO: Add other useful fields. apiVersion, kind, uid?'
                          type: string
                      required:
                      - name
                      type: object
                    url:
                      description: URL is the base URL for Venafi Cloud
                      type: string
                  required:
                  - apiTokenSecretRef
                  type: object
                tpp:
                  description: TPP specifies Trust Protection Platform configuration
                    settings. Only one of TPP or Cloud may be specified.
                  properties:
                    caBundle:
                      description: CABundle is a PEM encoded TLS certificate to use
                        to verify connections to the TPP instance. If specified, system
                        roots will not be used and the issuing CA for the TPP instance
                        must be verifiable using the provided root. If not specified,
                        the connection will be verified using the cert-manager system
                        root certificates.
                      format: byte
                      type: string
                    credentialsRef:
                      description: CredentialsRef is a reference to a Secret containing
                        the username and password for the TPP server. The secret must
                        contain two keys, 'username' and 'password'.
                      properties:
                        name:
                          description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                            TODO: Add other useful fields. apiVersion, kind, uid?'
                          type: string
                      required:
                      - name
                      type: object
                    url:
                      description: URL is the base URL for the Venafi TPP instance
                      type: string
                  required:
                  - credentialsRef
                  - url
                  type: object
                zone:
                  description: Zone is the Venafi Policy Zone to use for this issuer.
                    All requests made to the Venafi platform will be restricted by
                    the named zone policy. This field is required.
                  type: string
              required:
              - zone
              type: object
          type: object
        status:
          description: IssuerStatus contains status information about an Issuer
          properties:
            acme:
              properties:
                lastRegisteredEmail:
                  description: LastRegisteredEmail is the email associated with the
                    latest registered ACME account, in order to track changes made
                    to registered account associated with the  Issuer
                  type: string
                uri:
                  description: URI is the unique account identifier, which can also
                    be used to retrieve account details from the CA
                  type: string
              type: object
            conditions:
              items:
                description: IssuerCondition contains condition information for an
                  Issuer.
                properties:
                  lastTransitionTime:
                    description: LastTransitionTime is the timestamp corresponding
                      to the last status change of this condition.
                    format: date-time
                    type: string
                  message:
                    description: Message is a human readable description of the details
                      of the last transition, complementing reason.
                    type: string
                  reason:
                    description: Reason is a brief machine readable explanation for
                      the condition's last transition.
                    type: string
                  status:
                    description: Status of the condition, one of ('True', 'False',
                      'Unknown').
                    enum:
                    - "True"
                    - "False"
                    - Unknown
                    type: string
                  type:
                    description: Type of the condition, currently ('Ready').
                    type: string
                required:
                - status
                - type
                type: object
              type: array
          type: object
  versions:
  - name: v1alpha2
    served: true
    storage: true
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    cert-manager.io/inject-ca-from-secret: 'kube-system/cert-manager-webhook-ca'
  labels:
    app: 'cert-manager'
    app.kubernetes.io/instance: 'cert-manager'
    app.kubernetes.io/managed-by: 'Tiller'
    app.kubernetes.io/name: 'cert-manager'
    helm.sh/chart: 'cert-manager-v0.15.1'
  name: orders.acme.cert-manager.io
spec:
  additionalPrinterColumns:
  - JSONPath: .status.state
    name: State
    type: string
  - JSONPath: .spec.issuerRef.name
    name: Issuer
    priority: 1
    type: string
  - JSONPath: .status.reason
    name: Reason
    priority: 1
    type: string
  - JSONPath: .metadata.creationTimestamp
    description: CreationTimestamp is a timestamp representing the server time when
      this object was created. It is not guaranteed to be set in happens-before order
      across separate operations. Clients may not set this value. It is represented
      in RFC3339 form and is in UTC.
    name: Age
    type: date
  group: acme.cert-manager.io
  names:
    kind: Order
    listKind: OrderList
    plural: orders
    singular: order
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Order is a type to represent an Order with an ACME server
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            commonName:
              description: CommonName is the common name as specified on the DER encoded
                CSR. If CommonName is not specified, the first DNSName specified will
                be used as the CommonName. At least one of CommonName or a DNSNames
                must be set. This field must match the corresponding field on the
                DER encoded CSR.
              type: string
            csr:
              description: Certificate signing request bytes in DER encoding. This
                will be used when finalizing the order. This field must be set on
                the order.
              format: byte
              type: string
            dnsNames:
              description: DNSNames is a list of DNS names that should be included
                as part of the Order validation process. If CommonName is not specified,
                the first DNSName specified will be used as the CommonName. At least
                one of CommonName or a DNSNames must be set. This field must match
                the corresponding field on the DER encoded CSR.
              items:
                type: string
              type: array
            issuerRef:
              description: IssuerRef references a properly configured ACME-type Issuer
                which should be used to create this Order. If the Issuer does not
                exist, processing will be retried. If the Issuer is not an 'ACME'
                Issuer, an error will be returned and the Order will be marked as
                failed.
              properties:
                group:
                  type: string
                kind:
                  type: string
                name:
                  type: string
              required:
              - name
              type: object
          required:
          - csr
          - issuerRef
          type: object
        status:
          properties:
            authorizations:
              description: Authorizations contains data returned from the ACME server
                on what authorizations must be completed in order to validate the
                DNS names specified on the Order.
              items:
                description: ACMEAuthorization contains data returned from the ACME
                  server on an authorization that must be completed in order validate
                  a DNS name on an ACME Order resource.
                properties:
                  challenges:
                    description: Challenges specifies the challenge types offered
                      by the ACME server. One of these challenge types will be selected
                      when validating the DNS name and an appropriate Challenge resource
                      will be created to perform the ACME challenge process.
                    items:
                      description: Challenge specifies a challenge offered by the
                        ACME server for an Order. An appropriate Challenge resource
                        can be created to perform the ACME challenge process.
                      properties:
                        token:
                          description: Token is the token that must be presented for
                            this challenge. This is used to compute the 'key' that
                            must also be presented.
                          type: string
                        type:
                          description: Type is the type of challenge being offered,
                            e.g. http-01, dns-01
                          type: string
                        url:
                          description: URL is the URL of this challenge. It can be
                            used to retrieve additional metadata about the Challenge
                            from the ACME server.
                          type: string
                      required:
                      - token
                      - type
                      - url
                      type: object
                    type: array
                  identifier:
                    description: Identifier is the DNS name to be validated as part
                      of this authorization
                    type: string
                  initialState:
                    description: InitialState is the initial state of the ACME authorization
                      when first fetched from the ACME server. If an Authorization
                      is already 'valid', the Order controller will not create a Challenge
                      resource for the authorization. This will occur when working
                      with an ACME server that enables 'authz reuse' (such as Let's
                      Encrypt's production endpoint). If not set and 'identifier'
                      is set, the state is assumed to be pending and a Challenge will
                      be created.
                    enum:
                    - valid
                    - ready
                    - pending
                    - processing
                    - invalid
                    - expired
                    - errored
                    type: string
                  url:
                    description: URL is the URL of the Authorization that must be
                      completed
                    type: string
                  wildcard:
                    description: Wildcard will be true if this authorization is for
                      a wildcard DNS name. If this is true, the identifier will be
                      the *non-wildcard* version of the DNS name. For example, if
                      '*.example.com' is the DNS name being validated, this field
                      will be 'true' and the 'identifier' field will be 'example.com'.
                    type: boolean
                required:
                - url
                type: object
              type: array
            certificate:
              description: Certificate is a copy of the PEM encoded certificate for
                this Order. This field will be populated after the order has been
                successfully finalized with the ACME server, and the order has transitioned
                to the 'valid' state.
              format: byte
              type: string
            failureTime:
              description: FailureTime stores the time that this order failed. This
                is used to influence garbage collection and back-off.
              format: date-time
              type: string
            finalizeURL:
              description: FinalizeURL of the Order. This is used to obtain certificates
                for this order once it has been completed.
              type: string
            reason:
              description: Reason optionally provides more information about a why
                the order is in the current state.
              type: string
            state:
              description: State contains the current state of this Order resource.
                States 'success' and 'expired' are 'final'
              enum:
              - valid
              - ready
              - pending
              - processing
              - invalid
              - expired
              - errored
              type: string
            url:
              description: URL of the Order. This will initially be empty when the
                resource is first created. The Order controller will populate this
                field when the Order is first processed. This field will be immutable
                after it is initially set.
              type: string
          type: object
      required:
      - metadata
  versions:
  - name: v1alpha2
    served: true
    storage: true

---
# Source: cert-manager/templates/cainjector-rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: cert-manager-cainjector
  labels:
    app: cainjector
    app.kubernetes.io/name: cainjector
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "cainjector"
    helm.sh/chart: cert-manager-v0.15.1
rules:
  - apiGroups: ["cert-manager.io"]
    resources: ["certificates"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["get", "create", "update", "patch"]
  - apiGroups: ["admissionregistration.k8s.io"]
    resources: ["validatingwebhookconfigurations", "mutatingwebhookconfigurations"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["apiregistration.k8s.io"]
    resources: ["apiservices"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["apiextensions.k8s.io"]
    resources: ["customresourcedefinitions"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["auditregistration.k8s.io"]
    resources: ["auditsinks"]
    verbs: ["get", "list", "watch", "update"]
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: cert-manager-cainjector
  labels:
    app: cainjector
    app.kubernetes.io/name: cainjector
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "cainjector"
    helm.sh/chart: cert-manager-v0.15.1
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cert-manager-cainjector
subjects:
  - name: cert-manager-cainjector
    namespace: kube-system
    kind: ServiceAccount

---
# leader election rules
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: cert-manager-cainjector:leaderelection
  namespace: kube-system
  labels:
    app: cainjector
    app.kubernetes.io/name: cainjector
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "cainjector"
    helm.sh/chart: cert-manager-v0.15.1
rules:
  # Used for leader election by the controller
  # cert-manager-cainjector-leader-election is used by the CertificateBased injector controller
  #   see cmd/cainjector/start.go#L113
  # cert-manager-cainjector-leader-election-core is used by the SecretBased injector controller
  #   see cmd/cainjector/start.go#L137
  - apiGroups: [""]
    resources: ["configmaps"]
    resourceNames: ["cert-manager-cainjector-leader-election", "cert-manager-cainjector-leader-election-core"]
    verbs: ["get", "update", "patch"]
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["create"]
---

# grant cert-manager permission to manage the leaderelection configmap in the
# leader election namespace
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: cert-manager-cainjector:leaderelection
  namespace: kube-system
  labels:
    app: cainjector
    app.kubernetes.io/name: cainjector
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "cainjector"
    helm.sh/chart: cert-manager-v0.15.1
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cert-manager-cainjector:leaderelection
subjects:
  - kind: ServiceAccount
    name: cert-manager-cainjector
    namespace: kube-system
---
# Source: cert-manager/templates/rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: cert-manager:leaderelection
  namespace: kube-system
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
rules:
  # Used for leader election by the controller
  - apiGroups: [""]
    resources: ["configmaps"]
    resourceNames: ["cert-manager-controller"]
    verbs: ["get", "update", "patch"]
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["create"]

---

# grant cert-manager permission to manage the leaderelection configmap in the
# leader election namespace
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: cert-manager:leaderelection
  namespace: kube-system
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cert-manager:leaderelection
subjects:
  - apiGroup: ""
    kind: ServiceAccount
    name: cert-manager
    namespace: kube-system

---

# Issuer controller role
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: cert-manager-controller-issuers
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
rules:
  - apiGroups: ["cert-manager.io"]
    resources: ["issuers", "issuers/status"]
    verbs: ["update"]
  - apiGroups: ["cert-manager.io"]
    resources: ["issuers"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list", "watch", "create", "update", "delete"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "patch"]

---

# ClusterIssuer controller role
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: cert-manager-controller-clusterissuers
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
rules:
  - apiGroups: ["cert-manager.io"]
    resources: ["clusterissuers", "clusterissuers/status"]
    verbs: ["update"]
  - apiGroups: ["cert-manager.io"]
    resources: ["clusterissuers"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list", "watch", "create", "update", "delete"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "patch"]

---

# Certificates controller role
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: cert-manager-controller-certificates
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
rules:
  - apiGroups: ["cert-manager.io"]
    resources: ["certificates", "certificates/status", "certificaterequests", "certificaterequests/status"]
    verbs: ["update"]
  - apiGroups: ["cert-manager.io"]
    resources: ["certificates", "certificaterequests", "clusterissuers", "issuers"]
    verbs: ["get", "list", "watch"]
  # We require these rules to support users with the OwnerReferencesPermissionEnforcement
  # admission controller enabled:
  # https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#ownerreferencespermissionenforcement
  - apiGroups: ["cert-manager.io"]
    resources: ["certificates/finalizers", "certificaterequests/finalizers"]
    verbs: ["update"]
  - apiGroups: ["acme.cert-manager.io"]
    resources: ["orders"]
    verbs: ["create", "delete", "get", "list", "watch"]
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list", "watch", "create", "update", "delete"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "patch"]

---

# Orders controller role
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: cert-manager-controller-orders
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
rules:
  - apiGroups: ["acme.cert-manager.io"]
    resources: ["orders", "orders/status"]
    verbs: ["update"]
  - apiGroups: ["acme.cert-manager.io"]
    resources: ["orders", "challenges"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["cert-manager.io"]
    resources: ["clusterissuers", "issuers"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["acme.cert-manager.io"]
    resources: ["challenges"]
    verbs: ["create", "delete"]
  # We require these rules to support users with the OwnerReferencesPermissionEnforcement
  # admission controller enabled:
  # https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#ownerreferencespermissionenforcement
  - apiGroups: ["acme.cert-manager.io"]
    resources: ["orders/finalizers"]
    verbs: ["update"]
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "patch"]

---

# Challenges controller role
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: cert-manager-controller-challenges
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
rules:
  # Use to update challenge resource status
  - apiGroups: ["acme.cert-manager.io"]
    resources: ["challenges", "challenges/status"]
    verbs: ["update"]
  # Used to watch challenge resources
  - apiGroups: ["acme.cert-manager.io"]
    resources: ["challenges"]
    verbs: ["get", "list", "watch"]
  # Used to watch challenges, issuer and clusterissuer resources
  - apiGroups: ["cert-manager.io"]
    resources: ["issuers", "clusterissuers"]
    verbs: ["get", "list", "watch"]
  # Need to be able to retrieve ACME account private key to complete challenges
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list", "watch"]
  # Used to create events
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "patch"]
  # HTTP01 rules
  - apiGroups: [""]
    resources: ["pods", "services"]
    verbs: ["get", "list", "watch", "create", "delete"]
  - apiGroups: ["extensions"]
    resources: ["ingresses"]
    verbs: ["get", "list", "watch", "create", "delete", "update"]
  # We require the ability to specify a custom hostname when we are creating
  # new ingress resources.
  # See: https://github.com/openshift/origin/blob/21f191775636f9acadb44fa42beeb4f75b255532/pkg/route/apiserver/admission/ingress_admission.go#L84-L148
  - apiGroups: ["route.openshift.io"]
    resources: ["routes/custom-host"]
    verbs: ["create"]
  # We require these rules to support users with the OwnerReferencesPermissionEnforcement
  # admission controller enabled:
  # https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#ownerreferencespermissionenforcement
  - apiGroups: ["acme.cert-manager.io"]
    resources: ["challenges/finalizers"]
    verbs: ["update"]
  # DNS01 rules (duplicated above)
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list", "watch"]

---

# ingress-shim controller role
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: cert-manager-controller-ingress-shim
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
rules:
  - apiGroups: ["cert-manager.io"]
    resources: ["certificates", "certificaterequests"]
    verbs: ["create", "update", "delete"]
  - apiGroups: ["cert-manager.io"]
    resources: ["certificates", "certificaterequests", "issuers", "clusterissuers"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["extensions"]
    resources: ["ingresses"]
    verbs: ["get", "list", "watch"]
  # We require these rules to support users with the OwnerReferencesPermissionEnforcement
  # admission controller enabled:
  # https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#ownerreferencespermissionenforcement
  - apiGroups: ["extensions"]
    resources: ["ingresses/finalizers"]
    verbs: ["update"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "patch"]

---

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: cert-manager-controller-issuers
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cert-manager-controller-issuers
subjects:
  - name: cert-manager
    namespace: kube-system
    kind: ServiceAccount

---

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: cert-manager-controller-clusterissuers
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cert-manager-controller-clusterissuers
subjects:
  - name: cert-manager
    namespace: kube-system
    kind: ServiceAccount

---

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: cert-manager-controller-certificates
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cert-manager-controller-certificates
subjects:
  - name: cert-manager
    namespace: kube-system
    kind: ServiceAccount

---

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: cert-manager-controller-orders
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cert-manager-controller-orders
subjects:
  - name: cert-manager
    namespace: kube-system
    kind: ServiceAccount

---

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: cert-manager-controller-challenges
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cert-manager-controller-challenges
subjects:
  - name: cert-manager
    namespace: kube-system
    kind: ServiceAccount

---

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: cert-manager-controller-ingress-shim
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cert-manager-controller-ingress-shim
subjects:
  - name: cert-manager
    namespace: kube-system
    kind: ServiceAccount

---

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cert-manager-view
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
    rbac.authorization.k8s.io/aggregate-to-view: "true"
    rbac.authorization.k8s.io/aggregate-to-edit: "true"
    rbac.authorization.k8s.io/aggregate-to-admin: "true"
rules:
  - apiGroups: ["cert-manager.io"]
    resources: ["certificates", "certificaterequests", "issuers"]
    verbs: ["get", "list", "watch"]

---

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cert-manager-edit
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
    rbac.authorization.k8s.io/aggregate-to-edit: "true"
    rbac.authorization.k8s.io/aggregate-to-admin: "true"
rules:
  - apiGroups: ["cert-manager.io"]
    resources: ["certificates", "certificaterequests", "issuers"]
    verbs: ["create", "delete", "deletecollection", "patch", "update"]

---
# Source: cert-manager/templates/webhook-rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: cert-manager-webhook:dynamic-serving
  namespace: kube-system
  labels:
    app: webhook
    app.kubernetes.io/name: webhook
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "webhook"
    helm.sh/chart: cert-manager-v0.15.1
rules:
- apiGroups: [""]
  resources: ["secrets"]
  resourceNames:
  - 'cert-manager-webhook-ca'
  verbs: ["get", "list", "watch", "update"]
# It's not possible to grant CREATE permission on a single resourceName.
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["create"]
---

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: cert-manager-webhook:dynamic-serving
  namespace: kube-system
  labels:
    app: webhook
    app.kubernetes.io/name: webhook
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "webhook"
    helm.sh/chart: cert-manager-v0.15.1
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cert-manager-webhook:dynamic-serving
subjects:
- apiGroup: ""
  kind: ServiceAccount
  name: cert-manager-webhook
  namespace: kube-system
---
# Source: cert-manager/templates/service.yaml

apiVersion: v1
kind: Service
metadata:
  name: cert-manager
  namespace: kube-system
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
spec:
  type: ClusterIP
  ports:
    - protocol: TCP
      port: 9402
      targetPort: 9402
  selector:
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/component: "controller"

---
# Source: cert-manager/templates/webhook-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: cert-manager-webhook
  namespace: kube-system
  labels:
    app: webhook
    app.kubernetes.io/name: webhook
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "webhook"
    helm.sh/chart: cert-manager-v0.15.1
spec:
  type: ClusterIP
  ports:
  - name: https
    port: 443
    targetPort: 10250
  selector:
    app.kubernetes.io/name: webhook
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/component: "webhook"

---
# Source: cert-manager/templates/cainjector-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cert-manager-cainjector
  namespace: kube-system
  labels:
    app: cainjector
    app.kubernetes.io/name: cainjector
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "cainjector"
    helm.sh/chart: cert-manager-v0.15.1
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: cainjector
      app.kubernetes.io/instance: cert-manager
      app.kubernetes.io/component: "cainjector"
  template:
    metadata:
      labels:
        app: cainjector
        app.kubernetes.io/name: cainjector
        app.kubernetes.io/instance: cert-manager
        app.kubernetes.io/managed-by: Tiller
        app.kubernetes.io/component: "cainjector"
        helm.sh/chart: cert-manager-v0.15.1
    spec:
      serviceAccountName: cert-manager-cainjector
      containers:
        - name: cert-manager
          image: {{.CertManagerCAInjectorImage}}
          imagePullPolicy: IfNotPresent
          args:
          - --v=2
          - --leader-election-namespace=kube-system
          env:
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          resources:
            {}
            
---
# Source: cert-manager/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cert-manager
  namespace: kube-system
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "controller"
    helm.sh/chart: cert-manager-v0.15.1
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: cert-manager
      app.kubernetes.io/instance: cert-manager
      app.kubernetes.io/component: "controller"
  template:
    metadata:
      labels:
        app: cert-manager
        app.kubernetes.io/name: cert-manager
        app.kubernetes.io/instance: cert-manager
        app.kubernetes.io/component: "controller"
        app.kubernetes.io/managed-by: Tiller
        helm.sh/chart: cert-manager-v0.15.1
      annotations:
        prometheus.io/path: "/metrics"
        prometheus.io/scrape: 'true'
        prometheus.io/port: '9402'
    spec:
      serviceAccountName: cert-manager
      containers:
        - name: cert-manager
          image: {{.GetCertManagerControllerImage}}
          imagePullPolicy: IfNotPresent
          args:
          - --v=2
          - --cluster-resource-namespace=$(POD_NAMESPACE)
          - --leader-election-namespace=kube-system
          ports:
          - containerPort: 9402
            protocol: TCP
          env:
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          resources:
            {}
            

---
# Source: cert-manager/templates/webhook-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cert-manager-webhook
  namespace: kube-system
  labels:
    app: webhook
    app.kubernetes.io/name: webhook
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "webhook"
    helm.sh/chart: cert-manager-v0.15.1
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: webhook
      app.kubernetes.io/instance: cert-manager
      app.kubernetes.io/component: "webhook"
  template:
    metadata:
      labels:
        app: webhook
        app.kubernetes.io/name: webhook
        app.kubernetes.io/instance: cert-manager
        app.kubernetes.io/managed-by: Tiller
        app.kubernetes.io/component: "webhook"
        helm.sh/chart: cert-manager-v0.15.1
    spec:
      serviceAccountName: cert-manager-webhook
      containers:
        - name: cert-manager
          image: {{.GetCertManagerWebhookImage}}
          imagePullPolicy: IfNotPresent
          args:
          - --v=2
          - --secure-port=10250
          - --dynamic-serving-ca-secret-namespace=kube-system
          - --dynamic-serving-ca-secret-name=cert-manager-webhook-ca
          - --dynamic-serving-dns-names=cert-manager-webhook,cert-manager-webhook.kube-system,cert-manager-webhook.kube-system.svc
          ports:
          - name: https
            containerPort: 10250
          livenessProbe:
            httpGet:
              path: /livez
              port: 6080
              scheme: HTTP
            initialDelaySeconds: 60
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /healthz
              port: 6080
              scheme: HTTP
            initialDelaySeconds: 5
            periodSeconds: 5
          env:
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          resources:
            {}
            


---
# Source: cert-manager/templates/webhook-mutating-webhook.yaml
apiVersion: admissionregistration.k8s.io/v1beta1
kind: MutatingWebhookConfiguration
metadata:
  name: cert-manager-webhook
  labels:
    app: webhook
    app.kubernetes.io/name: webhook
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "webhook"
    helm.sh/chart: cert-manager-v0.15.1
  annotations:
    cert-manager.io/inject-ca-from-secret: "kube-system/cert-manager-webhook-ca"
webhooks:
  - name: webhook.cert-manager.io
    rules:
      - apiGroups:
          - "cert-manager.io"
          - "acme.cert-manager.io"
        apiVersions:
          - v1alpha2
          - v1alpha3
        operations:
          - CREATE
          - UPDATE
        resources:
          - "*/*"
    failurePolicy: Fail
    # Only include 'sideEffects' field in Kubernetes 1.12+
    sideEffects: None
    clientConfig:
      service:
        name: cert-manager-webhook
        namespace: kube-system
        path: /mutate

---
# Source: cert-manager/templates/cainjector-psp-clusterrole.yaml


---
# Source: cert-manager/templates/cainjector-psp-clusterrolebinding.yaml


---
# Source: cert-manager/templates/cainjector-psp.yaml

---
# Source: cert-manager/templates/crds.yaml


---
# Source: cert-manager/templates/psp-clusterrole.yaml


---
# Source: cert-manager/templates/psp-clusterrolebinding.yaml


---
# Source: cert-manager/templates/psp.yaml


---
# Source: cert-manager/templates/servicemonitor.yaml


---
# Source: cert-manager/templates/webhook-psp-clusterrole.yaml
 

---
# Source: cert-manager/templates/webhook-psp-clusterrolebinding.yaml


---
# Source: cert-manager/templates/webhook-psp.yaml


---
# Source: cert-manager/templates/webhook-validating-webhook.yaml
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  name: cert-manager-webhook
  labels:
    app: webhook
    app.kubernetes.io/name: webhook
    app.kubernetes.io/instance: cert-manager
    app.kubernetes.io/managed-by: Tiller
    app.kubernetes.io/component: "webhook"
    helm.sh/chart: cert-manager-v0.15.1
  annotations:
    cert-manager.io/inject-ca-from-secret: "kube-system/cert-manager-webhook-ca"
webhooks:
  - name: webhook.cert-manager.io
    namespaceSelector:
      matchExpressions:
      - key: "cert-manager.io/disable-validation"
        operator: "NotIn"
        values:
        - "true"
      - key: "name"
        operator: "NotIn"
        values:
        - kube-system
    rules:
      - apiGroups:
          - "cert-manager.io"
          - "acme.cert-manager.io"
        apiVersions:
          - v1alpha2
          - v1alpha3
        operations:
          - CREATE
          - UPDATE
        resources:
          - "*/*"
    failurePolicy: Fail
    # Only include 'sideEffects' field in Kubernetes 1.12+
    sideEffects: None
    clientConfig:
      service:
        name: cert-manager-webhook
        namespace: kube-system
        path: /validate
`
)

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

package auth_test

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/SUSE/skuba/pkg/skuba/actions/auth"
)

func startServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", auth.OpenIDHandler())
	mux.HandleFunc("/auth", auth.AuthHandler())
	mux.HandleFunc("/auth/local", auth.AuthLocalHandler())
	mux.HandleFunc("/token", auth.TokenHandler())
	mux.HandleFunc("/keys", auth.KeysHandler())
	mux.HandleFunc("/approval", auth.ApprovalHandler())

	srv := httptest.NewUnstartedServer(mux)
	cert, _ := tls.X509KeyPair(auth.LocalhostCert, auth.LocalhostKey)
	srv.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	srv.StartTLS()
	return srv
}

func Test_Login(t *testing.T) {
	testServer := startServer()
	defer testServer.Close()

	tests := []struct {
		name               string
		cfg                auth.LoginConfig
		expectedKubeConfCb func() *clientcmdapi.Config
		expectedError      bool
	}{
		{
			name: "secure ssl/tls",
			cfg: auth.LoginConfig{
				DexServer:   testServer.URL,
				Username:    auth.MockDefaultUsername,
				Password:    auth.MockDefaultPassword,
				RootCAPath:  "testdata/localhost.crt",
				ClusterName: "test-cluster-name",
			},
			expectedKubeConfCb: func() *clientcmdapi.Config {
				url, _ := url.Parse(testServer.URL)

				kubeConfig := clientcmdapi.NewConfig()
				kubeConfig.Clusters["test-cluster-name"] = &clientcmdapi.Cluster{
					Server:                   fmt.Sprintf("%s://%s:%s", url.Scheme, url.Hostname(), auth.DefaultPortAPIServer),
					CertificateAuthorityData: auth.LocalhostCert,
				}
				kubeConfig.Contexts["test-cluster-name"] = &clientcmdapi.Context{
					Cluster:  "test-cluster-name",
					AuthInfo: auth.MockDefaultUsername,
				}
				kubeConfig.CurrentContext = "test-cluster-name"
				kubeConfig.AuthInfos[auth.MockDefaultUsername] = &clientcmdapi.AuthInfo{
					AuthProvider: &clientcmdapi.AuthProviderConfig{
						Name: "oidc",
						Config: map[string]string{
							"idp-issuer-url": testServer.URL,
							"client-id":      "oidc-cli",
							"client-secret":  "swac7qakes7AvucH8bRucucH",
							"id-token":       auth.MockIDToken,
							"refresh-token":  auth.MockRefreshToken,
						},
					},
				}
				return kubeConfig
			},
		},
		{
			name: "insecure ssl/tls",
			cfg: auth.LoginConfig{
				DexServer:          testServer.URL,
				Username:           auth.MockDefaultUsername,
				Password:           auth.MockDefaultPassword,
				InsecureSkipVerify: true,
				ClusterName:        "test-cluster-name",
			},
			expectedKubeConfCb: func() *clientcmdapi.Config {
				url, _ := url.Parse(testServer.URL)

				kubeConfig := clientcmdapi.NewConfig()
				kubeConfig.Clusters["test-cluster-name"] = &clientcmdapi.Cluster{
					Server:                fmt.Sprintf("%s://%s:%s", url.Scheme, url.Hostname(), auth.DefaultPortAPIServer),
					InsecureSkipTLSVerify: true,
				}
				kubeConfig.Contexts["test-cluster-name"] = &clientcmdapi.Context{
					Cluster:  "test-cluster-name",
					AuthInfo: auth.MockDefaultUsername,
				}
				kubeConfig.CurrentContext = "test-cluster-name"
				kubeConfig.AuthInfos[auth.MockDefaultUsername] = &clientcmdapi.AuthInfo{
					AuthProvider: &clientcmdapi.AuthProviderConfig{
						Name: "oidc",
						Config: map[string]string{
							"idp-issuer-url": testServer.URL,
							"client-id":      "oidc-cli",
							"client-secret":  "swac7qakes7AvucH8bRucucH",
							"id-token":       "eyJhbGciOiJSUzI1NiIsImtpZCI6ImFlNjg5YTI1OWJkYjRjMWZiZDZmZGFjMzg0OTk5YTJhNWNlNmRmOGEifQ.eyJpc3MiOiJodHRwczovLzEwLjg2LjAuMTE2OjMyMDAwIiwic3ViIjoiQ2loamJqMW9aV3hzYjNkdmNteGtMRzkxUFhWelpYSnpMR1JqUFdWNFlXMXdiR1VzWkdNOWIzSm5FZ1JzWkdGdyIsImF1ZCI6WyJvaWRjIiwib2lkYy1jbGkiXSwiZXhwIjoxNTY1MTY1NjE0LCJpYXQiOjE1NjUwNzkyMTQsImF6cCI6Im9pZGMtY2xpIiwiYXRfaGFzaCI6ImJ2cml3TEthMzZuTjExUG5Jc3RRSmciLCJlbWFpbCI6ImhlbGxvQHN1c2UuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImdyb3VwcyI6WyJkZXZlbG9wZXJzIl0sIm5hbWUiOiJIZWxsbyBXb3JsZCJ9.AYN9cbk2hS6S8ZbQSZ4yoGksPJJ9qzbK8iXCoB6XXmhc5AUlwxnXQ-vzcp1u6h8AtY3iJX0s5ZwH3BthKEBlj6Aad6v5qp62Ws0Wb1-RY6TcCNQv4AdpBuFlJtJIxp7wI33bR0gpLOMsjYJRgKuLvQ1Dn7tipT62CPhqwA91lT613_yByLC8ek1Qy3RSwJIA_hkJT0H-yMHM2JC5WuB3P0MEURfl2QIXaWDjoV5RcL0dh_dkwy2v6zxgCPu0gFvL2BOrcHPjv6k6kphMnQ8uCbQaEfNxuMYr7zDRWBcNSpfjhbbYRAjNBHbpMorM3mT83GB76cxdUWCW2q69nM1B_w",
							"refresh-token":  "ChludG1ncnh1aHQ1a3F0dG83enRvYmtlc2hiEhlsNWNvM3V6cmVjb2FxYW1maHZqa2F5azJh",
						},
					},
				}
				return kubeConfig
			},
		},
		{
			name: "with apiserver url",
			cfg: auth.LoginConfig{
				DexServer:          testServer.URL,
				Username:           auth.MockDefaultUsername,
				Password:           auth.MockDefaultPassword,
				InsecureSkipVerify: true,
				ClusterName:        "test-cluster-name",
				KubeAPIServer:      "https://10.2.3.4:6443",
			},
			expectedKubeConfCb: func() *clientcmdapi.Config {
				kubeConfig := clientcmdapi.NewConfig()
				kubeConfig.Clusters["test-cluster-name"] = &clientcmdapi.Cluster{
					Server:                "https://10.2.3.4:6443",
					InsecureSkipTLSVerify: true,
				}
				kubeConfig.Contexts["test-cluster-name"] = &clientcmdapi.Context{
					Cluster:  "test-cluster-name",
					AuthInfo: auth.MockDefaultUsername,
				}
				kubeConfig.CurrentContext = "test-cluster-name"
				kubeConfig.AuthInfos[auth.MockDefaultUsername] = &clientcmdapi.AuthInfo{
					AuthProvider: &clientcmdapi.AuthProviderConfig{
						Name: "oidc",
						Config: map[string]string{
							"idp-issuer-url": testServer.URL,
							"client-id":      "oidc-cli",
							"client-secret":  "swac7qakes7AvucH8bRucucH",
							"id-token":       "eyJhbGciOiJSUzI1NiIsImtpZCI6ImFlNjg5YTI1OWJkYjRjMWZiZDZmZGFjMzg0OTk5YTJhNWNlNmRmOGEifQ.eyJpc3MiOiJodHRwczovLzEwLjg2LjAuMTE2OjMyMDAwIiwic3ViIjoiQ2loamJqMW9aV3hzYjNkdmNteGtMRzkxUFhWelpYSnpMR1JqUFdWNFlXMXdiR1VzWkdNOWIzSm5FZ1JzWkdGdyIsImF1ZCI6WyJvaWRjIiwib2lkYy1jbGkiXSwiZXhwIjoxNTY1MTY1NjE0LCJpYXQiOjE1NjUwNzkyMTQsImF6cCI6Im9pZGMtY2xpIiwiYXRfaGFzaCI6ImJ2cml3TEthMzZuTjExUG5Jc3RRSmciLCJlbWFpbCI6ImhlbGxvQHN1c2UuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImdyb3VwcyI6WyJkZXZlbG9wZXJzIl0sIm5hbWUiOiJIZWxsbyBXb3JsZCJ9.AYN9cbk2hS6S8ZbQSZ4yoGksPJJ9qzbK8iXCoB6XXmhc5AUlwxnXQ-vzcp1u6h8AtY3iJX0s5ZwH3BthKEBlj6Aad6v5qp62Ws0Wb1-RY6TcCNQv4AdpBuFlJtJIxp7wI33bR0gpLOMsjYJRgKuLvQ1Dn7tipT62CPhqwA91lT613_yByLC8ek1Qy3RSwJIA_hkJT0H-yMHM2JC5WuB3P0MEURfl2QIXaWDjoV5RcL0dh_dkwy2v6zxgCPu0gFvL2BOrcHPjv6k6kphMnQ8uCbQaEfNxuMYr7zDRWBcNSpfjhbbYRAjNBHbpMorM3mT83GB76cxdUWCW2q69nM1B_w",
							"refresh-token":  "ChludG1ncnh1aHQ1a3F0dG83enRvYmtlc2hiEhlsNWNvM3V6cmVjb2FxYW1maHZqa2F5azJh",
						},
					},
				}
				return kubeConfig
			},
		},
		{
			name: "with invalid apiserver url",
			cfg: auth.LoginConfig{
				DexServer:          testServer.URL,
				Username:           auth.MockDefaultUsername,
				Password:           auth.MockDefaultPassword,
				InsecureSkipVerify: true,
				KubeAPIServer:      ".10.2.3.4:6443",
			},
			expectedError: true,
		},
		{
			name: "cert file not exist",
			cfg: auth.LoginConfig{
				DexServer:   testServer.URL,
				Username:    auth.MockDefaultUsername,
				Password:    auth.MockDefaultPassword,
				RootCAPath:  "testdata/nonexist.crt",
				ClusterName: "test-cluster-name",
			},
			expectedError: true,
		},
		{
			name: "auth failed",
			cfg: auth.LoginConfig{
				DexServer:          testServer.URL,
				Username:           "admin",
				Password:           "1234",
				InsecureSkipVerify: true,
			},
			expectedError: true,
		},
		{
			name: "invalid url",
			cfg: auth.LoginConfig{
				DexServer:          ".1.2.3.4",
				Username:           auth.MockDefaultUsername,
				Password:           auth.MockDefaultPassword,
				InsecureSkipVerify: true,
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotKubeConfig, err := auth.Login(tt.cfg)

			if tt.expectedError {
				if err == nil {
					t.Errorf("error expected on %s, but no error reported", tt.name)
				}
				return
			} else if !tt.expectedError && err != nil {
				t.Errorf("error not expected on %s, but an error was reported (%v)", tt.name, err)
				return
			}

			if !reflect.DeepEqual(gotKubeConfig, tt.expectedKubeConfCb()) {
				t.Errorf("got %v, want %v", gotKubeConfig, tt.expectedKubeConfCb())
				return
			}
		})
	}
}

func Test_SaveKubeconfig(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		kubeConfig    *clientcmdapi.Config
		expectedError bool
	}{
		{
			name:       "success output",
			kubeConfig: clientcmdapi.NewConfig(),
		},
		{
			name:          "open file failed",
			filename:      "path/to/kubeconfig",
			kubeConfig:    clientcmdapi.NewConfig(),
			expectedError: true,
		},
		{
			name:          "encode failed",
			kubeConfig:    nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var path string
			if tt.filename != "" {
				path = filepath.Join("testdata", tt.filename+".golden")
			} else {
				path = filepath.Join("testdata", tt.name+".golden")
			}
			err := auth.SaveKubeconfig(path, tt.kubeConfig)

			if tt.expectedError && err == nil {
				t.Errorf("error expected on %s, but no error reported", tt.name)
				return
			} else if !tt.expectedError && err != nil {
				t.Errorf("error not expected on %s, but an error was reported (%v)", tt.name, err)
				return
			}
		})
	}
}

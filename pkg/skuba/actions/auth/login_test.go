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

package auth

import (
	"reflect"
	"testing"
)

func Test_processURL(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		defaultScheme string
		defaultPort   string
		expectedURL   string
		expectedError bool
	}{
		{
			name:          "ip",
			url:           "10.1.2.3",
			defaultScheme: "https",
			defaultPort:   "32000",
			expectedURL:   "https://10.1.2.3:32000",
		},
		{
			name:          "https+ip",
			url:           "https://10.1.2.3",
			defaultScheme: "https",
			defaultPort:   "32000",
			expectedURL:   "https://10.1.2.3:32000",
		},
		{
			name:          "ip+port",
			url:           "10.1.2.3:32000",
			defaultScheme: "https",
			defaultPort:   "32000",
			expectedURL:   "https://10.1.2.3:32000",
		},
		{
			name:          "https+ip+port",
			url:           "https://10.1.2.3:32000",
			defaultScheme: "https",
			defaultPort:   "32000",
			expectedURL:   "https://10.1.2.3:32000",
		},
		{
			name:          "http+host+port",
			url:           "http://10.1.2.3:32000",
			defaultScheme: "https",
			defaultPort:   "32000",
			expectedError: true,
		},
		{
			name:          "https+port",
			url:           "https://:32000",
			defaultScheme: "https",
			defaultPort:   "32000",
			expectedError: true,
		},
		{
			name:          "host",
			url:           "localhost.net",
			defaultScheme: "https",
			defaultPort:   "32000",
			expectedURL:   "https://localhost.net:32000",
		},
		{
			name:          "https+host",
			url:           "https://localhost.net",
			defaultScheme: "https",
			defaultPort:   "32000",
			expectedURL:   "https://localhost.net:32000",
		},
		{
			name:          "host+port",
			url:           "localhost.net:32000",
			defaultScheme: "https",
			defaultPort:   "32000",
			expectedURL:   "https://localhost.net:32000",
		},
		{
			name:          "https+host+port",
			url:           "https://localhost.net:32000",
			defaultScheme: "https",
			defaultPort:   "32000",
			expectedURL:   "https://localhost.net:32000",
		},
		{
			name:          "http+host+port",
			url:           "http://localhost.net:32000",
			defaultScheme: "https",
			defaultPort:   "32000",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotURL, err := processURL(tt.url, tt.defaultScheme, tt.defaultPort)

			if tt.expectedError && err == nil {
				t.Errorf("error expected on %s, but no error reported", tt.name)
			} else if !tt.expectedError && err != nil {
				t.Errorf("error not expected on %s, but an error was reported (%v)", tt.name, err)
				return
			}

			if !reflect.DeepEqual(gotURL, tt.expectedURL) {
				t.Errorf("got %v, want %v", gotURL, tt.expectedURL)
			}
		})
	}
}

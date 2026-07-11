/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flowcontrol

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDebugDumpHandlersSetResponseHeaders(t *testing.T) {
	cfgCtlr := &configController{
		priorityLevelStates: map[string]*priorityLevelState{},
	}

	tests := []struct {
		name    string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{
			name:    "dump priority levels",
			handler: cfgCtlr.dumpPriorityLevels,
		},
		{
			name:    "dump queues",
			handler: cfgCtlr.dumpQueues,
		},
		{
			name:    "dump requests",
			handler: cfgCtlr.dumpRequests,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/debug/api_priority_and_fairness", nil)
			resp := httptest.NewRecorder()

			tc.handler(resp, req)

			if got := resp.Header().Get("Content-Type"); got != debugResponseContentType {
				t.Fatalf("expected Content-Type %q, got %q", debugResponseContentType, got)
			}
			if got := resp.Header().Get("X-Content-Type-Options"); got != "nosniff" {
				t.Fatalf("expected X-Content-Type-Options header %q, got %q", "nosniff", got)
			}
		})
	}
}

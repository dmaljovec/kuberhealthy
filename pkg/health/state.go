// Copyright 2018 Comcast Cable Communications Management, LLC
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package health

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	khstatev1 "github.com/Comcast/kuberhealthy/v2/pkg/apis/khstate/v1"
)

// State represents the results of all checks being managed along with a top-level OK and Error state. This is displayed
// on the kuberhealthy status page as JSON
type State struct {
	OK            bool
	Errors        []string
	CheckDetails  map[string]khstatev1.WorkloadDetails // map of check names to last run timestamp
	JobDetails    map[string]khstatev1.WorkloadDetails // map of job names to last run timestamp
	CurrentMaster string
}

// AddError adds new errors to State
func (h *State) AddError(s ...string) {
	for _, str := range s {
		if len(s) == 0 {
			log.Warningln("AddError was called but the error was blank so it was skipped.")
			continue
		}
		log.Debugln("Appending error:", str)
		h.Errors = append(h.Errors, str)
	}
}

// WriteHTTPStatusResponse writes a response to an http response writer
func (h *State) WriteHTTPStatusResponse(w http.ResponseWriter) error {

	currentStatus := *h

	// marshal the health check results into a json blob of bytes
	b, err := json.MarshalIndent(currentStatus, "", "  ")
	if err != nil {
		log.Warningln("Error marshaling health check json for caller:", err)
		return err
	}

	// write the output to the caller
	_, err = w.Write(b)
	if err != nil {
		log.Errorln("Error writing response to caller:", err)
		return err
	}

	return err
}

// NewState creates a new health check result response
func NewState() State {
	s := State{}
	s.OK = true
	s.Errors = []string{}
	s.CheckDetails = make(map[string]khstatev1.WorkloadDetails)
	s.JobDetails = make(map[string]khstatev1.WorkloadDetails)
	return s
}

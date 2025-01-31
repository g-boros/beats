// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package snapshot

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	// "strconv"
	// "path/filepath"
	"testing"

	"github.com/elastic/beats/v7/metricbeat/module/elasticsearch"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mbtest "github.com/elastic/beats/v7/metricbeat/mb/testing"
)

var info = elasticsearch.Info{
	ClusterID:   "1234",
	ClusterName: "helloworld",
}

func TestNumOfEvents(t *testing.T) {
	input, err := ioutil.ReadFile("./_meta/test/snapshot.714.json")
	require.NoError(t, err)

	reporter := &mbtest.CapturingReporterV2{}
	eventsMapping(reporter, info, input)

	require.True(t, len(reporter.GetEvents()) >= 1)
	require.Equal(t, 0, len(reporter.GetErrors()))
}

func TestSchema(t *testing.T) {
	input, err := ioutil.ReadFile("./_meta/test/snapshot.714.json")
	require.NoError(t, err)

	reporter := &mbtest.CapturingReporterV2{}
	eventsMapping(reporter, info, input)

	require.GreaterOrEqual(t, len(reporter.GetEvents()), 1)
	version, err := reporter.GetEvents()[0].MetricSetFields.GetValue("version")
	if assert.NoError(t, err) {
		versionInt, ok := version.(int64)
		require.True(t, ok)
		require.GreaterOrEqual(t, versionInt, int64(1))
	}

	name, err := reporter.GetEvents()[0].MetricSetFields.GetValue("name")
	if assert.NoError(t, err) {
		assert.Greater(t, len(fmt.Sprintf("%v", name)), 0)
	}
}

func TestData(t *testing.T) {
	mux := http.NewServeMux()

	mux.Handle("/_slm/policy", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input, _ := ioutil.ReadFile(("./_meta/test/snapshot.714.json"))
		w.Write(input)
	}))
	server := httptest.NewServer(mux)
	defer server.Close()

	ms := mbtest.NewReportingMetricSetV2Error(t, getConfig(server.URL))
	if err := mbtest.WriteEventsReporterV2Error(ms, t, ""); err != nil {
		t.Fatal("write", err)
	}
}

func getConfig(host string) map[string]interface{} {
	return map[string]interface{}{
		"module":     elasticsearch.ModuleName,
		"metricsets": []string{"snapshot"},
		"hosts":      []string{host},
	}
}

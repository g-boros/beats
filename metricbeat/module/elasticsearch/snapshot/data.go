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
	"encoding/json"
	"fmt"

	"github.com/joeshaw/multierror"
	"github.com/pkg/errors"

	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/common"
	s "github.com/elastic/beats/v7/libbeat/common/schema"
	c "github.com/elastic/beats/v7/libbeat/common/schema/mapstriface"
	"github.com/elastic/beats/v7/metricbeat/mb"
	"github.com/elastic/beats/v7/metricbeat/module/elasticsearch"
)

var (
 schema = s.Schema{
	// "name":         c.Str("name"),
	"version":            c.Int("version"),
	"modified_date_millis": c.Int("modified_date_millis"),
	"policy":	c.Dict("policy", s.Schema {
		"name":  c.Str("name"),
		"schedule":  c.Str("schedule"),
		"repository":  c.Str("repository"),
		"config": c.Ifc("config"),
		"retention": c.Dict("retention", s.Schema{
			"expire_after": c.Str("expire_after"),
			"min_count": c.Int("min_count"),
			"max_count": c.Int("max_count"),
		}),
	}),
	"last_success": c.Dict("last_success", s.Schema{
		"snapshot_name": c.Str("snapshot_name"),
		"time": c.Int("time"),
	}, c.DictOptional),
	"last_failure": c.Dict("last_failure", s.Schema{
		"snapshot_name": c.Str("snapshot_name"),
		"time": c.Int("time"),
		// intentionally skip as it can be too verbose and not a metric
		// "details": c.Str("details")
	}, c.DictOptional),
	"next_execution_millis": c.Int("next_execution_millis"),
	"stats": c.Dict("stats", s.Schema{
		"policy": c.Str("policy"),
		"snapshots_taken": c.Int("snapshots_taken"),
		"snapshots_failed": c.Int("snapshots_failed"),
		"snapshots_deleted": c.Int("snapshots_deleted"),
		"snapshot_deletion_failures": c.Int("snapshot_deletion_failures"),
	}),
})

func eventsMapping(r mb.ReporterV2, info elasticsearch.Info, content []byte) error {
	var errs multierror.Errors

	// Using two-phase processing as policy name is unavailable beforehand
	var snapshotPoliciesMap map[string]json.RawMessage
	err := json.Unmarshal(content, &snapshotPoliciesMap)
	if err != nil {
		errs = append(errs, errors.Wrap(err, "failure mapping Elasticsearch Snapshot API response to JSON"))
	}

	var event mb.Event
	event.RootFields = common.MapStr{}
	event.RootFields.Put("service.name", elasticsearch.ModuleName)

	event.ModuleFields = common.MapStr{}
	event.ModuleFields.Put("cluster.name", info.ClusterName)
	event.ModuleFields.Put("cluster.id", info.ClusterID)

	for key, value := range snapshotPoliciesMap {
		var snapshotPolicy = make(map[string]interface{})
		if err := json.Unmarshal(value, &snapshotPolicy); err != nil {
			errs = append(errs, errors.Wrap(err, "failure mapping Elasticsearch Snapshot API response"))
			continue
		}
		logp.L().Debugf("snapshot API response: %v", fmt.Sprint(snapshotPolicy))

		if event.MetricSetFields, err = schema.Apply(snapshotPolicy); err != nil {
			errs = append(errs, errors.Wrap(err, "failure applying snapshot schema"))
			continue
		}

		if _, err := event.MetricSetFields.Put("name", key); err != nil {
			errs = append(errs, errors.Wrap(err, "failure extending snapshot schema with name field"))
			continue
		}
	}
	r.Event(event)
	return errs.Err()
}

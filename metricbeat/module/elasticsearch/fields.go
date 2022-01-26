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

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package elasticsearch

import (
	"github.com/elastic/beats/v7/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "elasticsearch", asset.ModuleFieldsPri, AssetElasticsearch); err != nil {
		panic(err)
	}
}

// AssetElasticsearch returns asset data.
// This is the base64 encoded zlib format compressed contents of module/elasticsearch.
func AssetElasticsearch() string {
	return "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"
}

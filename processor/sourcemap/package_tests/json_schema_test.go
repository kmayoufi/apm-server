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

package package_tests

import (
	"testing"

	sm "github.com/elastic/apm-server/processor/sourcemap"
	"github.com/elastic/apm-server/processor/sourcemap/generated/schema"
	"github.com/elastic/apm-server/tests"
)

//Check whether attributes are added to the example payload but not to the schema
func TestPayloadAttributesInSchema(t *testing.T) {

	tests.TestPayloadAttributesInSchema(t,
		"sourcemap",
		tests.NewSet("sourcemap", "sourcemap.file", "sourcemap.names", "sourcemap.sources", "sourcemap.sourceRoot",
			"sourcemap.mappings", "sourcemap.sourcesContent", "sourcemap.version"),
		schema.PayloadSchema)
}

var (
	procSetup = tests.ProcessorSetup{
		Proc:            sm.NewProcessor(),
		FullPayloadPath: "../testdata/sourcemap/payload.json",
		TemplatePaths:   []string{"../_meta/fields.yml"},
	}
)

func TestAttributesPresenceInSourcemap(t *testing.T) {
	requiredKeys := tests.NewSet("service_name", "service_version",
		"bundle_filepath", "sourcemap")
	procSetup.AttrsPresence(t, requiredKeys, nil)
}

func TestKeywordLimitationOnSourcemapAttributes(t *testing.T) {
	mapping := map[string]string{
		"sourcemap.service.name":    "service_name",
		"sourcemap.service.version": "service_version",
		"sourcemap.bundle_filepath": "bundle_filepath",
	}
	procSetup.KeywordLimitation(t, tests.NewSet(), mapping)
}

func TestPayloadDataForSourcemap(t *testing.T) {
	type val = []interface{}
	payloadData := []tests.SchemaTestData{
		// add test data for testing
		// * specific edge cases
		// * multiple allowed dataypes
		// * regex pattern, time formats
		// * length restrictions, other than keyword length restrictions

		{Key: "sourcemap", Invalid: []tests.Invalid{
			{Msg: `error validating sourcemap`, Values: val{""}},
			{Msg: `sourcemap not in expected format`, Values: val{[]byte{}}}}},
		{Key: "service_name", Valid: val{tests.Str1024},
			Invalid: []tests.Invalid{
				{Msg: `service_name/minlength`, Values: val{""}},
				{Msg: `service_name/maxlength`, Values: val{tests.Str1025}},
				{Msg: `service_name/pattern`, Values: val{tests.Str1024Special}}}},
		{Key: "service_version", Valid: val{tests.Str1024},
			Invalid: []tests.Invalid{{Msg: `service_version/minlength`, Values: val{""}}}},
		{Key: "bundle_filepath", Valid: []interface{}{tests.Str1024},
			Invalid: []tests.Invalid{{Msg: `bundle_filepath/minlength`, Values: val{""}}}},
	}
	procSetup.DataValidation(t, payloadData)
}

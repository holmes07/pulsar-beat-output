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

package testing

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/mb"
)

var (
	// Use `go test -data` to update files.
	dataFlag = flag.Bool("data", false, "Write updated data.json files")
)

// WriteEvent fetches a single event writes the output to a ./_meta/data.json
// file.
func WriteEvent(f mb.EventFetcher, t testing.TB) error {
	if !*dataFlag {
		t.Skip("skip data generation tests")
	}

	event, err := f.Fetch()
	if err != nil {
		return err
	}

	fullEvent := CreateFullEvent(f, event)
	WriteEventToDataJSON(t, fullEvent, ".")
	return nil
}

// WriteEvents fetches events and writes the first event to a ./_meta/data.json
// file.
func WriteEvents(f mb.EventsFetcher, t testing.TB) error {
	return WriteEventsCond(f, t, nil)

}

// WriteEventsCond fetches events and writes the first event that matches the condition
// to a ./_meta/data.json file.
func WriteEventsCond(f mb.EventsFetcher, t testing.TB, cond func(e common.MapStr) bool) error {
	if !*dataFlag {
		t.Skip("skip data generation tests")
	}

	events, err := f.Fetch()
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return fmt.Errorf("no events were generated")
	}

	event, err := SelectEvent(events, cond)
	if err != nil {
		return err
	}

	fullEvent := CreateFullEvent(f, event)
	WriteEventToDataJSON(t, fullEvent, "")
	return nil
}

// WriteEventsReporterV2 fetches events and writes the first event to a ./_meta/data.json
// file.
func WriteEventsReporterV2(f mb.ReportingMetricSetV2, t testing.TB, path string) error {
	return WriteEventsReporterV2Cond(f, t, path, nil)
}

// WriteEventsReporterV2Cond fetches events and writes the first event that matches
// the condition to a file.
func WriteEventsReporterV2Cond(f mb.ReportingMetricSetV2, t testing.TB, path string, cond func(common.MapStr) bool) error {
	if !*dataFlag {
		t.Skip("skip data generation tests")
	}

	events, errs := ReportingFetchV2(f)
	if len(errs) > 0 {
		return errs[0]
	}

	if len(events) == 0 {
		return fmt.Errorf("no events were generated")
	}

	match, err := SelectEventV2(f, events, cond)
	if err != nil {
		return err
	}

	e := StandardizeEvent(f, match, mb.AddMetricSetInfo)

	WriteEventToDataJSON(t, e, path)
	return nil
}

// CreateFullEvent builds a full event given the data generated by a MetricSet.
// This simulates the output of Metricbeat as if it were
// 2016-05-23T08:05:34.853Z and the hostname is host.example.com.
func CreateFullEvent(ms mb.MetricSet, metricSetData common.MapStr) beat.Event {
	return StandardizeEvent(
		ms,
		mb.TransformMapStrToEvent(ms.Module().Name(), metricSetData, nil),
		mb.AddMetricSetInfo,
	)
}

// StandardizeEvent builds a beat.Event given the data generated by a MetricSet.
// This simulates the output as if it were 2016-05-23T08:05:34.853Z and the
// hostname is host.example.com and the RTT is 155us.
func StandardizeEvent(ms mb.MetricSet, e mb.Event, modifiers ...mb.EventModifier) beat.Event {
	startTime, err := time.Parse(time.RFC3339Nano, "2017-10-12T08:05:34.853Z")
	if err != nil {
		panic(err)
	}

	e.Timestamp = startTime
	e.Took = 115 * time.Microsecond
	e.Host = ms.Host()
	if e.Namespace == "" {
		e.Namespace = ms.Registration().Namespace
	}

	fullEvent := e.BeatEvent(ms.Module().Name(), ms.Name(), modifiers...)

	fullEvent.Fields["agent"] = common.MapStr{
		"name":     "host.example.com",
		"hostname": "host.example.com",
	}

	return fullEvent
}

// WriteEventToDataJSON writes the given event as "pretty" JSON to
// a ./_meta/data.json file. If the -data CLI flag is unset or false then the
// method is a no-op.
func WriteEventToDataJSON(t testing.TB, fullEvent beat.Event, postfixPath string) {
	if !*dataFlag {
		return
	}

	p, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if stat, err := os.Stat(postfixPath); postfixPath == "" || (err == nil && stat.IsDir()) {
		p = path.Join(p, postfixPath, "_meta", "data.json")
	} else {
		p = postfixPath
	}

	fields := fullEvent.Fields
	fields["@timestamp"] = fullEvent.Timestamp

	output, err := json.MarshalIndent(&fullEvent.Fields, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	if err = ioutil.WriteFile(p, output, 0644); err != nil {
		t.Fatal(err)
	}
}

// SelectEvent selects the first event that matches an specific condition
func SelectEvent(events []common.MapStr, cond func(e common.MapStr) bool) (common.MapStr, error) {
	if cond == nil && len(events) > 0 {
		return events[0], nil
	}
	for _, e := range events {
		if cond(e) {
			return e, nil
		}
	}
	return nil, fmt.Errorf("no events satisfied the condition")
}

// SelectEventV2 selects the first event that matches an specific condition
func SelectEventV2(f mb.ReportingMetricSetV2, events []mb.Event, cond func(e common.MapStr) bool) (mb.Event, error) {
	if cond == nil && len(events) > 0 {
		return events[0], nil
	}
	for _, e := range events {
		fields := StandardizeEvent(f, e, mb.AddMetricSetInfo).Fields
		if cond(fields) {
			return e, nil
		}
	}
	return mb.Event{}, fmt.Errorf("no events satisfied the condition")
}
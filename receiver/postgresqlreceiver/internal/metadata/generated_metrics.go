// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/model/pdata"
)

// Type is the component type name.
const Type config.Type = "postgresql"

// MetricIntf is an interface to generically interact with generated metric.
type MetricIntf interface {
	Name() string
	New() pdata.Metric
	Init(metric pdata.Metric)
}

// Intentionally not exposing this so that it is opaque and can change freely.
type metricImpl struct {
	name     string
	initFunc func(pdata.Metric)
}

// Name returns the metric name.
func (m *metricImpl) Name() string {
	return m.name
}

// New creates a metric object preinitialized.
func (m *metricImpl) New() pdata.Metric {
	metric := pdata.NewMetric()
	m.Init(metric)
	return metric
}

// Init initializes the provided metric object.
func (m *metricImpl) Init(metric pdata.Metric) {
	m.initFunc(metric)
}

type metricStruct struct {
	PostgresqlBackends   MetricIntf
	PostgresqlBlocksRead MetricIntf
	PostgresqlCommits    MetricIntf
	PostgresqlDbSize     MetricIntf
	PostgresqlOperations MetricIntf
	PostgresqlRollbacks  MetricIntf
	PostgresqlRows       MetricIntf
}

// Names returns a list of all the metric name strings.
func (m *metricStruct) Names() []string {
	return []string{
		"postgresql.backends",
		"postgresql.blocks_read",
		"postgresql.commits",
		"postgresql.db_size",
		"postgresql.operations",
		"postgresql.rollbacks",
		"postgresql.rows",
	}
}

var metricsByName = map[string]MetricIntf{
	"postgresql.backends":    Metrics.PostgresqlBackends,
	"postgresql.blocks_read": Metrics.PostgresqlBlocksRead,
	"postgresql.commits":     Metrics.PostgresqlCommits,
	"postgresql.db_size":     Metrics.PostgresqlDbSize,
	"postgresql.operations":  Metrics.PostgresqlOperations,
	"postgresql.rollbacks":   Metrics.PostgresqlRollbacks,
	"postgresql.rows":        Metrics.PostgresqlRows,
}

func (m *metricStruct) ByName(n string) MetricIntf {
	return metricsByName[n]
}

// Metrics contains a set of methods for each metric that help with
// manipulating those metrics.
var Metrics = &metricStruct{
	&metricImpl{
		"postgresql.backends",
		func(metric pdata.Metric) {
			metric.SetName("postgresql.backends")
			metric.SetDescription("The number of backends.")
			metric.SetUnit("1")
			metric.SetDataType(pdata.MetricDataTypeSum)
			metric.Sum().SetIsMonotonic(false)
			metric.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
		},
	},
	&metricImpl{
		"postgresql.blocks_read",
		func(metric pdata.Metric) {
			metric.SetName("postgresql.blocks_read")
			metric.SetDescription("The number of blocks read.")
			metric.SetUnit("1")
			metric.SetDataType(pdata.MetricDataTypeSum)
			metric.Sum().SetIsMonotonic(true)
			metric.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
		},
	},
	&metricImpl{
		"postgresql.commits",
		func(metric pdata.Metric) {
			metric.SetName("postgresql.commits")
			metric.SetDescription("The number of commits.")
			metric.SetUnit("1")
			metric.SetDataType(pdata.MetricDataTypeSum)
			metric.Sum().SetIsMonotonic(true)
			metric.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
		},
	},
	&metricImpl{
		"postgresql.db_size",
		func(metric pdata.Metric) {
			metric.SetName("postgresql.db_size")
			metric.SetDescription("The database disk usage.")
			metric.SetUnit("By")
			metric.SetDataType(pdata.MetricDataTypeSum)
			metric.Sum().SetIsMonotonic(false)
			metric.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
		},
	},
	&metricImpl{
		"postgresql.operations",
		func(metric pdata.Metric) {
			metric.SetName("postgresql.operations")
			metric.SetDescription("The number of db row operations.")
			metric.SetUnit("1")
			metric.SetDataType(pdata.MetricDataTypeSum)
			metric.Sum().SetIsMonotonic(true)
			metric.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
		},
	},
	&metricImpl{
		"postgresql.rollbacks",
		func(metric pdata.Metric) {
			metric.SetName("postgresql.rollbacks")
			metric.SetDescription("The number of rollbacks.")
			metric.SetUnit("1")
			metric.SetDataType(pdata.MetricDataTypeSum)
			metric.Sum().SetIsMonotonic(true)
			metric.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
		},
	},
	&metricImpl{
		"postgresql.rows",
		func(metric pdata.Metric) {
			metric.SetName("postgresql.rows")
			metric.SetDescription("The number of rows in the database.")
			metric.SetUnit("1")
			metric.SetDataType(pdata.MetricDataTypeSum)
			metric.Sum().SetIsMonotonic(false)
			metric.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
		},
	},
}

// M contains a set of methods for each metric that help with
// manipulating those metrics. M is an alias for Metrics
var M = Metrics

// Attributes contains the possible metric attributes that can be used.
var Attributes = struct {
	// Database (The name of the database.)
	Database string
	// Operation (The database operation.)
	Operation string
	// Source (The block read source type.)
	Source string
	// State (The tuple (row) state.)
	State string
	// Table (The schema name followed by the table name.)
	Table string
}{
	"database",
	"operation",
	"source",
	"state",
	"table",
}

// A is an alias for Attributes.
var A = Attributes

// AttributeOperation are the possible values that the attribute "operation" can have.
var AttributeOperation = struct {
	Ins    string
	Upd    string
	Del    string
	HotUpd string
}{
	"ins",
	"upd",
	"del",
	"hot_upd",
}

// AttributeSource are the possible values that the attribute "source" can have.
var AttributeSource = struct {
	HeapRead  string
	HeapHit   string
	IdxRead   string
	IdxHit    string
	ToastRead string
	ToastHit  string
	TidxRead  string
	TidxHit   string
}{
	"heap_read",
	"heap_hit",
	"idx_read",
	"idx_hit",
	"toast_read",
	"toast_hit",
	"tidx_read",
	"tidx_hit",
}

// AttributeState are the possible values that the attribute "state" can have.
var AttributeState = struct {
	Dead string
	Live string
}{
	"dead",
	"live",
}

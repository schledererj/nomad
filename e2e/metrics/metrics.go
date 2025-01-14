package metrics

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/nomad/e2e/e2eutil"
	"github.com/hashicorp/nomad/e2e/framework"
	"github.com/hashicorp/nomad/helper/uuid"
	"github.com/prometheus/common/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MetricsTest struct {
	framework.TC
	jobIDs       []string
	prometheusID string
	fabioID      string
	fabioAddress string
}

func init() {
	framework.AddSuites(&framework.TestSuite{
		Component:   "Metrics",
		CanRunLocal: true,
		Cases: []framework.TestCase{
			new(MetricsTest),
		},
	})
}

// Stand up prometheus to collect metrics from all clients and allocs,
// with fabio as a system job in front of it so that we don't need to
// have prometheus use host networking
func (tc *MetricsTest) BeforeAll(f *framework.F) {
	t := f.T()
	e2eutil.WaitForLeader(t, tc.Nomad())
	e2eutil.WaitForNodesReady(t, tc.Nomad(), 1)
	err := tc.setUpPrometheus(f)
	require.Nil(t, err)
}

// Clean up the target jobs after each test case, but keep fabio/prometheus
// for reuse between the two test cases (Windows vs Linux)
func (tc *MetricsTest) AfterEach(f *framework.F) {
	if os.Getenv("NOMAD_TEST_SKIPCLEANUP") == "1" {
		return
	}
	for _, jobID := range tc.jobIDs {
		tc.Nomad().Jobs().Deregister(jobID, true, nil)
	}
	tc.jobIDs = []string{}
	tc.Nomad().System().GarbageCollect()
}

// Clean up fabio/prometheus
func (tc *MetricsTest) AfterAll(f *framework.F) {
	if os.Getenv("NOMAD_TEST_SKIPCLEANUP") == "1" {
		return
	}
	tc.tearDownPrometheus(f)
}

// TestMetricsLinux runs a collection of jobs that exercise alloc metrics.
// Then we query prometheus to verify we're collecting client and alloc metrics
// and correctly presenting them to the prometheus scraper.
func (tc *MetricsTest) TestMetricsLinux(f *framework.F) {
	t := f.T()
	clientNodes, err := e2eutil.ListLinuxClientNodes(tc.Nomad())
	require.Nil(t, err)
	if len(clientNodes) == 0 {
		t.Skip("no Linux clients")
	}

	workloads := map[string]string{
		"cpustress":  "nomad_client_allocs_cpu_user",
		"diskstress": "nomad_client_allocs_memory_rss", // TODO(tgross): do we have disk stats?
		"helloworld": "nomad_client_allocs_cpu_allocated",
		"memstress":  "nomad_client_allocs_memory_usage",
		"simpleweb":  "nomad_client_allocs_memory_rss",
	}

	tc.runWorkloads(t, workloads)
	tc.queryClientMetrics(t, clientNodes)
	tc.queryAllocMetrics(t, workloads)
}

// Run workloads from and wait for allocations
func (tc *MetricsTest) runWorkloads(t *testing.T, workloads map[string]string) {
	for jobName := range workloads {
		uuid := uuid.Generate()
		jobID := "metrics-" + jobName + "-" + uuid[0:8]
		tc.jobIDs = append(tc.jobIDs, jobID)
		file := "metrics/input/" + jobName + ".nomad"
		allocs := e2eutil.RegisterAndWaitForAllocs(t, tc.Nomad(), file, jobID)
		if len(allocs) == 0 {
			t.Fatalf("failed to register %s", jobID)
		}
	}
}

// query prometheus to verify that metrics are being collected
// from clients
func (tc *MetricsTest) queryClientMetrics(t *testing.T, clientNodes []string) {
	metrics := []string{
		"nomad_client_allocated_memory",
		"nomad_client_host_cpu_user",
		"nomad_client_host_disk_available",
		"nomad_client_host_memory_used",
		"nomad_client_uptime",
	}
	// we start with a very long timeout here because it takes a while for
	// prometheus to be live and for jobs to initially register metrics.
	timeout := 60 * time.Second

	for _, metric := range metrics {
		var results model.Vector
		var err error
		ok := assert.Eventually(t, func() bool {
			results, err = tc.promQuery(metric)
			if err != nil {
				return false
			}
			instances := make(map[model.LabelValue]struct{})
			for _, result := range results {
				instances[result.Metric["instance"]] = struct{}{}
			}
			if len(instances) != len(clientNodes) {
				err = fmt.Errorf("expected metric '%s' for all clients. got:\n%v",
					metric, results)
				return false
			}
			return true
		}, timeout, 1*time.Second)
		require.Truef(t, ok, "prometheus query failed: %v", err)

		// shorten the timeout after the first workload is successfully
		// queried so that we don't hang the whole test run if something's
		// wrong with only one of the jobs
		timeout = 10 * time.Second
	}
}

// query promtheus to verify that metrics are being collected
// from allocations
func (tc *MetricsTest) queryAllocMetrics(t *testing.T, workloads map[string]string) {
	// we start with a very long timeout here because it takes a while for
	// prometheus to be live and for jobs to initially register metrics.
	timeout := 60 * time.Second
	for jobName, metric := range workloads {
		query := fmt.Sprintf("%s{exported_job=\"%s\"}", metric, jobName)
		var results model.Vector
		var err error
		ok := assert.Eventually(t, func() bool {
			results, err = tc.promQuery(query)
			if err != nil {
				return false
			}
			return true
		}, timeout, 1*time.Second)
		require.Truef(t, ok, "prometheus query failed: %v", err)

		// make sure we didn't just collect a bunch of zero metrics
		lastResult := results[len(results)-1]
		require.Greaterf(t, float64(lastResult.Value), 0.0,
			"expected non-zero metrics: %v", results,
		)

		// shorten the timeout after the first workload is successfully
		// queried so that we don't hang the whole test run if something's
		// wrong with only one of the jobs
		timeout = 10 * time.Second
	}
}

// TestMetricsWindows runs a collection of jobs that exercise alloc metrics.
// Then we query prometheus to verify we're collecting client and alloc metrics
// and correctly presenting them to the prometheus scraper.
func (tc *MetricsTest) TestMetricsWindows(f *framework.F) {
	t := f.T()
	clientNodes, err := e2eutil.ListWindowsClientNodes(tc.Nomad())
	require.Nil(t, err)
	if len(clientNodes) == 0 {
		t.Skip("no Windows clients")
	}

	// TODO(tgross): run metrics on Windows, too
}

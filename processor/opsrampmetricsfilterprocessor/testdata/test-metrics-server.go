package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Metrics that should pass through the filter (based on alert definitions)
	nodeCPUSecondsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "node_cpu_seconds_total",
			Help: "Seconds the CPUs spent in each mode",
		},
		[]string{"mode"},
	)

	nodeMemoryAvailableBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "node_memory_MemAvailable_bytes",
			Help: "Available memory in bytes",
		},
	)

	nodeMemoryTotalBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "node_memory_MemTotal_bytes",
			Help: "Total memory in bytes",
		},
	)

	up = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "up",
			Help: "Whether the target is up",
		},
		[]string{"job"},
	)

	httpRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
	)

	// Additional metrics from alert definitions
	apiserverRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "apiserver_request_total",
			Help: "Total API server requests",
		},
		[]string{"verb", "code"},
	)

	k8sNodeConditionReady = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_node_condition_ready",
			Help: "Node ready condition",
		},
		[]string{"node"},
	)

	k8sNodeConditionDiskPressure = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_node_condition_disk_pressure",
			Help: "Node disk pressure condition",
		},
		[]string{"node"},
	)

	k8sNodeConditionMemoryPressure = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_node_condition_memory_pressure",
			Help: "Node memory pressure condition",
		},
		[]string{"node"},
	)

	k8sNodeConditionNetworkUnavailable = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_node_condition_network_unavailable",
			Help: "Node network unavailable condition",
		},
		[]string{"node"},
	)

	k8sNodeConditionPidPressure = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_node_condition_pid_pressure",
			Help: "Node PID pressure condition",
		},
		[]string{"node"},
	)

	k8sPodCPULimitUtilizationRatio = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_pod_cpu_limit_utilization_ratio",
			Help: "Pod CPU limit utilization ratio",
		},
		[]string{"pod", "namespace"},
	)

	k8sPodMemoryLimitUtilizationRatio = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_pod_memory_limit_utilization_ratio",
			Help: "Pod memory limit utilization ratio",
		},
		[]string{"pod", "namespace"},
	)

	k8sNodeMemoryWorkingSet = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_node_memory_working_set",
			Help: "Node memory working set",
		},
		[]string{"k8s_node_name"},
	)

	k8sNodeMemoryAvailable = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_node_memory_available",
			Help: "Node memory available",
		},
		[]string{"k8s_node_name"},
	)

	// Additional metrics needed for alert definitions
	k8sPodMemoryUsageBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_pod_memory_usage_bytes",
			Help: "Pod memory usage in bytes",
		},
		[]string{"pod", "namespace"},
	)

	k8sNodeCPUUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_node_cpu_usage",
			Help: "Node CPU usage",
		},
		[]string{"k8s_node_name"},
	)

	k8sNodeAllocatableCPU = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_node_allocatable_cpu",
			Help: "Node allocatable CPU",
		},
		[]string{"k8s_node_name"},
	)

	k8sPodCPULimitUtilization = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_pod_cpu_limit_utilization",
			Help: "Pod CPU limit utilization",
		},
		[]string{"pod", "namespace"},
	)

	// Metrics that should be filtered OUT (not in alert definitions)
	unrelatedMetric = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "unrelated_metric",
			Help: "This should be filtered out",
		},
	)

	anotherMetric = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "another_metric",
			Help: "This should also be filtered out",
		},
	)

	randomMetric = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "random_metric",
			Help: "Random metric that should be filtered out",
		},
	)
)

func init() {
	// Register all metrics
	prometheus.MustRegister(
		nodeCPUSecondsTotal,
		nodeMemoryAvailableBytes,
		nodeMemoryTotalBytes,
		up,
		httpRequestsTotal,
		apiserverRequestTotal,
		k8sNodeConditionReady,
		k8sNodeConditionDiskPressure,
		k8sNodeConditionMemoryPressure,
		k8sNodeConditionNetworkUnavailable,
		k8sNodeConditionPidPressure,
		k8sPodCPULimitUtilizationRatio,
		k8sPodMemoryLimitUtilizationRatio,
		k8sNodeMemoryWorkingSet,
		k8sNodeMemoryAvailable,
		k8sPodMemoryUsageBytes,
		k8sNodeCPUUsage,
		k8sNodeAllocatableCPU,
		k8sPodCPULimitUtilization,
		unrelatedMetric,
		anotherMetric,
		randomMetric,
	)
}

func updateMetrics() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Update metrics that should pass through the filter
			nodeCPUSecondsTotal.WithLabelValues("idle").Add(rand.Float64())
			nodeCPUSecondsTotal.WithLabelValues("user").Add(rand.Float64() * 0.5)
			nodeCPUSecondsTotal.WithLabelValues("system").Add(rand.Float64() * 0.3)

			nodeMemoryAvailableBytes.Set(float64(rand.Intn(4000000000) + 1000000000))
			nodeMemoryTotalBytes.Set(8000000000)

			up.WithLabelValues("apiserver").Set(float64(rand.Intn(2)))
			up.WithLabelValues("prometheus").Set(1)

			httpRequestsTotal.Add(float64(rand.Intn(10)))

			// API server metrics
			apiserverRequestTotal.WithLabelValues("GET", "200").Add(float64(rand.Intn(50)))
			apiserverRequestTotal.WithLabelValues("POST", "201").Add(float64(rand.Intn(20)))
			apiserverRequestTotal.WithLabelValues("GET", "404").Add(float64(rand.Intn(5)))

			// Kubernetes node conditions
			nodes := []string{"node1", "node2", "node3"}
			for _, node := range nodes {
				k8sNodeConditionReady.WithLabelValues(node).Set(1)
				k8sNodeConditionDiskPressure.WithLabelValues(node).Set(float64(rand.Intn(2)))
				k8sNodeConditionMemoryPressure.WithLabelValues(node).Set(float64(rand.Intn(2)))
				k8sNodeConditionNetworkUnavailable.WithLabelValues(node).Set(float64(rand.Intn(2)))
				k8sNodeConditionPidPressure.WithLabelValues(node).Set(float64(rand.Intn(2)))

				// Node memory metrics
				k8sNodeMemoryWorkingSet.WithLabelValues(node).Set(float64(rand.Intn(4000000000) + 1000000000))
				k8sNodeMemoryAvailable.WithLabelValues(node).Set(float64(rand.Intn(2000000000) + 500000000))

				// Node CPU metrics
				k8sNodeCPUUsage.WithLabelValues(node).Set(rand.Float64() * 100)
				k8sNodeAllocatableCPU.WithLabelValues(node).Set(float64(rand.Intn(8) + 1))
			}

			// Pod metrics
			pods := []string{"pod1", "pod2", "pod3"}
			namespaces := []string{"default", "kube-system"}
			for _, pod := range pods {
				for _, ns := range namespaces {
					k8sPodCPULimitUtilizationRatio.WithLabelValues(pod, ns).Set(rand.Float64())
					k8sPodCPULimitUtilization.WithLabelValues(pod, ns).Set(rand.Float64())
					k8sPodMemoryLimitUtilizationRatio.WithLabelValues(pod, ns).Set(rand.Float64())
					k8sPodMemoryUsageBytes.WithLabelValues(pod, ns).Set(float64(rand.Intn(500000000) + 100000000))
				}
			}

			// Update metrics that should be filtered out
			unrelatedMetric.Set(rand.Float64() * 100)
			anotherMetric.Add(rand.Float64() * 5)
			randomMetric.Set(rand.Float64() * 1000)

			log.Printf("Updated metrics at %v", time.Now().Format(time.RFC3339))
		}
	}
}

func main() {
	log.Println("Starting Go metrics server on :8080...")
	log.Println("Metrics endpoint: http://localhost:8080/metrics")
	log.Println()
	log.Println("Expected metrics to PASS through filter (from alert definitions):")
	log.Println("  ✓ apiserver_request_total")
	log.Println("  ✓ k8s_pod_cpu_limit_utilization_ratio")
	log.Println("  ✓ k8s_pod_memory_usage_bytes (missing from server)")
	log.Println("  ✓ k8s_node_condition_ready")
	log.Println("  ✓ k8s_node_condition_disk_pressure")
	log.Println("  ✓ k8s_node_condition_memory_pressure")
	log.Println("  ✓ k8s_node_condition_network_unavailable")
	log.Println("  ✓ k8s_node_condition_pid_pressure")
	log.Println("  ✓ k8s_node_cpu_usage (missing from server)")
	log.Println("  ✓ k8s_node_allocatable_cpu (missing from server)")
	log.Println()
	log.Println("Expected metrics to be FILTERED OUT (not in alert definitions):")
	log.Println("  ✗ node_cpu_seconds_total")
	log.Println("  ✗ node_memory_MemAvailable_bytes")
	log.Println("  ✗ node_memory_MemTotal_bytes")
	log.Println("  ✗ up")
	log.Println("  ✗ http_requests_total")
	log.Println("  ✗ k8s_pod_memory_limit_utilization_ratio")
	log.Println("  ✗ k8s_node_memory_working_set")
	log.Println("  ✗ k8s_node_memory_available")
	log.Println("  ✗ unrelated_metric")
	log.Println("  ✗ another_metric")
	log.Println("  ✗ random_metric")
	log.Println()

	// Start updating metrics in the background
	go updateMetrics()

	// Set up HTTP server
	http.Handle("/metrics", promhttp.Handler())

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

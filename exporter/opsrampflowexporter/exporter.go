// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampflowexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opsrampflowexporter"

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	aflow "opsrampflowexporter/internal/aflow/v1"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type flowExporter struct {
	client   *grpc.ClientConn
	endpoint string
	token    string
}

func newFlowExporter(cfg *Config, set exporter.Settings) (*flowExporter, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &flowExporter{
		endpoint: cfg.Endpoint,
		token:    cfg.AuthToken,
	}, nil
}

func (l *flowExporter) Start(ctx context.Context, host component.Host) error {

	conn, err := grpc.NewClient(
		l.endpoint+":443",
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName:         l.endpoint + ":443",
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		})),
	)
	if err != nil {
		log.Fatalf("Failed to create gRPC connection: %v", err)
		return err
	}
	l.client = conn

	return nil
}

func (l *flowExporter) consumeLogs(ctx context.Context, ld plog.Logs) error {
	log.Printf("Flow exporter consumeLogs not implemented yet")
	return nil
}

func (l *flowExporter) consumeMetrics(ctx context.Context, md pmetric.Metrics) error {

	aggregatedFlows := convertMetricsToAggregatedFlows(md)

	client := aflow.NewNetflowLogsClient(l.client)

	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		"Authorization": "Bearer " + l.token,
	}))

	res, err := client.AddNetflowDetails(ctx, &aggregatedFlows)
	if err != nil {
		log.Printf("aaaaa%v", err)
		return fmt.Errorf("failed to send flows: %v", err)
	}
	fmt.Println("response is %v and err %v", res, err)

	return nil
}

func convertMetricsToAggregatedFlows(md pmetric.Metrics) aflow.AggregatedFlows {
	fmt.Println("convertMetricsToAggregatedFlows")
	flows := aflow.AggregatedFlows{
		AggregationInterval: 60,
		CollectorFlows: []*aflow.CollectorFlows{
			{
				Collector: &aflow.Collector{
					Hostname:  "ebpf-collector.example.com",
					Ip:        "192.168.1.1",
					Uuid:      "ebpf-collector-123",
					SysUptime: 3600,
					SchemaUrl: "http://example.com/schema/ebpf-collector",
				},
				ExporterFlows: []*aflow.ExporterFlows{
					{
						Exporter: &aflow.Exporter{
							Hostname:  "ebpf-exporter.example.com",
							Ip:        "192.168.1.2",
							Uuid:      "ebpf-exporter-456",
							SchemaUrl: "http://example.com/schema/ebpf-exporter",
						},
						Flows: []*aflow.Flow{},
					},
				},
			},
		},
	}

	flowMap := make(map[string]*aflow.Flow)
	fmt.Printf("%v\n", flows)
	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		rm := rms.At(i)
		sms := rm.ScopeMetrics()
		for j := 0; j < sms.Len(); j++ {
			sm := sms.At(j)
			metrics := sm.Metrics()

			for k := 0; k < metrics.Len(); k++ {
				metric := metrics.At(k)

				if metric.Type() == pmetric.MetricTypeGauge {
					dps := metric.Gauge().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dp := dps.At(l)
						attrs := dp.Attributes()

						srcIP, _ := attrs.Get("src_ip")
						dstIP, _ := attrs.Get("dst_ip")
						srcPort, _ := attrs.Get("src_port")
						dstPort, _ := attrs.Get("dst_port")
						protocol, _ := attrs.Get("protocol")

						flowKey := fmt.Sprintf("%s-%s-%d-%d-%d",
							srcIP.Str(),
							dstIP.Str(),
							srcPort.Int(),
							dstPort.Int(),
							protocol.Int())

						flow, exists := flowMap[flowKey]
						if !exists {
							flow = &aflow.Flow{
								Kind:         aflow.Flow_KIND_NETFLOW_9,
								SamplingRate: 1,
								SampleCount:  1,
								IpVersion:    4,
								Protocol:     uint32(protocol.Int()),
								EthernetType: "IPv4",
								SrcAddr:      srcIP.Str(),
								DstAddr:      dstIP.Str(),
								SrcPort:      uint32(srcPort.Int()),
								DstPort:      uint32(dstPort.Int()),
								Meta: &aflow.Meta{
									Metrics: &aflow.Meta_Metrics{},
									Native: &aflow.Meta_NativeInfo{
										InIfName:   "eth0",
										InIfIndex:  1,
										OutIfName:  "eth1",
										OutIfIndex: 2,
									},
								},
							}
							flowMap[flowKey] = flow
						}

						switch metric.Name() {
						case "network.packets":
							flow.Meta.Metrics.Pkts = uint32(dp.IntValue())
						case "network.bytes":
							flow.Meta.Metrics.Bytes = uint64(dp.IntValue())
							fmt.Printf("%v\n", flow.Meta.Metrics.Bytes)
						}
					}
				}
			}
		}
	}
	var flowSlice []*aflow.Flow
	for _, flow := range flowMap {
		flowSlice = append(flowSlice, flow)
	}

	if len(flowSlice) > 0 {
		flows.CollectorFlows[0].ExporterFlows[0].Flows = flowSlice
	}
	fmt.Printf("flows: %v", flows)
	return flows
}

package v1alpha1_test

import (
	envoy_resource "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	core_plugins "github.com/kumahq/kuma/pkg/core/plugins"
	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	"github.com/kumahq/kuma/pkg/core/xds"
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	api "github.com/kumahq/kuma/pkg/plugins/policies/meshaccesslog/api/v1alpha1"
	plugin "github.com/kumahq/kuma/pkg/plugins/policies/meshaccesslog/plugin/v1alpha1"
	test_model "github.com/kumahq/kuma/pkg/test/resources/model"
	util_proto "github.com/kumahq/kuma/pkg/util/proto"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
	envoy_common "github.com/kumahq/kuma/pkg/xds/envoy"
	. "github.com/kumahq/kuma/pkg/xds/envoy/listeners"
	"github.com/kumahq/kuma/pkg/xds/generator"
)

var _ = Describe("MeshAccessLog", func() {
	type testCase struct {
		resources         []core_xds.Resource
		toRules           core_xds.ToRules
		fromRules         core_xds.FromRules
		expectedListeners []string
	}
	DescribeTable("should generate proper Envoy config",
		func(given testCase) {
			resources := core_xds.NewResourceSet()
			for _, res := range given.resources {
				r := res
				resources.Add(&r)
			}

			context := xds_context.Context{}
			proxy := xds.Proxy{
				Dataplane: &core_mesh.DataplaneResource{
					Meta: &test_model.ResourceMeta{
						Mesh: "default",
						Name: "backend",
					},
					Spec: &mesh_proto.Dataplane{
						Networking: &mesh_proto.Dataplane_Networking{
							Inbound: []*mesh_proto.Dataplane_Networking_Inbound{
								{
									Tags: map[string]string{
										mesh_proto.ServiceTag: "backend",
									},
									Address: "127.0.0.1",
									Port:    17777,
								},
							},
							Outbound: []*mesh_proto.Dataplane_Networking_Outbound{
								{
									Address: "127.0.0.1",
									Port:    27777,
									Tags: map[string]string{
										mesh_proto.ServiceTag: "other-service",
									},
								},
							},
						},
					},
				},
				Policies: xds.MatchedPolicies{
					Dynamic: map[core_model.ResourceType]xds.TypedMatchingPolicies{
						api.MeshAccessLogType: {
							Type:      api.MeshAccessLogType,
							ToRules:   given.toRules,
							FromRules: given.fromRules,
						},
					},
				},
			}
			plugin := plugin.NewPlugin().(core_plugins.PolicyPlugin)

			Expect(plugin.Apply(resources, context, &proxy)).To(Succeed())

			for i, r := range resources.ListOf(envoy_resource.ListenerType) {
				actual, err := util_proto.ToYAML(r.Resource)
				Expect(err).ToNot(HaveOccurred())

				Expect(actual).To(MatchYAML(given.expectedListeners[i]))
			}
		},
		Entry("basic outbound route", testCase{
			resources: []core_xds.Resource{{
				Name:   "outbound",
				Origin: generator.OriginOutbound,
				Resource: NewListenerBuilder(envoy_common.APIV3).
					Configure(OutboundListener("outbound:127.0.0.1:27777", "127.0.0.1", 27777, core_xds.SocketAddressProtocolTCP)).
					Configure(FilterChain(NewFilterChainBuilder(envoy_common.APIV3).
						Configure(HttpConnectionManager("127.0.0.1:27777", false)).
						Configure(
							HttpOutboundRoute(
								"backend",
								envoy_common.Routes{{
									Clusters: []envoy_common.Cluster{envoy_common.NewCluster(
										envoy_common.WithService("backend"),
										envoy_common.WithWeight(100),
									)},
								}},
								map[string]map[string]bool{
									"kuma.io/service": {
										"web": true,
									},
								},
							),
						),
					)).MustBuild(),
			}},
			toRules: core_xds.ToRules{
				Rules: []*core_xds.Rule{
					{
						Subset: []core_xds.Tag{},
						Conf: &api.MeshAccessLog_Conf{
							Backends: []*api.MeshAccessLog_Backend{{
								File: &api.MeshAccessLog_FileBackend{
									Format: &api.MeshAccessLog_Format{
										Plain: "",
									},
									Path: "/tmp/log",
								},
							}},
						},
					},
				}},
			expectedListeners: []string{`
            address:
              socketAddress:
                address: 127.0.0.1
                portValue: 27777
            filterChains:
            - filters:
              - name: envoy.filters.network.http_connection_manager
                typedConfig:
                  '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                  accessLog:
                  - name: envoy.access_loggers.file
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.access_loggers.file.v3.FileAccessLog
                      logFormat:
                          textFormatSource:
                              inlineString: |
                                [%START_TIME%] default "%REQ(:method)% %REQ(x-envoy-original-path?:path)% %PROTOCOL%" %RESPONSE_CODE% %RESPONSE_FLAGS% %BYTES_RECEIVED% %BYTES_SENT% %DURATION% %RESP(x-envoy-upstream-service-time)% "%REQ(x-forwarded-for)%" "%REQ(user-agent)%" "%REQ(x-b3-traceid?x-datadog-traceid)%" "%REQ(x-request-id)%" "%REQ(:authority)%" "backend" "other-service" "" "%UPSTREAM_HOST%"
                      path: /tmp/log
                  httpFilters:
                  - name: envoy.filters.http.router
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
                  routeConfig:
                    name: outbound:backend
                    validateClusters: false
                    requestHeadersToAdd:
                    - header:
                        key: x-kuma-tags
                        value: '&kuma.io/service=web&'
                    virtualHosts:
                    - domains:
                      - '*'
                      name: backend
                      routes:
                      - match:
                          prefix: /
                        route:
                          cluster: backend
                          timeout: 0s
                  statPrefix: "127_0_0_1_27777"
            name: outbound:127.0.0.1:27777
            trafficDirection: OUTBOUND`,
			}}),
		Entry("basic outbound tcp", testCase{
			resources: []core_xds.Resource{{
				Name:   "outbound",
				Origin: generator.OriginOutbound,
				Resource: NewListenerBuilder(envoy_common.APIV3).
					Configure(OutboundListener("outbound:127.0.0.1:27777", "127.0.0.1", 27777, core_xds.SocketAddressProtocolTCP)).
					Configure(FilterChain(NewFilterChainBuilder(envoy_common.APIV3).
						Configure(TcpProxy(
							"127.0.0.1:27777",
							envoy_common.NewCluster(
								envoy_common.WithService("backend"),
								envoy_common.WithWeight(100),
							),
						)),
					)).MustBuild(),
			}},
			toRules: core_xds.ToRules{
				Rules: []*core_xds.Rule{
					{
						Subset: []core_xds.Tag{},
						Conf: &api.MeshAccessLog_Conf{
							Backends: []*api.MeshAccessLog_Backend{{
								File: &api.MeshAccessLog_FileBackend{
									Format: &api.MeshAccessLog_Format{
										Plain: "",
									},
									Path: "/tmp/log",
								},
							}},
						},
					},
				}},
			expectedListeners: []string{`
            address:
              socketAddress:
                address: 127.0.0.1
                portValue: 27777
            filterChains:
                - filters:
                    - name: envoy.filters.network.tcp_proxy
                      typedConfig:
                        '@type': type.googleapis.com/envoy.extensions.filters.network.tcp_proxy.v3.TcpProxy
                        accessLog:
                            - name: envoy.access_loggers.file
                              typedConfig:
                                '@type': type.googleapis.com/envoy.extensions.access_loggers.file.v3.FileAccessLog
                                logFormat:
                                    textFormatSource:
                                        inlineString: |
                                            [%START_TIME%] %RESPONSE_FLAGS% default (backend)->%UPSTREAM_HOST%(other-service) took %DURATION%ms, sent %BYTES_SENT% bytes, received: %BYTES_RECEIVED% bytes
                                path: /tmp/log
                        cluster: backend
                        statPrefix: "127_0_0_1_27777"
            name: outbound:127.0.0.1:27777
            trafficDirection: OUTBOUND`,
			}}),
		Entry("basic outbound route without match", testCase{
			resources: []core_xds.Resource{{
				Name:   "outbound",
				Origin: generator.OriginOutbound,
				Resource: NewListenerBuilder(envoy_common.APIV3).
					Configure(OutboundListener("outbound:127.0.0.1:27777", "127.0.0.1", 27777, core_xds.SocketAddressProtocolTCP)).
					Configure(FilterChain(NewFilterChainBuilder(envoy_common.APIV3).
						Configure(HttpConnectionManager("127.0.0.1:27777", false)).
						Configure(
							HttpOutboundRoute(
								"backend",
								envoy_common.Routes{
									{
										Clusters: []envoy_common.Cluster{envoy_common.NewCluster(
											envoy_common.WithService("backend"),
											envoy_common.WithWeight(100),
										)},
									},
								},
								map[string]map[string]bool{
									"kuma.io/service": {
										"web": true,
									},
								},
							),
						),
					)).MustBuild(),
			}},
			toRules: core_xds.ToRules{
				Rules: []*core_xds.Rule{
					{
						Subset: []core_xds.Tag{{
							Key:   mesh_proto.ServiceTag,
							Value: "other",
						}},
						Conf: &api.MeshAccessLog_Conf{
							Backends: []*api.MeshAccessLog_Backend{{
								File: &api.MeshAccessLog_FileBackend{
									Format: &api.MeshAccessLog_Format{
										Plain: "",
									},
									Path: "/tmp/log",
								},
							}},
						},
					},
				}},
			expectedListeners: []string{`
            address:
              socketAddress:
                address: 127.0.0.1
                portValue: 27777
            filterChains:
            - filters:
              - name: envoy.filters.network.http_connection_manager
                typedConfig:
                  '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                  httpFilters:
                  - name: envoy.filters.http.router
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
                  routeConfig:
                    name: outbound:backend
                    validateClusters: false
                    requestHeadersToAdd:
                    - header:
                        key: x-kuma-tags
                        value: '&kuma.io/service=web&'
                    virtualHosts:
                    - domains:
                      - '*'
                      name: backend
                      routes:
                      - match:
                          prefix: /
                        route:
                          cluster: backend
                          timeout: 0s
                  statPrefix: "127_0_0_1_27777"
            name: outbound:127.0.0.1:27777
            trafficDirection: OUTBOUND`,
			}}),
		Entry("basic inbound route", testCase{
			resources: []core_xds.Resource{{
				Name:   "inbound",
				Origin: generator.OriginInbound,
				Resource: NewListenerBuilder(envoy_common.APIV3).
					Configure(InboundListener("inbound:127.0.0.1:17777", "127.0.0.1", 17777, core_xds.SocketAddressProtocolTCP)).
					Configure(FilterChain(NewFilterChainBuilder(envoy_common.APIV3).
						Configure(HttpConnectionManager("127.0.0.1:17777", false)).
						Configure(
							HttpInboundRoutes(
								"backend",
								envoy_common.Routes{
									{
										Clusters: []envoy_common.Cluster{envoy_common.NewCluster(
											envoy_common.WithService("backend"),
											envoy_common.WithWeight(100),
										)},
									},
								},
							),
						),
					)).MustBuild(),
			}},
			fromRules: core_xds.FromRules{
				Rules: map[xds.InboundListener]xds.Rules{
					{Address: "127.0.0.1", Port: 17777}: {{
						Subset: []core_xds.Tag{},
						Conf: &api.MeshAccessLog_Conf{
							Backends: []*api.MeshAccessLog_Backend{{
								File: &api.MeshAccessLog_FileBackend{
									Format: &api.MeshAccessLog_Format{
										Plain: "",
									},
									Path: "/tmp/log",
								},
							}},
						},
					}},
				}},
			expectedListeners: []string{`
            address:
              socketAddress:
                address: 127.0.0.1
                portValue: 17777
            enableReusePort: false
            filterChains:
            - filters:
              - name: envoy.filters.network.http_connection_manager
                typedConfig:
                  '@type': type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                  accessLog:
                  - name: envoy.access_loggers.file
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.access_loggers.file.v3.FileAccessLog
                      logFormat:
                          textFormatSource:
                              inlineString: |
                                [%START_TIME%] default "%REQ(:method)% %REQ(x-envoy-original-path?:path)% %PROTOCOL%" %RESPONSE_CODE% %RESPONSE_FLAGS% %BYTES_RECEIVED% %BYTES_SENT% %DURATION% %RESP(x-envoy-upstream-service-time)% "%REQ(x-forwarded-for)%" "%REQ(user-agent)%" "%REQ(x-b3-traceid?x-datadog-traceid)%" "%REQ(x-request-id)%" "%REQ(:authority)%" "unknown" "backend" "" "%UPSTREAM_HOST%"
                      path: /tmp/log
                  httpFilters:
                  - name: envoy.filters.http.router
                    typedConfig:
                      '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
                  routeConfig:
                    name: inbound:backend
                    validateClusters: false
                    requestHeadersToRemove:
                    - x-kuma-tags
                    virtualHosts:
                    - domains:
                      - '*'
                      name: backend
                      routes:
                      - match:
                          prefix: /
                        route:
                          cluster: backend
                          timeout: 0s
                  statPrefix: "127_0_0_1_17777"
            name: inbound:127.0.0.1:17777
            trafficDirection: INBOUND`,
			}}),
	)
})
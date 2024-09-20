// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package sidecar

import (
	"context"
	"fmt"

	_struct "github.com/golang/protobuf/ptypes/struct"
	"go.uber.org/zap"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	istionet "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/5GSEC/SentryFlow/pkg/util"
)

const (
	FilterName                        = "http-filter"
	UpstreamAndClusterName            = "sentryflow"
	ApiPath                           = "/api/v1/events"
	FilterURI                         = "https://raw.githubusercontent.com/anurag-rajawat/envoy-wasm-filters/main/httpfilters.wasm"
	RemoteWasmFilterClusterName       = "remote_wasm"
	FilterSha256                      = "714be6a76e8853fa331a285c8d420a740675708f1503df88370a30197f8b6e37"
	Timeout                           = "5s"
	SentryFlowDefaultFilterServerPort = 8081
)

// StartMonitoring begins monitoring API calls within the Istio (sidecar based)
// service mesh deployed in a Kubernetes cluster. It achieves this by creating a
// custom EnvoyFilter resource in Kubernetes.
func StartMonitoring(ctx context.Context, k8sClient client.Client) {
	logger := util.LoggerFromCtx(ctx).Named("istio-sidecar")

	logger.Info("Starting istio sidecar mesh monitoring")

	if err := createEnvoyFilter(ctx, k8sClient); err != nil {
		logger.Errorf("Failed to create EnvoyFilter. Stopping istio sidecar mesh monitoring, error: %v", err)
		return
	}
	logger.Info("Started istio sidecar mesh monitoring")

	<-ctx.Done()
	logger.Info("Shutting down istio sidecar mesh monitoring")
	if err := deleteEnvoyFilter(logger, k8sClient); err != nil {
		logger.Errorf("Failed to delete EnvoyFilter, error: %v", err)
	}

	logger.Info("Stopped istio sidecar mesh monitoring")
}

func deleteEnvoyFilter(logger *zap.SugaredLogger, k8sClient client.Client) error {
	existingFilter := &istionet.EnvoyFilter{}
	if err := k8sClient.Get(context.Background(), types.NamespacedName{Name: FilterName, Namespace: "istio-system"}, existingFilter); err != nil {
		return err
	}

	if err := k8sClient.Delete(context.Background(), existingFilter); err != nil {
		return err
	}
	logger.Infow("Deleted EnvoyFilter", "name", FilterName, "namespace", "istio-system")

	return nil
}

func createEnvoyFilter(ctx context.Context, k8sClient client.Client) error {
	logger := util.LoggerFromCtx(ctx)

	var configVal = fmt.Sprintf(`{"upstream_name": "%v", "authority": "%v", "api_path": "%v"}
`, UpstreamAndClusterName, UpstreamAndClusterName, ApiPath)

	filter := &istionet.EnvoyFilter{
		TypeMeta: v1.TypeMeta{
			Kind:       "EnvoyFilter",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: FilterName,
			// Deploy the filter to whatever istio considers its "root" namespace so that we
			// don't have to create the ConfigMap(s) containing the WASM filter binary, and
			// the associated annotations/configuration for the Istio sidecar(s).
			// https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#MeshConfig:~:text=No-,rootNamespace,-string
			Namespace: "istio-system",
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "sentryflow",
			},
		},
		Spec: networkingv1alpha3.EnvoyFilter{
			ConfigPatches: []*networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{
				{
					ApplyTo: networkingv1alpha3.EnvoyFilter_HTTP_FILTER,
					Match: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
						Context: networkingv1alpha3.EnvoyFilter_ANY,
						ObjectTypes: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
							Listener: &networkingv1alpha3.EnvoyFilter_ListenerMatch{
								FilterChain: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterChainMatch{
									Filter: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterMatch{
										Name: "envoy.filters.network.http_connection_manager",
										SubFilter: &networkingv1alpha3.EnvoyFilter_ListenerMatch_SubFilterMatch{
											Name: "envoy.filters.http.router",
										},
									},
								},
							},
						},
					},
					Patch: &networkingv1alpha3.EnvoyFilter_Patch{
						Operation: networkingv1alpha3.EnvoyFilter_Patch_INSERT_BEFORE,
						Value: &_struct.Struct{
							Fields: map[string]*_struct.Value{
								"name": {
									Kind: &_struct.Value_StringValue{
										StringValue: "envoy.filters.http.wasm",
									},
								},
								"typedConfig": {
									Kind: &_struct.Value_StructValue{
										StructValue: &_struct.Struct{
											Fields: map[string]*_struct.Value{
												"@type": {
													Kind: &_struct.Value_StringValue{
														StringValue: "type.googleapis.com/udpa.type.v1.TypedStruct",
													},
												},
												"typeUrl": {
													Kind: &_struct.Value_StringValue{
														StringValue: "type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm",
													},
												},
												"value": {
													Kind: &_struct.Value_StructValue{
														StructValue: &_struct.Struct{
															Fields: map[string]*_struct.Value{
																"config": {
																	Kind: &_struct.Value_StructValue{
																		StructValue: &_struct.Struct{
																			Fields: map[string]*_struct.Value{
																				"name": {
																					Kind: &_struct.Value_StringValue{
																						StringValue: FilterName,
																					},
																				},
																				"rootId": {
																					Kind: &_struct.Value_StringValue{
																						StringValue: FilterName,
																					},
																				},
																				"configuration": {
																					Kind: &_struct.Value_StructValue{
																						StructValue: &_struct.Struct{
																							Fields: map[string]*_struct.Value{
																								"@type": {
																									Kind: &_struct.Value_StringValue{
																										StringValue: "type.googleapis.com/google.protobuf.StringValue",
																									},
																								},
																								"value": {
																									Kind: &_struct.Value_StringValue{
																										StringValue: configVal,
																									},
																								},
																							},
																						},
																					},
																				},
																				"vmConfig": {
																					Kind: &_struct.Value_StructValue{
																						StructValue: &_struct.Struct{
																							Fields: map[string]*_struct.Value{
																								"code": {
																									Kind: &_struct.Value_StructValue{
																										StructValue: &_struct.Struct{
																											Fields: map[string]*_struct.Value{
																												"remote": {
																													Kind: &_struct.Value_StructValue{
																														StructValue: &_struct.Struct{
																															Fields: map[string]*_struct.Value{
																																"http_uri": {
																																	Kind: &_struct.Value_StructValue{
																																		StructValue: &_struct.Struct{
																																			Fields: map[string]*_struct.Value{
																																				"uri": {
																																					Kind: &_struct.Value_StringValue{
																																						StringValue: FilterURI,
																																					},
																																				},
																																				"timeout": {
																																					Kind: &_struct.Value_StringValue{
																																						StringValue: Timeout,
																																					},
																																				},
																																				"cluster": {
																																					Kind: &_struct.Value_StringValue{
																																						StringValue: RemoteWasmFilterClusterName,
																																					},
																																				},
																																			},
																																		},
																																	},
																																},
																																"sha256": {
																																	Kind: &_struct.Value_StringValue{
																																		StringValue: FilterSha256,
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"runtime": {
																									Kind: &_struct.Value_StringValue{
																										StringValue: "envoy.wasm.runtime.v8",
																									},
																								},
																								"vmId": {
																									Kind: &_struct.Value_StringValue{
																										StringValue: FilterName,
																									},
																								},
																								"allow_precompiled": {
																									Kind: &_struct.Value_BoolValue{
																										BoolValue: true,
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					ApplyTo: networkingv1alpha3.EnvoyFilter_CLUSTER,
					Match: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
						Context: networkingv1alpha3.EnvoyFilter_SIDECAR_OUTBOUND,
					},
					Patch: &networkingv1alpha3.EnvoyFilter_Patch{
						Operation: networkingv1alpha3.EnvoyFilter_Patch_ADD,
						Value: &_struct.Struct{
							Fields: map[string]*_struct.Value{
								"name": {
									Kind: &_struct.Value_StringValue{
										StringValue: UpstreamAndClusterName,
									},
								},
								"type": {
									Kind: &_struct.Value_StringValue{
										StringValue: "LOGICAL_DNS",
									},
								},
								"connect_timeout": {
									Kind: &_struct.Value_StringValue{
										StringValue: "1s",
									},
								},
								"lb_policy": {
									Kind: &_struct.Value_StringValue{
										StringValue: "ROUND_ROBIN",
									},
								},
								"load_assignment": {
									Kind: &_struct.Value_StructValue{
										StructValue: &_struct.Struct{
											Fields: map[string]*_struct.Value{
												"cluster_name": {
													Kind: &_struct.Value_StringValue{
														StringValue: UpstreamAndClusterName,
													},
												},
												"endpoints": {
													Kind: &_struct.Value_ListValue{
														ListValue: &_struct.ListValue{
															Values: []*_struct.Value{
																{
																	Kind: &_struct.Value_StructValue{
																		StructValue: &_struct.Struct{
																			Fields: map[string]*_struct.Value{
																				"lb_endpoints": {
																					Kind: &_struct.Value_ListValue{
																						ListValue: &_struct.ListValue{
																							Values: []*_struct.Value{{
																								Kind: &_struct.Value_StructValue{
																									StructValue: &_struct.Struct{
																										Fields: map[string]*_struct.Value{
																											"endpoint": {
																												Kind: &_struct.Value_StructValue{
																													StructValue: &_struct.Struct{
																														Fields: map[string]*_struct.Value{
																															"address": {
																																Kind: &_struct.Value_StructValue{
																																	StructValue: &_struct.Struct{
																																		Fields: map[string]*_struct.Value{
																																			"socket_address": {
																																				Kind: &_struct.Value_StructValue{
																																					StructValue: &_struct.Struct{
																																						Fields: map[string]*_struct.Value{
																																							"protocol": {
																																								Kind: &_struct.Value_StringValue{
																																									StringValue: "TCP",
																																								},
																																							},
																																							"address": {
																																								Kind: &_struct.Value_StringValue{
																																									StringValue: UpstreamAndClusterName + "." + UpstreamAndClusterName,
																																								},
																																							},
																																							"port_value": {
																																								Kind: &_struct.Value_NumberValue{
																																									NumberValue: SentryFlowDefaultFilterServerPort,
																																								},
																																							},
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																							}},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					ApplyTo: networkingv1alpha3.EnvoyFilter_CLUSTER,
					Match: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
						Context: networkingv1alpha3.EnvoyFilter_SIDECAR_OUTBOUND,
					},
					Patch: &networkingv1alpha3.EnvoyFilter_Patch{
						Operation: networkingv1alpha3.EnvoyFilter_Patch_ADD,
						Value: &_struct.Struct{
							Fields: map[string]*_struct.Value{
								"name": {
									Kind: &_struct.Value_StringValue{
										StringValue: RemoteWasmFilterClusterName,
									},
								},
								"type": {
									Kind: &_struct.Value_StringValue{
										StringValue: "STRICT_DNS",
									},
								},
								"connect_timeout": {
									Kind: &_struct.Value_StringValue{
										StringValue: "1s",
									},
								},
								"dns_refresh_rate": {
									Kind: &_struct.Value_StringValue{
										StringValue: Timeout,
									},
								},
								"dns_lookup_family": {
									Kind: &_struct.Value_StringValue{
										StringValue: "V4_ONLY",
									},
								},
								"lb_policy": {
									Kind: &_struct.Value_StringValue{
										StringValue: "ROUND_ROBIN",
									},
								},
								"load_assignment": {
									Kind: &_struct.Value_StructValue{
										StructValue: &_struct.Struct{
											Fields: map[string]*_struct.Value{
												"cluster_name": {
													Kind: &_struct.Value_StringValue{
														StringValue: RemoteWasmFilterClusterName,
													},
												},
												"endpoints": {
													Kind: &_struct.Value_ListValue{
														ListValue: &_struct.ListValue{
															Values: []*_struct.Value{
																{
																	Kind: &_struct.Value_StructValue{
																		StructValue: &_struct.Struct{
																			Fields: map[string]*_struct.Value{
																				"lb_endpoints": {
																					Kind: &_struct.Value_ListValue{
																						ListValue: &_struct.ListValue{
																							Values: []*_struct.Value{
																								{
																									Kind: &_struct.Value_StructValue{
																										StructValue: &_struct.Struct{
																											Fields: map[string]*_struct.Value{
																												"endpoint": {
																													Kind: &_struct.Value_StructValue{
																														StructValue: &_struct.Struct{
																															Fields: map[string]*_struct.Value{
																																"address": {
																																	Kind: &_struct.Value_StructValue{
																																		StructValue: &_struct.Struct{
																																			Fields: map[string]*_struct.Value{
																																				"socket_address": {
																																					Kind: &_struct.Value_StructValue{
																																						StructValue: &_struct.Struct{
																																							Fields: map[string]*_struct.Value{
																																								"address": {
																																									Kind: &_struct.Value_StringValue{
																																										StringValue: "raw.githubusercontent.com",
																																									},
																																								},
																																								"port_value": {
																																									Kind: &_struct.Value_NumberValue{
																																										NumberValue: 443,
																																									},
																																								},
																																							},
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								"transport_socket": {
									Kind: &_struct.Value_StructValue{
										StructValue: &_struct.Struct{
											Fields: map[string]*_struct.Value{
												"name": {
													Kind: &_struct.Value_StringValue{
														StringValue: "envoy.transport_sockets.tls",
													},
												},
												"typed_config": {
													Kind: &_struct.Value_StructValue{
														StructValue: &_struct.Struct{
															Fields: map[string]*_struct.Value{
																"@type": {
																	Kind: &_struct.Value_StringValue{
																		StringValue: "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext",
																	},
																},
																"sni": {
																	Kind: &_struct.Value_StringValue{
																		StringValue: "raw.githubusercontent.com",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	existingFilter := &istionet.EnvoyFilter{}
	if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(filter), existingFilter); err != nil {
		if errors.IsNotFound(err) {
			if err := k8sClient.Create(ctx, filter); err != nil {
				return err
			}
			logger.Infow("Created EnvoyFilter", "name", filter.Name, "namespace", filter.Namespace)
			return nil
		}
		return err
	}
	logger.Infow("Found EnvoyFilter", "name", filter.Name, "namespace", filter.Namespace)
	return nil
}

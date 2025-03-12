# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: sentryflow.proto
# Protobuf Python Version: 5.29.3
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    29,
    3,
    '',
    'sentryflow.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x10sentryflow.proto\x12\x08protobuf\"1\n\nClientInfo\x12\x10\n\x08hostName\x18\x01 \x01(\t\x12\x11\n\tIPAddress\x18\x02 \x01(\t\"\xe7\x03\n\x06\x41PILog\x12\n\n\x02id\x18\x01 \x01(\x04\x12\x11\n\ttimeStamp\x18\x02 \x01(\t\x12\x14\n\x0csrcNamespace\x18\x0b \x01(\t\x12\x0f\n\x07srcName\x18\x0c \x01(\t\x12\x30\n\x08srcLabel\x18\r \x03(\x0b\x32\x1e.protobuf.APILog.SrcLabelEntry\x12\x0f\n\x07srcType\x18\x15 \x01(\t\x12\r\n\x05srcIP\x18\x16 \x01(\t\x12\x0f\n\x07srcPort\x18\x17 \x01(\t\x12\x14\n\x0c\x64stNamespace\x18\x1f \x01(\t\x12\x0f\n\x07\x64stName\x18  \x01(\t\x12\x30\n\x08\x64stLabel\x18! \x03(\x0b\x32\x1e.protobuf.APILog.DstLabelEntry\x12\x0f\n\x07\x64stType\x18) \x01(\t\x12\r\n\x05\x64stIP\x18* \x01(\t\x12\x0f\n\x07\x64stPort\x18+ \x01(\t\x12\x10\n\x08protocol\x18\x33 \x01(\t\x12\x0e\n\x06method\x18\x34 \x01(\t\x12\x0c\n\x04path\x18\x35 \x01(\t\x12\x14\n\x0cresponseCode\x18\x36 \x01(\x05\x1a/\n\rSrcLabelEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\x1a/\n\rDstLabelEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01:\x02\x18\x01\"\xd9\x01\n\x08\x41PIEvent\x12$\n\x08metadata\x18\x01 \x01(\x0b\x32\x12.protobuf.Metadata\x12\"\n\x06source\x18\x03 \x01(\x0b\x32\x12.protobuf.Workload\x12\'\n\x0b\x64\x65stination\x18\x04 \x01(\x0b\x32\x12.protobuf.Workload\x12\"\n\x07request\x18\x05 \x01(\x0b\x32\x11.protobuf.Request\x12$\n\x08response\x18\x06 \x01(\x0b\x32\x12.protobuf.Response\x12\x10\n\x08protocol\x18\x07 \x01(\t\"\xa1\x01\n\x08Metadata\x12\x12\n\ncontext_id\x18\x01 \x01(\r\x12\x11\n\ttimestamp\x18\x02 \x01(\x04\x12\x19\n\ristio_version\x18\x03 \x01(\tB\x02\x18\x01\x12\x0f\n\x07mesh_id\x18\x04 \x01(\t\x12\x11\n\tnode_name\x18\x05 \x01(\t\x12\x15\n\rreceiver_name\x18\x06 \x01(\t\x12\x18\n\x10receiver_version\x18\x07 \x01(\t\"E\n\x08Workload\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x11\n\tnamespace\x18\x02 \x01(\t\x12\n\n\x02ip\x18\x03 \x01(\t\x12\x0c\n\x04port\x18\x04 \x01(\x05\"x\n\x07Request\x12/\n\x07headers\x18\x01 \x03(\x0b\x32\x1e.protobuf.Request.HeadersEntry\x12\x0c\n\x04\x62ody\x18\x02 \x01(\t\x1a.\n\x0cHeadersEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\"\x9c\x01\n\x08Response\x12\x30\n\x07headers\x18\x01 \x03(\x0b\x32\x1f.protobuf.Response.HeadersEntry\x12\x0c\n\x04\x62ody\x18\x02 \x01(\t\x12 \n\x18\x62\x61\x63kend_latency_in_nanos\x18\x03 \x01(\x04\x1a.\n\x0cHeadersEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\"\x7f\n\nAPIMetrics\x12<\n\x0cperAPICounts\x18\x01 \x03(\x0b\x32&.protobuf.APIMetrics.PerAPICountsEntry\x1a\x33\n\x11PerAPICountsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\x04:\x02\x38\x01\"l\n\x0bMetricValue\x12/\n\x05value\x18\x01 \x03(\x0b\x32 .protobuf.MetricValue.ValueEntry\x1a,\n\nValueEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\"\xb5\x02\n\x0c\x45nvoyMetrics\x12\x11\n\ttimeStamp\x18\x01 \x01(\t\x12\x11\n\tnamespace\x18\x0b \x01(\t\x12\x0c\n\x04name\x18\x0c \x01(\t\x12\x11\n\tIPAddress\x18\r \x01(\t\x12\x32\n\x06labels\x18\x0e \x03(\x0b\x32\".protobuf.EnvoyMetrics.LabelsEntry\x12\x34\n\x07metrics\x18\x15 \x03(\x0b\x32#.protobuf.EnvoyMetrics.MetricsEntry\x1a-\n\x0bLabelsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\x1a\x45\n\x0cMetricsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12$\n\x05value\x18\x02 \x01(\x0b\x32\x15.protobuf.MetricValue:\x02\x38\x01\x32\xbd\x02\n\nSentryFlow\x12:\n\tGetAPILog\x12\x14.protobuf.ClientInfo\x1a\x10.protobuf.APILog\"\x03\x88\x02\x01\x30\x01\x12\x39\n\x0bGetAPIEvent\x12\x14.protobuf.ClientInfo\x1a\x12.protobuf.APIEvent0\x01\x12\x36\n\x0cSendAPIEvent\x12\x12.protobuf.APIEvent\x1a\x12.protobuf.APIEvent\x12=\n\rGetAPIMetrics\x12\x14.protobuf.ClientInfo\x1a\x14.protobuf.APIMetrics0\x01\x12\x41\n\x0fGetEnvoyMetrics\x12\x14.protobuf.ClientInfo\x1a\x16.protobuf.EnvoyMetrics0\x01\x42-Z+github.com/5GSEC/SentryFlow/protobuf/golangb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'sentryflow_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z+github.com/5GSEC/SentryFlow/protobuf/golang'
  _globals['_APILOG_SRCLABELENTRY']._loaded_options = None
  _globals['_APILOG_SRCLABELENTRY']._serialized_options = b'8\001'
  _globals['_APILOG_DSTLABELENTRY']._loaded_options = None
  _globals['_APILOG_DSTLABELENTRY']._serialized_options = b'8\001'
  _globals['_APILOG']._loaded_options = None
  _globals['_APILOG']._serialized_options = b'\030\001'
  _globals['_METADATA'].fields_by_name['istio_version']._loaded_options = None
  _globals['_METADATA'].fields_by_name['istio_version']._serialized_options = b'\030\001'
  _globals['_REQUEST_HEADERSENTRY']._loaded_options = None
  _globals['_REQUEST_HEADERSENTRY']._serialized_options = b'8\001'
  _globals['_RESPONSE_HEADERSENTRY']._loaded_options = None
  _globals['_RESPONSE_HEADERSENTRY']._serialized_options = b'8\001'
  _globals['_APIMETRICS_PERAPICOUNTSENTRY']._loaded_options = None
  _globals['_APIMETRICS_PERAPICOUNTSENTRY']._serialized_options = b'8\001'
  _globals['_METRICVALUE_VALUEENTRY']._loaded_options = None
  _globals['_METRICVALUE_VALUEENTRY']._serialized_options = b'8\001'
  _globals['_ENVOYMETRICS_LABELSENTRY']._loaded_options = None
  _globals['_ENVOYMETRICS_LABELSENTRY']._serialized_options = b'8\001'
  _globals['_ENVOYMETRICS_METRICSENTRY']._loaded_options = None
  _globals['_ENVOYMETRICS_METRICSENTRY']._serialized_options = b'8\001'
  _globals['_SENTRYFLOW'].methods_by_name['GetAPILog']._loaded_options = None
  _globals['_SENTRYFLOW'].methods_by_name['GetAPILog']._serialized_options = b'\210\002\001'
  _globals['_CLIENTINFO']._serialized_start=30
  _globals['_CLIENTINFO']._serialized_end=79
  _globals['_APILOG']._serialized_start=82
  _globals['_APILOG']._serialized_end=569
  _globals['_APILOG_SRCLABELENTRY']._serialized_start=469
  _globals['_APILOG_SRCLABELENTRY']._serialized_end=516
  _globals['_APILOG_DSTLABELENTRY']._serialized_start=518
  _globals['_APILOG_DSTLABELENTRY']._serialized_end=565
  _globals['_APIEVENT']._serialized_start=572
  _globals['_APIEVENT']._serialized_end=789
  _globals['_METADATA']._serialized_start=792
  _globals['_METADATA']._serialized_end=953
  _globals['_WORKLOAD']._serialized_start=955
  _globals['_WORKLOAD']._serialized_end=1024
  _globals['_REQUEST']._serialized_start=1026
  _globals['_REQUEST']._serialized_end=1146
  _globals['_REQUEST_HEADERSENTRY']._serialized_start=1100
  _globals['_REQUEST_HEADERSENTRY']._serialized_end=1146
  _globals['_RESPONSE']._serialized_start=1149
  _globals['_RESPONSE']._serialized_end=1305
  _globals['_RESPONSE_HEADERSENTRY']._serialized_start=1100
  _globals['_RESPONSE_HEADERSENTRY']._serialized_end=1146
  _globals['_APIMETRICS']._serialized_start=1307
  _globals['_APIMETRICS']._serialized_end=1434
  _globals['_APIMETRICS_PERAPICOUNTSENTRY']._serialized_start=1383
  _globals['_APIMETRICS_PERAPICOUNTSENTRY']._serialized_end=1434
  _globals['_METRICVALUE']._serialized_start=1436
  _globals['_METRICVALUE']._serialized_end=1544
  _globals['_METRICVALUE_VALUEENTRY']._serialized_start=1500
  _globals['_METRICVALUE_VALUEENTRY']._serialized_end=1544
  _globals['_ENVOYMETRICS']._serialized_start=1547
  _globals['_ENVOYMETRICS']._serialized_end=1856
  _globals['_ENVOYMETRICS_LABELSENTRY']._serialized_start=1740
  _globals['_ENVOYMETRICS_LABELSENTRY']._serialized_end=1785
  _globals['_ENVOYMETRICS_METRICSENTRY']._serialized_start=1787
  _globals['_ENVOYMETRICS_METRICSENTRY']._serialized_end=1856
  _globals['_SENTRYFLOW']._serialized_start=1859
  _globals['_SENTRYFLOW']._serialized_end=2176
# @@protoc_insertion_point(module_scope)

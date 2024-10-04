from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ClientInfo(_message.Message):
    __slots__ = ("hostName", "IPAddress")
    HOSTNAME_FIELD_NUMBER: _ClassVar[int]
    IPADDRESS_FIELD_NUMBER: _ClassVar[int]
    hostName: str
    IPAddress: str
    def __init__(self, hostName: _Optional[str] = ..., IPAddress: _Optional[str] = ...) -> None: ...

class APILog(_message.Message):
    __slots__ = ("id", "timeStamp", "srcNamespace", "srcName", "srcLabel", "srcType", "srcIP", "srcPort", "dstNamespace", "dstName", "dstLabel", "dstType", "dstIP", "dstPort", "protocol", "method", "path", "responseCode")
    class SrcLabelEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    class DstLabelEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    ID_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    SRCNAMESPACE_FIELD_NUMBER: _ClassVar[int]
    SRCNAME_FIELD_NUMBER: _ClassVar[int]
    SRCLABEL_FIELD_NUMBER: _ClassVar[int]
    SRCTYPE_FIELD_NUMBER: _ClassVar[int]
    SRCIP_FIELD_NUMBER: _ClassVar[int]
    SRCPORT_FIELD_NUMBER: _ClassVar[int]
    DSTNAMESPACE_FIELD_NUMBER: _ClassVar[int]
    DSTNAME_FIELD_NUMBER: _ClassVar[int]
    DSTLABEL_FIELD_NUMBER: _ClassVar[int]
    DSTTYPE_FIELD_NUMBER: _ClassVar[int]
    DSTIP_FIELD_NUMBER: _ClassVar[int]
    DSTPORT_FIELD_NUMBER: _ClassVar[int]
    PROTOCOL_FIELD_NUMBER: _ClassVar[int]
    METHOD_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    RESPONSECODE_FIELD_NUMBER: _ClassVar[int]
    id: int
    timeStamp: str
    srcNamespace: str
    srcName: str
    srcLabel: _containers.ScalarMap[str, str]
    srcType: str
    srcIP: str
    srcPort: str
    dstNamespace: str
    dstName: str
    dstLabel: _containers.ScalarMap[str, str]
    dstType: str
    dstIP: str
    dstPort: str
    protocol: str
    method: str
    path: str
    responseCode: int
    def __init__(self, id: _Optional[int] = ..., timeStamp: _Optional[str] = ..., srcNamespace: _Optional[str] = ..., srcName: _Optional[str] = ..., srcLabel: _Optional[_Mapping[str, str]] = ..., srcType: _Optional[str] = ..., srcIP: _Optional[str] = ..., srcPort: _Optional[str] = ..., dstNamespace: _Optional[str] = ..., dstName: _Optional[str] = ..., dstLabel: _Optional[_Mapping[str, str]] = ..., dstType: _Optional[str] = ..., dstIP: _Optional[str] = ..., dstPort: _Optional[str] = ..., protocol: _Optional[str] = ..., method: _Optional[str] = ..., path: _Optional[str] = ..., responseCode: _Optional[int] = ...) -> None: ...

class APIEvent(_message.Message):
    __slots__ = ("metadata", "source", "destination", "request", "response", "protocol")
    METADATA_FIELD_NUMBER: _ClassVar[int]
    SOURCE_FIELD_NUMBER: _ClassVar[int]
    DESTINATION_FIELD_NUMBER: _ClassVar[int]
    REQUEST_FIELD_NUMBER: _ClassVar[int]
    RESPONSE_FIELD_NUMBER: _ClassVar[int]
    PROTOCOL_FIELD_NUMBER: _ClassVar[int]
    metadata: Metadata
    source: Workload
    destination: Workload
    request: Request
    response: Response
    protocol: str
    def __init__(self, metadata: _Optional[_Union[Metadata, _Mapping]] = ..., source: _Optional[_Union[Workload, _Mapping]] = ..., destination: _Optional[_Union[Workload, _Mapping]] = ..., request: _Optional[_Union[Request, _Mapping]] = ..., response: _Optional[_Union[Response, _Mapping]] = ..., protocol: _Optional[str] = ...) -> None: ...

class Metadata(_message.Message):
    __slots__ = ("context_id", "timestamp", "istio_version", "mesh_id", "node_name")
    CONTEXT_ID_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    ISTIO_VERSION_FIELD_NUMBER: _ClassVar[int]
    MESH_ID_FIELD_NUMBER: _ClassVar[int]
    NODE_NAME_FIELD_NUMBER: _ClassVar[int]
    context_id: int
    timestamp: int
    istio_version: str
    mesh_id: str
    node_name: str
    def __init__(self, context_id: _Optional[int] = ..., timestamp: _Optional[int] = ..., istio_version: _Optional[str] = ..., mesh_id: _Optional[str] = ..., node_name: _Optional[str] = ...) -> None: ...

class Workload(_message.Message):
    __slots__ = ("name", "namespace", "ip", "port")
    NAME_FIELD_NUMBER: _ClassVar[int]
    NAMESPACE_FIELD_NUMBER: _ClassVar[int]
    IP_FIELD_NUMBER: _ClassVar[int]
    PORT_FIELD_NUMBER: _ClassVar[int]
    name: str
    namespace: str
    ip: str
    port: int
    def __init__(self, name: _Optional[str] = ..., namespace: _Optional[str] = ..., ip: _Optional[str] = ..., port: _Optional[int] = ...) -> None: ...

class Request(_message.Message):
    __slots__ = ("headers", "body")
    class HeadersEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    HEADERS_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    headers: _containers.ScalarMap[str, str]
    body: str
    def __init__(self, headers: _Optional[_Mapping[str, str]] = ..., body: _Optional[str] = ...) -> None: ...

class Response(_message.Message):
    __slots__ = ("headers", "body")
    class HeadersEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    HEADERS_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    headers: _containers.ScalarMap[str, str]
    body: str
    def __init__(self, headers: _Optional[_Mapping[str, str]] = ..., body: _Optional[str] = ...) -> None: ...

class APIMetrics(_message.Message):
    __slots__ = ("perAPICounts",)
    class PerAPICountsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: int
        def __init__(self, key: _Optional[str] = ..., value: _Optional[int] = ...) -> None: ...
    PERAPICOUNTS_FIELD_NUMBER: _ClassVar[int]
    perAPICounts: _containers.ScalarMap[str, int]
    def __init__(self, perAPICounts: _Optional[_Mapping[str, int]] = ...) -> None: ...

class MetricValue(_message.Message):
    __slots__ = ("value",)
    class ValueEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    VALUE_FIELD_NUMBER: _ClassVar[int]
    value: _containers.ScalarMap[str, str]
    def __init__(self, value: _Optional[_Mapping[str, str]] = ...) -> None: ...

class EnvoyMetrics(_message.Message):
    __slots__ = ("timeStamp", "namespace", "name", "IPAddress", "labels", "metrics")
    class LabelsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    class MetricsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MetricValue
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MetricValue, _Mapping]] = ...) -> None: ...
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    NAMESPACE_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    IPADDRESS_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    METRICS_FIELD_NUMBER: _ClassVar[int]
    timeStamp: str
    namespace: str
    name: str
    IPAddress: str
    labels: _containers.ScalarMap[str, str]
    metrics: _containers.MessageMap[str, MetricValue]
    def __init__(self, timeStamp: _Optional[str] = ..., namespace: _Optional[str] = ..., name: _Optional[str] = ..., IPAddress: _Optional[str] = ..., labels: _Optional[_Mapping[str, str]] = ..., metrics: _Optional[_Mapping[str, MetricValue]] = ...) -> None: ...

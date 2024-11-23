package apispec

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

func MergePathItems(dst, src *openapi3.PathItem) *openapi3.PathItem {
	dst.Get, _ = mergeOperation(dst.Get, src.Get)
	dst.Put, _ = mergeOperation(dst.Put, src.Put)
	dst.Post, _ = mergeOperation(dst.Post, src.Post)
	dst.Delete, _ = mergeOperation(dst.Delete, src.Delete)
	dst.Options, _ = mergeOperation(dst.Options, src.Options)
	dst.Head, _ = mergeOperation(dst.Head, src.Head)
	dst.Patch, _ = mergeOperation(dst.Patch, src.Patch)

	// TODO what about merging parameters?

	return dst
}

func CopyPathItemWithNewOperation(item *openapi3.PathItem, method string, operation *openapi3.Operation) *openapi3.PathItem {
	// TODO - do we want to do : ret = *item?
	ret := openapi3.PathItem{}
	ret.Get = item.Get
	ret.Put = item.Put
	ret.Patch = item.Patch
	ret.Post = item.Post
	ret.Head = item.Head
	ret.Delete = item.Delete
	ret.Options = item.Options
	ret.Parameters = item.Parameters

	AddOperationToPathItem(&ret, method, operation)
	return &ret
}

func GetOperationFromPathItem(item *openapi3.PathItem, method string) *openapi3.Operation {
	switch method {
	case http.MethodGet:
		return item.Get
	case http.MethodDelete:
		return item.Delete
	case http.MethodOptions:
		return item.Options
	case http.MethodPatch:
		return item.Patch
	case http.MethodHead:
		return item.Head
	case http.MethodPost:
		return item.Post
	case http.MethodPut:
		return item.Put
	}
	return nil
}

func AddOperationToPathItem(item *openapi3.PathItem, method string, operation *openapi3.Operation) {
	switch method {
	case http.MethodGet:
		item.Get = operation
	case http.MethodDelete:
		item.Delete = operation
	case http.MethodOptions:
		item.Options = operation
	case http.MethodPatch:
		item.Patch = operation
	case http.MethodHead:
		item.Head = operation
	case http.MethodPost:
		item.Post = operation
	case http.MethodPut:
		item.Put = operation
	}
}

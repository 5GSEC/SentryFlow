package apispec

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var defaultIgnoredHeaders = []string{
	contentTypeHeaderName,
	acceptTypeHeaderName,
	authorizationTypeHeaderName,
}

func createHeadersToIgnore(headers []string) map[string]struct{} {
	ret := make(map[string]struct{})

	for _, header := range append(defaultIgnoredHeaders, headers...) {
		ret[strings.ToLower(header)] = struct{}{}
	}

	return ret
}

func shouldIgnoreHeader(headerToIgnore map[string]struct{}, headerKey string) bool {
	_, ok := headerToIgnore[strings.ToLower(headerKey)]
	return ok
}

func (o *OperationGenerator) addResponseHeader(response *openapi3.Response, headerKey, headerValue string) *openapi3.Response {
	if shouldIgnoreHeader(o.ResponseHeadersToIgnore, headerKey) {
		return response
	}

	if response.Headers == nil {
		response.Headers = make(openapi3.Headers)
	}

	response.Headers[headerKey] = &openapi3.HeaderRef{
		Value: &openapi3.Header{
			Parameter: openapi3.Parameter{
				Schema: openapi3.NewSchemaRef("",
					getSchemaFromValue(headerValue, true, openapi3.ParameterInHeader)),
			},
		},
	}

	return response
}

// https://swagger.io/docs/specification/describing-parameters/#header-parameters
func (o *OperationGenerator) addHeaderParam(operation *openapi3.Operation, headerKey, headerValue string) *openapi3.Operation {
	if shouldIgnoreHeader(o.RequestHeadersToIgnore, headerKey) {
		return operation
	}

	headerParam := openapi3.NewHeaderParameter(headerKey).
		WithSchema(getSchemaFromValue(headerValue, true, openapi3.ParameterInHeader))
	operation.AddParameter(headerParam)

	return operation
}

// https://swagger.io/docs/specification/describing-parameters/#cookie-parameters
func (o *OperationGenerator) addCookieParam(operation *openapi3.Operation, headerValue string) *openapi3.Operation {
	// Multiple cookie parameters are sent in the same header, separated by a semicolon and space.
	for _, cookie := range strings.Split(headerValue, "; ") {
		cookieKeyAndValue := strings.Split(cookie, "=")
		if len(cookieKeyAndValue) != 2 { // nolint:gomnd
			logger.Warnf("unsupported cookie param. %v", cookie)
			continue
		}
		key, value := cookieKeyAndValue[0], cookieKeyAndValue[1]
		// Cookie parameters can be primitive values, arrays and objects.
		// Arrays and objects are serialized using the form style.
		headerParam := openapi3.NewCookieParameter(key).WithSchema(getSchemaFromValue(value, true, openapi3.ParameterInCookie))
		operation.AddParameter(headerParam)
	}

	return operation
}

func ConvertHeadersToMap(headers []*Header) map[string]string {
	headersMap := make(map[string]string)

	for _, header := range headers {
		headersMap[strings.ToLower(header.Key)] = header.Value
	}

	return headersMap
}

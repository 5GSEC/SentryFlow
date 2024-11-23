package apispec

import (
	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/utils/field"

	"github.com/5gsec/sentryflow/speculator/pkg/util"
)

var supportedParametersInTypes = []string{openapi3.ParameterInHeader, openapi3.ParameterInQuery, openapi3.ParameterInPath, openapi3.ParameterInCookie}

func mergeOperation(operation, operation2 *openapi3.Operation) (*openapi3.Operation, []conflict) {
	if op, shouldReturn := shouldReturnIfNil(operation, operation2); shouldReturn {
		return op.(*openapi3.Operation), nil
	}

	var requestBodyConflicts, paramConflicts, resConflicts []conflict

	ret := openapi3.NewOperation()

	ret.RequestBody, requestBodyConflicts = mergeRequestBody(operation.RequestBody, operation2.RequestBody,
		field.NewPath("requestBody"))
	ret.Parameters, paramConflicts = mergeParameters(operation.Parameters, operation2.Parameters,
		field.NewPath("parameters"))
	ret.Responses, resConflicts = mergeResponses(*operation.Responses, *operation2.Responses,
		field.NewPath("responses"))

	ret.Security = mergeOperationSecurity(operation.Security, operation2.Security)

	conflicts := append(paramConflicts, resConflicts...)
	conflicts = append(conflicts, requestBodyConflicts...)

	if len(conflicts) > 0 {
		logger.Warnf("found conflicts while merging operation: %v and operation: %v. conflicts: %v", operation, operation2, conflicts)
	}

	return ret, conflicts
}

func mergeOperationSecurity(security, security2 *openapi3.SecurityRequirements) *openapi3.SecurityRequirements {
	if s, shouldReturn := shouldReturnIfNil(security, security2); shouldReturn {
		return s.(*openapi3.SecurityRequirements)
	}

	var mergedSecurity openapi3.SecurityRequirements

	ignoreSecurityKeyMap := map[string]bool{}

	for _, securityMap := range *security {
		mergedSecurity, ignoreSecurityKeyMap = appendSecurityIfNeeded(securityMap, mergedSecurity, ignoreSecurityKeyMap)
	}
	for _, securityMap := range *security2 {
		mergedSecurity, ignoreSecurityKeyMap = appendSecurityIfNeeded(securityMap, mergedSecurity, ignoreSecurityKeyMap)
	}

	return &mergedSecurity
}

func appendSecurityIfNeeded(securityMap openapi3.SecurityRequirement, mergedSecurity openapi3.SecurityRequirements, ignoreSecurityKeyMap map[string]bool) (openapi3.SecurityRequirements, map[string]bool) {
	for key, values := range securityMap {
		// ignore if already appended the exact security key
		if ignoreSecurityKeyMap[key] {
			continue
		}
		// https://swagger.io/docs/specification/authentication/
		// We will treat multiple authentication types as an OR
		// (Security schemes combined via OR are alternatives – any one can be used in the given context)
		mergedSecurity = append(mergedSecurity, map[string][]string{key: values})
		ignoreSecurityKeyMap[key] = true
	}

	return mergedSecurity, ignoreSecurityKeyMap
}

func mergeRequestBody(body, body2 *openapi3.RequestBodyRef, path *field.Path) (*openapi3.RequestBodyRef, []conflict) {
	if p, shouldReturn := shouldReturnIfEmptyRequestBody(body, body2); shouldReturn {
		return p, nil
	}

	content, conflicts := mergeContent(body.Value.Content, body2.Value.Content, path.Child("content"))

	return &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithContent(content),
	}, conflicts
}

func shouldReturnIfEmptyRequestBody(body, body2 *openapi3.RequestBodyRef) (*openapi3.RequestBodyRef, bool) {
	if isEmptyRequestBody(body) {
		return body2, true
	}

	if isEmptyRequestBody(body2) {
		return body, true
	}

	return nil, false
}

func isEmptyRequestBody(body *openapi3.RequestBodyRef) bool {
	return body == nil || body.Value == nil || len(body.Value.Content) == 0
}

func mergeParameters(parameters, parameters2 openapi3.Parameters, path *field.Path) (openapi3.Parameters, []conflict) {
	if p, shouldReturn := shouldReturnIfEmptyParameters(parameters, parameters2); shouldReturn {
		return p, nil
	}

	var retParameters openapi3.Parameters
	var retConflicts []conflict

	parametersByIn := getParametersByIn(parameters)
	parameters2ByIn := getParametersByIn(parameters2)
	for _, inType := range supportedParametersInTypes {
		mergedParameters, conflicts := mergeParametersByInType(parametersByIn[inType], parameters2ByIn[inType], path)
		retParameters = append(retParameters, mergedParameters...)
		retConflicts = append(retConflicts, conflicts...)
	}

	return retParameters, retConflicts
}

func getParametersByIn(parameters openapi3.Parameters) map[string]openapi3.Parameters {
	ret := make(map[string]openapi3.Parameters)

	for i, parameter := range parameters {
		if parameter.Value == nil {
			continue
		}

		switch parameter.Value.In {
		case openapi3.ParameterInCookie, openapi3.ParameterInHeader, openapi3.ParameterInQuery, openapi3.ParameterInPath:
			ret[parameter.Value.In] = append(ret[parameter.Value.In], parameters[i])
		default:
			logger.Warnf("in parameter not supported. %v", parameter.Value.In)
		}
	}

	return ret
}

func mergeParametersByInType(parameters, parameters2 openapi3.Parameters, path *field.Path) (openapi3.Parameters, []conflict) {
	if p, shouldReturn := shouldReturnIfEmptyParameters(parameters, parameters2); shouldReturn {
		return p, nil
	}

	var retParameters openapi3.Parameters
	var retConflicts []conflict

	parametersMapByName := makeParametersMapByName(parameters)
	parameters2MapByName := makeParametersMapByName(parameters2)

	// go over first parameters list
	// 1. merge mutual parameters
	// 2. add non-mutual parameters
	for name, param := range parametersMapByName {
		if param2, ok := parameters2MapByName[name]; ok {
			mergedParameter, conflicts := mergeParameter(param.Value, param2.Value, path.Child(name))
			retConflicts = append(retConflicts, conflicts...)
			retParameters = append(retParameters, &openapi3.ParameterRef{Value: mergedParameter})
		} else {
			retParameters = append(retParameters, param)
		}
	}

	// add non-mutual parameters from the second list
	for name, param := range parameters2MapByName {
		if _, ok := parametersMapByName[name]; !ok {
			retParameters = append(retParameters, param)
		}
	}

	return retParameters, retConflicts
}

func makeParametersMapByName(parameters openapi3.Parameters) map[string]*openapi3.ParameterRef {
	ret := make(map[string]*openapi3.ParameterRef)

	for i := range parameters {
		ret[parameters[i].Value.Name] = parameters[i]
	}

	return ret
}

func mergeParameter(parameter, parameter2 *openapi3.Parameter, path *field.Path) (*openapi3.Parameter, []conflict) {
	if p, shouldReturn := shouldReturnIfEmptyParameter(parameter, parameter2); shouldReturn {
		return p, nil
	}

	type1, type2 := parameter.Schema.Value.Type, parameter2.Schema.Value.Type
	switch conflictSolver(type1, type2) {
	case NoConflict, PreferType1:
		// do nothing, parameter is used.
	case PreferType2:
		// use parameter2.
		type1 = type2
		parameter = parameter2
	case ConflictUnresolved:
		return parameter, []conflict{
			{
				path: path,
				obj1: parameter,
				obj2: parameter2,
				msg:  createConflictMsg(path, type1, type2),
			},
		}
	}

	if type1.Includes(openapi3.TypeBoolean) || type1.Includes(openapi3.TypeInteger) || type1.Includes(openapi3.TypeNumber) || type1.Includes(openapi3.TypeString) {
		schema, conflicts := mergeSchema(parameter.Schema.Value, parameter2.Schema.Value, path)
		return parameter.WithSchema(schema), conflicts
	}
	if type1.Includes(openapi3.TypeArray) {
		items, conflicts := mergeSchemaItems(parameter.Schema.Value.Items, parameter2.Schema.Value.Items, path)
		return parameter.WithSchema(openapi3.NewArraySchema().WithItems(items.Value)), conflicts
	}
	if type1.Includes(openapi3.TypeObject) || type1.Includes("") {
		// when type is missing it is probably an object - we should try and merge the parameter schema
		schema, conflicts := mergeSchema(parameter.Schema.Value, parameter2.Schema.Value, path.Child("schema"))
		return parameter.WithSchema(schema), conflicts
	}
	logger.Warnf("unsupported schema type in parameter: %v", type1)

	return parameter, nil
}

func mergeSchemaItems(items, items2 *openapi3.SchemaRef, path *field.Path) (*openapi3.SchemaRef, []conflict) {
	if s, shouldReturn := shouldReturnIfNil(items, items2); shouldReturn {
		return s.(*openapi3.SchemaRef), nil
	}
	schema, conflicts := mergeSchema(items.Value, items2.Value, path.Child("items"))
	return &openapi3.SchemaRef{Value: schema}, conflicts
}

func mergeSchema(schema, schema2 *openapi3.Schema, path *field.Path) (*openapi3.Schema, []conflict) {
	if s, shouldReturn := shouldReturnIfNil(schema, schema2); shouldReturn {
		return s.(*openapi3.Schema), nil
	}

	if s, shouldReturn := shouldReturnIfEmptySchemaType(schema, schema2); shouldReturn {
		return s, nil
	}

	switch conflictSolver(schema.Type, schema2.Type) {
	case NoConflict, PreferType1:
		// do nothing, schema is used.
	case PreferType2:
		// use schema2.
		schema = schema2
	case ConflictUnresolved:
		return schema, []conflict{
			{
				path: path,
				obj1: schema,
				obj2: schema2,
				msg:  createConflictMsg(path, schema.Type, schema2.Type),
			},
		}
	}

	if schema.Type.Includes(openapi3.TypeBoolean) || schema.Type.Includes(openapi3.TypeInteger) || schema.Type.Includes(openapi3.TypeNumber) {
		return schema, nil
	}
	if schema.Type.Includes(openapi3.TypeString) {
		// Ignore format only if both schemas are string type and formats are different.
		if schema2.Type.Includes(openapi3.TypeString) && schema.Format != schema2.Format {
			schema.Format = ""
		}
		return schema, nil
	}

	return schema, nil
}

func mergeProperties(properties, properties2 openapi3.Schemas, path *field.Path) (openapi3.Schemas, []conflict) {
	retProperties := make(openapi3.Schemas)
	var retConflicts []conflict

	// go over first properties list
	// 1. merge mutual properties
	// 2. add non-mutual properties
	for key := range properties {
		schema := properties[key]
		if schema2, ok := properties2[key]; ok {
			mergedSchema, conflicts := mergeSchema(schema.Value, schema2.Value, path.Child(key))
			retConflicts = append(retConflicts, conflicts...)
			retProperties[key] = &openapi3.SchemaRef{Value: mergedSchema}
		} else {
			retProperties[key] = schema
		}
	}

	// add non-mutual properties from the second list
	for key, schema := range properties2 {
		if _, ok := properties[key]; !ok {
			retProperties[key] = schema
		}
	}

	return retProperties, retConflicts
}

func mergeResponses(responses, responses2 openapi3.Responses, path *field.Path) (*openapi3.Responses, []conflict) {
	if r, shouldReturn := shouldReturnIfEmptyResponses(&responses, &responses2); shouldReturn {
		return r, nil
	}

	var retConflicts []conflict

	retResponses := openapi3.NewResponses()

	// go over first responses list
	// 1. merge mutual response code responses
	// 2. add non-mutual response code responses
	for code, response := range responses.Map() {
		if response2 := responses2.Value(code); response2 != nil {
			mergedResponse, conflicts := mergeResponse(response.Value, response2.Value, path.Child(code))
			retConflicts = append(retConflicts, conflicts...)
			retResponses.Set(code, &openapi3.ResponseRef{Value: mergedResponse})
		} else {
			retResponses.Set(code, responses.Value(code))
		}
	}

	// add non-mutual parameters from the second list
	for code := range responses2.Map() {
		if val := responses.Value(code); val != nil {
			retResponses.Set(code, responses2.Value(code))
		}
	}

	return retResponses, retConflicts
}

func mergeResponse(response, response2 *openapi3.Response, path *field.Path) (*openapi3.Response, []conflict) {
	var retConflicts []conflict
	retResponse := openapi3.NewResponse()
	if response.Description != nil {
		retResponse = retResponse.WithDescription(*response.Description)
	} else if response2.Description != nil {
		retResponse = retResponse.WithDescription(*response2.Description)
	}

	content, conflicts := mergeContent(response.Content, response2.Content, path.Child("content"))
	if len(content) > 0 {
		retResponse = retResponse.WithContent(content)
	}
	retConflicts = append(retConflicts, conflicts...)

	headers, conflicts := mergeResponseHeader(response.Headers, response2.Headers, path.Child("headers"))
	if len(headers) > 0 {
		retResponse.Headers = headers
	}
	retConflicts = append(retConflicts, conflicts...)

	return retResponse, retConflicts
}

func mergeContent(content openapi3.Content, content2 openapi3.Content, path *field.Path) (openapi3.Content, []conflict) {
	var retConflicts []conflict
	retContent := openapi3.NewContent()

	// go over first content list
	// 1. merge mutual content media type
	// 2. add non-mutual content media type
	for name, mediaType := range content {
		if mediaType2, ok := content2[name]; ok {
			mergedSchema, conflicts := mergeSchema(mediaType.Schema.Value, mediaType2.Schema.Value, path.Child(name))
			// TODO: handle mediaType.Encoding
			retConflicts = append(retConflicts, conflicts...)
			retContent[name] = openapi3.NewMediaType().WithSchema(mergedSchema)
		} else {
			retContent[name] = content[name]
		}
	}

	// add non-mutual content media type from the second list
	for name := range content2 {
		if _, ok := content[name]; !ok {
			retContent[name] = content2[name]
		}
	}

	return retContent, retConflicts
}

func mergeResponseHeader(headers, headers2 openapi3.Headers, path *field.Path) (openapi3.Headers, []conflict) {
	var retConflicts []conflict
	retHeaders := make(openapi3.Headers)

	// go over first headers list
	// 1. merge mutual headers
	// 2. add non-mutual headers
	for name, header := range headers {
		if header2, ok := headers2[name]; ok {
			mergedHeader, conflicts := mergeHeader(header.Value, header2.Value, path.Child(name))
			retConflicts = append(retConflicts, conflicts...)
			retHeaders[name] = &openapi3.HeaderRef{Value: mergedHeader}
		} else {
			retHeaders[name] = headers[name]
		}
	}

	// add non-mutual headers from the second list
	for name := range headers2 {
		if _, ok := headers[name]; !ok {
			retHeaders[name] = headers2[name]
		}
	}

	return retHeaders, retConflicts
}

func mergeHeader(header, header2 *openapi3.Header, path *field.Path) (*openapi3.Header, []conflict) {
	if h, shouldReturn := shouldReturnIfEmptyHeader(header, header2); shouldReturn {
		return h, nil
	}

	if header.In != header2.In {
		return header, []conflict{
			{
				path: path,
				obj1: header,
				obj2: header2,
				msg:  createHeaderInConflictMsg(path, header.In, header2.In),
			},
		}
	}

	schema, conflicts := mergeSchema(header.Schema.Value, header2.Schema.Value, path)
	header.Parameter = *header.WithSchema(schema)

	return header, conflicts
}

func shouldReturnIfEmptyParameter(param, param2 *openapi3.Parameter) (*openapi3.Parameter, bool) {
	if isEmptyParameter(param) {
		return param2, true
	}

	if isEmptyParameter(param2) {
		return param, true
	}

	return nil, false
}

func isEmptyParameter(param *openapi3.Parameter) bool {
	return param == nil || isEmptySchemaRef(param.Schema)
}

func shouldReturnIfEmptyHeader(header, header2 *openapi3.Header) (*openapi3.Header, bool) {
	if isEmptyHeader(header) {
		return header2, true
	}

	if isEmptyHeader(header2) {
		return header, true
	}

	return nil, false
}

func isEmptyHeader(header *openapi3.Header) bool {
	return header == nil || isEmptySchemaRef(header.Schema)
}

func isEmptySchemaRef(schemaRef *openapi3.SchemaRef) bool {
	return schemaRef == nil || schemaRef.Value == nil
}

func shouldReturnIfEmptyResponses(r, r2 *openapi3.Responses) (*openapi3.Responses, bool) {
	if r.Len() == 0 {
		return r2, true
	}
	if r2.Len() == 0 {
		return r, true
	}
	// both are not empty
	return nil, false
}

func shouldReturnIfEmptyParameters(parameters, parameters2 openapi3.Parameters) (openapi3.Parameters, bool) {
	if len(parameters) == 0 {
		return parameters2, true
	}
	if len(parameters2) == 0 {
		return parameters, true
	}
	// both are not empty
	return nil, false
}

func shouldReturnIfEmptySchemaType(s, s2 *openapi3.Schema) (*openapi3.Schema, bool) {
	if len(s.Type.Slice()) == 0 {
		return s2, true
	}
	if len(s2.Type.Slice()) == 0 {
		return s, true
	}
	// both are not empty
	return nil, false
}

// used only with pointers.
func shouldReturnIfNil(a, b interface{}) (interface{}, bool) {
	if util.IsNil(a) {
		return b, true
	}
	if util.IsNil(b) {
		return a, true
	}
	// both are not nil
	return nil, false
}

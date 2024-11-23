package apispec

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/url"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
	"github.com/xeipuuv/gojsonschema"

	"github.com/5gsec/sentryflow/speculator/pkg/util"
)

func getSchema(value any) (schema *openapi3.Schema, err error) {
	switch value.(type) {
	case bool:
		schema = openapi3.NewBoolSchema()
	case string:
		schema = getStringSchema(value)
	case json.Number:
		schema = getNumberSchema(value)
	case map[string]any:
		schema, err = getObjectSchema(value)
		if err != nil {
			return nil, err
		}
	case []any:
		schema, err = getArraySchema(value)
		if err != nil {
			return nil, err
		}
	case nil:
		// TODO: Not sure how to handle null. ex: {"size":3,"err":null}
		schema = openapi3.NewStringSchema()
	default:
		// TODO:
		// I've tested additionalProperties and it seems like properties - we will might have problems in the diff logic
		// openapi3.MapProperty()
		// openapi3.RefProperty()
		// openapi3.RefSchema()
		// openapi3.ComposedSchema() - discriminator?
		return nil, fmt.Errorf("unexpected value type. value=%v, type=%T", value, value)
	}

	return schema, nil
}

func getStringSchema(value any) (schema *openapi3.Schema) {
	return openapi3.NewStringSchema().WithFormat(getStringFormat(value))
}

func getNumberSchema(value any) (schema *openapi3.Schema) {
	// https://swagger.io/docs/specification/data-models/data-types/#numbers

	// It is important to try first convert it to int
	if _, err := value.(json.Number).Int64(); err != nil {
		// if failed to convert to int it's a double
		// TODO: we will set a 'double' and not a 'float' - is that ok?
		schema = openapi3.NewFloat64Schema()
	} else {
		schema = openapi3.NewInt64Schema()
	}
	// TODO: Format
	// openapi3.Int8Property()
	// openapi3.Int16Property()
	// openapi3.Int32Property()
	// openapi3.Float64Property()
	// openapi3.Float32Property()
	return schema /*.WithExample(value)*/
}

func getObjectSchema(value any) (schema *openapi3.Schema, err error) {
	schema = openapi3.NewObjectSchema()
	stringMapE, err := cast.ToStringMapE(value)
	if err != nil {
		return nil, fmt.Errorf("failed to cast to string map. value=%v: %w", value, err)
	}

	for key, val := range stringMapE {
		if s, err := getSchema(val); err != nil {
			return nil, fmt.Errorf("failed to get schema from string map. key=%v, value=%v: %w", key, val, err)
		} else {
			schema = schema.WithProperty(escapeString(key), s)
		}
	}

	return schema, nil
}

func escapeString(key string) string {
	// need to escape double quotes if exists
	if strings.Contains(key, "\"") {
		key = strings.ReplaceAll(key, "\"", "\\\"")
	}
	return key
}

func getArraySchema(value any) (schema *openapi3.Schema, err error) {
	sliceE, err := cast.ToSliceE(value)
	if err != nil {
		return nil, fmt.Errorf("failed to cast to slice. value=%v: %w", value, err)
	}

	// in order to support mixed type array we will map all schemas by schema type
	schemaTypeToSchema := make(map[string]*openapi3.Schema)
	for i := range sliceE {
		item, err := getSchema(sliceE[i])
		if err != nil {
			return nil, fmt.Errorf("failed to get items schema from slice. value=%v: %w", sliceE[i], err)
		}
		if len(item.Type.Slice()) > 0 {
			if _, ok := schemaTypeToSchema[item.Type.Slice()[0]]; !ok {
				schemaTypeToSchema[item.Type.Slice()[0]] = item
			}
		}
	}

	switch len(schemaTypeToSchema) {
	case 0:
		// array is empty, but we can't create an empty array property (Schemas with 'type: array', require a sibling 'items:' field)
		// we will create string type items as a default value
		schema = openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())
	case 1:
		for _, s := range schemaTypeToSchema {
			schema = openapi3.NewArraySchema().WithItems(s)
			break
		}
	default:
		// oneOf
		// https://swagger.io/docs/specification/data-models/oneof-anyof-allof-not/
		var schemas []*openapi3.Schema
		for _, s := range schemaTypeToSchema {
			schemas = append(schemas, s)
		}
		schema = openapi3.NewOneOfSchema(schemas...)
	}

	return schema, nil
}

type HTTPInteractionData struct {
	ReqBody, RespBody       string
	ReqHeaders, RespHeaders map[string]string
	QueryParams             url.Values
	statusCode              int
}

func (h *HTTPInteractionData) getReqContentType() string {
	return h.ReqHeaders[contentTypeHeaderName]
}

func (h *HTTPInteractionData) getRespContentType() string {
	return h.RespHeaders[contentTypeHeaderName]
}

type OperationGeneratorConfig struct {
	ResponseHeadersToIgnore []string
	RequestHeadersToIgnore  []string
}

type OperationGenerator struct {
	ResponseHeadersToIgnore map[string]struct{}
	RequestHeadersToIgnore  map[string]struct{}
}

func NewOperationGenerator(config OperationGeneratorConfig) *OperationGenerator {
	return &OperationGenerator{
		ResponseHeadersToIgnore: createHeadersToIgnore(config.ResponseHeadersToIgnore),
		RequestHeadersToIgnore:  createHeadersToIgnore(config.RequestHeadersToIgnore),
	}
}

// Note: SecuritySchemes might be updated.
func (o *OperationGenerator) GenerateSpecOperation(data *HTTPInteractionData, securitySchemes openapi3.SecuritySchemes) (*openapi3.Operation, error) {
	operation := openapi3.NewOperation()

	if len(data.ReqBody) > 0 {
		reqContentType := data.getReqContentType()
		if reqContentType == "" {
			logger.Infof("Missing Content-Type header, ignoring request body. (%v)", data.ReqBody)
		} else {
			mediaType, mediaTypeParams, err := mime.ParseMediaType(reqContentType)
			if err != nil {
				return nil, fmt.Errorf("failed to parse request media type. Content-Type=%v: %w", reqContentType, err)
			}
			switch true {
			case util.IsApplicationJSONMediaType(mediaType):
				reqBodyJSON, err := gojsonschema.NewStringLoader(data.ReqBody).LoadJSON()
				if err != nil {
					return nil, fmt.Errorf("failed to load json from request body. body=%v: %w", data.ReqBody, err)
				}

				reqSchema, err := getSchema(reqBodyJSON)
				if err != nil {
					return nil, fmt.Errorf("failed to get schema from request body. body=%v: %w", data.ReqBody, err)
				}

				operationSetRequestBody(operation, openapi3.NewRequestBody().WithJSONSchema(reqSchema))
			case mediaType == mediaTypeApplicationForm:
				operation, securitySchemes, err = handleApplicationFormURLEncodedBody(operation, securitySchemes, data.ReqBody)
				if err != nil {
					return nil, fmt.Errorf("failed to handle %s body: %v", mediaTypeApplicationForm, err)
				}
			case mediaType == mediaTypeMultipartFormData:
				// Multipart requests combine one or more sets of data into a single body, separated by boundaries.
				// You typically use these requests for file uploads and for transferring data of several types
				// in a single request (for example, a file along with a JSON object).
				// https://swagger.io/docs/specification/describing-request-body/multipart-requests/
				schema, err := getMultipartFormDataSchema(data.ReqBody, mediaTypeParams)
				if err != nil {
					return nil, fmt.Errorf("failed to get multipart form-data schema from request body. body=%v: %v", data.ReqBody, err)
				}
				operationSetRequestBody(operation, openapi3.NewRequestBody().WithFormDataSchema(schema))
			default:
				logger.Infof("Treating %v as default request content type (no schema)", reqContentType)
			}
		}
	}

	for key, value := range data.ReqHeaders {
		lowerKey := strings.ToLower(key)
		if lowerKey == authorizationTypeHeaderName {
			// https://datatracker.ietf.org/doc/html/rfc6750#section-2.1
			operation, securitySchemes = handleAuthReqHeader(operation, securitySchemes, value)
		} else if APIKeyNames[lowerKey] {
			schemeKey := APIKeyAuthSecuritySchemeKey
			operation = addSecurity(operation, schemeKey)
			securitySchemes = updateSecuritySchemes(securitySchemes, schemeKey, NewAPIKeySecuritySchemeInHeader(key))
		} else if lowerKey == cookieTypeHeaderName {
			operation = o.addCookieParam(operation, value)
		} else {
			operation = o.addHeaderParam(operation, key, value)
		}
	}

	for key, values := range data.QueryParams {
		lowerKey := strings.ToLower(key)
		if lowerKey == AccessTokenParamKey {
			// https://datatracker.ietf.org/doc/html/rfc6750#section-2.3
			operation, securitySchemes = handleAuthQueryParam(operation, securitySchemes, values)
		} else if APIKeyNames[lowerKey] {
			schemeKey := APIKeyAuthSecuritySchemeKey
			operation = addSecurity(operation, schemeKey)
			securitySchemes = updateSecuritySchemes(securitySchemes, schemeKey, NewAPIKeySecuritySchemeInQuery(key))
		} else {
			operation = addQueryParam(operation, key, values)
		}
	}

	// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#responseObject
	// REQUIRED. A short description of the response.
	response := openapi3.NewResponse().WithDescription("response")
	if len(data.RespBody) > 0 {
		respContentType := data.getRespContentType()
		if respContentType == "" {
			logger.Infof("Missing Content-Type header, ignoring response body. (%v)", data.RespBody)
		} else {
			mediaType, _, err := mime.ParseMediaType(respContentType)
			if err != nil {
				return nil, fmt.Errorf("failed to parse response media type. Content-Type=%v: %w", respContentType, err)
			}
			switch true {
			case util.IsApplicationJSONMediaType(mediaType):
				respBodyJSON, err := gojsonschema.NewStringLoader(data.RespBody).LoadJSON()
				if err != nil {
					return nil, fmt.Errorf("failed to load json from response body. body=%v: %w", data.RespBody, err)
				}

				respSchema, err := getSchema(respBodyJSON)
				if err != nil {
					return nil, fmt.Errorf("failed to get schema from response body. body=%v: %w", respBodyJSON, err)
				}

				response = response.WithJSONSchema(respSchema)
			default:
				logger.Infof("Treating %v as default response content type (no schema)", respContentType)
			}
		}
	}

	for key, value := range data.RespHeaders {
		response = o.addResponseHeader(response, key, value)
	}

	operation.AddResponse(data.statusCode, response)
	operation.AddResponse(0 /*"default"*/, openapi3.NewResponse().WithDescription("default"))

	return operation, nil
}

func operationSetRequestBody(operation *openapi3.Operation, reqBody *openapi3.RequestBody) {
	operation.RequestBody = &openapi3.RequestBodyRef{Value: reqBody}
}

func CloneOperation(op *openapi3.Operation) (*openapi3.Operation, error) {
	var out openapi3.Operation

	opB, err := json.Marshal(op)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operation (%+v): %v", op, err)
	}

	if err := json.Unmarshal(opB, &out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}

	return &out, nil
}

func getBearerAuthClaims(bearerToken string) (claims jwt.MapClaims, found bool) {
	if len(bearerToken) == 0 {
		logger.Warnf("authZ token provided with no value.")
		return nil, false
	}

	// Parse the claims without validating (since we don't want to bother downloading a key)
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(bearerToken, jwt.MapClaims{})
	if err != nil {
		logger.Warnf("authZ token is not a JWT.")
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Infof("authZ token had unintelligble claims.")
		return nil, false
	}

	return claims, true
}

func generateBearerAuthScheme(operation *openapi3.Operation, claims jwt.MapClaims, key string) (*openapi3.Operation, *openapi3.SecurityScheme) {
	switch key {
	case BearerAuthSecuritySchemeKey:
		// https://swagger.io/docs/specification/authentication/bearer-authentication/
		return addSecurity(operation, key), openapi3.NewJWTSecurityScheme()
	case OAuth2SecuritySchemeKey:
		// https://swagger.io/docs/specification/authentication/oauth2/
		// we can't know the flow type (implicit, password, clientCredentials or authorizationCode) so we choose authorizationCode for now
		scopes := getScopesFromJWTClaims(claims)
		oAuth2SecurityScheme := NewOAuth2SecurityScheme(scopes)
		return addSecurity(operation, key, scopes...), oAuth2SecurityScheme
	default:
		logger.Warnf("Unsupported BearerAuth key: %v", key)
		return operation, nil
	}
}

func getScopesFromJWTClaims(claims jwt.MapClaims) []string {
	var scopes []string
	if claims == nil {
		return scopes
	}

	if scope, ok := claims["scope"]; ok {
		scopes = strings.Split(scope.(string), " ")
		logger.Debugf("found OAuth token scopes: %v", scopes)
	} else {
		logger.Warnf("no scopes defined in this token")
	}
	return scopes
}

func handleAuthQueryParam(operation *openapi3.Operation, securitySchemes openapi3.SecuritySchemes, values []string) (*openapi3.Operation, openapi3.SecuritySchemes) {
	if len(values) > 1 {
		// RFC 6750 does not prohibit multiple tokens, but we do not know whether
		// they would be AND or OR so we just pick the latest.
		logger.Warnf("Found %v tokens in query parameters, using only the last", len(values))
		values = values[len(values)-1:]
	}

	// Use scheme as security scheme name
	securitySchemeKey := OAuth2SecuritySchemeKey
	claims, _ := getBearerAuthClaims(values[0])

	if hasSecurity(operation, securitySchemeKey) {
		// RFC 6750 states multiple methods (form, uri query, header) cannot be used.
		logger.Errorf("OAuth tokens supplied with multiple methods, ignoring query param")
		return operation, securitySchemes
	}

	var scheme *openapi3.SecurityScheme
	operation, scheme = generateBearerAuthScheme(operation, claims, securitySchemeKey)
	if scheme != nil {
		securitySchemes = updateSecuritySchemes(securitySchemes, securitySchemeKey, scheme)
	}
	return operation, securitySchemes
}

func handleAuthReqHeader(operation *openapi3.Operation, securitySchemes openapi3.SecuritySchemes, value string) (*openapi3.Operation, openapi3.SecuritySchemes) {
	if strings.HasPrefix(value, BasicAuthPrefix) {
		// https://swagger.io/docs/specification/authentication/basic-authentication/
		// Use scheme as security scheme name
		key := BasicAuthSecuritySchemeKey
		operation = addSecurity(operation, key)
		securitySchemes = updateSecuritySchemes(securitySchemes, key, NewBasicAuthSecurityScheme())
	} else if strings.HasPrefix(value, BearerAuthPrefix) {
		// https://swagger.io/docs/specification/authentication/bearer-authentication/
		// https://datatracker.ietf.org/doc/html/rfc6750#section-2.1
		// Use scheme as security scheme name. For OAuth, we should consider checking
		// supported scopes to allow multiple defs.
		key := BearerAuthSecuritySchemeKey
		claims, found := getBearerAuthClaims(strings.TrimPrefix(value, BearerAuthPrefix))
		if found {
			key = OAuth2SecuritySchemeKey
		}

		if hasSecurity(operation, key) {
			// RFC 6750 states multiple methods (form, uri query, header) cannot be used.
			logger.Error("OAuth tokens supplied with multiple methods, ignoring header")
			return operation, securitySchemes
		}

		var scheme *openapi3.SecurityScheme
		operation, scheme = generateBearerAuthScheme(operation, claims, key)
		if scheme != nil {
			securitySchemes = updateSecuritySchemes(securitySchemes, key, scheme)
		}
	} else {
		logger.Warnf("ignoring unknown authorization header value (%v)", value)
	}
	return operation, securitySchemes
}

func addSecurity(op *openapi3.Operation, name string, scopes ...string) *openapi3.Operation {
	// https://swagger.io/docs/specification/authentication/
	// We will treat multiple authentication types as an OR
	// (Security schemes combined via OR are alternatives – any one can be used in the given context)
	securityRequirement := openapi3.NewSecurityRequirement()

	if len(scopes) > 0 {
		securityRequirement[name] = scopes
	} else {
		// We must use an empty array as the scopes, otherwise it will create invalid swagger
		securityRequirement[name] = []string{}
	}

	if op.Security == nil {
		op.Security = openapi3.NewSecurityRequirements()
	}
	op.Security.With(securityRequirement)

	return op
}

func hasSecurity(op *openapi3.Operation, name string) bool {
	if op.Security == nil {
		return false
	}

	for _, securityScheme := range *op.Security {
		if _, ok := securityScheme[name]; ok {
			return true
		}
	}
	return false
}

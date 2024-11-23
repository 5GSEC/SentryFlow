package apispec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/yudai/gojsondiff"
)

const (
	schemasRefPrefix    = "#/components/schemas/"
	maxSchemaToRefDepth = 20
)

// will return a map of SchemaRef and update the operation accordingly.
func updateSchemas(schemas openapi3.Schemas, op *openapi3.Operation) (retSchemas openapi3.Schemas, retOperation *openapi3.Operation) {
	if op == nil {
		return schemas, op
	}

	for i, response := range op.Responses.Map() {
		if response.Value == nil {
			continue
		}
		for content, mediaType := range response.Value.Content {
			schemas, mediaType.Schema = schemaToRef(schemas, mediaType.Schema.Value, "", 0)
			op.Responses.Value(i).Value.Content[content] = mediaType
		}
	}

	for i, parameter := range op.Parameters {
		if parameter.Value == nil {
			continue
		}
		for content, mediaType := range parameter.Value.Content {
			schemas, mediaType.Schema = schemaToRef(schemas, mediaType.Schema.Value, "", 0)
			op.Parameters[i].Value.Content[content] = mediaType
		}
	}

	if op.RequestBody != nil && op.RequestBody.Value != nil {
		for content, mediaType := range op.RequestBody.Value.Content {
			schemas, mediaType.Schema = schemaToRef(schemas, mediaType.Schema.Value, "", 0)
			op.RequestBody.Value.Content[content] = mediaType
		}
	}

	return schemas, op
}

func schemaToRef(schemas openapi3.Schemas, schema *openapi3.Schema, schemeNameHint string, depth int) (retSchemes openapi3.Schemas, schemaRef *openapi3.SchemaRef) {
	if schema == nil {
		return schemas, nil
	}

	if depth >= maxSchemaToRefDepth {
		logger.Warnf("Maximum depth was reached")
		return schemas, openapi3.NewSchemaRef("", schema)
	}

	if schema.Type.Is(openapi3.TypeArray) {
		if schema.Items == nil {
			// no need to create definition for an empty array
			return schemas, openapi3.NewSchemaRef("", schema)
		}
		// remove plural from def name hint when it's an array type (if exist)
		schemas, schema.Items = schemaToRef(schemas, schema.Items.Value, strings.TrimSuffix(schemeNameHint, "s"), depth+1)
		return schemas, openapi3.NewSchemaRef("", schema)
	}

	if !schema.Type.Is(openapi3.TypeObject) {
		return schemas, openapi3.NewSchemaRef("", schema)
	}

	if schema.Properties == nil || len(schema.Properties) == 0 {
		// no need to create ref for an empty object
		return schemas, openapi3.NewSchemaRef("", schema)
	}

	// go over all properties in the object and convert each one to ref if needed
	var propNames []string
	for propName := range schema.Properties {
		var ref *openapi3.SchemaRef
		schemas, ref = schemaToRef(schemas, schema.Properties[propName].Value, propName, depth+1)
		if ref != nil {
			schema.Properties[propName] = ref
			propNames = append(propNames, propName)
		}
	}

	// look for schema in schemas with identical schema
	schemeName, exist := findScheme(schemas, schema)
	if !exist {
		// generate new definition
		schemeName = schemeNameHint
		if schemeName == "" {
			schemeName = generateDefNameFromPropNames(propNames)
		}
		if schemas == nil {
			schemas = make(openapi3.Schemas)
		}
		if existingSchema, ok := schemas[schemeName]; ok {
			logger.Debugf("Security scheme name exist with different schema. existingSchema=%+v, schema=%+v", existingSchema, schema)
			schemeName = getUniqueSchemeName(schemas, schemeName)
		}
		schemas[schemeName] = openapi3.NewSchemaRef("", schema)
	}

	return schemas, openapi3.NewSchemaRef(schemasRefPrefix+schemeName, nil)
}

func generateDefNameFromPropNames(propNames []string) string {
	// generate name based on properties names when 'defNameHint' is missing
	// sort the slice to get more stable test results
	sort.Strings(propNames)
	propString := strings.Join(propNames, "_")
	return regexp.MustCompile(`[^a-zA-Z0-9._-]+`).ReplaceAllString(propString, "")
}

func getUniqueSchemeName(schemes openapi3.Schemas, name string) string {
	counter := 0
	for {
		suggestedName := fmt.Sprintf("%s_%d", name, counter)
		if _, ok := schemes[suggestedName]; !ok {
			// found a unique name
			return suggestedName
		}
		// suggestedName already exist - increase counter and look again
		counter++
	}
}

// will look for identical scheme in schemes map.
func findScheme(schemas openapi3.Schemas, schema *openapi3.Schema) (schemeName string, exist bool) {
	schemaBytes, _ := json.Marshal(schema)
	differ := gojsondiff.New()
	for name, defSchema := range schemas {
		defSchemaBytes, _ := json.Marshal(defSchema)
		diff, err := differ.Compare(defSchemaBytes, schemaBytes)
		if err != nil {
			logger.Errorf("Failed to compare schemas: %v", err)
			continue
		}
		if !diff.Modified() {
			logger.Debugf("Schema was found in schemas. schema=%+v, def name=%v", schema, name)
			return name, true
		}
	}

	logger.Debugf("Schema was not found in schemas. schema=%+v", schema)
	return "", false
}

package apispec

import (
	"encoding/json"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

var (
	stringNumberObject = createObjectSchema(map[string]*openapi3.Schema{
		openapi3.TypeString: openapi3.NewStringSchema(),
		openapi3.TypeNumber: openapi3.NewFloat64Schema(),
	})
	stringBooleanObject = createObjectSchema(map[string]*openapi3.Schema{
		openapi3.TypeString:  openapi3.NewStringSchema(),
		openapi3.TypeBoolean: openapi3.NewBoolSchema(),
	})
	stringIntegerObject = createObjectSchema(map[string]*openapi3.Schema{
		openapi3.TypeString:  openapi3.NewStringSchema(),
		openapi3.TypeInteger: openapi3.NewInt64Schema(),
	})
)

func marshal(obj interface{}) string {
	objB, _ := json.Marshal(obj)
	return string(objB)
}

func createObjectSchema(properties map[string]*openapi3.Schema) *openapi3.Schema {
	return openapi3.NewObjectSchema().WithProperties(properties)
}

func createObjectSchemaWithRef(properties map[string]*openapi3.SchemaRef) *openapi3.Schema {
	objectSchema := openapi3.NewObjectSchema()
	for name, ref := range properties {
		objectSchema.WithPropertyRef(name, ref)
	}

	return objectSchema
}

func Test_findDefinition(t *testing.T) {
	type args struct {
		schemas openapi3.Schemas
		schema  *openapi3.Schema
	}
	tests := []struct {
		name        string
		args        args
		wantDefName string
		wantExist   bool
	}{
		{
			name: "identical string schema exist",
			args: args{
				schemas: openapi3.Schemas{
					"string": &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
				schema: openapi3.NewStringSchema(),
			},
			wantDefName: "string",
			wantExist:   true,
		},
		{
			name: "identical string schema does not exist",
			args: args{
				schemas: openapi3.Schemas{
					"string": &openapi3.SchemaRef{Value: openapi3.NewStringSchema().WithFormat("format")},
				},
				schema: openapi3.NewStringSchema(),
			},
			wantDefName: "",
			wantExist:   false,
		},
		{
			name: "identical object schema exist (object order is different)",
			args: args{
				schemas: openapi3.Schemas{
					"object": &openapi3.SchemaRef{Value: openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
						openapi3.TypeObject: stringIntegerObject,
						openapi3.TypeString: openapi3.NewStringSchema(),
					})},
				},
				schema: createObjectSchema(
					map[string]*openapi3.Schema{
						openapi3.TypeString: openapi3.NewStringSchema(),
						openapi3.TypeObject: createObjectSchema(
							map[string]*openapi3.Schema{
								openapi3.TypeInteger: openapi3.NewInt64Schema(),
								openapi3.TypeString:  openapi3.NewStringSchema(),
							},
						),
					},
				),
			},
			wantDefName: "object",
			wantExist:   true,
		},
		{
			name: "identical object schema does not exist",
			args: args{
				schemas: openapi3.Schemas{
					"object": &openapi3.SchemaRef{Value: openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
						openapi3.TypeString: openapi3.NewStringSchema(),
						openapi3.TypeObject: stringIntegerObject,
					})},
				},
				schema: createObjectSchema(
					map[string]*openapi3.Schema{
						openapi3.TypeString: openapi3.NewStringSchema(),
						openapi3.TypeObject: stringNumberObject,
					},
				),
			},
			wantDefName: "",
			wantExist:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDefName, gotExist := findScheme(tt.args.schemas, tt.args.schema)
			if gotDefName != tt.wantDefName {
				t.Errorf("findScheme() gotDefName = %v, want %v", gotDefName, tt.wantDefName)
			}
			if gotExist != tt.wantExist {
				t.Errorf("findScheme() gotExist = %v, want %v", gotExist, tt.wantExist)
			}
		})
	}
}

func Test_getUniqueDefName(t *testing.T) {
	type args struct {
		schemas openapi3.Schemas
		name    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "name does not exist",
			args: args{
				schemas: openapi3.Schemas{
					"string": &openapi3.SchemaRef{Value: stringIntegerObject},
				},
				name: "no-test",
			},
			want: "no-test_0",
		},
		{
			name: "name exist once",
			args: args{
				schemas: openapi3.Schemas{
					"test_0": &openapi3.SchemaRef{Value: stringIntegerObject},
				},
				name: "test",
			},
			want: "test_1",
		},
		{
			name: "name exist multiple times",
			args: args{
				schemas: openapi3.Schemas{
					"test":   &openapi3.SchemaRef{Value: stringIntegerObject},
					"test_0": &openapi3.SchemaRef{Value: stringNumberObject},
					"test_1": &openapi3.SchemaRef{Value: stringBooleanObject},
				},
				name: "test",
			},
			want: "test_2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUniqueSchemeName(tt.args.schemas, tt.args.name); got != tt.want {
				t.Errorf("getUniqueSchemeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createArraySchemaWithRefItems(name string) *openapi3.Schema {
	arraySchemaWithRefItems := openapi3.NewArraySchema()
	arraySchemaWithRefItems.Items = openapi3.NewSchemaRef(schemasRefPrefix+name, nil)
	return arraySchemaWithRefItems
}

func Test_schemaToRef(t *testing.T) {
	arraySchemaWithNilItems := openapi3.NewArraySchema()
	arraySchemaWithNilItems.Items = nil
	type args struct {
		schemas     openapi3.Schemas
		schema      *openapi3.Schema
		defNameHint string
		depth       int
	}
	tests := []struct {
		name           string
		args           args
		wantRetSchemas openapi3.Schemas
		wantRetSchema  *openapi3.SchemaRef
	}{
		{
			name: "nil schema",
			args: args{
				schemas: openapi3.Schemas{
					"test": &openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
				},
				schema:      nil,
				defNameHint: "",
			},
			wantRetSchemas: openapi3.Schemas{
				"test": &openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
			},
			wantRetSchema: nil,
		},
		{
			name: "array schema with nil items",
			args: args{
				schemas: openapi3.Schemas{
					"test": &openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
				},
				schema:      openapi3.NewArraySchema(),
				defNameHint: "",
			},
			wantRetSchemas: openapi3.Schemas{
				"test": &openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
			},
			wantRetSchema: openapi3.NewSchemaRef("", openapi3.NewArraySchema()),
		},
		{
			name: "array schema with non object items - no change for schemas",
			args: args{
				schemas: openapi3.Schemas{
					"test": &openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
				},
				schema:      openapi3.NewArraySchema().WithItems(openapi3.NewBoolSchema()),
				defNameHint: "",
			},
			wantRetSchemas: openapi3.Schemas{
				"test": &openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
			},
			wantRetSchema: openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewBoolSchema())),
		},
		{
			name: "array schema with object items - use hint name",
			args: args{
				schemas: openapi3.Schemas{
					"test": &openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
				},
				schema:      openapi3.NewArraySchema().WithItems(stringNumberObject),
				defNameHint: "hint",
			},
			wantRetSchemas: openapi3.Schemas{
				"test": &openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
				"hint": &openapi3.SchemaRef{Value: stringNumberObject},
			},
			wantRetSchema: openapi3.NewSchemaRef("", createArraySchemaWithRefItems("hint")),
		},
		{
			name: "array schema with object items - hint name already exist",
			args: args{
				schemas: openapi3.Schemas{
					"hint": &openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
				},
				schema:      openapi3.NewArraySchema().WithItems(stringNumberObject),
				defNameHint: "hint",
			},
			wantRetSchemas: openapi3.Schemas{
				"hint":   &openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
				"hint_0": &openapi3.SchemaRef{Value: stringNumberObject},
			},
			wantRetSchema: openapi3.NewSchemaRef("", createArraySchemaWithRefItems("hint_0")),
		},
		{
			name: "primitive type",
			args: args{
				schemas: openapi3.Schemas{
					"test": &openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
				},
				schema: openapi3.NewInt64Schema(),
			},
			wantRetSchemas: openapi3.Schemas{
				"test": &openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
			},
			wantRetSchema: openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
		},
		{
			name: "empty object - no new schemas",
			args: args{
				schemas: openapi3.Schemas{
					"test": &openapi3.SchemaRef{Value: stringNumberObject},
				},
				schema: openapi3.NewObjectSchema(),
			},
			wantRetSchemas: openapi3.Schemas{
				"test": &openapi3.SchemaRef{Value: stringNumberObject},
			},
			wantRetSchema: openapi3.NewSchemaRef("", openapi3.NewObjectSchema()),
		},
		{
			name: "object - definition exist",
			args: args{
				schemas: openapi3.Schemas{
					"test": &openapi3.SchemaRef{Value: stringNumberObject},
				},
				schema: stringNumberObject,
			},
			wantRetSchemas: openapi3.Schemas{
				"test": &openapi3.SchemaRef{Value: stringNumberObject},
			},
			wantRetSchema: openapi3.NewSchemaRef(schemasRefPrefix+"test", nil),
		},
		{
			name: "object - definition does not exist",
			args: args{
				schemas: openapi3.Schemas{
					"test": &openapi3.SchemaRef{Value: stringBooleanObject},
				},
				schema: stringNumberObject,
			},
			wantRetSchemas: openapi3.Schemas{
				"test":          &openapi3.SchemaRef{Value: stringBooleanObject},
				"number_string": &openapi3.SchemaRef{Value: stringNumberObject},
			},
			wantRetSchema: openapi3.NewSchemaRef(schemasRefPrefix+"number_string", nil),
		},
		{
			name: "object - definition does not exist - use hint",
			args: args{
				schemas: openapi3.Schemas{
					"test": &openapi3.SchemaRef{Value: stringBooleanObject},
				},
				schema:      stringNumberObject,
				defNameHint: "hint",
			},
			wantRetSchemas: openapi3.Schemas{
				"test": &openapi3.SchemaRef{Value: stringBooleanObject},
				"hint": &openapi3.SchemaRef{Value: stringNumberObject},
			},
			wantRetSchema: openapi3.NewSchemaRef(schemasRefPrefix+"hint", nil),
		},
		{
			name: "object in object",
			args: args{
				schemas: nil,
				schema: createObjectSchema(
					map[string]*openapi3.Schema{
						openapi3.TypeString: openapi3.NewStringSchema(),
						openapi3.TypeObject: stringNumberObject,
					},
				),
			},
			wantRetSchemas: openapi3.Schemas{
				"object": openapi3.NewSchemaRef("", stringNumberObject),
				"object_string": openapi3.NewSchemaRef("", createObjectSchemaWithRef(
					map[string]*openapi3.SchemaRef{
						openapi3.TypeObject: openapi3.NewSchemaRef(schemasRefPrefix+"object", nil),
						openapi3.TypeString: openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
					},
				)),
			},
			wantRetSchema: openapi3.NewSchemaRef(schemasRefPrefix+"object_string", nil),
		},
		{
			name: "array of object in an object",
			args: args{
				schemas: nil,
				schema: createObjectSchema(
					map[string]*openapi3.Schema{
						openapi3.TypeBoolean: openapi3.NewBoolSchema(),

						/*use plural to check the removal of the "s"*/
						"objects": openapi3.NewArraySchema().WithItems(stringNumberObject),
					},
				),
			},
			wantRetSchemas: openapi3.Schemas{
				"object": openapi3.NewSchemaRef("", stringNumberObject),
				"boolean_objects": openapi3.NewSchemaRef("", createObjectSchemaWithRef(
					map[string]*openapi3.SchemaRef{
						openapi3.TypeBoolean: openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
						"objects":            openapi3.NewSchemaRef("", createArraySchemaWithRefItems("object")),
					},
				)),
			},
			wantRetSchema: openapi3.NewSchemaRef(schemasRefPrefix+"boolean_objects", nil),
		},
		{
			name: "object in object in object - max depth was reached after 1 object - ref was not created",
			args: args{
				schemas: nil,
				schema: createObjectSchema(
					map[string]*openapi3.Schema{
						"obj1": createObjectSchema(
							map[string]*openapi3.Schema{
								"obj2": stringNumberObject,
							},
						),
					},
				),
				depth: maxSchemaToRefDepth - 1,
			},
			wantRetSchemas: openapi3.Schemas{
				"obj1": openapi3.NewSchemaRef("", createObjectSchema(
					map[string]*openapi3.Schema{
						"obj1": createObjectSchema(
							map[string]*openapi3.Schema{
								"obj2": stringNumberObject,
							},
						),
					},
				),
				),
			},
			wantRetSchema: openapi3.NewSchemaRef(schemasRefPrefix+"obj1", nil),
		},
		{
			name: "object in object in object - max depth was reached after 2 objects - ref was not created",
			args: args{
				schemas: nil,
				schema: createObjectSchema(
					map[string]*openapi3.Schema{
						"obj1": createObjectSchema(
							map[string]*openapi3.Schema{
								"obj2":   stringNumberObject,
								"string": openapi3.NewStringSchema(),
							},
						),
					},
				),
				depth: maxSchemaToRefDepth - 2,
			},
			wantRetSchemas: openapi3.Schemas{
				"obj1": openapi3.NewSchemaRef("", createObjectSchema(
					map[string]*openapi3.Schema{
						"obj2":   stringNumberObject,
						"string": openapi3.NewStringSchema(),
					},
				),
				),
				"obj1_0": openapi3.NewSchemaRef("", createObjectSchemaWithRef(
					map[string]*openapi3.SchemaRef{
						"obj1": openapi3.NewSchemaRef(schemasRefPrefix+"obj1", nil),
					},
				),
				),
			},
			wantRetSchema: openapi3.NewSchemaRef(schemasRefPrefix+"obj1_0", nil),
		},
		{
			name: "max depth was reached - ref was not created",
			args: args{
				schemas: nil,
				schema: createObjectSchema(
					map[string]*openapi3.Schema{
						openapi3.TypeBoolean: openapi3.NewBoolSchema(),

						/*use plural to check the removal of the "s"*/
						"objects": openapi3.NewArraySchema().WithItems(stringNumberObject),
					},
				),
				depth: maxSchemaToRefDepth,
			},
			wantRetSchemas: nil,
			wantRetSchema: openapi3.NewSchemaRef("", createObjectSchema(
				map[string]*openapi3.Schema{
					openapi3.TypeBoolean: openapi3.NewBoolSchema(),

					/*use plural to check the removal of the "s"*/
					"objects": openapi3.NewArraySchema().WithItems(stringNumberObject),
				},
			)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRetSchemas, gotRetSchema := schemaToRef(tt.args.schemas, tt.args.schema, tt.args.defNameHint, tt.args.depth)
			assertEqual(t, gotRetSchemas, tt.wantRetSchemas)
			assertEqual(t, gotRetSchema, tt.wantRetSchema)
		})
	}
}

var interactionReqBody = `{"active":true,
"certificateVersion":"86eb5278-676a-3b7c-b29d-4a57007dc7be",
"controllerInstanceInfo":{"replicaId":"portshift-agent-66fc77c848-tmmk8"},
"policyAndAppVersion":1621477900361,
"version":"1.147.1"}`

var interactionRespBody = `{"cvss":[{"score":7.8,"vector":"AV:L/AC:L/PR:N/UI:R/S:U/C:H/I:H/A:H"}]}`

var interaction = &HTTPInteractionData{
	ReqBody:  interactionReqBody,
	RespBody: interactionRespBody,
	ReqHeaders: map[string]string{
		contentTypeHeaderName: mediaTypeApplicationJSON,
	},
	RespHeaders: map[string]string{
		contentTypeHeaderName: mediaTypeApplicationJSON,
	},
	statusCode: 200,
}

func createArraySchemaWithRef(ref string) *openapi3.Schema {
	arraySchema := openapi3.NewArraySchema()
	arraySchema.Items = &openapi3.SchemaRef{Ref: ref}
	return arraySchema
}

func Test_updateSchemas(t *testing.T) {
	op := NewOperation(t, interaction).Op
	retOp := NewOperation(t, interaction).Op
	retOp.RequestBody = &openapi3.RequestBodyRef{Value: openapi3.NewRequestBody().WithJSONSchemaRef(&openapi3.SchemaRef{
		Ref: schemasRefPrefix + "active_certificateVersion_controllerInstanceInfo_policyAndAppVersion_version",
	})}
	retOp.Responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().WithDescription("response").WithJSONSchemaRef(&openapi3.SchemaRef{
			Ref: schemasRefPrefix + "cvss",
		}),
	})

	type args struct {
		schemas openapi3.Schemas
		op      *openapi3.Operation
	}
	tests := []struct {
		name             string
		args             args
		wantRetSchemas   openapi3.Schemas
		wantRetOperation *openapi3.Operation
	}{
		{
			name: "sanity",
			args: args{
				schemas: nil,
				op:      op,
			},
			wantRetSchemas: openapi3.Schemas{
				"controllerInstanceInfo": openapi3.NewSchemaRef("", createObjectSchema(
					map[string]*openapi3.Schema{
						"replicaId": openapi3.NewStringSchema(),
					},
				)),
				"active_certificateVersion_controllerInstanceInfo_policyAndAppVersion_version": openapi3.NewSchemaRef("", createObjectSchemaWithRef(
					map[string]*openapi3.SchemaRef{
						"active":                 openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
						"certificateVersion":     openapi3.NewSchemaRef("", openapi3.NewUUIDSchema()),
						"controllerInstanceInfo": openapi3.NewSchemaRef(schemasRefPrefix+"controllerInstanceInfo", nil),
						"policyAndAppVersion":    openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
						"version":                openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
					},
				)),
				"cvs": openapi3.NewSchemaRef("", createObjectSchema(
					map[string]*openapi3.Schema{
						"score":  openapi3.NewFloat64Schema(),
						"vector": openapi3.NewStringSchema(),
					},
				)),
				"cvss": openapi3.NewSchemaRef("", createObjectSchema(
					map[string]*openapi3.Schema{
						"cvss": createArraySchemaWithRef(schemasRefPrefix + "cvs"),
					},
				)),
			},
			wantRetOperation: retOp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRetSchemas, gotRetOperation := updateSchemas(tt.args.schemas, tt.args.op)
			assertEqual(t, gotRetSchemas, tt.wantRetSchemas)
			assertEqual(t, gotRetOperation, tt.wantRetOperation)
		})
	}
}

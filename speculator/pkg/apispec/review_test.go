package apispec

//import (
//	"net/http"
//	"reflect"
//	"sort"
//	"sync"
//	"testing"
//
//	"github.com/getkin/kin-openapi/openapi3"
//	"github.com/gofrs/uuid"
//
//	"github.com/5gsec/sentryflow/speculator/pkg/pathtrie"
//)
//
//var dataCombined = &HTTPInteractionData{
//	ReqBody:  combinedReq,
//	RespBody: combinedRes,
//	ReqHeaders: map[string]string{
//		contentTypeHeaderName: mediaTypeApplicationJSON,
//	},
//	RespHeaders: map[string]string{
//		contentTypeHeaderName: mediaTypeApplicationJSON,
//	},
//	statusCode: 200,
//}
//
//func TestSpec_ApplyApprovedReview(t *testing.T) {
//	host := "host"
//	port := "8080"
//	uuidVar, _ := uuid.NewV4()
//
//	type fields struct {
//		Host         string
//		Port         string
//		ID           uuid.UUID
//		ApprovedSpec *ApprovedSpec
//		LearningSpec *LearningSpec
//		Mutex        sync.Mutex
//	}
//	type args struct {
//		approvedReviews *ApprovedSpecReview
//		specVersion     OASVersion
//	}
//	tests := []struct {
//		name     string
//		fields   fields
//		args     args
//		wantSpec *Spec
//		wantErr  bool
//	}{
//		{
//			name: "1 reviewed path item. modified path param. same path item. 2 Paths",
//			fields: fields{
//				Host: host,
//				Port: port,
//				ID:   uuidVar,
//				ApprovedSpec: &ApprovedSpec{
//					PathItems: map[string]*openapi3.PathItem{},
//				},
//				LearningSpec: &LearningSpec{
//					PathItems: map[string]*openapi3.PathItem{
//						"/api/1": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/api/2": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data2).Op).PathItem,
//					},
//				},
//			},
//			args: args{
//				specVersion: OASv3,
//				approvedReviews: &ApprovedSpecReview{
//					PathToPathItem: map[string]*openapi3.PathItem{
//						"/api/1": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/api/2": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data2).Op).PathItem,
//					},
//					PathItemsReview: []*ApprovedSpecReviewPathItem{
//						{
//							ReviewPathItem: ReviewPathItem{
//								ParameterizedPath: "/api/{param1}",
//								Paths: map[string]bool{
//									"/api/1": true,
//									"/api/2": true,
//								},
//							},
//							PathUUID: "1",
//						},
//					},
//				},
//			},
//			wantSpec: &Spec{
//				SpecInfo: SpecInfo{
//					ID:   uuidVar,
//					Host: host,
//					Port: port,
//					ApprovedSpec: &ApprovedSpec{
//						PathItems: map[string]*openapi3.PathItem{
//							"/api/{param1}": &NewTestPathItem().
//								WithOperation(http.MethodGet, NewOperation(t, dataCombined).Op).
//								WithPathParams("param1", openapi3.NewInt64Schema()).PathItem,
//						},
//						SpecVersion: OASv3,
//					},
//					LearningSpec: &LearningSpec{
//						PathItems: map[string]*openapi3.PathItem{},
//					},
//					ApprovedPathTrie: createPathTrie(map[string]string{
//						"/api/{param1}": "1",
//					}),
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "user took out one path out of the parameterized path, and also one more path has learned between review and approve (should ignore it and not delete)",
//			fields: fields{
//				Host: host,
//				Port: port,
//				ID:   uuidVar,
//				ApprovedSpec: &ApprovedSpec{
//					PathItems: map[string]*openapi3.PathItem{},
//				},
//				LearningSpec: &LearningSpec{
//					PathItems: map[string]*openapi3.PathItem{
//						"api/3/foo": &NewTestPathItem().PathItem,
//						"/api/1": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/api/2": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data2).Op).PathItem,
//					},
//				},
//			},
//			args: args{
//				specVersion: OASv2,
//				approvedReviews: &ApprovedSpecReview{
//					PathToPathItem: map[string]*openapi3.PathItem{
//						"/api/1": &NewTestPathItem().
//							WithOperation(http.MethodPost, NewOperation(t, Data).Op).PathItem,
//						"/api/2": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data2).Op).PathItem,
//					},
//					PathItemsReview: []*ApprovedSpecReviewPathItem{
//						{
//							ReviewPathItem: ReviewPathItem{
//								ParameterizedPath: "/api/{param1}",
//								Paths: map[string]bool{
//									"/api/2": true,
//								},
//							},
//							PathUUID: "1",
//						},
//						{
//							ReviewPathItem: ReviewPathItem{
//								ParameterizedPath: "/api/1",
//								Paths: map[string]bool{
//									"/api/1": true,
//								},
//							},
//							PathUUID: "2",
//						},
//					},
//				},
//			},
//			wantSpec: &Spec{
//				SpecInfo: SpecInfo{
//					Host: host,
//					Port: port,
//					ID:   uuidVar,
//					ApprovedSpec: &ApprovedSpec{
//						PathItems: map[string]*openapi3.PathItem{
//							"/api/{param1}": &NewTestPathItem().
//								WithOperation(http.MethodGet, NewOperation(t, Data2).Op).
//								WithPathParams("param1", openapi3.NewInt64Schema()).PathItem,
//							"/api/1": &NewTestPathItem().
//								WithOperation(http.MethodPost, NewOperation(t, Data).Op).PathItem,
//						},
//						SpecVersion: OASv2,
//					},
//					LearningSpec: &LearningSpec{
//						PathItems: map[string]*openapi3.PathItem{
//							"api/3/foo": &NewTestPathItem().PathItem,
//						},
//					},
//					ApprovedPathTrie: createPathTrie(map[string]string{
//						"/api/{param1}": "1",
//						"/api/1":        "2",
//					}),
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "multiple methods",
//			fields: fields{
//				Host: host,
//				Port: port,
//				ID:   uuidVar,
//				ApprovedSpec: &ApprovedSpec{
//					PathItems: map[string]*openapi3.PathItem{},
//				},
//				LearningSpec: &LearningSpec{
//					PathItems: map[string]*openapi3.PathItem{
//						"/anything": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).
//							WithOperation(http.MethodPost, NewOperation(t, Data).Op).PathItem,
//						"/headers": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/user-agent": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//					},
//				},
//			},
//			args: args{
//				specVersion: OASv3,
//				approvedReviews: &ApprovedSpecReview{
//					PathToPathItem: map[string]*openapi3.PathItem{
//						"/anything": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).
//							WithOperation(http.MethodPost, NewOperation(t, Data).Op).PathItem,
//						"/headers": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/user-agent": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//					},
//					PathItemsReview: []*ApprovedSpecReviewPathItem{
//						{
//							ReviewPathItem: ReviewPathItem{
//								ParameterizedPath: "/api/{test}",
//								Paths: map[string]bool{
//									"/anything":   true,
//									"/headers":    true,
//									"/user-agent": true,
//								},
//							},
//							PathUUID: "1",
//						},
//					},
//				},
//			},
//			wantSpec: &Spec{
//				SpecInfo: SpecInfo{
//					Host: host,
//					Port: port,
//					ID:   uuidVar,
//					ApprovedSpec: &ApprovedSpec{
//						PathItems: map[string]*openapi3.PathItem{
//							"/api/{test}": &NewTestPathItem().
//								WithOperation(http.MethodPost, NewOperation(t, Data).Op).
//								WithOperation(http.MethodGet, NewOperation(t, Data).Op).
//								WithPathParams("test", openapi3.NewStringSchema()).PathItem,
//						},
//						SpecVersion: OASv3,
//					},
//					LearningSpec: &LearningSpec{
//						PathItems: map[string]*openapi3.PathItem{},
//					},
//					ApprovedPathTrie: createPathTrie(map[string]string{
//						"/api/{test}": "1",
//					}),
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "new parameterized path, unmerge of path item",
//			fields: fields{
//				Host: host,
//				Port: port,
//				ID:   uuidVar,
//				ApprovedSpec: &ApprovedSpec{
//					PathItems: map[string]*openapi3.PathItem{},
//				},
//				LearningSpec: &LearningSpec{
//					PathItems: map[string]*openapi3.PathItem{
//						"/api/1": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/api/2": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/api/foo": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/user/1/bar/2": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//					},
//				},
//			},
//			args: args{
//				specVersion: OASv3,
//				approvedReviews: &ApprovedSpecReview{
//					PathToPathItem: map[string]*openapi3.PathItem{
//						"/api/1": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/api/2": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/api/foo": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/user/1/bar/2": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//					},
//					PathItemsReview: []*ApprovedSpecReviewPathItem{
//						{
//							ReviewPathItem: ReviewPathItem{
//								ParameterizedPath: "/api/{param1}",
//								Paths: map[string]bool{
//									"/api/1": true,
//									"/api/2": true,
//								},
//							},
//							PathUUID: "1",
//						},
//						{
//							ReviewPathItem: ReviewPathItem{
//								ParameterizedPath: "/api/foo",
//								Paths: map[string]bool{
//									"/api/foo": true,
//								},
//							},
//							PathUUID: "2",
//						},
//						{
//							ReviewPathItem: ReviewPathItem{
//								ParameterizedPath: "/user/{param1}/bar/{param2}",
//								Paths: map[string]bool{
//									"/user/1/bar/2": true,
//								},
//							},
//							PathUUID: "3",
//						},
//					},
//				},
//			},
//			wantSpec: &Spec{
//				SpecInfo: SpecInfo{
//					Host: host,
//					Port: port,
//					ID:   uuidVar,
//					ApprovedSpec: &ApprovedSpec{
//						PathItems: map[string]*openapi3.PathItem{
//							"/api/{param1}": &NewTestPathItem().
//								WithOperation(http.MethodGet, NewOperation(t, Data).Op).
//								WithPathParams("param1", openapi3.NewInt64Schema()).PathItem,
//							"/api/foo": &NewTestPathItem().
//								WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//							"/user/{param1}/bar/{param2}": &NewTestPathItem().
//								WithOperation(http.MethodGet, NewOperation(t, Data).Op).
//								WithPathParams("param1", openapi3.NewInt64Schema()).
//								WithPathParams("param2", openapi3.NewInt64Schema()).PathItem,
//						},
//						SpecVersion: OASv3,
//					},
//					LearningSpec: &LearningSpec{
//						PathItems: map[string]*openapi3.PathItem{},
//					},
//					ApprovedPathTrie: createPathTrie(map[string]string{
//						"/api/{param1}":               "1",
//						"/api/foo":                    "2",
//						"/user/{param1}/bar/{param2}": "3",
//					}),
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "new parameterized path, unmerge of path item with security",
//			fields: fields{
//				Host: host,
//				Port: port,
//				ID:   uuidVar,
//				ApprovedSpec: &ApprovedSpec{
//					PathItems:       map[string]*openapi3.PathItem{},
//					SecuritySchemes: openapi3.SecuritySchemes{},
//				},
//				LearningSpec: &LearningSpec{
//					PathItems: map[string]*openapi3.PathItem{
//						"/api/1": &NewTestPathItem().WithOperation(http.MethodGet,
//							NewOperation(t, Data).WithSecurityRequirement(openapi3.SecurityRequirement{BasicAuthSecuritySchemeKey: {}}).Op).PathItem,
//						"/api/2": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/api/foo": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/user/1/bar/2": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//					},
//				},
//			},
//			args: args{
//				specVersion: OASv3,
//				approvedReviews: &ApprovedSpecReview{
//					PathToPathItem: map[string]*openapi3.PathItem{
//						"/api/1": &NewTestPathItem().WithOperation(http.MethodGet,
//							NewOperation(t, Data).WithSecurityRequirement(openapi3.SecurityRequirement{BasicAuthSecuritySchemeKey: {}}).Op).PathItem,
//						"/api/2": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/api/foo": &NewTestPathItem().
//							WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/user/1/bar/2": &NewTestPathItem().WithOperation(http.MethodGet,
//							NewOperation(t, Data).WithSecurityRequirement(openapi3.SecurityRequirement{OAuth2SecuritySchemeKey: {}}).Op).PathItem,
//					},
//					PathItemsReview: []*ApprovedSpecReviewPathItem{
//						{
//							ReviewPathItem: ReviewPathItem{
//								ParameterizedPath: "/api/{param1}",
//								Paths: map[string]bool{
//									"/api/1": true,
//									"/api/2": true,
//								},
//							},
//							PathUUID: "1",
//						},
//						{
//							ReviewPathItem: ReviewPathItem{
//								ParameterizedPath: "/api/foo",
//								Paths: map[string]bool{
//									"/api/foo": true,
//								},
//							},
//							PathUUID: "2",
//						},
//						{
//							ReviewPathItem: ReviewPathItem{
//								ParameterizedPath: "/user/{param1}/bar/{param2}",
//								Paths: map[string]bool{
//									"/user/1/bar/2": true,
//								},
//							},
//							PathUUID: "3",
//						},
//					},
//				},
//			},
//			wantSpec: &Spec{
//				SpecInfo: SpecInfo{
//					Host: host,
//					Port: port,
//					ID:   uuidVar,
//					ApprovedSpec: &ApprovedSpec{
//						PathItems: map[string]*openapi3.PathItem{
//							"/api/{param1}": &NewTestPathItem().WithOperation(http.MethodGet,
//								NewOperation(t, Data).
//									WithSecurityRequirement(openapi3.SecurityRequirement{BasicAuthSecuritySchemeKey: {}}).Op).
//								WithPathParams("param1", openapi3.NewInt64Schema()).PathItem,
//							"/api/foo": &NewTestPathItem().WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//							"/user/{param1}/bar/{param2}": &NewTestPathItem().WithOperation(http.MethodGet,
//								NewOperation(t, Data).
//									WithSecurityRequirement(openapi3.SecurityRequirement{OAuth2SecuritySchemeKey: {}}).Op).
//								WithPathParams("param1", openapi3.NewInt64Schema()).
//								WithPathParams("param2", openapi3.NewInt64Schema()).PathItem,
//						},
//						SecuritySchemes: openapi3.SecuritySchemes{
//							BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
//							OAuth2SecuritySchemeKey:    &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
//						},
//						SpecVersion: OASv3,
//					},
//					LearningSpec: &LearningSpec{
//						PathItems: map[string]*openapi3.PathItem{},
//					},
//					ApprovedPathTrie: createPathTrie(map[string]string{
//						"/api/{param1}":               "1",
//						"/api/foo":                    "2",
//						"/user/{param1}/bar/{param2}": "3",
//					}),
//				},
//			},
//			wantErr: false,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Spec{
//				SpecInfo: SpecInfo{
//					Host:             tt.fields.Host,
//					Port:             tt.fields.Port,
//					ID:               tt.fields.ID,
//					ApprovedSpec:     tt.fields.ApprovedSpec,
//					LearningSpec:     tt.fields.LearningSpec,
//					ApprovedPathTrie: pathtrie.New(),
//				},
//			}
//			err := s.ApplyApprovedReview(tt.args.approvedReviews, tt.args.specVersion)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Error response not as expected. want error: %v. error: %v", tt.wantErr, err)
//				return
//			}
//
//			//assert.DeepEqual(t, s, tt.wantSpec, cmpopts.IgnoreUnexported(openapi3.Schema{}, Spec{}), cmpopts.IgnoreTypes(openapi3.ExtensionProps{}))
//			assertEqual(t, s, tt.wantSpec)
//		})
//	}
//}
//
//func TestSpec_CreateSuggestedReview(t *testing.T) {
//	type fields struct {
//		ID                        uuid.UUID
//		ApprovedSpec              *ApprovedSpec
//		LearningSpec              *LearningSpec
//		LearningParametrizedPaths *LearningParametrizedPaths
//		Mutex                     sync.Mutex
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   *SuggestedSpecReview
//	}{
//		{
//			name: "2 paths - map to one parameterized path",
//			fields: fields{
//				ID: uuid.UUID{},
//				LearningSpec: &LearningSpec{
//					PathItems: map[string]*openapi3.PathItem{
//						"/api/1": &NewTestPathItem().WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/api/2": &NewTestPathItem().WithOperation(http.MethodGet, NewOperation(t, Data2).Op).PathItem,
//					},
//				},
//				LearningParametrizedPaths: &LearningParametrizedPaths{
//					Paths: map[string]map[string]bool{
//						"/api/{param1}": {"/api/1": true, "/api/2": true},
//					},
//				},
//			},
//			want: &SuggestedSpecReview{
//				PathToPathItem: map[string]*openapi3.PathItem{
//					"/api/1": &NewTestPathItem().WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//					"/api/2": &NewTestPathItem().WithOperation(http.MethodGet, NewOperation(t, Data2).Op).PathItem,
//				},
//				PathItemsReview: []*SuggestedSpecReviewPathItem{
//					{
//						ReviewPathItem: ReviewPathItem{
//							ParameterizedPath: "/api/{param1}",
//							Paths: map[string]bool{
//								"/api/1": true,
//								"/api/2": true,
//							},
//						},
//					},
//				},
//			},
//		},
//		{
//			name: "4 paths - 2 under one parameterized path with one param, one is not parameterized, one is parameterized path with 2 params",
//			fields: fields{
//				ID: uuid.UUID{},
//				LearningSpec: &LearningSpec{
//					PathItems: map[string]*openapi3.PathItem{
//						"/api/1":           &NewTestPathItem().WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/api/2":           &NewTestPathItem().WithOperation(http.MethodGet, NewOperation(t, Data2).Op).PathItem,
//						"/api/foo":         &NewTestPathItem().WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//						"/api/foo/1/bar/2": &NewTestPathItem().WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//					},
//				},
//				LearningParametrizedPaths: &LearningParametrizedPaths{
//					Paths: map[string]map[string]bool{
//						"/api/{param1}":                  {"/api/1": true, "/api/2": true},
//						"/api/foo/{param1}/bar/{param2}": {"/api/foo/1/bar/2": true},
//						"/api/foo":                       {"/api/foo": true},
//					},
//				},
//			},
//			want: &SuggestedSpecReview{
//				PathToPathItem: map[string]*openapi3.PathItem{
//					"/api/1":           &NewTestPathItem().WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//					"/api/2":           &NewTestPathItem().WithOperation(http.MethodGet, NewOperation(t, Data2).Op).PathItem,
//					"/api/foo":         &NewTestPathItem().WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//					"/api/foo/1/bar/2": &NewTestPathItem().WithOperation(http.MethodGet, NewOperation(t, Data).Op).PathItem,
//				},
//				PathItemsReview: []*SuggestedSpecReviewPathItem{
//					{
//						ReviewPathItem: ReviewPathItem{
//							ParameterizedPath: "/api/{param1}",
//							Paths: map[string]bool{
//								"/api/1": true,
//								"/api/2": true,
//							},
//						},
//					},
//					{
//						ReviewPathItem: ReviewPathItem{
//							ParameterizedPath: "/api/foo/{param1}/bar/{param2}",
//							Paths: map[string]bool{
//								"/api/foo/1/bar/2": true,
//							},
//						},
//					},
//					{
//						ReviewPathItem: ReviewPathItem{
//							ParameterizedPath: "/api/foo",
//							Paths: map[string]bool{
//								"/api/foo": true,
//							},
//						},
//					},
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Spec{
//				SpecInfo: SpecInfo{
//					ID:           tt.fields.ID,
//					ApprovedSpec: tt.fields.ApprovedSpec,
//					LearningSpec: tt.fields.LearningSpec,
//				},
//			}
//			got := s.CreateSuggestedReview()
//			sort.Slice(got.PathItemsReview, func(i, j int) bool {
//				return got.PathItemsReview[i].ParameterizedPath > got.PathItemsReview[j].ParameterizedPath
//			})
//			sort.Slice(tt.want.PathItemsReview, func(i, j int) bool {
//				return tt.want.PathItemsReview[i].ParameterizedPath > tt.want.PathItemsReview[j].ParameterizedPath
//			})
//			gotB := marshal(got)
//			wantB := marshal(tt.want)
//			if gotB != wantB {
//				t.Errorf("CreateSuggestedReview() got = %v, want %v", gotB, wantB)
//			}
//		})
//	}
//}
//
//func TestSpec_createLearningParametrizedPaths(t *testing.T) {
//	type fields struct {
//		Host         string
//		ID           uuid.UUID
//		ApprovedSpec *ApprovedSpec
//		LearningSpec *LearningSpec
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   *LearningParametrizedPaths
//	}{
//		{
//			name: "",
//			fields: fields{
//				LearningSpec: &LearningSpec{
//					PathItems: map[string]*openapi3.PathItem{
//						"/api/1": &NewTestPathItem().PathItem,
//					},
//				},
//			},
//			want: &LearningParametrizedPaths{
//				Paths: map[string]map[string]bool{
//					"/api/{param1}": {"/api/1": true},
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Spec{
//				SpecInfo: SpecInfo{
//					Host:         tt.fields.Host,
//					ID:           tt.fields.ID,
//					ApprovedSpec: tt.fields.ApprovedSpec,
//					LearningSpec: tt.fields.LearningSpec,
//				},
//			}
//			if got := s.createLearningParametrizedPaths(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("createLearningParametrizedPaths() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_addPathParamsToPathItem(t *testing.T) {
//	type args struct {
//		pathItem      *openapi3.PathItem
//		suggestedPath string
//		paths         map[string]bool
//	}
//	tests := []struct {
//		name         string
//		args         args
//		wantPathItem *openapi3.PathItem
//	}{
//		{
//			name: "1 param",
//			args: args{
//				pathItem:      &NewTestPathItem().PathItem,
//				suggestedPath: "/api/{param1}/foo",
//				paths: map[string]bool{
//					"api/1/foo": true,
//					"api/2/foo": true,
//				},
//			},
//			wantPathItem: &NewTestPathItem().WithPathParams("param1", openapi3.NewInt64Schema()).PathItem,
//		},
//		{
//			name: "2 params",
//			args: args{
//				pathItem:      &NewTestPathItem().PathItem,
//				suggestedPath: "/api/{param1}/foo/{param2}",
//				paths: map[string]bool{
//					"api/1/foo/2":   true,
//					"api/2/foo/345": true,
//				},
//			},
//			wantPathItem: &NewTestPathItem().
//				WithPathParams("param1", openapi3.NewInt64Schema()).
//				WithPathParams("param2", openapi3.NewInt64Schema()).PathItem,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			addPathParamsToPathItem(tt.args.pathItem, tt.args.suggestedPath, tt.args.paths)
//			assertEqual(t, tt.args.pathItem, tt.wantPathItem)
//		})
//	}
//}
//
//func Test_updateSecurityDefinitionsFromPathItem(t *testing.T) {
//	type args struct {
//		securitySchemes openapi3.SecuritySchemes
//		item            *openapi3.PathItem
//	}
//	tests := []struct {
//		name string
//		args args
//		want openapi3.SecuritySchemes
//	}{
//		{
//			name: "Get operation",
//			args: args{
//				securitySchemes: openapi3.SecuritySchemes{},
//				item: &openapi3.PathItem{
//					Get: createOperationWithSecurity(&openapi3.SecurityRequirements{
//						{
//							BasicAuthSecuritySchemeKey: {},
//						},
//					}),
//				},
//			},
//			want: openapi3.SecuritySchemes{
//				BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
//			},
//		},
//		{
//			name: "Put operation",
//			args: args{
//				securitySchemes: openapi3.SecuritySchemes{},
//				item: &openapi3.PathItem{
//					Put: createOperationWithSecurity(&openapi3.SecurityRequirements{
//						{
//							OAuth2SecuritySchemeKey: {"admin"},
//						},
//						{
//							BasicAuthSecuritySchemeKey: {},
//						},
//					}),
//				},
//			},
//			want: openapi3.SecuritySchemes{
//				OAuth2SecuritySchemeKey:    &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
//				BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
//			},
//		},
//		{
//			name: "Post operation",
//			args: args{
//				securitySchemes: openapi3.SecuritySchemes{},
//				item: &openapi3.PathItem{
//					Post: createOperationWithSecurity(&openapi3.SecurityRequirements{
//						{
//							OAuth2SecuritySchemeKey: {"admin"},
//						},
//						{
//							BasicAuthSecuritySchemeKey: {},
//						},
//					}),
//				},
//			},
//			want: openapi3.SecuritySchemes{
//				OAuth2SecuritySchemeKey:    &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
//				BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
//			},
//		},
//		{
//			name: "Delete operation",
//			args: args{
//				securitySchemes: openapi3.SecuritySchemes{},
//				item: &openapi3.PathItem{
//					Delete: createOperationWithSecurity(&openapi3.SecurityRequirements{
//						{
//							OAuth2SecuritySchemeKey: {"admin"},
//						},
//						{
//							BasicAuthSecuritySchemeKey: {},
//						},
//					}),
//				},
//			},
//			want: openapi3.SecuritySchemes{
//				OAuth2SecuritySchemeKey:    &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
//				BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
//			},
//		},
//		{
//			name: "Options operation",
//			args: args{
//				securitySchemes: openapi3.SecuritySchemes{},
//				item: &openapi3.PathItem{
//					Options: createOperationWithSecurity(&openapi3.SecurityRequirements{
//						{
//							OAuth2SecuritySchemeKey: {"admin"},
//						},
//						{
//							BasicAuthSecuritySchemeKey: {},
//						},
//					}),
//				},
//			},
//			want: openapi3.SecuritySchemes{
//				OAuth2SecuritySchemeKey:    &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
//				BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
//			},
//		},
//		{
//			name: "Head operation",
//			args: args{
//				securitySchemes: openapi3.SecuritySchemes{},
//				item: &openapi3.PathItem{
//					Head: createOperationWithSecurity(&openapi3.SecurityRequirements{
//						{
//							OAuth2SecuritySchemeKey: {"admin"},
//						},
//						{
//							BasicAuthSecuritySchemeKey: {},
//						},
//					}),
//				},
//			},
//			want: openapi3.SecuritySchemes{
//				OAuth2SecuritySchemeKey:    &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
//				BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
//			},
//		},
//		{
//			name: "Patch operation",
//			args: args{
//				securitySchemes: openapi3.SecuritySchemes{},
//				item: &openapi3.PathItem{
//					Patch: createOperationWithSecurity(&openapi3.SecurityRequirements{
//						{
//							OAuth2SecuritySchemeKey: {"admin"},
//						},
//						{
//							BasicAuthSecuritySchemeKey: {},
//						},
//					}),
//				},
//			},
//			want: openapi3.SecuritySchemes{
//				OAuth2SecuritySchemeKey:    &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
//				BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
//			},
//		},
//		{
//			name: "Multiple operations",
//			args: args{
//				securitySchemes: openapi3.SecuritySchemes{},
//				item: &openapi3.PathItem{
//					Get: createOperationWithSecurity(&openapi3.SecurityRequirements{
//						{
//							BasicAuthSecuritySchemeKey: {},
//						},
//					}),
//					Put: createOperationWithSecurity(&openapi3.SecurityRequirements{
//						{
//							OAuth2SecuritySchemeKey: {"read"},
//						},
//					}),
//					Post: createOperationWithSecurity(&openapi3.SecurityRequirements{
//						{
//							"unsupported": {"read"},
//						},
//					}),
//					Delete: createOperationWithSecurity(&openapi3.SecurityRequirements{
//						{
//							OAuth2SecuritySchemeKey: {"admin"},
//						},
//						{
//							BasicAuthSecuritySchemeKey: {},
//						},
//					}),
//					Options: createOperationWithSecurity(nil),
//				},
//			},
//			want: openapi3.SecuritySchemes{
//				OAuth2SecuritySchemeKey:    &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
//				BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := updateSecuritySchemesFromPathItem(tt.args.securitySchemes, tt.args.item); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("updateSecuritySchemesFromPathItem() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

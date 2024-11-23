package apispec

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/5gsec/sentryflow/speculator/pkg/pathtrie"
)

func CreateDefaultSpec(host string, port string, config OperationGeneratorConfig) *Spec {
	return &Spec{
		SpecInfo: SpecInfo{
			Host: host,
			Port: port,
			LearningSpec: &LearningSpec{
				PathItems:       map[string]*openapi3.PathItem{},
				SecuritySchemes: openapi3.SecuritySchemes{},
			},
			ApprovedSpec: &ApprovedSpec{
				PathItems:       map[string]*openapi3.PathItem{},
				SecuritySchemes: openapi3.SecuritySchemes{},
			},
			ApprovedPathTrie: pathtrie.New(),
			ProvidedPathTrie: pathtrie.New(),
		},
		OpGenerator: NewOperationGenerator(config),
	}
}

func createDefaultSwaggerInfo() *openapi3.Info {
	return &openapi3.Info{
		Description:    "This is a generated Open API Spec",
		Title:          "Swagger",
		TermsOfService: "https://swagger.io/terms/",
		Contact: &openapi3.Contact{
			Email: "apiteam@swagger.io",
		},
		License: &openapi3.License{
			Name: "Apache 2.0",
			URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
		},
		Version: "1.0.0",
	}
}

package apispec

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/golang-jwt/jwt/v5"
)

type OAuth2Claims struct {
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}

const (
	BasicAuthSecuritySchemeKey  = "BasicAuth"
	APIKeyAuthSecuritySchemeKey = "ApiKeyAuth"
	OAuth2SecuritySchemeKey     = "OAuth2"
	BearerAuthSecuritySchemeKey = "BearerAuth"

	BearerAuthPrefix = "Bearer "
	BasicAuthPrefix  = "Basic "

	AccessTokenParamKey = "access_token"

	tknURL           = "https://example.com/oauth2/token"
	authorizationURL = "https://example.com/oauth2/authorize"

	apiKeyType      = "apiKey"
	basicAuthType   = "http"
	basicAuthScheme = "basic"
	oauth2Type      = "oauth2"
)

// APIKeyNames is set of names of headers or query params defining API keys.
// This should be runtime configurable, of course.
// Note: keys should be lowercase.
var APIKeyNames = map[string]bool{
	"key":     true, // Google
	"api_key": true,
}

func newAPIKeySecurityScheme(name string) *openapi3.SecurityScheme {
	// https://swagger.io/docs/specification/authentication/api-keys/
	return &openapi3.SecurityScheme{
		Type: apiKeyType,
		Name: name,
	}
}

func NewAPIKeySecuritySchemeInHeader(name string) *openapi3.SecurityScheme {
	return newAPIKeySecurityScheme(name).WithIn(openapi3.ParameterInHeader)
}

func NewAPIKeySecuritySchemeInQuery(name string) *openapi3.SecurityScheme {
	return newAPIKeySecurityScheme(name).WithIn(openapi3.ParameterInQuery)
}

func NewBasicAuthSecurityScheme() *openapi3.SecurityScheme {
	// https://swagger.io/docs/specification/authentication/basic-authentication/
	return &openapi3.SecurityScheme{
		Type:   basicAuthType,
		Scheme: basicAuthScheme,
	}
}

func NewOAuth2SecurityScheme(scopes []string) *openapi3.SecurityScheme {
	// https://swagger.io/docs/specification/authentication/oauth2/
	// we can't know the flow type (implicit, password, clientCredentials or authorizationCode)
	// so we choose authorizationCode for now
	return &openapi3.SecurityScheme{
		Type: oauth2Type,
		Flows: &openapi3.OAuthFlows{
			AuthorizationCode: &openapi3.OAuthFlow{
				AuthorizationURL: authorizationURL,
				TokenURL:         tknURL,
				Scopes:           createOAuthFlowScopes(scopes, []string{}),
			},
		},
	}
}

func updateSecuritySchemesFromOperation(securitySchemes openapi3.SecuritySchemes, op *openapi3.Operation) openapi3.SecuritySchemes {
	if op == nil || op.Security == nil {
		return securitySchemes
	}

	// Note: usage goes in the other direction; i.e., the security schemes do contain more detail, and operations
	// (security requirements) reference those schemes. The reference is required to be valid (i.e., the
	// name in the operation MUST be present in the security schemes) for OAuth openapi3 v2.0.  Here we assume
	// schemes are generic to push the operation's security requirements into the general security schemes.
	for _, securityGroup := range *op.Security {
		for key := range securityGroup {
			var scheme *openapi3.SecurityScheme
			switch key {
			case BasicAuthSecuritySchemeKey:
				scheme = NewBasicAuthSecurityScheme()
			case OAuth2SecuritySchemeKey:
				// we can't know the flow type (implicit, password, clientCredentials or authorizationCode) so
				// we choose authorizationCode for now
				scheme = NewOAuth2SecurityScheme(nil)
			case BearerAuthSecuritySchemeKey:
				scheme = openapi3.NewJWTSecurityScheme()
			case APIKeyAuthSecuritySchemeKey:
				// Use random key since it is not specified
				for apiKeyName := range APIKeyNames {
					scheme = NewAPIKeySecuritySchemeInHeader(apiKeyName)
					break
				}
			default:
				logger.Warnf("Unsupported security definition key: %v", key)
			}
			securitySchemes = updateSecuritySchemes(securitySchemes, key, scheme)
		}
	}

	return securitySchemes
}

func updateSecuritySchemes(securitySchemes openapi3.SecuritySchemes, key string, securityScheme *openapi3.SecurityScheme) openapi3.SecuritySchemes {
	// we can override SecuritySchemes if exists since it has the same key and value
	switch key {
	case BasicAuthSecuritySchemeKey, OAuth2SecuritySchemeKey, APIKeyAuthSecuritySchemeKey:
		securitySchemes[key] = &openapi3.SecuritySchemeRef{Value: securityScheme}
	default:
		logger.Warnf("Unsupported security definition key: %v", key)
	}

	return securitySchemes
}

func createOAuthFlowScopes(scopes []string, descriptions []string) map[string]string {
	flowScopes := make(map[string]string)
	if len(descriptions) > 0 {
		if len(descriptions) < len(scopes) {
			logger.Errorf("too few descriptions (%v) supplied for security scheme scopes (%v)", len(descriptions), len(scopes))
		}
		for idx, scope := range scopes {
			if idx < len(descriptions) {
				flowScopes[scope] = descriptions[idx]
			} else {
				flowScopes[scope] = ""
			}
		}
	} else {
		for _, scope := range scopes {
			flowScopes[scope] = ""
		}
	}
	return flowScopes
}

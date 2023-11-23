package security

import (
	"fmt"
	"strings"

	"github.com/hotmall/go-raml/raml"
	"sort"
)

const (
	Oauth2              = "OAuth 2.0"
	BasicAuthentication = "Basic Authentication"
	PassThrough         = "Pass Through"
)

// security define a security scheme, we only support oauth2 now.
// we generate middleware that checking for oauth2 credential
type Security struct {
	*raml.SecurityScheme
	Name        string
	PackageName string
	Headers     []*raml.Header
	QueryParams []*raml.NamedParameter
	//apiDef      *raml.APIDefinition
}

// New creates a security struct
func New(ss raml.SecurityScheme, name, packageName string) Security {
	sd := Security{
		SecurityScheme: &ss,
	}
	sd.Name = SecuritySchemeName(name)
	sd.PackageName = packageName

	// assign header, if any
	for k, v := range sd.DescribedBy.Headers {
		header := v
		header.Name = string(k)
		sd.Headers = append(sd.Headers, &header)
	}
	sort.SliceStable(sd.Headers, func(i, j int) bool { return sd.Headers[i].Name < sd.Headers[j].Name })

	// assign query params if any
	for k, v := range sd.DescribedBy.QueryParameters {
		queryParam := v
		queryParam.Name = string(k)
		sd.QueryParams = append(sd.QueryParams, &queryParam)
	}
	sort.SliceStable(sd.QueryParams, func(i, j int) bool { return sd.QueryParams[i].Name < sd.QueryParams[j].Name })

	return sd
}

// Supported returns true if the security scheme is supported by go-raml
func Supported(ss raml.SecurityScheme) bool {
	switch ss.Type {
	case Oauth2:
		_, ok := ss.Settings["accessTokenUri"]
		return ok
	default:
		return false
	}
}

// GetMethodSecuredBy get SecuredBy field of a method
func GetMethodSecuredBy(apiDef *raml.APIDefinition, r *raml.Resource, m *raml.Method) []raml.DefinitionChoice {
	// get the secured by
	securedBy := func() []raml.DefinitionChoice {
		if len(m.SecuredBy) > 0 {
			return m.SecuredBy
		} else if sb := FindResourceSecuredBy(r); len(sb) > 0 {
			return sb
		}
		return apiDef.SecuredBy
	}()

	// filter only security scheme we supported
	var filtered []raml.DefinitionChoice
	for _, sb := range securedBy {
		ss, ok := apiDef.SecuritySchemes[sb.Name]
		if !ok {
			continue
		}
		if !Supported(ss) {
			continue
		}
		filtered = append(filtered, sb)
	}
	return filtered
}

// get array of security scopes in the form of quoted string
func GetQuotedScopes(ss raml.DefinitionChoice) ([]string, error) {
	var quoted []string
	scopes, err := getSecurityScopes(ss)
	if err != nil {
		return quoted, err
	}
	for _, s := range scopes {
		quoted = append(quoted, fmt.Sprintf(`"%v"`, s))
	}
	return quoted, nil
}

// get scopes of a security scheme as []string
func getSecurityScopes(ss raml.DefinitionChoice) ([]string, error) {
	scopes := []string{}

	// check if there is scopes
	v, ok := ss.Parameters["scopes"]
	if !ok {
		return scopes, nil
	}

	// cast it to []string, return error if failed
	scopesArr, ok := v.([]interface{})
	if !ok {
		return scopes, fmt.Errorf("scopes must be array")
	}

	// build []string
	for _, s := range scopesArr {
		scopes = append(scopes, s.(string))
	}
	return scopes, nil
}

// return security scheme name that could be used in code
func SecuritySchemeName(name string) string {
	return strings.Replace(name, " ", "", -1)
}

// validate security scheme:
// - not empty
// - not 'null'
// - oauth2 -> we only support oauth2 now
func ValidateScheme(name string, apiDef *raml.APIDefinition) bool {
	if name == "" || name == "null" {
		return false
	}
	if ss, ok := apiDef.GetSecurityScheme(name); ok {
		return ss.Type == Oauth2
	}
	return false
}

// find resource's securedBy recursively
func FindResourceSecuredBy(r *raml.Resource) []raml.DefinitionChoice {
	if len(r.SecuredBy) > 0 {
		return r.SecuredBy
	}
	if r.Parent == nil {
		return []raml.DefinitionChoice{}
	}
	return FindResourceSecuredBy(r.Parent)
}

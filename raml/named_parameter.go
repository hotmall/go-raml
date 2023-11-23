package raml

// NamedParameter is collection of named parameters
// The RAML Specification uses collections of named parameters for the
// following properties: URI parameters, query string parameters, form
// parameters, request bodies (depending on the media type), and request
// and response headers.
//
// Some fields are pointers to distinguish Zero values and no values
type NamedParameter struct {

	// NOTE: We currently do not support Named Parameters With Multiple Types.
	// TODO: Add support for Named Parameters With Multiple Types. Should be
	// done sort of like the DefinitionChoice type.

	// The name of the Parameter, as defined by the type containing it.
	Name string
	// TODO: Fill this during the post-processing phase

	// A friendly name used only for display or documentation purposes.
	// If displayName is not specified, it defaults to the property's key
	DisplayName string `yaml:"displayName"` // TODO: Auto-fill this

	// The intended use or meaning of the parameter
	Description string `yaml:"description"`

	// The primitive type of the parameter's resolved value. Can be:
	//
	// Type	Description
	// string	- Value MUST be a string.
	// number	- Value MUST be a number. Indicate floating point numbers as defined by YAML.
	// integer	- Value MUST be an integer. Floating point numbers are not allowed. The integer type is a subset of the number type.
	// date		- Value MUST be a string representation of a date as defined in RFC2616 Section 3.3 [RFC2616]. See Date Representations.
	// boolean	- Value MUST be either the string "true" or "false" (without the quotes).
	// file		- (Applicable only to Form properties) Value is a file. Client generators SHOULD use this type to handle file uploads correctly.
	Type string
	// TODO: Verify the enum options

	// If the enum attribute is defined, API clients and servers MUST verify
	// that a parameter's value matches a value in the enum array
	// Enum parsing is currently disabled because of:
	// https://github.com/hotmall/go-raml/issues/99
	//Enum []Any `yaml:",flow"`

	// The pattern attribute is a regular expression that a parameter of type
	// string MUST match. Regular expressions MUST follow the regular
	// expression specification from ECMA 262/Perl 5. (string only)
	Pattern *string

	// The minLength attribute specifies the parameter value's minimum number
	// of characters (string only)
	MinLength *int `yaml:"minLength"`
	// TODO: go-yaml doesn't raise an error when the minLength isn't an integer!
	// find out why and fix it.

	// The maxLength attribute specifies the parameter value's maximum number
	// of characters (string only)
	MaxLength *int `yaml:"maxLength"`

	// The minimum attribute specifies the parameter's minimum value. (numbers
	// only)
	Minimum *float64

	// The maximum attribute specifies the parameter's maximum value. (numbers
	// only)
	Maximum *float64

	// An example value for the property. This can be used, e.g., by
	// documentation generators to generate sample values for the property.
	Example interface{}

	// The repeat attribute specifies that the parameter can be repeated,
	// i.e. the parameter can be used multiple times
	Repeat *bool // TODO: What does this mean?

	// Whether the parameter and its value MUST be present when a call is made.
	// In general, parameters are optional unless the required attribute is
	// included and its value set to 'true'.
	// For a URI parameter, its default value is 'true'.
	Required bool

	// The default value to use for the property if the property is omitted or
	// its value is not specified
	Default Any

	format Any `ramlFormat:"Named parameters must be mappings. Example: userId: {displayName: 'User ID', description: 'Used to identify the user.', type: 'integer', minimum: 1, example: 5}"`
}

// check if an element exist in enum field
/*func (np *NamedParameter) existInEnum(elem Any) bool {
	for _, e := range np.Enum {
		if reflect.DeepEqual(e, elem) {
			return true
		}
	}
	return false
}
*/

func (np *NamedParameter) inherit(parent NamedParameter, dicts map[string]interface{}) {
	np.Name = substituteParams(np.Name, parent.Name, dicts)
	np.DisplayName = substituteParams(np.DisplayName, parent.DisplayName, dicts)
	np.Description = substituteParams(np.Description, parent.Description, dicts)
	np.Type = parent.Type

	/*
		for _, elem := range parent.Enum {
			if !np.existInEnum(elem) {
				np.Enum = append(np.Enum, elem)
			}
		}
	*/
	np.Pattern = inheritStringPointer(np.Pattern, parent.Pattern, dicts)
	np.MinLength = inheritIntPointer(np.MinLength, parent.MinLength)
	np.MaxLength = inheritIntPointer(np.MaxLength, parent.MaxLength)
	if parent.Maximum != nil {
		np.Maximum = parent.Maximum
	}
	if parent.Minimum != nil {
		np.Minimum = parent.Minimum
	}
	if parent.Repeat != nil {
		np.Repeat = parent.Repeat
	}
	if parent.Required {
		np.Required = true
	}
}

func inheritStringPointer(val, parent *string, dicts map[string]interface{}) *string {
	if parent == nil {
		return val
	}
	if val == nil {
		val = new(string)
	}
	*val = substituteParams(*val, *parent, dicts)
	return val
}

func inheritIntPointer(val, parent *int) *int {
	if parent == nil {
		return val
	}
	return parent
}

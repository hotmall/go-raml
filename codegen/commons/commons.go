package commons

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	log "github.com/Sirupsen/logrus"

	"github.com/Jumpscale/go-raml/codegen/templates"
	"github.com/Jumpscale/go-raml/raml"
	"github.com/Jumpscale/go-raml/utils"
)

var (
	regValidIdentifier = regexp.MustCompile(`[^A-Za-z0-9\[\]]+`)
)

const (
	underscore = "_"
)

const (
	// ReqBodySuffix is suffix name for request body object
	ReqBodySuffix = "ReqBody"

	// RespBodySuffix is suffix name for response body object
	RespBodySuffix = "RespBody"

	LangGo     = "go"
	LangPython = "python"
)

// doNormalizeURI removes `{`, `}`, and `/` from an URI
func doNormalizeURI(URI string) string {
	s := strings.Replace(URI, "/", " ", -1)
	s = strings.Replace(s, "{", "", -1)
	return strings.Replace(s, "}", "", -1)
}

// NormalizeURI removes `{`, `}`, `/`, and space from an URI
func NormalizeURI(URI string) string {
	return strings.Replace(doNormalizeURI(URI), " ", "", -1)
}

// NormalizeURITitle does NormalizeURI with first character in upper case
func NormalizeURITitle(URI string) string {
	s := strings.Title(doNormalizeURI(URI))
	return strings.Replace(s, " ", "", -1)

}

// ParseDescription create string slice from an RAML description.
// each element is a  description line
func ParseDescription(desc string) []string {
	if desc == "" {
		return nil
	}
	// we need to trim it because our parser usually give
	// space after last newline
	desc = strings.TrimSpace(desc)

	if desc == "" {
		return []string{}
	}

	return strings.Split(desc, "\n")
}

// GenerateFile generates file from a template.
// if file already exist and override==false, file won't be regenerated
// funcMap = this parameter is used for passing go function to the template
func GenerateFile(data interface{}, tmplFile, tmplName, filename string, override bool) error {
	if !override && isFileExist(filename) {
		log.Infof("file %v already exist and override=false, no need to regenerate", filename)
		return nil
	}

	if err := checkCreateDir(filepath.Dir(filename)); err != nil {
	}

	// pass Go function to template
	funcMap := template.FuncMap{
		"ToLower":           strings.ToLower,
		"Title":             strings.Title,
		"Camelize":          utils.Camelize,
		"CamelizeDownFirst": utils.CamelizeDownFirst,
	}

	// all template files path is relative to current directory (./)
	// while go-bindata files exist in templates directory
	tmplFile = strings.Replace(tmplFile, "./", "", -1)

	byteData, err := templates.Asset(tmplFile)
	if err != nil {
		return err
	}

	t, err := template.New(tmplName).Funcs(funcMap).Parse(string(byteData))
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Infof("generating file %v", filename)
	if err := t.ExecuteTemplate(f, tmplName, data); err != nil {
		return err
	}

	switch {
	case strings.HasSuffix(filename, ".go"):
		return runGoFmt(filename)
	case strings.HasSuffix(filename, ".py"):
		return runAutoPep8(filename)
	}
	return nil
}

// MapToSortedStrings returns sorted string arrays from a map
func MapToSortedStrings(m map[string]struct{}) []string {
	ss := []string{}
	for k := range m {
		ss = append(ss, k)
	}
	sort.Strings(ss)
	return ss
}

// cek if a file exist
func isFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// ParamizingURI creates parameterized URI
// Input : raw string, ex : /users/{userId}/address/{addressId}
// Output : "/users/"+userId+"/address/"+addressId
// TODO : optimize with regex
func ParamizingURI(URI, sep string) string {
	uri := `"` + URI + `"`
	// replace { with "+
	uri = strings.Replace(uri, "{", `"`+sep, -1)

	// if ended with }/" or }", remove trailing "
	if strings.HasSuffix(uri, `}/"`) || strings.HasSuffix(uri, `}"`) {
		uri = uri[:len(uri)-1]
	}

	// replace } with +"
	uri = strings.Replace(uri, "}", sep+`"`, -1)

	// clean trailing +"
	if strings.HasSuffix(uri, sep+`"`) {
		uri = uri[:len(uri)-2]
	}
	return postProcessParamURI(uri, sep)
}

// run `go fmt` command to a file
func runGoFmt(filePath string) error {
	args := []string{
		"-w",
		filePath,
	}

	if out, err := exec.Command("gofmt", args...).CombinedOutput(); err != nil {
		log.Errorf("Error running go fmt on '%s' failed:\n%s", filePath, string(out))
		return errors.New("go fmt failed")
	}
	return nil
}

func runAutoPep8(filename string) error {
	args := []string{
		"-a",
		"-a",
		"--in-place",
		"--max-line-length",
		"120",
		filename,
	}
	if out, err := exec.Command("autopep8", args...).CombinedOutput(); err != nil {
		log.Errorf("Error running autopep8 on '%s' failed:\n%s",
			filename, string(out))
		return errors.New("autopep8 failed: make sure you have it installed")
	}
	return nil
}

// normalize package name because not all characters can be used as package name
func NormalizePkgName(name string) string {
	return strings.Replace(name, "-", "_", -1)
}

// inverse of normalizePkgName
func DenormalizePkgName(name string) string {
	return strings.Replace(name, "_", "-", -1)
}

// AtoiOrPanic convert a string to int and panic if failed
func AtoiOrPanic(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("%v is not valid integer string. err = %v", str, err)
	}
	return i
}

// create directory if not exist
func checkCreateDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0777)
	}
	return nil
}

// IsStrInArray check if a string `str` is part of array `arr`
func IsStrInArray(arr []string, str string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}
	return false
}

// NormalizeIdentifier change invalid character in identifier
// to `_`.
// Edge cases:
// - If started with invalid char, we prepend with `The_`
// - don't replace `.` if it means a library
func NormalizeIdentifier(s string) string {
	if s == "" {
		return s
	}

	str := regValidIdentifier.ReplaceAllString(s, underscore)

	// it needs to be started with letter
	// if not, prepend `The_`
	startedWithLetter := func() bool {
		r := []rune(str)
		return unicode.IsLetter(r[0])
	}()
	if !startedWithLetter && !strings.HasPrefix(str, "[]") {
		str = "The_" + str
	}

	return strings.Trim(str, "_")
}

func NormalizeIdentifierWithLib(s string, apiDef *raml.APIDefinition) string {
	if strings.Index(s, ".") < 0 {
		return NormalizeIdentifier(s)
	}
	splitted := strings.Split(s, ".")
	if len(splitted) != 2 {
		log.Fatalf("libImportPath invalid:%v", s)
	}

	if _, libFile := apiDef.FindLibFile(DenormalizePkgName(splitted[0])); libFile == "" {
		// the '.' doesn't mean library
		return NormalizeIdentifier(s)
	}
	return splitted[0] + "." + NormalizeIdentifier(splitted[1])
}

func DisplayNameToFuncName(str string) string {
	str = strings.Replace(str, " ", "", -1) // remove the space
	return NormalizeIdentifier(str)         // change the other to _
}

func SnackCaseServerMethodName(displayName, verb string, resource *raml.Resource) string {
	if len(displayName) > 0 {
		return DisplayNameToFuncName(displayName)
	}
	return NormalizeIdentifier(snakeCaseResourceURI(resource) + "_" + strings.ToLower(verb))
}

// create snake case function name from a resource URI
func snakeCaseResourceURI(r *raml.Resource) string {
	return _snakeCaseResourceURI(r, "")
}

func _snakeCaseResourceURI(r *raml.Resource, completeURI string) string {
	if r == nil {
		return completeURI
	}
	var snake string
	if len(r.URI) > 0 {
		uri := NormalizeURI(r.URI)
		if len(uri) > 0 {
			if r.Parent != nil { // not root resource, need to add "_"
				snake = "_"
			}

			if strings.HasPrefix(r.URI, "/{") {
				snake += "by" + strings.ToUpper(uri[:1])
			} else {
				snake += strings.ToLower(uri[:1])
			}

			if len(uri) > 1 { // append with the rest of uri
				snake += uri[1:]
			}
		}
	}
	return _snakeCaseResourceURI(r.Parent, snake+completeURI)
}

func postProcessParamURI(uri, sep string) string {
	arrs := strings.Split(uri, sep)
	for i, elem := range arrs {
		if strings.HasPrefix(elem, `"`) {
			continue
		}
		arrs[i] = utils.CamelizeDownFirst(elem)
	}
	return strings.Join(arrs, sep)
}

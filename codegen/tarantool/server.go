package tarantool

import (
	"os/exec"
	"path"

	"fmt"

	"path/filepath"

	"os"
	"strings"

	"github.com/hotmall/go-raml/codegen/capnp"
	"github.com/hotmall/go-raml/codegen/commons"
	"github.com/hotmall/go-raml/raml"
)

const (
	handlersDir = "handlers"
)

// Server represents a tarantool server
type Server struct {
	apiDef     *raml.APIDefinition
	apiDocsDir string // apidocs directory. apidocs won't be generated if it is empty
	TargetDir  string // root directory of the generated code
	Resources  []tarantoolResource
}

// NewServer creates a new tarantool server
func NewServer(apiDef *raml.APIDefinition, apiDocsDir string, targetDir string) *Server {
	resources := getServerResourcesDefs(apiDef)
	return &Server{
		apiDef:     apiDef,
		apiDocsDir: apiDocsDir,
		TargetDir:  targetDir,
		Resources:  resources,
	}
}

// APIDocsDir implements generator.Server.APIDocsDir interface
func (s *Server) APIDocsDir() string {
	return s.apiDocsDir
}

// Generate implements generator.Server.Generate interface
func (s *Server) Generate() error {
	if err := s.generateSchemas(); err != nil {
		return err
	}

	if err := s.generateMain(); err != nil {
		return err
	}

	return s.generateResources()
}

func (s *Server) generateSchemas() error {
	schemasPath := path.Join(s.TargetDir, handlersDir, "schemas")
	err := capnp.GenerateCapnp(s.apiDef, schemasPath, "plain", "")
	if err != nil {
		return err
	}

	args, err := filepath.Glob(path.Join(schemasPath, "*.capnp"))
	if err != nil {
		return err
	}

	// lua generator used the name of the first file in the list to determine the name of the compiled lua file
	compiledFile := strings.Replace(args[0], ".capnp", "_capnp.lua", 1)

	args = append([]string{"compile", "-olua"}, args...)
	cmd := exec.Command("capnp", args...)
	err = cmd.Run()
	if err != nil {
		return err
	}

	if _, err = os.Stat(compiledFile); os.IsNotExist(err) {
		return fmt.Errorf(
			"Can't find the compiled lua schema file expected at %v", compiledFile)
	}

	// rename the compiled lua file to a more generic name
	err = os.Rename(compiledFile, path.Join(
		schemasPath, "schema.lua"))
	if err != nil {
		return err
	}

	return nil
}

// generateMain generates the main server file
func (s *Server) generateMain() error {
	filename := path.Join(s.TargetDir, "main.lua")
	return commons.GenerateFile(s, "./templates/tarantool/server_main.tmpl", "server_main", filename, true)
}

// generateResources generates all resource apis and end point handlers
func (s *Server) generateResources() error {
	var allMethods []method

	for _, resource := range s.Resources {
		filename := path.Join(s.TargetDir, fmt.Sprintf("%v_api", strings.ToLower(resource.Name))+".lua")

		if err := commons.GenerateFile(resource, "./templates/tarantool/server_resource_api.tmpl",
			"server_resource_api", filename, true); err != nil {
			return err
		}

		for _, method := range resource.Methods {
			allMethods = append(allMethods, *method)
			filename := path.Join(s.TargetDir, handlersDir, method.Handler()+".lua")

			// generate method handler file
			if err := commons.GenerateFile(method, "./templates/tarantool/server_method_handler.tmpl",
				"server_method_handler", filename, false); err != nil {
				return err
			}
		}
	}

	methods := map[string]interface{}{
		"Methods": allMethods,
	}

	// Generate handlers file
	filename := path.Join(s.TargetDir, handlersDir, "handlers.lua")
	return commons.GenerateFile(methods, "./templates/tarantool/server_handlers.tmpl", "server_handlers", filename, true)
}

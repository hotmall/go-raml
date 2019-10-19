package gorestful

import (
	"path/filepath"
	"strings"

	"github.com/Jumpscale/go-raml/codegen/commons"
)

func (gr *goResource) ServiceImporters() []string {
	ip := map[string]struct{}{}
	ip[`"`+globRootImportPath+"/types"+`"`] = struct{}{}
	ip[`"`+globRootImportPath+"/delegate"+`"`] = struct{}{}
	ip[`"github.com/emicklei/go-restful"`] = struct{}{}

	isStrConv := false
	for _, m := range gr.Methods {
		for _, p := range m.QueryParameters {
			if p.Type != "string" {
				isStrConv = true
				break
			}
		}

		if !isStrConv {
			for _, p := range m.URIParameters {
				if p.Type != "string" {
					isStrConv = true
					break
				}
			}
		}
	}

	if isStrConv {
		ip[`"strconv"`] = struct{}{}
	}
	return commons.MapToSortedStrings(ip)
}

func (gr *goResource) DelegateImporters() []string {
	ip := map[string]struct{}{}
	ip[`"`+globRootImportPath+"/types"+`"`] = struct{}{}
	return commons.MapToSortedStrings(ip)
}

func (gr *goResource) ResourceImporters() []string {
	ip := map[string]struct{}{}
	ip[`"`+globRootImportPath+"/types"+`"`] = struct{}{}
	ip[`"`+globRootImportPath+"/service"+`"`] = struct{}{}
	ip[`"github.com/emicklei/go-restful"`] = struct{}{}
	return commons.MapToSortedStrings(ip)
}

func (gr *goResource) generateService(dir string) error {
	// Generate method context
	for _, m := range gr.Methods {
		gc := newGoContext(m.method)
		if err := gc.generate(dir); err != nil {
			return err
		}
	}

	// Generate delegate
	filename := filepath.Join(dir+"/"+delegateDir, strings.ToLower(gr.Name)+"_if") + ".go"
	if err := commons.GenerateFile(gr, "./templates/golang/gorestful_delegate.tmpl",
		"gorestful_delegate_template", filename, true); err != nil {
		return err
	}

	// Generate service
	filename = filepath.Join(dir+"/"+serviceDir, strings.ToLower(gr.Name)+"_service") + ".go"
	if err := commons.GenerateFile(gr, "./templates/golang/gorestful_service.tmpl",
		"gorestful_service_template", filename, true); err != nil {
		return err
	}

	// Generate resource
	filename = filepath.Join(dir+"/"+resourceDir, strings.ToLower(gr.Name)+"_resource") + ".go"
	if err := commons.GenerateFile(gr, "./templates/golang/gorestful_resource.tmpl",
		"gorestful_resource_template", filename, false); err != nil {
		return err
	}

	return nil
}

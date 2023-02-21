package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

const (
	//defaultProjectBaseDir = "test"
	defaultProjectBaseDir    = ".."
	defaultTemplatesDir      = "templates"
	defaultTemplateExtension = ".tmpl"

	defaultBasePackage = "github.com/gostevedore/stevedore"

	cliComponentName         = "cli"
	entrypointComponentName  = "entrypoint"
	applicationComponentName = "application"
	handlerComponentName     = "handler"

	defaultCliRelativeDir         = "internal/infrastructure/cli"
	defaultEntrypointRelativeDir  = "internal/entrypoint"
	defaultApplicationRelativeDir = "internal/application"
	defaultHandlerRelativeDir     = "internal/handler"
)

func main() {

	var err error

	if os.Args == nil || len(os.Args) < 2 {
		fmt.Println("Please provide a use case [example: get/images]")
		os.Exit(1)
	}

	newUseCase := os.Args[1]

	newUseCaseCliComponent := NewComponent(cliComponentName, newUseCase, defaultCliRelativeDir, defaultTemplatesDir, defaultTemplateExtension)
	newUseCaseEntrypointComponent := NewComponent(entrypointComponentName, newUseCase, defaultEntrypointRelativeDir, defaultTemplatesDir, defaultTemplateExtension)
	newUseCaseApplicationComponent := NewComponent(applicationComponentName, newUseCase, defaultApplicationRelativeDir, defaultTemplatesDir, defaultTemplateExtension)
	newUseCaseHandlerComponent := NewComponent(handlerComponentName, newUseCase, defaultHandlerRelativeDir, defaultTemplatesDir, defaultTemplateExtension)

	err = newUseCaseCliComponent.Do()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = newUseCaseEntrypointComponent.Do()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = newUseCaseApplicationComponent.Do()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = newUseCaseHandlerComponent.Do()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

type component struct {
	componentType string
	useCase       string
	relativeDir   string
	tmplDir       string
	tmplExt       string
}

type ComponentData struct {
	ObjectNameBase string
	FileNameBase   string
	UseCase        string
	PackageName    string

	CliPackageURL         string
	EntrypointPackageURL  string
	ApplicationPackageURL string
	HandlerPackageURL     string

	CliObject         string
	EntrypointObject  string
	ApplicationObject string
	HandlerObject     string
}

func NewComponent(componentType, useCase, relativeDir, tmplDir, tmplExt string) *component {

	return &component{
		componentType: componentType,
		useCase:       useCase,
		relativeDir:   relativeDir,
		tmplDir:       tmplDir,
		tmplExt:       tmplExt,
	}
}

func (c *component) Do() error {

	fmt.Println("Creating code for '" + c.componentType + "'")
	err := c.generateCodeFromTemplates()
	fmt.Println()
	if err != nil {
		return err
	}

	return nil
}

func (c *component) createUseCaseNormalizedName() string {

	normalized := ""

	splitedName := strings.Split(c.useCase, string(filepath.Separator))
	for _, name := range splitedName {
		normalized += strings.Title(name)
	}

	return normalized
}

func (c *component) generateDestinationFileName(file string) string {
	var dest string

	switch filepath.Base(file) {
	case fmt.Sprintf("%s.go", strings.ToLower(c.componentType)):
		normalizedName := c.createUseCaseNormalizedName()
		runed := []rune(normalizedName)
		runed[0] = unicode.ToLower(runed[0])
		dest = string(runed)
		dest = fmt.Sprintf("%s.go", string(runed))
	case fmt.Sprintf("%s_test.go", strings.ToLower(c.componentType)):
		normalizedName := c.createUseCaseNormalizedName()
		runed := []rune(normalizedName)
		runed[0] = unicode.ToLower(runed[0])
		dest = fmt.Sprintf("%s_test.go", string(runed))
	default:
		dest = filepath.Base(file)
	}

	return filepath.Join(defaultProjectBaseDir, c.relativeDir, c.useCase, dest)
}

func (c *component) createComponentData() *ComponentData {

	objectNameBase := c.createUseCaseNormalizedName()

	runed := []rune(objectNameBase)
	runed[0] = unicode.ToLower(runed[0])
	fileNameBase := string(runed)

	return &ComponentData{
		ObjectNameBase: objectNameBase,
		FileNameBase:   fileNameBase,
		UseCase:        c.useCase,
		PackageName:    filepath.Base(c.useCase),

		CliPackageURL:         filepath.Join(defaultBasePackage, defaultCliRelativeDir, c.useCase),
		EntrypointPackageURL:  filepath.Join(defaultBasePackage, defaultEntrypointRelativeDir, c.useCase),
		ApplicationPackageURL: filepath.Join(defaultBasePackage, defaultApplicationRelativeDir, c.useCase),
		HandlerPackageURL:     filepath.Join(defaultBasePackage, defaultHandlerRelativeDir, c.useCase),

		CliObject:         c.createUseCaseNormalizedName(),
		EntrypointObject:  c.createUseCaseNormalizedName() + "Entrypoint",
		ApplicationObject: c.createUseCaseNormalizedName() + "Application",
		HandlerObject:     c.createUseCaseNormalizedName() + "Handler",
	}
}

func (c *component) generateCodeFromTemplates() error {

	templatesDir := filepath.Join(c.tmplDir, c.componentType)

	data := c.createComponentData()

	err := filepath.Walk(templatesDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				dest := path[:len(path)-len(c.tmplExt)]

				_, err := os.Stat(c.generateDestinationFileName(dest))
				if errors.Is(err, os.ErrNotExist) {

					os.MkdirAll(filepath.Dir(c.generateDestinationFileName(dest)), 0755)

					fmt.Printf("  Creating file %s from %s... ", c.generateDestinationFileName(dest), path)
					fileTemplate, err := template.ParseFiles(path)
					if err != nil {
						fmt.Println("failed")
						return err
					}

					file, err := os.Create(c.generateDestinationFileName(dest))
					if err != nil {
						fmt.Println("failed")
						return err
					}

					err = fileTemplate.Execute(file, data)
					if err != nil {
						fmt.Println("failed")
						return err
					}
					fmt.Println("done")
				}
			}

			return nil
		})
	if err != nil {
		return err
	}

	return nil
}

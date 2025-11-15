package project

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type Analyzer struct {
	finder *Finder
}

func NewAnalyzer(finder *Finder) *Analyzer {
	return &Analyzer{
		finder: finder,
	}
}

func (a *Analyzer) FindExistingEntities(domainPath string) ([]string, error) {
	root, err := a.finder.FindRoot()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(root, domainPath)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var entities []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		if strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}

		filePath := filepath.Join(fullPath, entry.Name())
		entityNames, err := a.extractStructNames(filePath)
		if err != nil {
			continue
		}

		entities = append(entities, entityNames...)
	}

	return entities, nil
}

func (a *Analyzer) extractStructNames(filePath string) ([]string, error) {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var names []string

	for _, decl := range node.Decls {

		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				names = append(names, typeSpec.Name.Name)
			}
		}
	}

	return names, nil
}

func (a *Analyzer) FindExistingRepositories(repoPath string) ([]string, error) {
	root, err := a.finder.FindRoot()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(root, repoPath)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var repos []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		if strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".go")
		name = strings.TrimSuffix(name, "_repository")
		name = strings.TrimSuffix(name, "_repo")

		parts := strings.Split(name, "_")
		for i, part := range parts {
			if len(part) > 0 {
				parts[i] = strings.ToUpper(part[:1]) + part[1:]
			}
		}

		repoName := strings.Join(parts, "")
		if repoName != "" {
			repos = append(repos, repoName)
		}
	}

	return repos, nil
}

func (a *Analyzer) FileExists(path string) bool {
	root, err := a.finder.FindRoot()
	if err != nil {
		return false
	}

	fullPath := filepath.Join(root, path)
	_, err = os.Stat(fullPath)
	return err == nil
}

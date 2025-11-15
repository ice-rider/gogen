package dependency

import (
	"fmt"
	"strings"

	"gogen/pkg/models"
)

type Graph struct {
	nodes map[string]*Node
	edges map[string][]string
}

type Node struct {
	Name      string
	Type      models.ComponentType
	Component models.Component
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]*Node),
		edges: make(map[string][]string),
	}
}

func (g *Graph) AddNode(name string, nodeType models.ComponentType, component models.Component) {
	g.nodes[name] = &Node{
		Name:      name,
		Type:      nodeType,
		Component: component,
	}
}

func (g *Graph) AddEdge(from, to string) {
	if g.edges[from] == nil {
		g.edges[from] = make([]string, 0)
	}
	g.edges[from] = append(g.edges[from], to)
}

func (g *Graph) BuildFromPlan(plan *models.GenerationPlan) {

	for i := range plan.Entities {
		entity := &plan.Entities[i]
		g.AddNode(entity.Name, models.ComponentTypeEntity, entity)
	}

	for i := range plan.Repositories {
		repo := &plan.Repositories[i]
		repoName := repo.Name + "Repository"
		g.AddNode(repoName, models.ComponentTypeRepository, repo)

		g.AddEdge(repoName, repo.Entity)
	}

	for i := range plan.UseCases {
		uc := &plan.UseCases[i]
		ucName := uc.Name + "UseCase"
		g.AddNode(ucName, models.ComponentTypeUseCase, uc)

		for _, dep := range uc.Dependencies {
			g.AddEdge(ucName, dep.Name)
		}
	}
}

func (g *Graph) TopologicalSort() ([]string, error) {
	visited := make(map[string]bool)
	temp := make(map[string]bool)
	result := make([]string, 0)

	var visit func(string) error
	visit = func(node string) error {
		if temp[node] {
			return fmt.Errorf("circular dependency detected at %s", node)
		}

		if visited[node] {
			return nil
		}

		temp[node] = true

		for _, dep := range g.edges[node] {
			if err := visit(dep); err != nil {
				return err
			}
		}

		temp[node] = false
		visited[node] = true
		result = append(result, node)

		return nil
	}

	for node := range g.nodes {
		if !visited[node] {
			if err := visit(node); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func (g *Graph) DetectCycles() [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(string, []string) bool
	dfs = func(node string, path []string) bool {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)

		for _, neighbor := range g.edges[node] {
			if !visited[neighbor] {
				if dfs(neighbor, path) {
					return true
				}
			} else if recStack[neighbor] {

				cycleStart := 0
				for i, n := range path {
					if n == neighbor {
						cycleStart = i
						break
					}
				}
				cycle := append([]string{}, path[cycleStart:]...)
				cycles = append(cycles, cycle)
			}
		}

		recStack[node] = false
		return false
	}

	for node := range g.nodes {
		if !visited[node] {
			dfs(node, []string{})
		}
	}

	return cycles
}

func (g *Graph) Print() string {
	var sb strings.Builder

	sb.WriteString("Dependency Graph:\n")
	sb.WriteString("================\n\n")

	for node, deps := range g.edges {
		nodeInfo := g.nodes[node]
		sb.WriteString(fmt.Sprintf("%s (%s)\n", node, nodeInfo.Type))

		if len(deps) > 0 {
			sb.WriteString("  depends on:\n")
			for _, dep := range deps {
				depInfo := g.nodes[dep]
				sb.WriteString(fmt.Sprintf("    - %s (%s)\n", dep, depInfo.Type))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

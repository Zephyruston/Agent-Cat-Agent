package analysis

import (
	"bufio"
	"os"
	"strings"
)

type Dependency struct {
	Name    string
	Version string
}

func ParseGoMod(path string) ([]Dependency, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var deps []Dependency
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "require ") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				deps = append(deps, Dependency{Name: fields[1], Version: fields[2]})
			}
		}
	}
	return deps, scanner.Err()
}

func ParseRequirements(path string) ([]Dependency, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var deps []Dependency
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "==")
		if len(parts) == 2 {
			deps = append(deps, Dependency{Name: parts[0], Version: parts[1]})
		}
	}
	return deps, scanner.Err()
}

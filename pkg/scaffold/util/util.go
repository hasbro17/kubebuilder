/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v1"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/input"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/project"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/v1/resource"
)

func GetResourceInfo(r *resource.Resource, repo, domain string) (resourcePackage, groupDomain string) {
	// Use the k8s.io/api package for core resources
	coreGroups := map[string]string{
		"apps":                  "",
		"admission":             "k8s.io",
		"admissionregistration": "k8s.io",
		"auditregistration":     "k8s.io",
		"apiextensions":         "k8s.io",
		"authentication":        "k8s.io",
		"authorization":         "k8s.io",
		"autoscaling":           "",
		"batch":                 "",
		"certificates":          "k8s.io",
		"coordination":          "k8s.io",
		"core":                  "",
		"events":                "k8s.io",
		"extensions":            "",
		"imagepolicy":           "k8s.io",
		"networking":            "k8s.io",
		"node":                  "k8s.io",
		"metrics":               "k8s.io",
		"policy":                "",
		"rbac.authorization":    "k8s.io",
		"scheduling":            "k8s.io",
		"setting":               "k8s.io",
		"storage":               "k8s.io",
	}
	resourcePath := filepath.Join("api", r.Version, fmt.Sprintf("%s_types.go", strings.ToLower(r.Kind)))
	if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
		if domain, found := coreGroups[r.Group]; found {
			// TODO: support apiextensions.k8s.io and metrics.k8s.io.
			// apiextensions.k8s.io is in k8s.io/apiextensions-apiserver/pkg/apis/apiextensions
			// metrics.k8s.io is in k8s.io/metrics/pkg/apis/metrics
			resourcePackage := path.Join("k8s.io", "api", r.Group)
			groupDomain = r.Group
			if domain != "" {
				groupDomain = r.Group + "." + domain
			}
			return resourcePackage, groupDomain
		}
		// TODO: need to support '--resource-pkg-path' flag for specifying resourcePath
	}
	return path.Join(repo, "api"), r.Group + "." + domain
}

// LoadProjectFile reads the project file and deserializes it into a Project
func LoadProjectFile(path string) (input.ProjectFile, error) {
	in, err := ioutil.ReadFile(path) // nolint: gosec
	if err != nil {
		return input.ProjectFile{}, err
	}
	p := input.ProjectFile{}
	err = yaml.Unmarshal(in, &p)
	if err != nil {
		return input.ProjectFile{}, err
	}
	if p.Version == "" {
		// older kubebuilder project does not have scaffolding version
		// specified, so default it to Version1
		p.Version = project.Version1
	}
	return p, nil
}

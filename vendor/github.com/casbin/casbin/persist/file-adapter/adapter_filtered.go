// Copyright 2017 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fileadapter

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
)

// FilteredAdapter is the filtered file adapter for Casbin. It can load policy
// from file or save policy to file and supports loading of filtered policies.
type FilteredAdapter struct {
	*Adapter
	filtered bool
}

// Filter defines the filtering rules for a FilteredAdapter's policy. Empty values
// are ignored, but all others must match the filter.
type Filter struct {
	P []string
	G []string
}

// NewFilteredAdapter is the constructor for FilteredAdapter.
func NewFilteredAdapter(filePath string) *FilteredAdapter {
	a := FilteredAdapter{}
	a.filtered = true
	a.Adapter = NewAdapter(filePath)
	return &a
}

// LoadPolicy loads all policy rules from the storage.
func (a *FilteredAdapter) LoadPolicy(model model.Model) error {
	a.filtered = false
	return a.Adapter.LoadPolicy(model)
}

// LoadFilteredPolicy loads only policy rules that match the filter.
func (a *FilteredAdapter) LoadFilteredPolicy(model model.Model, filter interface{}) error {
	if filter == nil {
		return a.LoadPolicy(model)
	}
	if a.filePath == "" {
		return errors.New("invalid file path, file path cannot be empty")
	}

	filterValue, ok := filter.(*Filter)
	if !ok {
		return errors.New("invalid filter type")
	}
	err := a.loadFilteredPolicyFile(model, filterValue, persist.LoadPolicyLine)
	if err == nil {
		a.filtered = true
	}
	return err
}

func (a *FilteredAdapter) loadFilteredPolicyFile(model model.Model, filter *Filter, handler func(string, model.Model)) error {
	f, err := os.Open(a.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if filterLine(line, filter) {
			continue
		}

		handler(line, model)
	}
	return scanner.Err()
}

// IsFiltered returns true if the loaded policy has been filtered.
func (a *FilteredAdapter) IsFiltered() bool {
	return a.filtered
}

// SavePolicy saves all policy rules to the storage.
func (a *FilteredAdapter) SavePolicy(model model.Model) error {
	if a.filtered {
		return errors.New("cannot save a filtered policy")
	}
	return a.Adapter.SavePolicy(model)
}

func filterLine(line string, filter *Filter) bool {
	if filter == nil {
		return false
	}
	p := strings.Split(line, ",")
	if len(p) == 0 {
		return true
	}
	var filterSlice []string
	switch strings.TrimSpace(p[0]) {
	case "p":
		filterSlice = filter.P
	case "g":
		filterSlice = filter.G
	}
	return filterWords(p, filterSlice)
}

func filterWords(line []string, filter []string) bool {
	if len(line) < len(filter)+1 {
		return true
	}
	var skipLine bool
	for i, v := range filter {
		if len(v) > 0 && strings.TrimSpace(v) != strings.TrimSpace(line[i+1]) {
			skipLine = true
			break
		}
	}
	return skipLine
}

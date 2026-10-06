package validation_test

import (
	"strings"
	"testing"

	"github.com/beego/beego/v2/core/validation"
)

type recursiveErrorChild struct {
	Name string `valid:"UnknownRecursiveRule"`
}

type recursiveValidChild struct {
	Name string `valid:"Required"`
}

func TestRecursiveValidPreservesChildError(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input interface{}
	}{
		{"value siblings", struct {
			First recursiveErrorChild
			Last  recursiveValidChild
		}{Last: recursiveValidChild{Name: "valid"}}},
		{"pointer siblings", &struct {
			First *recursiveErrorChild
			Last  *recursiveValidChild
		}{First: &recursiveErrorChild{}, Last: &recursiveValidChild{Name: "valid"}}},
		{"nested siblings", struct {
			First struct{ Child recursiveErrorChild }
			Last  recursiveValidChild
		}{Last: recursiveValidChild{Name: "valid"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			valid := validation.Validation{}
			pass, err := valid.RecursiveValid(tc.input)
			if pass || err == nil {
				t.Fatalf("RecursiveValid = (%v, %v), want false and child rule error", pass, err)
			}
			if !strings.Contains(err.Error(), "UnknownRecursiveRule") {
				t.Fatalf("lost child rule diagnostic: %v", err)
			}
			if valid.HasErrors() {
				t.Fatal("an invalid rule should return an error, not a validation failure")
			}
		})
	}
}

func TestRecursiveValidSiblingControls(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    interface{}
		wantPass bool
	}{
		{"valid siblings", struct {
			First recursiveValidChild
			Last  recursiveValidChild
		}{recursiveValidChild{"first"}, recursiveValidChild{"last"}}, true},
		{"first validation failure", struct {
			First recursiveValidChild
			Last  recursiveValidChild
		}{Last: recursiveValidChild{"last"}}, false},
		{"last validation failure", struct {
			First recursiveValidChild
			Last  recursiveValidChild
		}{First: recursiveValidChild{"first"}}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			valid := validation.Validation{}
			pass, err := valid.RecursiveValid(tc.input)
			if err != nil || pass != tc.wantPass || valid.HasErrors() == tc.wantPass {
				t.Fatalf("RecursiveValid = (%v, %v), HasErrors = %v, want pass %v", pass, err, valid.HasErrors(), tc.wantPass)
			}
		})
	}
}

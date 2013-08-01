package beego

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildTemplate(t *testing.T) {
	dir := "_beeTmp"
	files := []string{
		"1.tpl",
		"2.html",
		"3.htmltpl",
		"4.mystyle",
	}
	if err := os.MkdirAll(dir, 0777); err != nil {
		t.Fatal(err)
	}
	for _, name := range files {
		if _, err := os.Create(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}
	if err := BuildTemplate(dir); err != nil {
		t.Fatal(err)
	}
	if len(BeeTemplates) != 1 {
		t.Fatalf("should be 1 but got %v", len(BeeTemplates))
	}
	for _, v := range BeeTemplates {
		if len(v.Templates()) != 3 {
			t.Errorf("should be 3 but got %v", len(v.Templates()))
		}
	}

	AddTemplateExt("mystyle")
	if err := BuildTemplate(dir); err != nil {
		t.Fatal(err)
	}
	if len(BeeTemplates) != 1 {
		t.Fatalf("should be 1 but got %v", len(BeeTemplates))
	}
	for _, v := range BeeTemplates {
		if len(v.Templates()) != 4 {
			t.Errorf("should be 4 but got %v", len(v.Templates()))
		}
	}
}

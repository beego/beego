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
		if f, err := os.Create(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		} else {
			f.Close()
		}
	}
	if err := BuildTemplate(dir); err != nil {
		t.Fatal(err)
	}
	if len(BeeTemplates) != 2 {
		t.Fatalf("should be 2 but got %v", len(BeeTemplates))
	}

	AddTemplateExt("mystyle")
	if err := BuildTemplate(dir); err != nil {
		t.Fatal(err)
	}
	if len(BeeTemplates) != 3 {
		t.Fatalf("should be 3 but got %v", len(BeeTemplates))
	}
	for _, name := range files {
		os.Remove(filepath.Join(dir, name))
	}
	os.Remove(dir)
}

package logs

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestArchiver(t *testing.T) {
	for name, ar := range SupportedFormats {
		name, ar := name, ar
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// skip RAR for now
			if _, ok := ar.(rarFormat); ok {
				t.Skip("not supported")
			}
			symmetricTest(t, name, ar)
		})
	}
}

// symmetricTest performs a symmetric test by using ar.Make to make an archive
// from the test corpus, then using ar.Open to open the archive and comparing
// the contents to ensure they are equal.
func symmetricTest(t *testing.T, name string, ar Archiver) {
	tmp, err := ioutil.TempDir("", "archiver")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	// Test creating archive
	outfile := filepath.Join(tmp, "test-"+name)
	err = ar.Make(outfile, []string{"testdata"})
	if err != nil {
		t.Fatalf("making archive: didn't expect an error, but got: %v", err)
	}

	if !ar.Match(outfile) {
		t.Fatalf("identifying format should be 'true', but got 'false'")
	}

	var expectedFileCount int
	filepath.Walk("testdata", func(fpath string, info os.FileInfo, err error) error {
		expectedFileCount++
		return nil
	})

	// Test extracting archive
	dest := filepath.Join(tmp, "extraction_test")
	os.Mkdir(dest, 0755)
	err = ar.Open(outfile, dest)
	if err != nil {
		t.Fatalf("extracting archive [%s -> %s]: didn't expect an error, but got: %v", outfile, dest, err)
	}

	// If outputs equals inputs, we're good; traverse output files
	// and compare file names, file contents, and file count.

	var actualFileCount int
	filepath.Walk(dest, func(fpath string, info os.FileInfo, err error) error {
		if fpath == dest {
			return nil
		}
		actualFileCount++

		origPath, err := filepath.Rel(dest, fpath)
		if err != nil {
			t.Fatalf("%s: Error inducing original file path: %v", fpath, err)
		}

		if info.IsDir() {
			// stat dir instead of read file
			_, err = os.Stat(origPath)
			if err != nil {
				t.Fatalf("%s: Couldn't stat original directory (%s): %v",
					fpath, origPath, err)
			}
			return nil
		}

		expectedFileInfo, err := os.Stat(origPath)
		if err != nil {
			t.Fatalf("%s: Error obtaining original file info: %v", fpath, err)
		}
		expected, err := ioutil.ReadFile(origPath)
		if err != nil {
			t.Fatalf("%s: Couldn't open original file (%s) from disk: %v",
				fpath, origPath, err)
		}

		actualFileInfo, err := os.Stat(fpath)
		if err != nil {
			t.Fatalf("%s: Error obtaining actual file info: %v", fpath, err)
		}
		actual, err := ioutil.ReadFile(fpath)
		if err != nil {
			t.Fatalf("%s: Couldn't open new file from disk: %v", fpath, err)
		}

		if actualFileInfo.Mode() != expectedFileInfo.Mode() {
			t.Fatalf("%s: File mode differed between on disk and compressed",
				expectedFileInfo.Mode().String()+" : "+actualFileInfo.Mode().String())
		}
		if !bytes.Equal(expected, actual) {
			t.Fatalf("%s: File contents differed between on disk and compressed", origPath)
		}

		return nil
	})

	if got, want := actualFileCount, expectedFileCount; got != want {
		t.Fatalf("Expected %d resulting files, got %d", want, got)
	}
}

func BenchmarkMake(b *testing.B) {
	tmp, err := ioutil.TempDir("", "archiver")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	for name, ar := range SupportedFormats {
		name, ar := name, ar
		b.Run(name, func(b *testing.B) {
			// skip RAR for now
			if _, ok := ar.(rarFormat); ok {
				b.Skip("not supported")
			}
			outfile := filepath.Join(tmp, "benchMake-"+name)
			for i := 0; i < b.N; i++ {
				err = ar.Make(outfile, []string{"testdata"})
				if err != nil {
					b.Fatalf("making archive: didn't expect an error, but got: %v", err)
				}
			}
		})
	}
}

func BenchmarkOpen(b *testing.B) {
	tmp, err := ioutil.TempDir("", "archiver")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	for name, ar := range SupportedFormats {
		name, ar := name, ar
		b.Run(name, func(b *testing.B) {
			// skip RAR for now
			if _, ok := ar.(rarFormat); ok {
				b.Skip("not supported")
			}
			// prepare a archive
			outfile := filepath.Join(tmp, "benchMake-"+name)
			err = ar.Make(outfile, []string{"testdata"})
			if err != nil {
				b.Fatalf("open archive: didn't expect an error, but got: %v", err)
			}
			// prepare extraction destination
			dest := filepath.Join(tmp, "extraction_test")
			os.Mkdir(dest, 0755)

			// let's go
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err = ar.Open(outfile, dest)
				if err != nil {
					b.Fatalf("open archive: didn't expect an error, but got: %v", err)
				}
			}
		})
	}
}

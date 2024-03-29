// Copyright ©2020 The go-fonts Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	log.SetPrefix("dejavu-gen: ")
	log.SetFlags(0)

	var (
		src = flag.String(
			"src",
			"https://github.com/dejavu-fonts/dejavu-fonts/releases/download/version_2_37/dejavu-fonts-ttf-2.37.zip",
			"remote ZIP file holding OTF files for DejaVu fonts",
		)
	)

	flag.Parse()

	tmp, err := os.MkdirTemp("", "go-fonts-dejavu-")
	if err != nil {
		log.Fatalf("could not create tmp dir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	var zr *zip.ReadCloser

	switch {
	case strings.HasPrefix(*src, "http://"),
		strings.HasPrefix(*src, "https://"):
		zr, err = fetch(tmp, *src)
		if err != nil {
			log.Fatalf("could not fetch DejaVu sources: %+v", err)
		}
	default:
		zr, err = zip.OpenReader(*src)
		if err != nil {
			log.Fatalf("could not open local DejaVu sources: %+v", err)
		}
	}
	defer zr.Close()

	for _, f := range zr.File {
		if !strings.HasSuffix(f.Name, suffix) {
			continue
		}
		err := gen(f)
		if err != nil {
			log.Fatalf("could not generate font: %+v", err)
		}
	}
}

func fetch(tmp, src string) (*zip.ReadCloser, error) {
	resp, err := http.Get(src)
	if err != nil {
		return nil, fmt.Errorf("could not GET %q: %w", src, err)
	}
	defer resp.Body.Close()

	f, err := os.Create(path.Join(tmp, "dejavu.zip"))
	if err != nil {
		return nil, fmt.Errorf("could not create zip file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not copy zip file: %w", err)
	}

	err = f.Close()
	if err != nil {
		return nil, fmt.Errorf("could not save zip file: %w", err)
	}

	zr, err := zip.OpenReader(f.Name())
	if err != nil {
		return nil, fmt.Errorf("could not open zip file: %w", err)
	}
	return zr, nil
}

func gen(f *zip.File) error {
	fname := path.Base(f.Name)
	log.Printf("generating fonts package for %q...", fname)

	r, err := f.Open()
	if err != nil {
		return fmt.Errorf("could not decompress zip file %q: %w", f.Name, err)
	}
	defer r.Close()

	raw := new(bytes.Buffer)

	_, err = io.Copy(raw, r)
	if err != nil {
		return fmt.Errorf("could not download TTF file: %w", err)
	}

	err = do(fname, raw.Bytes())
	if err != nil {
		return fmt.Errorf("could not generate package for %q: %w", fname, err)
	}

	return nil
}

func do(ttfName string, src []byte) error {
	fontName := fontName(ttfName)
	pkgName := pkgName(ttfName)
	if err := os.Mkdir(pkgName, 0777); err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create package dir %q: %w", pkgName, err)
	}

	b := new(bytes.Buffer)
	fmt.Fprintf(b, "// generated by go run gen-fonts.go; DO NOT EDIT\n\n")
	fmt.Fprintf(b, "// Package %s provides the %q TrueType font\n", pkgName, fontName)
	fmt.Fprintf(b, "// from the DejaVu font family.\n")
	fmt.Fprintf(b, "package %[1]s // import \"github.com/go-fonts/dejavu/%[1]s\"\n\n", pkgName)
	fmt.Fprintf(b, "import _ \"embed\"\n")
	fmt.Fprintf(b, "// TTF is the data for the %q TrueType font.\n", fontName)
	fmt.Fprintf(b, "//\n//go:embed %s\n", ttfName)
	fmt.Fprintf(b, "var TTF  []byte\n")

	dst, err := format.Source(b.Bytes())
	if err != nil {
		return fmt.Errorf("could not format source: %w", err)
	}

	err = os.WriteFile(filepath.Join(pkgName, "data.go"), dst, 0666)
	if err != nil {
		return fmt.Errorf("could not write package source file: %w", err)
	}

	err = os.WriteFile(filepath.Join(pkgName, ttfName), src, 0666)
	if err != nil {
		return fmt.Errorf("could not write package TTF file: %w", err)
	}

	return nil
}

const suffix = ".ttf"

// fontName maps "Go-Regular.ttf" to "Go Regular".
func fontName(ttfName string) string {
	s := ttfName[:len(ttfName)-len(suffix)]
	s = strings.Replace(s, "-", " ", -1)
	return s
}

// pkgName maps "Go-Regular.ttf" to "goregular".
func pkgName(ttfName string) string {
	s := ttfName[:len(ttfName)-len(suffix)]
	s = strings.Replace(s, "-", "", -1)
	s = strings.ToLower(s)
	return s
}

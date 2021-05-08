package common

import (
	"github.com/yuin/goldmark"
	"io/ioutil"
	"os"
	"path/filepath"
)

// MdRenderRecursively render all markdown files in src
// recursively to dst into html.
// Note that src can be a single markdown file or a directory.
// Contents of src would be rendered recursively if src is a directory.
// dst would be ignored if src is a file.
// dst would be created if not exists.
// File in dst would be overwritten if overWrite is true.
// File without .md ext would be copied to dst if copyOthers is true.
func MdRenderRecursively(src string, dst string, overWrite bool, copyOthers bool) error {
	si, err := os.Stat(src)
	if err != nil {
		return err
	}

	if si.IsDir() {
		if !DirectoryExist(dst) {
			err = os.MkdirAll(dst, os.ModePerm)
			if err != nil {
				return err
			}
		}

		entries, err := ioutil.ReadDir(src)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			srcPath := filepath.Join(src, entry.Name())
			dstPath := filepath.Join(dst, entry.Name())

			err = MdRenderRecursively(srcPath, dstPath, overWrite, copyOthers)
			if err != nil {
				return err
			}
		}
	} else {
		if filepath.Ext(src) != ".md" {
			if copyOthers {
				return CopyFile(src, dst)
			}
		} else {
			dstPath := ChExt(dst, ".html")
			return MdRenderFile(src, dstPath)
		}
	}
	return nil
}

// MdRenderFile render markdown file src to html file dst.
// dst would be src with extension replaced by ".html" if dst is "".
// dst would be overwritten if exists.
// dst would be created if not exists.
func MdRenderFile(src string, dst string) error {
	if dst == "" {
		// TODO: Is it safe to change string parameter directly?
		dst = ChExt(src, ".html")
	}
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	return goldmark.Convert(input, out)
}
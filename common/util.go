package common

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func FileExist(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return false
	}
	return true
}

func DirectoryExist(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil || !info.IsDir() {
		return false
	}
	return true
}

func PathExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err != nil
}

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
// From https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree.
// Source directory must exist.
// Destination directory would be created if not exists.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	if !DirectoryExist(src) {
		return &ErrDirectoryNotExists{Path: src}
	}

	if !DirectoryExist(dst) {
		si, err := os.Stat(src)
		if err != nil {
			return err
		}
		err = os.MkdirAll(dst, si.Mode())
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

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ChExt return a new path which replace the extension of srcPath with newExt.
// Note that newExt should contain ".".
func ChExt(srcPath string, newExt string) string {
	return strings.TrimSuffix(srcPath, filepath.Ext(srcPath)) + newExt
}
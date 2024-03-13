package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
)

type Args struct {
	Src string            `short:"s" long:"src" description:"Source file path" required:"true"`
	Dst map[string]string `short:"d" long:"dst" description:"Destination file path, format: ext1:dir1" required:"true"`
}

func checkArgs(args *Args) error {
	srcFile, err := os.Lstat(args.Src)
	if err != nil {
		return fmt.Errorf("lstat src file[%s] fail: %w", args.Src, err)
	}
	if !srcFile.IsDir() {
		return fmt.Errorf("src file[%s] is not a directory", args.Src)
	}
	for ext, dst := range args.Dst {
		log.Printf("ext: %s, dst: %s", ext, dst)
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	args := &Args{}
	_, err := flags.Parse(args)
	if err != nil {
		log.Fatalf("parse flags fail: %v", err)
	}
	if err := checkArgs(args); err != nil {
		log.Fatalf("check args fail: %v", err)
	}

	dstMap := make(map[string]string, len(args.Dst))
	for ext, dst := range args.Dst {
		dstMap[strings.ToLower(ext)] = dst
	}

	err = filepath.WalkDir(args.Src, func(path string, d os.DirEntry, err error) error {
		ext := filepath.Ext(path)
		if ext == "" {
			return nil
		}
		if strings.HasPrefix(ext, ".") {
			ext = ext[1:]
		}
		dst, ok := dstMap[strings.ToLower(ext)]
		if !ok {
			return nil
		}
		srcRel, err := filepath.Rel(args.Src, path)
		if err != nil {
			return fmt.Errorf("rel path fail: %w", err)
		}
		dstPath := filepath.Join(dst, srcRel)
		dstFile, err := os.Lstat(dstPath)
		if err == nil {
			// check file size and mod time
			srcFile, err := os.Lstat(path)
			if err != nil {
				return fmt.Errorf("lstat src file[%s] fail: %w", path, err)
			}
			if srcFile.Size() != dstFile.Size() || srcFile.ModTime() != dstFile.ModTime() {
				// copy file
				log.Printf("move file(recover): %s -> %s", path, dstPath)
			}
		} else {
			if os.IsNotExist(err) {
				// ok
				log.Printf("move file(rename): %s -> %s", path, dstPath)
			} else {
				return fmt.Errorf("lstat dst file[%s] fail: %w", dstPath, err)
			}
		}
		if err := os.Rename(path, dstPath); err != nil {
			if os.IsNotExist(err) {
				// parent not exist
				if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
					return fmt.Errorf("mkdir fail: %w", err)
				}
				// retry
				if err := os.Rename(path, dstPath); err != nil {
					return fmt.Errorf("move file[%s] to [%s] fail: %w", path, dstPath, err)
				}
			} else {
				return fmt.Errorf("move file[%s] to [%s] fail: %w", path, dstPath, err)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("walk dir fail: %v", err)
	}

}

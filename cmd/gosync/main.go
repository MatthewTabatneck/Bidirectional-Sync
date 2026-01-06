package main

import (
	"bidirectional-sync/internal/fs"
	"fmt"
	"log/slog"
)

func main() {
	root := "C:/Users/mmtab/Desktop/Linux-Laptop"

	fileHashes, failures, err := fs.ParseDirectory(root)
	if err != nil {
		slog.Error("directory scan failed", "root", root, "err", err)
		return
	}

	slog.Info("scan complete",
		"files_ok", len(fileHashes),
		"files_failed", len(failures),
	)

	for _, f := range failures {
		slog.Warn("file skipped",
			"path", f.Path,
			"stage", f.Stage,
			"err", f.Err,
		)
	}

	fmt.Println(fileHashes)
}

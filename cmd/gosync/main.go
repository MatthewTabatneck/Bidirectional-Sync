package main

import (
	// "bidirectional-sync/internal/fs"
	// "fmt"
	// "log/slog"
	"bidirectional-sync/db"
	"log"
	"time"
)

func main() {
	// root := "C:/Users/mmtab/Desktop/Linux-Laptop"

	// fileHashes, failures, err := fs.ParseDirectory(root)
	// if err != nil {
	// 	slog.Error("directory scan failed", "root", root, "err", err)
	// 	return
	// }

	// slog.Info("scan complete",
	// 	"files_ok", len(fileHashes),
	// 	"files_failed", len(failures),
	// )

	// for _, f := range failures {
	// 	slog.Warn("file skipped",
	// 		"path", f.Path,
	// 		"stage", f.Stage,
	// 		"err", f.Err,
	// 	)
	// }

	// fmt.Println(fileHashes)

	// 1. Initialize the DB
	store, err := db.NewStore("./gophersync.db")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	// 2. Example: Saving a file manually
	err = store.UpsertFile("tst.txt", "abc123hash", 1024, time.Now())
	if err != nil {
		log.Printf("Failed to save to DB: %v", err)
	}
}

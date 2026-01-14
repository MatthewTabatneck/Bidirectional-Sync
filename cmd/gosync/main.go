package main

import (
	"bidirectional-sync/internal/db"
	"bidirectional-sync/internal/fs"
	"fmt"
	"log"
)

func main() {
	// 1. Init DB
	store, _ := db.NewStore("./gophersync.db")
	defer store.Close()

	// 2. Run the scan
	fmt.Println("Scanning and hashing files...")
	failures, err := fs.ParseDirectory("C:/Users/mmtab/Desktop/Linux-Laptop", store)

	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}

	// 3. Report failures
	for _, f := range failures {
		fmt.Printf("Failed at stage [%s] for file %s: %v\n", f.Stage, f.Path, f.Err)
	}

	// 1. Initialize the DB
	// store, err := db.NewStore("./gophersync.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer store.Close()

	// // 2. Example: Saving a file manually
	// err = store.UpsertFile("tst.txt", "abc123hash", 1024, time.Now())
	// if err != nil {
	// 	log.Printf("Failed to save to DB: %v", err)
	// }
}

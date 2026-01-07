package fs

import (
	"bidirectional-sync/db"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Store of file information needed to compare with localdb and aws
type FileData struct {
	Path    string
	Hash    string
	Size    int64
	ModTime time.Time
}

// Store of failures that will be reran on increasing intervals and logged to fix issues
type FailureStage string

const (
	StageWalk FailureStage = "walk"
	StageStat FailureStage = "stat"
	StageHash FailureStage = "hash"
)

type FileFailure struct {
	Path  string
	Stage FailureStage
	Err   error
}

// ParseDirectory func will take chosen root directory, parse through each file and do the following
// - check and store file size
// - check and store last modified time
// - call and store hashFile
// All while storing any error the func comes across and moving onto the next file
func ParseDirectory(root string, store *db.Store) ([]FileFailure, error) {
	var failures []FileFailure

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			failures = append(failures, FileFailure{
				Path:  path,
				Stage: StageWalk,
				Err:   walkErr,
			})
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Base(path) == "gophersync.db" {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			failures = append(failures, FileFailure{
				Path:  path,
				Stage: StageStat,
				Err:   err,
			})
			return nil
		}

		hash, err := hashFile(path)
		if err != nil {
			failures = append(failures, FileFailure{
				Path:  path,
				Stage: StageHash,
				Err:   err,
			})
			return nil
		}

		// --- NEW: UPSERT TO DATABASE ---
		// Instead of results = append(results, ...), we save directly.
		err = store.UpsertFile(path, hash, info.Size(), info.ModTime())
		if err != nil {
			// We treat a DB failure as a major issue, but you could also
			// log it as a FailureStage "db" if you want to keep going.
			return fmt.Errorf("failed to save %s to db: %w", path, err)
		}

		return nil
	})

	return failures, err
}

// Function hashFile will take individual file locations and create a sha256 hash of its data
// This is done through io.Copy which chunks each file into 32KB pieces to avoid FILE TOO LARGE errors
func hashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	defer file.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("read file %q for hashing: %w", filePath, err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

package db

import (
	"os"
	"testing"
	"time"
)

func TestNeedsUpdate(t *testing.T) {
	// 1. Setup temporary database
	dbPath := "test_sync.db"
	defer os.Remove(dbPath) // Cleanup after test

	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// 2. Define test data
	testPath := "C:/Users/test/money.txt"
	testSize := int64(1024)
	testTime := time.Now().Truncate(time.Second) // Match SQLite second precision

	// --- TEST 1: NEW FILE (Should be TRUE) ---
	needs, err := store.NeedsUpdate(testPath, testSize, testTime)
	if err != nil || !needs {
		t.Errorf("Test 1 Failed: New file should need update. Got: %v, err: %v", needs, err)
	}

	// 3. Insert the file into DB to test existing records
	err = store.UpsertFile(testPath, "fakehash123", testSize, testTime)
	if err != nil {
		t.Fatalf("Failed to upsert: %v", err)
	}

	// --- TEST 2: UNCHANGED FILE (Should be FALSE) ---
	needs, err = store.NeedsUpdate(testPath, testSize, testTime)
	if err != nil || needs {
		t.Errorf("Test 2 Failed: Unchanged file should NOT need update. Got: %v", needs)
	}

	// --- TEST 3: MODIFIED SIZE (Should be TRUE) ---
	needs, err = store.NeedsUpdate(testPath, 2048, testTime)
	if err != nil || !needs {
		t.Errorf("Test 3 Failed: Changed size should need update. Got: %v", needs)
	}

	// --- TEST 4: MODIFIED TIME (Should be TRUE) ---
	newTime := testTime.Add(time.Hour)
	needs, err = store.NeedsUpdate(testPath, testSize, newTime)
	if err != nil || !needs {
		t.Errorf("Test 4 Failed: Changed time should need update. Got: %v", needs)
	}
}

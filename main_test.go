package main

import (
	"fmt"
	"os"
	"testing"
)

func forceBackup(force bool, storage_dir string, client string) bool {
	if force {
		newFile, err := os.Create(storage_dir + "/" + client + "/backup")
		if err != nil {
			fmt.Println("Error: ", err)
		}
		newFile.Close()
	}

	return force_manual_backup(storage_dir, client)
}

func TestForceBackup(t *testing.T) {
	storage_dir := "backup-storage"
	client := "travis-ci"
	backup_path := storage_dir + "/" + client
	os.MkdirAll(backup_path, 0700)

	cases := []struct {
		name        string
		force       bool
		storage_dir string
		client      string
		expected    bool
	}{
		{"do_not_force_backup", false, storage_dir, client, false},
		{"force_backup", true, storage_dir, client, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := forceBackup(tc.force, tc.storage_dir, tc.client)
			if result != tc.expected {
				t.Fatalf("expected result %v, but got %v", tc.expected, result)
			}
		})

	}

	os.RemoveAll(storage_dir)
}

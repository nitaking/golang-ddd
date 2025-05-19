package bolt

import (
	"go.etcd.io/bbolt"
	"log"
	"syscall"
	"time"
)

const PERMISSION = 0600 //Only owner can read/write
const DbName = "data/local.db"
const NoteBucketName = "notes"

func NewBboltDB() *bbolt.DB {
	option := bbolt.Options{
		Timeout:         1 * time.Second,
		MmapFlags:       syscall.MAP_PRIVATE,
		InitialMmapSize: 10 * 1024 * 1024, // NOTE: Initial mmap size: 10MB
		Mlock:           false,
		FreelistType:    bbolt.FreelistMapType,
	}

	db, err := bbolt.Open(DbName, PERMISSION, &option)
	if err != nil {
		log.Fatalf("failed to open bbolt: %v", err)
	}
	return db
}

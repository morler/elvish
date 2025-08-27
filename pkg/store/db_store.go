package store

import (
	"fmt"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
	"src.elv.sh/pkg/logutil"
	. "src.elv.sh/pkg/store/storedefs"
)

var logger = logutil.GetLogger("[store] ")
var initDB = map[string](func(*bolt.Tx) error){}

// DBStore is the permanent storage backend for elvish. It is not thread-safe.
// In particular, the store may be closed while another goroutine is still
// accessing the store. To prevent bad things from happening, every time the
// main goroutine spawns a new goroutine to operate on the store, it should call
// wg.Add(1) in the main goroutine before spawning another goroutine, and
// call wg.Done() in the spawned goroutine after the operation is finished.
type DBStore interface {
	Store
	Close() error
}

type dbStore struct {
	db *bolt.DB
	wg sync.WaitGroup // used for registering outstanding operations on the store
}

func dbWithDefaultOptions(dbname string) (*bolt.DB, error) {
	// Configure database options with Windows-friendly settings
	options := &bolt.Options{
		Timeout: 2 * time.Second, // Increased timeout for Windows file locking
	}
	
	// On Windows, try multiple times with exponential backoff for file locks
	var db *bolt.DB
	var err error
	var attempt int
	for attempt = 0; attempt < 3; attempt++ {
		db, err = bolt.Open(dbname, 0644, options)
		if err == nil {
			break
		}
		
		// Wait before retrying, with exponential backoff
		if attempt < 2 {
			waitTime := time.Duration(100*(attempt+1)) * time.Millisecond
			logger.Printf("database open attempt %d failed: %v, retrying after %v", attempt+1, err, waitTime)
			time.Sleep(waitTime)
		}
	}
	
	if err != nil {
		logger.Printf("failed to open database after 3 attempts: %v", err)
	} else if attempt > 0 {
		logger.Printf("database opened successfully on attempt %d", attempt+1)
	}
	
	return db, err
}

// NewStore creates a new Store from the given file.
func NewStore(dbname string) (DBStore, error) {
	db, err := dbWithDefaultOptions(dbname)
	if err != nil {
		return nil, err
	}
	return NewStoreFromDB(db)
}

// NewStoreFromDB creates a new Store from a bolt DB.
func NewStoreFromDB(db *bolt.DB) (DBStore, error) {
	logger.Println("initializing store")
	defer logger.Println("initialized store")
	st := &dbStore{
		db: db,
		wg: sync.WaitGroup{},
	}

	err := db.Update(func(tx *bolt.Tx) error {
		for name, fn := range initDB {
			err := fn(tx)
			if err != nil {
				return fmt.Errorf("failed to %s: %v", name, err)
			}
		}
		return nil
	})
	return st, err
}

// Close waits for all outstanding operations to finish, and closes the
// database.
func (s *dbStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	
	logger.Println("closing database store, waiting for operations to complete")
	s.wg.Wait()
	logger.Println("all operations completed, closing database")
	
	err := s.db.Close()
	if err != nil {
		logger.Printf("error closing database: %v", err)
		return err
	}
	
	logger.Println("database closed successfully")
	return nil
}

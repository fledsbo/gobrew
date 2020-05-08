package storage

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	"github.com/fledsbo/gobrew/fermentation"
	"github.com/fledsbo/gobrew/hwinterface"
)

// Storage holds database state
type Storage struct {
	db                     *bolt.DB
	MonitorController      *hwinterface.MonitorController
	OutletController       *hwinterface.OutletController
	FermentationController *fermentation.Controller
}

const fermTable = "fermentations"
const outletTable = "outlets"

// NewStorage creates a new storage object
func NewStorage() *Storage {
	s := &Storage{}

	var err error

	s.db, err = bolt.Open("gobrew.db", 0600, nil)

	if err != nil {
		panic(err)
	}

	log.Println("Opened database")
	return s
}

// StoreFermentations stores the list of fermentations
func (s *Storage) StoreFermentations() error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(fermTable))
		if err != nil {
			return err
		}
		for _, ferm := range s.FermentationController.Batches {
			encoded, err := json.Marshal(ferm)
			log.Print("Json: " + string(encoded))
			if err != nil {
				return err
			}
			err = b.Put([]byte(ferm.Name), encoded)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// RemoveFermentation removes a specific fermentation from store
func (s *Storage) RemoveFermentation(key string) error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(fermTable))
		if err != nil {
			return err
		}
		return b.Delete([]byte(key))
	})

	return err
}

// LoadFermentations loads the list of outlets
func (s *Storage) LoadFermentations() error {
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(fermTable))
		if b == nil {
			return nil
		}

		b.ForEach(func(k, v []byte) error {
			var fermentation fermentation.Batch
			json.Unmarshal(v, &fermentation)
			s.FermentationController.Batches = append(s.FermentationController.Batches, &fermentation)
			return nil
		})
		return nil
	})

	return err
}

// StoreOutlets stores the list of outlets
func (s *Storage) StoreOutlets() error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(outletTable))
		if err != nil {
			return err
		}
		for _, outlet := range s.OutletController.Outlets {
			encoded, err := json.Marshal(outlet)
			log.Print("Json: " + string(encoded))
			if err != nil {
				return err
			}
			err = b.Put([]byte(outlet.Name), encoded)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// LoadOutlets loads the list of outlets
func (s *Storage) LoadOutlets() error {

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(outletTable))
		if b == nil {
			return nil
		}

		b.ForEach(func(k, v []byte) error {
			var outlet hwinterface.Outlet
			json.Unmarshal(v, &outlet)
			s.OutletController.Outlets[outlet.Name] = outlet
			return nil
		})
		return nil
	})

	return err
}

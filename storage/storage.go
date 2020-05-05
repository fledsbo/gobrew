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
	db                *bolt.DB
	MonitorController *hwinterface.MonitorController
	OutletController  *hwinterface.OutletController
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
func (s *Storage) StoreFermentations(ferms []*fermentation.FermentationController) error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(fermTable))
		if err != nil {
			return err
		}
		for _, ferm := range ferms {
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

// LoadFermentations loads the list of outlets
func (s *Storage) LoadFermentations() ([]*fermentation.FermentationController, error) {
	out := make([]*fermentation.FermentationController, 0, 25)
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(fermTable))
		if b == nil {
			return nil
		}

		b.ForEach(func(k, v []byte) error {
			fermentation := fermentation.NewFermentationController(string(k), s.MonitorController, s.OutletController)
			json.Unmarshal(v, fermentation)
			out = append(out, fermentation)
			return nil
		})
		return nil
	})

	return out, err
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

package apis

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/fledsbo/gobrew/config"
	"github.com/fledsbo/gobrew/fermentation"
)

type Brewfather struct {
	FermentationController *fermentation.Controller
	Config                 *config.Config
	LastUpdate             time.Time

	SumTemp      float64
	CountTemp    int
	SumGravity   float64
	CountGravity int
}

type BrewfatherReport struct {
	Name string   `json:"name"`
	Temp *float64 `json:"temp"`
	//	AuxTemp      *float64 `json:"aux_temp"`
	//	ExtTemp      *float64 `json:"ext_temp"`
	TempUnit    *string  `json:"temp_unit"`
	Gravity     *float64 `json:"gravity"`
	GravityUnit *string  `json:"gravity_unit"`
	//	Pressure     *float64
	//	PressureUnit *string `json:"pressure_unit"`
	//	Ph           *float64
	//	Bpm          *float64
	//	Comment      *string
	//	Beer         *string
}

func (b *Brewfather) update() error {
	celcius := "C"
	gravity := "G"

	if b.Config.BrewfatherStreamURL == nil {
		return nil
	}

	if time.Now().Sub(b.LastUpdate) < 15*time.Minute {
		return nil
	}

	for _, batch := range b.FermentationController.Batches {
		state := b.FermentationController.GetBatchState(batch)
		report := BrewfatherReport{
			Name:        "gobrew-" + batch.Name,
			Temp:        state.Temperature,
			TempUnit:    &celcius,
			Gravity:     state.Gravity,
			GravityUnit: &gravity,
		}
		reportJSON, err := json.Marshal(report)
		if err != nil {
			panic(err)
		}

		log.Println("Sending to " + *b.Config.BrewfatherStreamURL)
		log.Println(string(reportJSON))
		req, err := http.NewRequest("POST", *b.Config.BrewfatherStreamURL, bytes.NewBuffer(reportJSON))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		log.Println("response Status:", resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("response Body:", string(body))
	}
	b.LastUpdate = time.Now()
	return nil
}

func (b *Brewfather) Run() {
	for {
		time.Sleep(time.Minute)
		b.update()
	}
}

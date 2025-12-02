package main

import (
	"log"
	"bytes"
	"net/http"
	"errors"
	"io"
	"fmt"
	"os"
)

type holeState int
type HoleToken string

const (
	Initializing holeState = iota
	Occupied
	Free 
)

type Hole struct {
	HoleState holeState
	HoleToken HoleToken
	OccupyingMole string
}

func GenerateToken() HoleToken {
	return HoleToken("test")
}

func NewHole() *Hole {
	t := GenerateToken()
	return &Hole{HoleState: Initializing, HoleToken: t}
}

func (h *Hole) RegisterToWatcher() error {
	url := os.Getenv("APP1_URL")
	if url == "" {
	    url = "http://localhost:42069" // fallback for local dev
	}
	req, err := http.NewRequest("GET", url + "/hole/add", nil)
	if err != nil {
	    log.Fatal(err)
	}
	req.Header.Set("HOLE-Token", "TESTAROONI")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	b, _ := io.ReadAll(resp.Body)
	fmt.Printf("MESSAGE: " + string(b) + "\n")
	//check return of header to verify hole added
	if bytes.Equal(b, []byte("Hole Added\n")) {
		h.HoleState = Free
		return nil
	} else {
		return errors.New("Hole Not Registered to Watcher")
	}
}

func (h *Hole) fillHole(w http.ResponseWriter, r *http.Request) {
	url := os.Getenv("APP1_URL")
	if url == "" {
	    url = "http://localhost:42071" // fallback for local dev
	}
	h.OccupyingMole = url
	h.HoleState = Occupied
}
func (h *Hole) userWhack(w http.ResponseWriter, r *http.Request) {
	//TODO: Process Whack
	if h.HoleState == Occupied {
		w.Write([]byte("Dead Mole!\n"))
		req, err := http.NewRequest("GET", h.OccupyingMole  + "/die", nil)
		if err != nil {
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		
	} else {
		//send a Miss Message to Watcher
		w.Write([]byte("Nothing Here!\n"))
		url := os.Getenv("APP1_URL")
		if url == "" {
		    url = "http://localhost:42069" // fallback for local dev
		}
		req, err := http.NewRequest("GET", url + "/mole/die", nil)
		if err != nil {
		    log.Fatal(err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
	}

}
func main() {
	h := NewHole()
	err := h.RegisterToWatcher()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/fill", h.fillHole)
	http.HandleFunc("/whack", h.userWhack)
	url := os.Getenv("APP1_URL")
	if url == "" {
	    url = ":42070" // fallback for local dev
	} else {
		url = ":8080"
	}
	log.Fatal(http.ListenAndServe(url, nil))
}

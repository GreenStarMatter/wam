package main

import (
	"os"
	"fmt"
	"io"
	"net/http"
	"log"
	"bytes"
	"errors"
	"strings"
)

type MoleState int
type MoleToken string
type Hole string

const (
	Initializing MoleState = iota
	Tunneling 
	Residing
	Dead
)

type Mole struct {
	MoleState MoleState
	MoleToken MoleToken
	Home Hole
}

func GenerateToken() MoleToken {
	return MoleToken("test")
}

func (h Hole) occupy() (bool, error) {
	req, err := http.NewRequest("GET", string(h) + "/fill", nil)
	if err != nil {
	    log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}

func NewMole() *Mole {
	t := GenerateToken()
	return &Mole{MoleState: Initializing, MoleToken: t}
}
func parseHoles(body []byte) []string {
	holes := strings.Split(string(body), "\n")
	if len(holes)<2 {
		return nil
	}
	return holes[1:]
}
func (m *Mole) searchHoles() error {
	url := os.Getenv("APP1_URL")
	if url == "" {
	    url = "http://localhost:42069" // fallback for local dev
	}
	req, err := http.NewRequest("GET", url + "/hole/check", nil)
	if err != nil {
	    log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	b, _ := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	holes := parseHoles(b)
	var hole Hole
	for _, h := range holes {
		hole = Hole(h)
		success, err := hole.occupy()
		if err != nil {
			return err
		}
		if success {
			m.Home = hole
			return nil
		}
	}
	return nil//m.searchHoles()
}



func (m *Mole) RegisterToWatcher() error {
	url := os.Getenv("APP1_URL")
	if url == "" {
	    url = "http://localhost:42069" // fallback for local dev
	}
	req, err := http.NewRequest("GET", url + "/mole/add", nil)
	if err != nil {
	    log.Fatal(err)
	}
	req.Header.Set("MOLE-Token", "TESTMOLE")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	b, _ := io.ReadAll(resp.Body)
	fmt.Printf("MESSAGE: " + string(b) + "\n")
	//check return of header to verify hole added
	if bytes.Equal(b, []byte("Mole Added\n")) {
		m.MoleState = Tunneling
		return nil
	} else {
		return errors.New("Mole Not Registered to Watcher")
	}
}


func (m *Mole) Die(w http.ResponseWriter, r *http.Request) {
	url := os.Getenv("APP1_URL")
	if url == "" {
	    url = "http://localhost:42069" // fallback for local dev
	}
	req, err := http.NewRequest("GET", url + "/mole/die", nil)
	if err != nil {
	    log.Fatal(err)
	}
	req.Header.Set("MOLE-Token", "TESTMOLE")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error message", http.StatusBadRequest)
	}
	defer resp.Body.Close()
	
	b, _ := io.ReadAll(resp.Body)
	fmt.Printf("MESSAGE: " + string(b) + "\n")
	//check return of header to verify hole added
	os.Exit(0)
	if bytes.Equal(b, []byte("Mole killed\n")) {
		m.MoleState = Dead
		os.Exit(0)
		return
	} else {
		http.Error(w, "Mole Death Not Registered to Watcher", http.StatusBadRequest)
		return
	}
}


func main() {

	m := NewMole()
	var err error
	//Register Handle for Dead
	http.HandleFunc("/die", m.Die)
	for {
		switch m.MoleState {
			case Initializing:
				err = m.RegisterToWatcher()
				if err != nil {
					log.Fatal(err)
				}
			case Tunneling:
				err := m.searchHoles()
				if err != nil {
					log.Fatal(err)
				}
				m.MoleState = Residing
			case Residing:
				url := os.Getenv("APP1_URL")
				if url == "" {
				    url = ":42071" // fallback for local dev
				} else {
					url = ":8080"
				}
				log.Fatal(http.ListenAndServe(url, nil))
			case Dead:
				os.Exit(0)
		}
	}
}

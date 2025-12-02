package main

import (
	"net/http"
	"sync"
	"log"
	"os"
)

var masterHoles []Hole
var masterMoles []Mole
var mu sync.Mutex

type Mole struct {
	name string
	address string
}

type Hole struct {
	address string//TODO: Replace address with token
}

func verifyUniqueHole(h Hole) bool {
	for _, h2 := range masterHoles {
		if h.address == h2.address {
			return false
		}
	}
	return true
}

func verifyUniqueMole(m Mole) bool {
	for _, m2 := range masterMoles {
		if m.address == m2.address {
			return false
		}
	}
	return true
}

func addHole(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("HOLE-Token")
	if token == "" {
		w.Write([]byte("BAD NAME\n"))
		return
	}
	url := os.Getenv("APP2_URL")
	if url == "" {
	    url = "http://localhost:42070" // fallback for local dev
	}
	newHole := Hole{address: url}//TODO: Hardcode hole address :42070
	mu.Lock()
	if !verifyUniqueHole(newHole) {
		w.Write([]byte("Hole token Already Exists\n"))//TODO: Make this a specific return message that tells the sender to kill itself
		mu.Unlock()
		return
	} else {
		masterHoles = append(masterHoles, newHole)
	}
	mu.Unlock()
	w.Write([]byte("Hole Added\n"))
}

func addMole(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("MOLE-Token")
	if token == "" {
		w.Write([]byte("BAD NAME\n"))
		return
	}
	url := os.Getenv("APP3_URL")
	if url == "" {
	    url = "http://localhost:42071" // fallback for local dev
	}
	newMole := Mole{address:url}//TODO: Hardcode mole address :42071
	mu.Lock()
	if !verifyUniqueMole(newMole) {
		w.Write([]byte("Mole token Already Exists\n"))//TODO: Make this a specific return message that tells the sender to kill itself
		mu.Unlock()
		return
	} else {
		masterMoles = append(masterMoles, newMole)
	}
	mu.Unlock()
	w.Write([]byte("Mole Added\n"))
}

func holeCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("*****Holes*****\n"))
	for _, h := range masterHoles {
		w.Write([]byte(h.address + "\n"))
	}
}

func moleCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("*****Moles*****\n"))
	for _, m := range masterMoles {
		w.Write([]byte(m.address + "\n"))
	}
}

func moleKill(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("One less mole!\n"))
}

func main() {
	http.HandleFunc("/hole/add", addHole)
	http.HandleFunc("/hole/check", holeCheck)
	http.HandleFunc("/mole/add", addMole)
	http.HandleFunc("/mole/check", moleCheck)
	http.HandleFunc("/mole/die", moleKill)
	url := os.Getenv("APP2_URL")
	if url == "" {
	    url = ":42069" // fallback for local dev
	} else {
		url = ":8080"
	}

	log.Fatal(http.ListenAndServe(url, nil))
	//user is going to mimic moles and holes sending information to this file
	//If I get moles and holes up and running and dockerized, then this is pretty much done

	//add a user to pipe all state changes to
	//handle all state changes of holes and moles
}

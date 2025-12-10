package main

import (
	"fmt"
)

type HoleState int
type MoleState int
type HoleFactory struct {
	HoleId  int
	HoleSet HoleSet
}

type MoleFactory struct {
	MoleId int
}

const (
	Unoccupied HoleState = iota
	Occupied
)

const (
	HidingAlive MoleState = iota
	ExposedAlive
	Dead
)

type Hole struct {
	ID            int
	State         HoleState
	OccupyingMole *Mole
	ParentHoleSet HoleSet
}

type HoleSet struct {
	Available   map[int]*Hole
	Unavailable map[int]*Hole
}

func (hs *HoleSet) addToMap(m map[int]*Hole, h *Hole) error {
	if _, ok := m[h.ID]; ok {
		return fmt.Errorf("hole %d already exists", h.ID)
	}
	m[h.ID] = h
	return nil
}

func (hs *HoleSet) AddAvailable(h *Hole) error {
	return hs.addToMap(hs.Available, h)
}

func (hs *HoleSet) RemoveAvailable(h *Hole) {
	delete(hs.Available, h.ID)
}

func (hs *HoleSet) AddUnavailable(h *Hole) error {
	return hs.addToMap(hs.Unavailable, h)
}

func (hs *HoleSet) RemoveUnavailable(h *Hole) {
	delete(hs.Unavailable, h.ID)
}

func (f *HoleFactory) NewHole() (*Hole, error) {
	f.HoleId++
	h := &Hole{ID: f.HoleId, State: Unoccupied, ParentHoleSet: f.HoleSet}
	err := f.HoleSet.AddAvailable(h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func NewHoleFactory() *HoleFactory {
	return &HoleFactory{
		HoleSet: HoleSet{
			Available:   make(map[int]*Hole),
			Unavailable: make(map[int]*Hole),
		},
	}
}

type Mole struct {
	ID    int
	State MoleState
}

func (f *MoleFactory) NewMole() *Mole {
	f.MoleId++
	return &Mole{ID: f.MoleId, State: HidingAlive}
}

func (h *Hole) TryOccupy(m *Mole) bool {
	if h.State == Occupied {
		return false
	}

	h.ParentHoleSet.RemoveAvailable(h)
	h.ParentHoleSet.AddUnavailable(h)
	h.OccupyingMole = m
	h.State = Occupied
	return true
}

func (m *Mole) Occupy(h *Hole) bool {
	return h.TryOccupy(m)
}

func (m *Mole) ToggleState() {
	switch m.State {
	case HidingAlive:
		m.State = ExposedAlive
	case ExposedAlive:
		m.State = HidingAlive
	}
}

func (h *Hole) TryWhack() string {
	if h.State == Unoccupied {
		return "whiff, no moles here!"
	}

	if h.OccupyingMole.TryWhack() {
		return "bonked out of existence!"
	}

	return "missed and now its laughing!"
}

func (m *Mole) TryWhack() bool {
	if m.State != ExposedAlive {
		return false
	}
	m.State = Dead
	return true
}


func (m *Mole) GetAvailableHole(hs *HoleSet) *Hole {
	if m.State == Dead {
		return nil
	}
	if len(hs.Available) < 1 {
		return nil
	}
	for _, h := range hs.Available {
		return h
	}
	return nil
}

func (m *Mole) TryOccupy(hs *HoleSet) bool {
	h := m.GetAvailableHole(hs)
	if h == nil {
		return false
	}
	return m.Occupy(h)
}

func main() {

}

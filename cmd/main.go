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
	MoleSet MoleSet
}

const (
	Unoccupied HoleState = iota
	Occupied
)

const (
	TunnelingAlive MoleState = iota
	HidingAlive
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

type MoleSet struct {
	Housed map[int]*Mole
	Unhoused map[int]*Mole
	Dead map[int]*Mole
}


func (ms *MoleSet) addToMap(m map[int]*Mole, mo *Mole) error {
	if _, ok := m[mo.ID]; ok {
		return fmt.Errorf("mole %d already exists", mo.ID)
	}
	m[mo.ID] = mo
	return nil
}

func (ms *MoleSet) AddHoused(m *Mole) error {
	return ms.addToMap(ms.Housed, m)
}

func (ms *MoleSet) RemoveHoused(m *Mole) {
	delete(ms.Housed, m.ID)
}

func (ms *MoleSet) AddUnhoused(m *Mole) error {
	return ms.addToMap(ms.Unhoused, m)
}

func (ms *MoleSet) RemoveUnhoused(m *Mole) {
	delete(ms.Unhoused, m.ID)
}

func (ms *MoleSet) AddDead(m *Mole) error {
	return ms.addToMap(ms.Dead, m)
}

func (ms *MoleSet) RemoveDead(m *Mole) {
	delete(ms.Dead, m.ID)
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

func NewMoleFactory() *MoleFactory {
	return &MoleFactory{
		MoleSet: MoleSet{
			Housed:   make(map[int]*Mole),
			Unhoused: make(map[int]*Mole),
			Dead: make(map[int]*Mole),
		},
	}
}

type Mole struct {
	ID    int
	State MoleState
	HoleOccupied *Hole
	ParentMoleSet MoleSet
}

func (f *MoleFactory) NewMole() (*Mole, error) {
	f.MoleId++
	m := &Mole{ID: f.MoleId, State: TunnelingAlive, ParentMoleSet: f.MoleSet}
	err := f.MoleSet.AddUnhoused(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (h *Hole) TryOccupy(m *Mole) bool {
	if h.State == Occupied {
		return false
	}

	h.ParentHoleSet.RemoveAvailable(h)
	h.ParentHoleSet.AddUnavailable(h)
	m.ParentMoleSet.RemoveUnhoused(m)
	m.ParentMoleSet.AddHoused(m)
	h.OccupyingMole = m
	m.HoleOccupied = h
	m.State = HidingAlive
	h.State = Occupied
	return true
}

func (h *Hole) Free() {
	h.ParentHoleSet.AddAvailable(h)
	h.ParentHoleSet.RemoveUnavailable(h)
	m := h.OccupyingMole
	m.ParentMoleSet.AddUnhoused(m)
	m.ParentMoleSet.RemoveHoused(m)
	h.OccupyingMole.HoleOccupied = nil
	h.OccupyingMole.State = TunnelingAlive
	h.OccupyingMole = nil
	h.State = Unoccupied
}

func (m *Mole) Occupy(h *Hole) bool {
	return h.TryOccupy(m)
}

func (m *Mole) Tunnel(hs *HoleSet) {
	if m.HoleOccupied != nil {
		m.HoleOccupied.Free()
	}
	m.State = TunnelingAlive
	m.TryOccupy(hs)
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
	m.ParentMoleSet.RemoveHoused(m)
	m.ParentMoleSet.AddDead(m)
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


type GameState int

const (
	Initializing GameState = iota
	Playing
	End
)

type Game struct {
	HoleFactory *HoleFactory
	MoleFactory *MoleFactory
	State GameState
}

//make holes
func (g *Game) MakeHoles(holes int) {
	for _ = range holes {
		_, _ = g.HoleFactory.NewHole()
	}
}

func (g *Game) MakeMoles(moles int) {
	for _ = range moles {
		_, _ = g.MoleFactory.NewMole()
	}
}


func (g *Game) HouseMoles() {
	for _, m := range g.MoleFactory.MoleSet.Unhoused {
		_ = m.TryOccupy(&g.HoleFactory.HoleSet)
	}
}

func NewGame() *Game {
	hf := &HoleFactory{}
	mf := &MoleFactory{}
	return &Game{HoleFactory: hf, MoleFactory: mf, State: Initializing}
}
func (g *Game) Init(holes int, moles int) {
	g.HoleFactory = NewHoleFactory()
	g.MakeHoles(holes)
	g.MoleFactory = NewMoleFactory()
	g.MakeMoles(moles)
	g.HouseMoles()
}

func (g *Game) CheckWin(moles int) bool {
	if moles == len(g.MoleFactory.MoleSet.Dead) {
		return true
	}

	return false
}

//TODO: Make play logic (User interaction)
//TODO: Make Moles Randomely Move (Make a rand function which fires and chooses moles to move or toggle exposure)
//TODO: Make Player Able to Send Kill Command




func main() {

}

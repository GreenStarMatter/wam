package main

type HoleState int
type MoleState int
type HoleFactory struct{
	HoleId int
}

type MoleFactory struct{
	MoleId int
}

const (
	Unoccupied HoleState = iota
	Occupied
)

const (
	Hiding MoleState = iota
	Exposed
)

type Hole struct{
	ID int
	State HoleState
	OccupyingMole *Mole
}

func (f *HoleFactory) NewHole() *Hole {
	f.HoleId++
	return &Hole{ID: f.HoleId, State:Unoccupied}
}

type Mole struct{
	ID int
	State MoleState
}

func (f *MoleFactory) NewMole() *Mole {
	f.MoleId++
	return &Mole{ID: f.MoleId, State:Hiding}
}

func (h *Hole) TryOccupy(m *Mole) bool {
	if h.State == Occupied {
		return false
	}
	h.OccupyingMole = m
	h.State = Occupied
	return true
}

func (m *Mole) Occupy(h *Hole) bool {
	return h.TryOccupy(m)
}

func (m *Mole) ToggleState() {
	if m.State == Hiding {
		m.State = Exposed
	} else if m.State == Exposed {
		m.State = Hiding
	}
}

func main() {

}

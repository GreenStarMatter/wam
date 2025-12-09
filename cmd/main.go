package main

type HoleState int
type MoleState int
type HoleFactory struct {
	HoleId int
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
}

func (f *HoleFactory) NewHole() *Hole {
	f.HoleId++
	return &Hole{ID: f.HoleId, State: Unoccupied}
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
	h.OccupyingMole = m
	h.State = Occupied
	return true
}

func (m *Mole) Occupy(h *Hole) bool {
	return h.TryOccupy(m)
}

func (m *Mole) ToggleState() {
	if m.State == HidingAlive {
		m.State = ExposedAlive
	} else if m.State == ExposedAlive {
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

func main() {

}

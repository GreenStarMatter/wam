package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestHole(t *testing.T) {
	f := NewHoleFactory()
	h, _ := f.NewHole()
	assert.Equal(t, h, &Hole{ID: 1, State: Unoccupied, ParentHoleSet: f.HoleSet})
	err := f.HoleSet.AddAvailable(h)
	require.Error(t, err)
}

func TestMole(t *testing.T) {
	f := NewMoleFactory()
	m, _ := f.NewMole()
	assert.Equal(t, m, &Mole{ID: 1, State: TunnelingAlive, ParentMoleSet: f.MoleSet})
	//Manually set mole to hiding
	m.State = HidingAlive
	m.ToggleState()
	assert.Equal(t, m, &Mole{ID: 1, State: ExposedAlive, ParentMoleSet: f.MoleSet})
	m.ToggleState()
	assert.Equal(t, m, &Mole{ID: 1, State: HidingAlive, ParentMoleSet: f.MoleSet})
	whacked := m.TryWhack()
	assert.Equal(t, whacked, false)
	m.ToggleState()
	whacked = m.TryWhack()
	assert.Equal(t, whacked, true)
	assert.Equal(t, m, &Mole{ID: 1, State: Dead, ParentMoleSet: f.MoleSet})
	m.ToggleState()
	assert.Equal(t, m, &Mole{ID: 1, State: Dead, ParentMoleSet: f.MoleSet})
}

func TestMoleOccupy(t *testing.T) {
	hf := NewHoleFactory()
	mf := NewMoleFactory()
	h, _ := hf.NewHole()
	assert.Equal(t, h, &Hole{ID: 1, State: Unoccupied, ParentHoleSet: hf.HoleSet})
	m, _ := mf.NewMole()
	assert.Equal(t, m, &Mole{ID: 1, State: TunnelingAlive, ParentMoleSet: mf.MoleSet})
	occupied := m.TryOccupy(&hf.HoleSet)
	assert.Equal(t, occupied, true)
	assert.Equal(t, m, &Mole{ID: 1, State: HidingAlive, HoleOccupied: h, ParentMoleSet: mf.MoleSet})
	assert.Equal(t, h.OccupyingMole, m)
	assert.Equal(t, h.State, Occupied)
	h.Free()
	//Broken state, should never have a free called without wrapper
	assert.Equal(t, h, &Hole{ID: 1, State: Unoccupied, ParentHoleSet: hf.HoleSet})
	assert.Equal(t, m, &Mole{ID: 1, State: TunnelingAlive, ParentMoleSet: mf.MoleSet})
}

func TestMoleTunnel(t *testing.T) {
	hf := NewHoleFactory()
	mf := NewMoleFactory()
	h, _ := hf.NewHole()
	m, _ := mf.NewMole()
	m2, _ := mf.NewMole()
	occupied := m.TryOccupy(&hf.HoleSet)
	assert.Equal(t, occupied, true)
	assert.Equal(t, h.OccupyingMole, m)
	assert.Equal(t, h.State, Occupied)
	occupied2 := m2.TryOccupy(&hf.HoleSet)
	assert.Equal(t, occupied2, false)
	assert.Equal(t, m, &Mole{ID: 1, State: HidingAlive, HoleOccupied: h, ParentMoleSet: mf.MoleSet})
	assert.Equal(t, m2, &Mole{ID: 2, State: TunnelingAlive, HoleOccupied: nil, ParentMoleSet: mf.MoleSet})

	h2, _ := hf.NewHole()
	m2.Tunnel(&hf.HoleSet)
	assert.Equal(t, m2, &Mole{ID: 2, State: HidingAlive, HoleOccupied: h2, ParentMoleSet: mf.MoleSet})
}

func TestGameSetup(t *testing.T) {
	var buf bytes.Buffer
	g := NewGame(&buf)
	g.Init(5, 3)
	assert.Equal(t, 2, len(g.HoleFactory.HoleSet.Available))
	assert.Equal(t, 3, len(g.HoleFactory.HoleSet.Unavailable))

	g = NewGame(&buf)
	g.Init(5, 7)
	assert.Equal(t, 0, len(g.HoleFactory.HoleSet.Available))
	assert.Equal(t, 5, len(g.HoleFactory.HoleSet.Unavailable))

	g = NewGame(&buf)
	g.Init(0, 7)
	assert.Equal(t, 0, len(g.HoleFactory.HoleSet.Available))
	assert.Equal(t, 0, len(g.HoleFactory.HoleSet.Unavailable))
}

func TestGameWin(t *testing.T) {
	var buf bytes.Buffer
	g := NewGame(&buf)
	g.Init(3, 3)

	assert.Equal(t, 0, len(g.MoleFactory.MoleSet.Dead))

	hs := g.MoleFactory.MoleSet.Housed
	hs[1].ToggleState()
	hs[2].ToggleState()
	hs[3].ToggleState()

	_ = hs[1].TryWhack()
	_ = hs[2].TryWhack()
	_ = hs[3].TryWhack()

	assert.Equal(t, 3, len(g.MoleFactory.MoleSet.Dead))
	assert.True(t, g.CheckWin(3))
}

func TestGameInputHandling(t *testing.T) {
	var buf bytes.Buffer
	g := NewGame(&buf)
	g.Init(3, 3)
	commands := make(chan string)
	input := strings.NewReader("moles\nholes\nhelp\nwhack 2\nquit\n")
	scanner := g.InitForPlayer(input)
	go g.ReadCommands(scanner, commands)
	for g.State != End {
		g.ProcessPlayerInput(<-commands)
	}
	assert.Contains(t, buf.String(), "hole: 1")
	assert.Contains(t, buf.String(), "Dead: 0")
	assert.Contains(t, buf.String(), "HELP HELP")
	assert.Contains(t, buf.String(), "SHLONK!")

}

package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHole(t *testing.T) {
	f := &HoleFactory{}
	h := f.NewHole()
	assert.Equal(t, h, &Hole{ID: 1, State: Unoccupied})
}

func TestMole(t *testing.T) {
	f := &MoleFactory{}
	m := f.NewMole()
	assert.Equal(t, m, &Mole{ID: 1, State: HidingAlive})

	m.ToggleState()
	assert.Equal(t, m, &Mole{ID: 1, State: ExposedAlive})
	m.ToggleState()
	assert.Equal(t, m, &Mole{ID: 1, State: HidingAlive})
	whacked := m.TryWhack()
	assert.Equal(t, whacked, false)
	m.ToggleState()
	whacked = m.TryWhack()
	assert.Equal(t, whacked, true)
	assert.Equal(t, m, &Mole{ID: 1, State: Dead})
	m.ToggleState()
	assert.Equal(t, m, &Mole{ID: 1, State: Dead})
}

func TestMoleOccupy(t *testing.T) {
	hf := &HoleFactory{}
	mf := &MoleFactory{}
	h := hf.NewHole()
	assert.Equal(t, h, &Hole{ID: 1, State: Unoccupied})
	m := mf.NewMole()
	assert.Equal(t, m, &Mole{ID: 1, State: HidingAlive})
	occupied := m.Occupy(h)
	assert.Equal(t, occupied, true)
	assert.Equal(t, h.OccupyingMole, m)
	assert.Equal(t, h.State, Occupied)
}

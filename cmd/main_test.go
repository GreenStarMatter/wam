package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestHole(t *testing.T) {
	f := &HoleFactory{}
	h := f.NewHole()
	assert.Equal(t, h, &Hole{ID: 1, State:Unoccupied})
}

func TestMole(t *testing.T) {
	f := &MoleFactory{}
	m := f.NewMole()
	assert.Equal(t, m, &Mole{ID: 1, State:Hiding})

	m.ToggleState()
	assert.Equal(t, m, &Mole{ID: 1, State:Exposed})
	m.ToggleState()
	assert.Equal(t, m, &Mole{ID: 1, State:Hiding})
}


func TestMoleOccupy(t *testing.T) {
	hf := &HoleFactory{}
	mf := &MoleFactory{}
	h := hf.NewHole()
	assert.Equal(t, h, &Hole{ID: 1, State:Unoccupied})
	m := mf.NewMole()
	assert.Equal(t, m, &Mole{ID: 1, State:Hiding})
	occupied := m.Occupy(h)
	assert.Equal(t, occupied, true)
	assert.Equal(t, h.OccupyingMole, m)
	assert.Equal(t, h.State, Occupied)
}

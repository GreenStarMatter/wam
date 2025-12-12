package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const WelcomeMessage = `
|||=======MOLES MOLES MOLES MOLES=======|||\n
Welcome to a wonderful game of moles. It's quite simple:\n
There are holes which can be whacked and there are moles which need to be whacked!\n
Whack all the moles! GO!!!!!\n\n
`

const HelpMessage = `
|||=======HELP HELP HELP HELP HELP=======|||\n
The name of the game is to whack all of the moles:

Commands:
- whack [#]
	Attempt to whack a mole on hole #.  If a mole is there and is exposed then the whack will be successful and the mole will be removed from the game.
- moles
	Survey the moles.  Returns information about how many moles are left.
- holes
	Survey the holes.  Returns information about all the spots that can be whacked.
- quit
	Quits the game.
- help
	You are here.  Type this again and you will be here again.
`

type HoleState int
type MoleState int
type HoleFactory struct {
	HoleId  int
	HoleSet HoleSet
}

type MoleFactory struct {
	MoleId  int
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
	Housed   map[int]*Mole
	Unhoused map[int]*Mole
	Dead     map[int]*Mole
}

func (hs *HoleSet) PrintHolesString() string {
	var b strings.Builder

	for _, ho := range hs.Available {
		fmt.Fprintf(&b, "hole: %d\n", ho.ID)
	}
	for _, ho := range hs.Unavailable {
		fmt.Fprintf(&b, "hole: %d\n", ho.ID)
	}

	return b.String()
}

func (ms *MoleSet) GetMoleStats() string {
	return fmt.Sprintf("Alive: %d\nDead: %d\n", len(ms.Housed)+len(ms.Unhoused), len(ms.Dead))
}

func (hs *HoleSet) GetHole(id int) *Hole {
	if h, ok := hs.Available[id]; ok {
		return h
	}

	if h, ok := hs.Unavailable[id]; ok {
		return h
	}

	return nil
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
			Dead:     make(map[int]*Mole),
		},
	}
}

type Mole struct {
	ID            int
	State         MoleState
	HoleOccupied  *Hole
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

func (h *Hole) TryWhack() (bool, string) {
	if h.State == Unoccupied {
		return false, "whiff, no moles here!\n"
	}

	if h.OccupyingMole.TryWhack() {
		return true, "bonked out of existence!\n"
	}

	return false, "missed and now its laughing!\n"
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
	HoleFactory  *HoleFactory
	MoleFactory  *MoleFactory
	State        GameState
	Output       io.Writer
	WinCondition int
}

// make holes
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

func NewGame(out io.Writer) *Game {
	hf := &HoleFactory{}
	mf := &MoleFactory{}
	return &Game{HoleFactory: hf, MoleFactory: mf, State: Initializing, Output: out}
}
func (g *Game) Init(holes int, moles int) {
	g.WinCondition = moles
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

func (g *Game) InitForPlayer(input io.Reader) *bufio.Scanner {
	fmt.Fprintf(g.Output, WelcomeMessage)
	scanner := bufio.NewScanner(input)
	fmt.Fprintf(g.Output, "> ")
	g.State = Playing
	return scanner
}

func (g *Game) ReadCommands(scanner *bufio.Scanner, commands chan string) {
	for scanner.Scan() {
		commands <- scanner.Text()
	}
	close(commands)
}

func (g *Game) winCheck() {
	if len(g.MoleFactory.MoleSet.Dead) == g.WinCondition {
		fmt.Fprintf(g.Output, "Moles eliminated, YOU WIN!!!!\n")
		g.State = End
	}
}

func (g *Game) handleWhack(hole string) {
	hId, err := strconv.Atoi(hole)
	fmt.Fprintf(g.Output, "SHLONK!\n")
	if err != nil {
		fmt.Fprintf(g.Output, "Hole ID not recognized, where are you aiming?!\n")
		return
	}
	h := g.HoleFactory.HoleSet.GetHole(hId)
	if h == nil {
		fmt.Fprintf(g.Output, "Hole ID not recognized, where are you aiming?!\n")
		return
	}
	hit, resp := h.TryWhack()
	if hit {
		g.winCheck()
	}
	fmt.Fprint(g.Output, resp)
}

func (g *Game) handleMoles() {
	msg := g.MoleFactory.MoleSet.GetMoleStats()
	fmt.Fprint(g.Output, msg)
}
func (g *Game) handleHoles() {
	msg := g.HoleFactory.HoleSet.PrintHolesString()
	fmt.Fprint(g.Output, msg)
}

func (g *Game) handleHelp() {
	fmt.Fprintf(g.Output, HelpMessage)
}

func (g *Game) handleQuit() {
	fmt.Fprintf(g.Output, "GOODBYE QUITTER!\n")
	g.State = End
	//os.Exit(0)
}

func (g *Game) ProcessPlayerInput(commands string) {
	input := commands
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "whack":
		if len(parts) > 1 {
			g.handleWhack(parts[1])
		} else {
			fmt.Fprintf(g.Output, "Hole ID Not Specified\n")
		}
	case "moles":
		g.handleMoles()
	case "holes":
		g.handleHoles()
	case "help":
		g.handleHelp()
	case "quit":
		g.handleQuit()
	default:
		fmt.Fprintf(g.Output, "unknown commands\n")
	}
	fmt.Fprintf(g.Output, "> ")

}

func (g *Game) RunPlayLoop(commands chan string) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	for g.State != End {
		select {
		case <-tick.C:
			g.ProcessMoleMoves(30)
		case cmd, ok := <-commands:
			if !ok {
				return
			}
			g.ProcessPlayerInput(cmd)
		}
	}
}

func (g *Game) ProcessMoleMoves(entropy int) {

	for _, m := range g.MoleFactory.MoleSet.Unhoused {
		m.Tunnel(&g.HoleFactory.HoleSet)
	}

	for _, m := range g.MoleFactory.MoleSet.Housed {
		if rand.Intn(100) < entropy {
			fmt.Fprintf(g.Output, "mole %d vanished!\n", m.ID)
			m.Tunnel(&g.HoleFactory.HoleSet)
		}
		if rand.Intn(100) < entropy {
			m.ToggleState()
			if m.State == HidingAlive {
				fmt.Fprintf(g.Output, "mole %d vanished!\n", m.ID)
			} else {
				fmt.Fprintf(g.Output, "mole %d appeared in hole %d!\n", m.ID, m.HoleOccupied.ID)
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(os.Stdout)
	g.Init(3, 3)
	commands := make(chan string)
	scanner := g.InitForPlayer(os.Stdin)
	go g.ReadCommands(scanner, commands)
	g.RunPlayLoop(commands)
}

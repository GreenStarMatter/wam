# WHACK-A-MOLE


## V1 Spike
======================================================================
This is a goofy project designed to practice creating docker containers, communicating between them, and communicating with them from the command line.  The idea is that the user will spawn n servers (or holes for the moles).  They will then spawn k (some number equal to or less than n) moles.  The moles will simply be a state of the server which is occupied or not.  Each mole will be in a state of exposed or not.  A Watcher will receive all of the state information when a mole moves from one hole to another or if they become exposed or not.  The user will send curl commands to each of the servers to "Whack" the moles.  If a mole is whacked then it will be removed from the game.  The server will return to the watcher the result of the attempt.  Stats will be returned such as the hit percentage and time taken.

The current implementation was a rough draft which intended to test a proof of concept using HTTP to communicate between each of the mole, watcher, and hole.  Each of these microservices were containerized with Docker, use the Docker network to communicate, and are orchestrated with Docker Compose.

The next version will intend to clean-up each of the services.  This will be done by manually iterating through each of the files, identifying each of the end points, and creating a second version of these which is more "idiomatic" Go.


### hole/hole.go
**main logic**
A hole is created and then registered to the watcher.  Handlers are defined and attached to endpoints.  The hole app then listens and waits to be triggered to either issue a fill command from a mole or a whack command from the user.
**type and constants**
- holeState is an int type which helps keep the state of the hole in an iota constant
- HoleToken is a string type which helps define a token which the watcher will pass back and forth
- Hole is a struct type which has a HoleState, HoleToken, and an OccupyingMole
**End Points**
- /fill
	- Calls fillHole method on a Hole
	- Simply passes the url handled by the Dockerfile env variable (or localhost) to the holes attributes
	- TODO: this needs to be updated to either pass the correct environment from the hole (maybe tricky, might need some sort of "address bus"
- /whack
	- Checks the holestate to see if it is occupied
	- If occupied then send a request to the mole (in the h.OccupyingMole attribute) on its /die endpoint
	- If not occupied send a message back referring to the miss
**Helper Functionality**
	- There is a RegisterToWatcher method which registers the hole to the watcher so that it can be added to a master list so the moles know where to look and the state of a hole.
	- There is a NewHole() function which creates a new hole
	- There is a GenerateToken which generates a HoleToken to uniquely identify the hole (TODO: Move this to watcher maybe?  Or maybe start with a generic token which the watcher will randomely assign a new token which is not already a taken name)
### Moles
**main logic**
A new mole is created and a handle function for the die command is attached to the /die method.  A loop is then entered for the different mole states.  When initializing, the mole uses a RegisterToWatcher to place itself on a master list and verify that it is not a duplicate mole.  Then the Tunneling state looks through available holes in search of a m.Home.  The Residing state turns into a server waiting for a request from the attached hole.  The process exits if it moves to the Dead state.
**type and constants**
- MoleState which is an int the is used in a iota constant to track the state of the mole
- MoleToken which is a string that helps disinguish the moles from each other
- Hole which is a string that represents attach the occupied hole to the occupying mole
- Mole which is a struct which has a MoleState, MoleToken, and a Home
**End Points**
- /die
	- This endpoint causes the mole to exit the process
**Helper Functionality**
- RegisterToWatcher is a mole method which adds the mole to the master list
- searchHoles looks through all available holes for a possible home
- parseHoles takes the return from the Watcher and returns an array of holes (if any are available)
- NewMole creates a Mole
- occupy sends a request to the hole to set the Mole as the occupying mole
### Watcher
**main logic**
Quite a few hole and mole funcitonality are tied to endpoints and then a listen and serve method just waits to be called upon.
**type and constants**
A master list of holes is declared
A maser list of moles is declared
A mutex is declared
A mole type is made which contains the moles name and address
A hole type is made which contains the holes address

**End Points**
- /hole/add
	- Adds a hole to the master list if the token is not already on it
- /mole/add
	- Adds a mole to the master list if the token is not already on it
- /hole/check
	- Returns the master list of holes
- /mole/check
	- Returns the master list of moles
- /mole/die
	- Currently just returns a message, but this would be a good place to do a DB write (when it exists)
**Helper Functionality**
- verifyUniqeHole checks if hole is already in master list
- verifyUniqeMole checks if mole is already in master list


## V2 Projected Goals
======================================================================
The first goal is to go back through the original code and see if it can be cleaned up.
The second goal is to add the functionality for the different apps to send actual port information (if this is possible with Docker).  I think I should use docker-compose up --scale hole=3 --scale mole=2.  This command will allow theoretically allow me to create multiple of each of the apps that can be aware of themselves on the Docker network by using os.Hostname().  I should test the scale-up command and the os.Hostname command before doing so.
Lastly I would like to use nginx to communicate from outside the Docker container environment using a Reverse Proxy to be able to locate the apps in a manner similar to DNS.  I will need to move away from my port exposing in this method.

## V2 Realized Goals
======================================================================
I decided to go back to basics and rewrite my V1 spike.  This massively simplified the code as I rewrote it to a monolith.  This definitely doesn't have all the capability of the original but simplified some of the design by getting rid of most of the complexities in Docker, concurrency, HTTP, and 3 separate apps.  The code is much better organized from a testability stand-point, but I think suffers from having too much information folded into one file.  I did still develop with V3 in mind where I added some abstractions which were a bit ahead of this version.  I think I will first break this into separate modules because the game is getting a bit more complex as I develop and I pushing back towards the capabilities I was looking at testing will only add to this complexity.  I am much happier with this design, but I did race at the end and I think my code got a bit more sloppy for the effort.

This version is run by just: **go run ./cmd**

There is some madness in this implementation because I couldn't quickly figure out how to have log lines overwrite.  I wanted information to show in the terminal when moles appeared or vanished so the user would have feedback on what to do but this creates havoc without having a clean UI to work with.  This is something I'll need to figure out for future versions as I still imagine the app having a HUD like display, but I kind of like this chaos right now.  Really makes you root against the moles.

I kept all of the old code because I didn't have any real reason to delete it yet outside of keeping the repo clean, but this is just a toy for now.  I'll probably get rid of the bloat in V3 when I do some of the other cleaning I had planned.

## V3 Projected Goals
======================================================================
My next iteration, I want to break the monolithic main.go into 4 basic packages: main, game, moles, and holes.  The logic is already written so this will mostly just be parsing out the functionality to each of the individual files.  I think I'll also go back and try to flatten out some of the io/game loop logic I added near the end because I just started rushing to finish the commit.  Likely one more day would have it where it needs to be so I'm not too worried.

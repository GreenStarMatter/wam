# WHACK-A-MOLE


This is a CLI whack-a-mole game.

This version is run by just: **go run ./cmd**.  Just start a game and type **whack x** where x is the number of the hole you are aiming for!

Other Commands:
- *moles*: Gives information about moles alive and dead
- *holes*: Tells what holes are available and unavailable (for moles)
- *help*: Prints out possible commands
- *quit*: Quits


Currently game sizes can be adjusted by changing the code within cmd/main.go in the Init function.

The original scope was to make a game that would help practice creating docker containers, communicating between them, and communicating with them from the command line.  Some of this version is included within the commit project history, but the design became very convoluted so I decided to focus on the CLI portion before adding back this complexity.  There is still some work to be done on the CLI portion with information being blasted at the user, but I enjoyed the chaos enough that I kept it in for now.  Really makes you root against the moles.

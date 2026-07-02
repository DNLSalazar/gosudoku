package flagargs

import "flag"

type Args struct {
	Mode         int
	Speed        int
	InitialCells int
}

const MODE_DESCRIPTION = `Mode of game

- 1 for server game
- 2 for CLI game
- 3 for solver`

const SPEED_DESCRIPTION = `Speed for Solver Mode
- 1 For slow 
- 2 For medium
- 3 For fast
`

const CELLS_DESCRIPTION = `Number of static cells`

func GetArgs() Args {
	var mode, speed, cells int

	flag.IntVar(&mode, "m", 0, MODE_DESCRIPTION)
	flag.IntVar(&speed, "s", 2, SPEED_DESCRIPTION)
	flag.IntVar(&cells, "c", 17, CELLS_DESCRIPTION)
	flag.Parse()

	mode -= 1
	if mode < 0 || mode > 2 {
		mode = 0
	}

	if cells < 1 || cells > 80 {
		cells = 17
	}

	return Args{
		Mode:         mode,
		Speed:        speed,
		InitialCells: cells,
	}
}

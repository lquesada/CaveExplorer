package helpers

import "cavernal.com/world"

func NewBasicPlatformRoomSet(o, u, w, h, x func() world.ITile) []*world.Room {
	return []*world.Room{
		&world.Room{
		Floor: [][]func() world.ITile{
			{h, h, h, h, x, h, h, h, h},
			{h, o, o, o, o, o, o, o, h},
			{h, o, o, o, o, o, o, o, h},
			{h, o, o, o, o, o, o, o, h},
			{x, o, o, o, u, o, o, o, x},
			{h, o, o, o, o, o, o, o, h},
			{h, o, o, o, o, o, o, o, h},
			{h, o, o, o, o, o, o, o, h},
			{h, h, h, h, x, h, h, h, h},
		},
		InterestX: 4,
		InterestZ: 4,
	},
	&world.Room{
		Floor: [][]func() world.ITile{
			{h, h, h, x, h, h, h},
			{h, o, o, o, o, o, h},
			{h, o, o, o, o, o, h},
			{x, o, o, u, o, o, x},
			{h, o, o, o, o, o, h},
			{h, o, o, o, o, o, h},
			{h, h, h, x, h, h, h},
		},
		InterestX: 3,
		InterestZ: 3,
	},
	}
}

func NewBasicWallRoomSet(o, u, w, h, x func() world.ITile) []*world.Room {
	return []*world.Room{
		&world.Room{
		Floor: [][]func() world.ITile{
			{w, w, w, w, x, w, w, w, w},
			{w, o, o, o, o, o, o, o, w},
			{w, o, o, o, o, o, o, o, w},
			{w, o, o, o, o, o, o, o, w},
			{x, o, o, o, u, o, o, o, x},
			{w, o, o, o, o, o, o, o, w},
			{w, o, o, o, o, o, o, o, w},
			{w, o, o, o, o, o, o, o, w},
			{w, w, w, w, x, w, w, w, w},
		},
		InterestX: 4,
		InterestZ: 4,
	},
	&world.Room{
		Floor: [][]func() world.ITile{
			{w, w, w, x, w, w, w},
			{w, o, o, o, o, o, w},
			{w, o, o, o, o, o, w},
			{x, o, o, u, o, o, x},
			{w, o, o, o, o, o, w},
			{w, o, o, o, o, o, w},
			{w, w, w, x, w, w, w},
		},
		InterestX: 3,
		InterestZ: 3,
	},
	}
}

func NewBasicStartOrEndRoom(o, u, w, h, x func() world.ITile) *world.Room {
	fx := func() world.ITile {
		t := x()
		t.SetConnect(true)
		t.SetSeeThrough(true)
		t.SetWalkThrough(true)
		t.SetFallThrough(false)
		return t
	}
	return &world.Room{
		Floor: [][]func() world.ITile{
			{fx},
		},
		InterestX: 0,
		InterestZ: 0,
	}
}

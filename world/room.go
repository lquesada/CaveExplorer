package world

type Room struct {
	Floor [][]func() ITile
	InterestX int
	InterestZ int
}

func (r *Room) Width() int {
        if r.Height() == 0 {
                return 0
        }
        return len(r.Floor[0])
}

func (r *Room) Height() int {
        return len(r.Floor)
}

func (r *Room) GetTile(x, z int) ITile {
	if z < 0 || z >= len(r.Floor) {
		return nil
	}
	if x < 0 || x >= len(r.Floor[0]) {
		return nil
	}
	if r.Floor[z][x] == nil {
		return nil
	}
	return r.Floor[z][x]()
}

func (r *Room) CanApply(w *World, x, z int) bool {
	for i := 0; i < r.Width(); i++ {
		for j := 0; j < r.Height(); j++ {
			wI := x+i
			wJ := z+j
			if w.GetTile(wI, wJ) != w.DefaultTile() {
				return false
			}
		}
	}
	return true
}

func (r *Room) Apply(w *World, x, z int, y float32) (min, max *Coords) {
	w.GrowZ(z+r.Height())
	w.GrowZ(z)
	w.GrowX(x+r.Width())
	w.GrowX(x)
	for i := 0; i < r.Width(); i++ {
		for j := 0; j < r.Height(); j++ {
			wI := x+i
			wJ := z+j
			occupied := false
			if w.GetTile(wI, wJ) != w.DefaultTile() {
				occupied = true
			}
			t := r.GetTile(i, j)
			if t != nil {
				w.SetTile(wI, wJ, t)
				if occupied {
					t.SetConnect(false)
				}
				t.SetY(t.Y()+y)
			}
		}
	}
	return &Coords{x, z}, &Coords{x+r.Width()-1, z+r.Height()-1}
}

func (r *Room) FindConnect() []*Coords {
	coords := []*Coords{}
	for x := 0; x < r.Width(); x++ {
		for z := 0; z < r.Height(); z++ {
			if t := r.GetTile(x, z); t != nil && t.Connect() {
				coords = append(coords, &Coords{X: x, Z: z})
			}
		}
	}
	return coords
}

func (r *Room) CanConnect(tileComparator func(ITile, ITile) bool, roomX, roomZ int, w *World, worldX, worldZ int) bool {
  if !w.GetTile(worldX, worldZ).Connect() || !r.GetTile(roomX, roomZ).Connect() {
  	return false
  }
  anyChange := false
  for i := 0; i < r.Width(); i++ {
	for j := 0; j < r.Height(); j++ {
		wI := worldX+i-roomX
		wJ := worldZ+j-roomZ
		if (wI != worldX || wJ != worldZ) && w.GetTile(wI, wJ) != w.DefaultTile() && r.GetTile(i, j) != nil && !tileComparator(w.GetTile(wI, wJ), r.GetTile(i, j)) {	
			return false
		}
		if w.GetTile(wI, wJ) != w.DefaultTile() && r.GetTile(i, j) != nil {
			anyChange = true
		}
	}
  }
  if anyChange {
  	return true
  }
  return false
}

func (r *Room) Connect(connectorTile func() ITile, roomX, roomZ int, w *World, worldX, worldZ int) (min, max *Coords) {
  tile := w.GetTile(worldX, worldZ)
  minCoords, maxCoords := r.Apply(w, worldX-roomX, worldZ-roomZ, tile.Y())
  if connectorTile != nil && (r.Width() > 1 || r.Height() > 1) {
  	connector := connectorTile()
  	connector.SetY(tile.Y())
	tile = w.GetTile(worldX, worldZ)
  	connector.SetCleanable(tile.Cleanable())
  	w.SetTile(worldX, worldZ, connector)
  }
  return minCoords, maxCoords

}

func (r *Room) ValidTiles() []*Coords {
	validTiles := []*Coords{}
	for x := 0; x < r.Width(); x++ {
		for z := 0; z < r.Height(); z++ {
			valid := true
			if x == r.InterestX && z == r.InterestZ {
				continue
			}
			for dx := -1; dx <= 1; dx++ {
				for dz := -1; dz <= 1; dz++ {
					xC := x+dx
					zC := z+dz
					tile := r.GetTile(xC, zC)
					if tile == nil || !tile.WalkThrough() || tile.FallThrough() {
						valid = false
					}
				}
			}
			if valid {
				validTiles = append(validTiles, &Coords{X: x, Z: z})
			}
		}
	}
	return validTiles
}
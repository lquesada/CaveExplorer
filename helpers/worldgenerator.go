package helpers

import "math/rand"
import "cavernal.com/world"
import "cavernal.com/entity"
import "cavernal.com/game"
import "cavernal.com/lib/g3n/engine/math32"

type WorldGenerator struct{
	rooms []*world.Room
	startRoom *world.Room
	endRoom *world.Room
	comparator func(world.ITile, world.ITile) bool
	connectorTile func() world.ITile
	rand *rand.Rand
}

type PlacedRoom struct{
	Room *world.Room
	Coords *world.Coords
}

func NewWorldGenerator(rooms []*world.Room, startRoom *world.Room, endRoom *world.Room, comparator func(world.ITile, world.ITile) bool, connectorTile func() world.ITile) *WorldGenerator {
	return &WorldGenerator{
		rooms: rooms,
		startRoom: startRoom,
		endRoom: endRoom,
		comparator: comparator,
		connectorTile: connectorTile,
	}
}

func (g *WorldGenerator) Generate(seed int, tileSize, gravity float32, criticalPathRoomCount, extraRoomCount int, clean bool) (w *world.World, start, end *PlacedRoom, criticalPath, extra []*PlacedRoom) {
	g.rand = rand.New(rand.NewSource(int64(seed)))
	w = world.Empty(tileSize, gravity)
	criticalPath = []*PlacedRoom{}
	extra = []*PlacedRoom{}
	minCoords, maxCoords := g.startRoom.Apply(w, -g.startRoom.InterestX, -g.startRoom.InterestZ, 0)
	start = &PlacedRoom{g.startRoom, minCoords}

	for i := 0; i < criticalPathRoomCount; i++ {
		roomPicker := g.rand.Perm(len(g.rooms))
		for rI := 0; rI < len(roomPicker); rI++ {
			room := g.rooms[roomPicker[rI]]
			if ok, minC, maxC := g.addConnectedRoom(w, minCoords, maxCoords, room); ok {
				minCoords, maxCoords = minC, maxC
				criticalPath = append(criticalPath, &PlacedRoom{room, minC})
				break
			}
		}
	}
	
	if ok, minC, _ := g.addConnectedRoom(w, minCoords, maxCoords, g.endRoom); ok {
		minCoords = minC
	} else if ok, minC, _ := g.addConnectedRoom(w, w.MinCoords(), w.MaxCoords(), g.endRoom); ok {
		minCoords = minC
	} else {
		wConnect := w.FindConnect(minCoords, maxCoords)
		rConnect := g.endRoom.FindConnect()
		if len(wConnect) > 0 && len(rConnect) > 0 {
			minCoords, _ = g.endRoom.Connect(g.connectorTile, rConnect[0].X, rConnect[0].Z, w, wConnect[0].X, wConnect[0].Z)
		}
	}
	end = &PlacedRoom{g.endRoom, minCoords}
	
	for i := 0; i < extraRoomCount; i++ {
		roomPicker := g.rand.Perm(len(g.rooms))
		for rI := 0; rI < len(roomPicker); rI++ {
			room := g.rooms[roomPicker[rI]]
			if ok, minC, _ := g.addConnectedRoom(w, w.MinCoords(), w.MaxCoords(), room); ok {
				extra = append(extra, &PlacedRoom{room, minC})
				break
			}
		}
	}

	if clean {
		minC := w.MinCoords()
		maxC := w.MaxCoords()

		for x := minC.X; x <= maxC.X; x++ {
			for z := minC.Z; z <= maxC.Z; z++ {
				t := w.GetTile(x, z)
				if t == w.DefaultTile() {
					continue
				}
				if !t.Cleanable() {
					continue
				}
				g.tryRemove(w, x, z)
			}
		}
	}
	return w, start, end, criticalPath, extra
}

func (g *WorldGenerator) tryRemove(w *world.World, x, z int) {
	walkableCount := 0
	dx := []int{-1, 1, 0, 0}
	dz := []int{0, 0, -1, 1}
	for k := 0; k < 4; k++ {
		i, j := dx[k], dz[k]
		wt := w.GetTile(x+i, z+j)
		if wt.WalkThrough() && !wt.FallThrough() {
			walkableCount++
		}
	}
	if walkableCount <= 1 {
		g.cleanRemove(w, x, z)
	}
}

func (g *WorldGenerator) cleanRemove(w *world.World, x, z int) {
	if w.GetTile(x, z) == w.DefaultTile() || !w.GetTile(x, z).Cleanable() {
		return
	}
	w.SetTile(x, z, nil)
	walkableCount := 0
	minC := w.MinCoords()
	maxC := w.MaxCoords()
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Ignore center
			if j == 0 && i == 0 {
				continue
			}
			if walkableCount <= 1 {
				if x+i >= minC.X && x+i <= maxC.X && z+j >= minC.Z && z+j <= maxC.Z {
					g.tryRemove(w, x+i, z+j)
				}
			}
		}
	}

	var replace world.ITile
	dx := []int{-1, 1, 0, 0}
	dz := []int{0, 0, -1, 1}
	for k := 0; k < 4; k++ {
		i, j := dx[k], dz[k]
		wt := w.GetTile(x+i, z+j)
		if wt != w.DefaultTile() && !(wt.WalkThrough() && !wt.FallThrough()) {
			replace = wt
		}
	}
	if replace != nil {
		w.SetTile(x, z, replace)
	}
}

func (g *WorldGenerator) addConnectedRoom(w *world.World, minCoords, maxCoords *world.Coords, room *world.Room) (ok bool, minC, maxC *world.Coords) {
	wConnect := w.FindConnect(minCoords, maxCoords)
	wPicker := g.rand.Perm(len(wConnect))
	for wI := 0; wI < len(wPicker); wI++ {
		wCo := wConnect[wPicker[wI]]
		rConnect := room.FindConnect()
		rPicker := g.rand.Perm(len(rConnect))
		for pI := 0; pI < len(rPicker); pI++ {
			rCo := rConnect[rPicker[pI]]
			if room.CanConnect(g.comparator, rCo.X, rCo.Z, w, wCo.X, wCo.Z) {
				minCoords, maxCoords = room.Connect(g.connectorTile, rCo.X, rCo.Z, w, wCo.X, wCo.Z)
				return true, minCoords, maxCoords
			}
		}
	}
	return false, nil, nil
}

func MatchObstacleComparator(t1, t2 world.ITile) bool {
	return !t1.WalkThrough() == !t2.WalkThrough() || t1.FallThrough() == t2.FallThrough()
}

// Good fit for platforms.
func NonObstacleOverlapComparator(t1, t2 world.ITile) bool {
	return (!t1.WalkThrough() && !t2.WalkThrough()) || (t1.FallThrough() && t2.FallThrough())
}

func MatchSameTypeOrWalkthroughComparator(t1, t2 world.ITile) bool {
	return t1.FallThrough() == t2.FallThrough() &&
		   (!t1.WalkThrough() || !t2.WalkThrough()) &&
		   t1.SeeThrough() == t2.SeeThrough()
}

func MatchSameTypeComparator(t1, t2 world.ITile) bool {
	return t1.FallThrough() == t2.FallThrough() &&
		   t1.WalkThrough() == t2.WalkThrough() &&
		   t1.SeeThrough() == t2.SeeThrough()
}

func MatchAnythingComparator(t1, t2 world.ITile) bool {
	return true
}

type EntitySet struct {
	minAmount	int
	maxAmount	int
	minCount	int
	maxCount	int
	entityGenerator  func() entity.IEntity
}

func NewEntitySetCountable(minAmount, maxAmount, minCount, maxCount int, entityGenerator func() entity.IEntity) *EntitySet {
	s := NewEntitySet(minAmount, maxAmount, entityGenerator)
	s.minCount = minCount
	s.maxCount = maxCount
	return s
}

func NewEntitySet(minAmount, maxAmount int, entityGenerator func() entity.IEntity) *EntitySet {
	return &EntitySet{minAmount: minAmount, maxAmount: maxAmount, entityGenerator: entityGenerator}
}

func (g *WorldGenerator) Populate(s *game.State, rs []*PlacedRoom, entities []*EntitySet) {
	totalSize := 0
	placement := map[*PlacedRoom][]entity.IEntity{}
	validTiles := map[*PlacedRoom]int{}
	for _, r := range rs {
		validTiles[r] = len(r.Room.ValidTiles())
		totalSize += validTiles[r]
	}
	for _, e := range entities {
		amount := e.minAmount+int(math32.Round(g.rand.Float32()*float32(e.maxAmount-e.minAmount)))
		for c := 0; c < amount; c++ {
			picker := g.rand.Float32()*float32(totalSize)
			for _, r := range rs {
				picker -= float32(validTiles[r])
				if picker > 0 {
					continue
				}
				if _, ok := placement[r]; !ok {
					placement[r] = []entity.IEntity{}
				}
				ent := e.entityGenerator()
				enti, isItem := ent.(entity.IItem)
				if isItem {
					count := 1
					if e.minCount != 0 || e.maxCount != 0 {
						count = e.minCount+int(math32.Round(g.rand.Float32()*float32(e.maxCount-e.minCount)))
					}
					if count == 0 {
						break
					}
					enti.SetCount(count)
				}
				placement[r] = append(placement[r], ent)
				break
			}
		}
	}
	for _, r := range rs {
		g.PopulateRoom(s, r, placement[r])
	}
}

func (g *WorldGenerator) PlaceImportantEntityInRoom(s *game.State, r *PlacedRoom, important entity.IEntity) {
	if important != nil {
		tX := r.Coords.X+r.Room.InterestX
		tZ := r.Coords.Z+r.Room.InterestZ
		wC := s.World.Center(tX, tZ)
		tile := s.World.GetTile(tX, tZ)
		important.Position().Set(wC.X, tile.Y(), wC.Z)
		if i, ok := important.(entity.IItem); ok {
			s.Items = append(s.Items, i)
		} else if i, ok := important.(entity.IEnemy); ok {
			s.Enemies = append(s.Enemies, i)
		} else if _, ok := important.(entity.IPlayer); ok {
		} else {
			s.Attrezzo = append(s.Attrezzo, important)
		}
	}
}

func (g *WorldGenerator) PopulateRoom(s *game.State, r *PlacedRoom, other []entity.IEntity) {
	if len(other) == 0 {
		return
	}
	validTiles := r.Room.ValidTiles()
	if len(validTiles) == 0 {
		validTiles = append(validTiles, &world.Coords{r.Room.InterestX, r.Room.InterestZ})
	}
	for _, o := range other {
		if o == nil {
			continue
		}
		pos := validTiles[g.rand.Intn(len(validTiles))]
		tX := r.Coords.X+pos.X
		tZ := r.Coords.Z+pos.Z
		wC := s.World.Center(tX, tZ)
		tile := s.World.GetTile(tX, tZ)
		o.Position().Set(wC.X+(-0.4+g.rand.Float32()*0.8)*s.World.TileSize(), tile.Y(), wC.Z+(-0.2+g.rand.Float32()*0.4)*s.World.TileSize())
		if i, ok := o.(entity.IItem); ok {
			s.Items = append(s.Items, i)
		} else if i, ok := o.(entity.IEnemy); ok {
			s.Enemies = append(s.Enemies, i)
		} else if _, ok := o.(entity.IPlayer); ok {
		} else {
			s.Attrezzo = append(s.Attrezzo, o)
		}
	}
}
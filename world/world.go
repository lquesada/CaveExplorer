package world

import "github.com/lquesada/cavernal/lib/g3n/engine/math32"
import "github.com/lquesada/cavernal/model"

type Coords struct {
	X int
	Z int
}

type World struct {
	floor [][]ITile
	tileSize float32
	zeroX int
	zeroZ int
	defaultTile ITile
	gravity float32

}

func Empty(tileSize, gravity float32) *World {
	defaultTile := NewTile("far void", nil)
    defaultTile.SetWalkThrough(true)
    defaultTile.SetSeeThrough(true)
    defaultTile.SetFallThrough(true)
	return &World{
		floor: [][]ITile{},
		tileSize: tileSize,
		zeroX: 0,
		zeroZ: 0,
		defaultTile: defaultTile,
		gravity: gravity,
	}
}

func New(floor [][]ITile, tileSize float32, zeroX, zeroZ int, gravity float32) *World {
	defaultTile := NewTile("far void", nil)
    defaultTile.SetWalkThrough(true)
    defaultTile.SetSeeThrough(true)
    defaultTile.SetFallThrough(true)
	return &World{
		floor: floor,
		tileSize: tileSize,
		zeroX: zeroX,
		zeroZ: zeroZ,
		defaultTile: defaultTile,
		gravity: gravity,
	}
}

func (w *World) Node() model.INode {
	nodes := []model.INode{}
	for z, vv := range(w.floor) {
		for x, v := range(vv) {
			if v == nil {
				continue
			}
			if n := v.Node(); n != nil {
				nodes = append(nodes, n.Transform(
					&model.Transform{
						Position: &math32.Vector3{float32(x-w.zeroX)*w.TileSize(), v.Y(), float32(z-w.zeroZ)*w.TileSize()},
						Scale: &math32.Vector3{w.TileSize(), 1, w.TileSize()},
						}))
			}
		}
	}
	return model.NewNode(nodes...)
}

func (w *World) GetTile(x, z int) ITile {
	x += w.zeroX
	z += w.zeroZ
	if z < 0 || z >= len(w.floor) {
		return w.defaultTile
	}
	if x < 0 || x >= len(w.floor[z]) {
		return w.defaultTile
	}
	if w.floor[z][x] == nil {
		return w.defaultTile
	}
	return w.floor[z][x]
}

func (w *World) SetTile(x, z int, t ITile) {
	w.GrowZ(z)
	w.GrowX(x)
	x += w.zeroX
	z += w.zeroZ
	if t == w.defaultTile {
		t = nil
	}
	w.floor[z][x] = t
}

func (w *World) DefaultTile() ITile {
	return w.defaultTile
}

func (w *World) Width() int {
        if w.Height() == 0 {
                return 0
        }
        return len(w.floor[0])
}

func (w *World) Height() int {
        return len(w.floor)
}

func (w *World) Gravity() float32 {
    return w.gravity
}

func (w *World) TileSize() float32 {
	return w.tileSize
}

func (w *World) WhereStanding(position *math32.Vector3) (tileX, tileZ int) {
	fx := int(math32.Round(position.X/w.TileSize()))
	fz := int(math32.Round(position.Z/w.TileSize()))
	return fx, fz
}

func (w *World) Center(x, z int) (position *math32.Vector3) {
	return &math32.Vector3{float32(x)*w.TileSize(), 0, float32(z)*w.TileSize()}
}

func (w *World) IsFallingToVoid(position *math32.Vector3, radius float32) (falling bool, fallPoint *math32.Vector3) {
	tileSpan := int(math32.Ceil(radius*2/w.TileSize()))
	
	tileX, tileZ := w.WhereStanding(position)
	tile := w.GetTile(tileX, tileZ)	
	var xCoordSum, zCoordSum float32
	var count int
	if tile.FallThrough() && position.Y <= tile.Y() {
		fitsHole := true
		deltaX, deltaZ, deltaCount := w.DeltaXZForTileCoverage(radius)
		for i := 0; i < deltaCount; i++ {			
			cTileX, cTileZ := w.WhereStanding(&math32.Vector3{position.X + deltaX[i], 0, position.Z + deltaZ[i]})
			cTile := w.GetTile(cTileX, cTileZ)
			if !cTile.FallThrough() {
				fitsHole = false
			}
		}
		if fitsHole {
			return true, position
		}

		for i := 0; i < tileSpan; i++ {
			nextTile: for j := 0; j < tileSpan; j++ {
				for x := 0; x < tileSpan; x++ {
					for z := 0; z < tileSpan; z++ {
						evalX := tileX-i+x
						evalZ := tileZ-j+z
						evalTile := w.GetTile(evalX, evalZ)
						if !evalTile.FallThrough() || position.Y > tile.Y() {
							continue nextTile
						}
					}
				}
				tileCenter := w.Center(tileX-i, tileZ-j)
				xCoord := tileCenter.X
				zCoord := tileCenter.Z
				margin := w.TileSize()/2
				if tileCenter.X > 0 {
					xCoord += margin
				} else if tileCenter.X < 0 {
					xCoord -= margin
				}
				if tileCenter.Z > 0 {
					zCoord += margin
				} else if tileCenter.Z < 0 {
					zCoord -= margin
				}
				xCoordSum += xCoord
				zCoordSum += zCoord
				count++
			}
		}
		if count > 0 {
			centerX := xCoordSum/float32(count)
			centerZ := zCoordSum/float32(count)
			var fallX, fallZ float32
			if centerX < position.X + radius || centerX > position.X - radius{
				fallX = centerX
			}
			if centerZ < position.Z + radius || centerZ > position.Z - radius {
				fallZ = centerZ
			}
			return true, &math32.Vector3{fallX, 0, fallZ}
		}			
	}
	return false, nil
}

func (w *World) EvaluateRay(source, destination *math32.Vector3, ok func(ITile) bool) (success bool, fragmentSuccess float32) {
	angle := math32.Atan2(destination.Z-source.Z, destination.X-source.X)
	incX := math32.Cos(angle)
	incZ := math32.Sin(angle)
	current := math32.NewVector3(source.X, source.Y, source.Z)
	destinationX, destinationZ := w.WhereStanding(destination)
	currentX, currentZ := w.WhereStanding(current)
	for current.DistanceTo(destination) > 0.001 {
		if currentX == destinationX && currentZ == destinationZ {
			if !ok(w.GetTile(currentX, currentZ)) {
				return false, source.DistanceTo(current)
			} else {
				current.Copy(destination)
				break
			}
		}
		nextX := currentX
		if currentX != destinationX {
			if incX > 0 {
				nextX++
			} else if incX < 0 {
				nextX--
			}
		}
		nextZ := currentZ
		if currentZ != destinationZ {
			if incZ > 0 {
				nextZ++
			} else if incZ < 0 {
				nextZ--
			}
		}
		currentCenter := w.Center(currentX, currentZ)
		nextCenter := w.Center(nextX, nextZ)
		cutX := (currentCenter.X+nextCenter.X)/2
		cutZ := (currentCenter.Z+nextCenter.Z)/2
		deltaX := cutX - current.X
		deltaZ := cutZ - current.Z
		var needToCutX, needToCutZ float32
		if math32.Abs(incX) > 0.001 {
			needToCutX = deltaX/incX
		}
		if math32.Abs(incZ) > 0.001 {
			needToCutZ = deltaZ/incZ
		}

		var willProgress float32
		var fixX, fixZ bool
		if nextX != currentX && (needToCutX <= needToCutZ || currentZ == nextZ) {
			willProgress = needToCutX
			fixX = true
		}
		if nextZ != currentZ && (needToCutZ <= needToCutX || currentX == nextX) {
			willProgress = needToCutZ
			fixZ = true
		}

		progressX := willProgress * incX
		progressZ := willProgress * incZ

		current.X += progressX
		if fixX {
			current.X = cutX + 0.0001 * incX/math32.Abs(incX)
		}

		current.Z += progressZ
		if fixZ {
			current.Z = cutZ + 0.0001 * incZ/math32.Abs(incZ)
		}

		currentX, currentZ = w.WhereStanding(current)

		if !ok(w.GetTile(currentX, currentZ)) {
			return false, source.DistanceTo(current)
		}

	}
	return true, source.DistanceTo(current)
}

func (w *World) EvaluateRayRadius(source, destination *math32.Vector3, radius float32, ok func(ITile) bool) (success bool, fragmentSuccess float32) {
	deltaX, deltaZ, deltaCount := w.DeltaXZForBorderCoverage(radius)
	var failed bool
	var fragment float32
	for i := 0; i < deltaCount; i++ {
		dx := deltaX[i]
		dz := deltaZ[i]
		sourceCurrent := &math32.Vector3{source.X + dx, source.Y, source.Z + dz}
		destinationCurrent := &math32.Vector3{destination.X + dx, source.Y, destination.Z + dz}
		
		if suc, frag := w.EvaluateRay(sourceCurrent, destinationCurrent, ok); !suc {
			if !failed {
				failed = true
				fragment = frag
			} else {
				fragment = math32.Min(fragment, frag)
			}
		}
	}
	return !failed, fragment
}

func (w *World) DeltaXZForBorderCoverage(radius float32) (deltaX, deltaZ []float32, deltaCount int) {
	deltaX = []float32{}
	deltaZ = []float32{}

	for r := radius; r >= 0; r -= w.TileSize()/2 {
		deltaX = append(deltaX, r)
		deltaZ = append(deltaZ, -radius)

		deltaX = append(deltaX, r)
		deltaZ = append(deltaZ, radius)

		deltaX = append(deltaX, -r)
		deltaZ = append(deltaZ, -radius)

		deltaX = append(deltaX, -r)
		deltaZ = append(deltaZ, radius)

		deltaZ = append(deltaZ, r)
		deltaX = append(deltaX, -radius)

		deltaZ = append(deltaZ, r)
		deltaX = append(deltaX, radius)

		deltaZ = append(deltaZ, -r)
		deltaX = append(deltaX, -radius)

		deltaZ = append(deltaZ, -r)
		deltaX = append(deltaX, radius)
	}

	return deltaX, deltaZ, len(deltaX)
}

func (w *World) DeltaXZForTileCoverage(radius float32) (deltaX, deltaZ []float32, deltaCount int) {
	deltaX = []float32{}
	deltaZ = []float32{}

	for r1 := radius; r1 >= 0; r1 -= w.TileSize()/2 {
		for r2 := radius; r2 >= 0; r2 -= w.TileSize()/2 {
			deltaX = append(deltaX, r1)
			deltaZ = append(deltaZ, r2)

			deltaX = append(deltaX, -r1)
			deltaZ = append(deltaZ, -r2)

			deltaX = append(deltaX, r1)
			deltaZ = append(deltaZ, -r2)

			deltaX = append(deltaX, -r1)
			deltaZ = append(deltaZ, r2)
		}
	}
	return deltaX, deltaZ, len(deltaX)
}

func (w *World) GrowX(x int) {
  xZ := x+w.zeroX
  if xZ < 0 {
	  incX := -xZ
	  for rI, r := range w.floor {
	  	newLen := len(r)+incX
	  	rN := make([]ITile, newLen, newLen)
	  	for i := 0; i < incX; i++ {
	  		rN[i] = w.defaultTile
	  	}
	  	for i := incX; i < newLen; i++ {
	  		rN[i] = r[i-incX]
	  	}
	  	w.floor[rI] = rN
	  }
	  w.zeroX += incX
  }
  if xZ > w.Width() {
	  incX := xZ-w.Width()
	  for rI, r := range w.floor {
	  	newLen := len(r)+incX
	  	rN := make([]ITile, newLen, newLen)
	  	for i := 0; i < len(r); i++ {
	  		rN[i] = r[i]
	  	}
	  	for i := len(r); i < newLen; i++ {
	  		rN[i] = w.defaultTile
	  	}
	  	w.floor[rI] = rN
	  }
	}
}

func (w *World) GrowZ(z int) {
  zZ := z+w.zeroZ
  width := w.Width()
  if zZ < 0 {
	  incZ := -zZ
	  newLen := len(w.floor)+incZ
	  fN := make([][]ITile, newLen, newLen)
	  	for j := 0; j < incZ; j++ {
	  		fN[j] = make([]ITile, width, width)
	  		for i := 0; i < width; i++ {
	  			fN[j][i] = w.defaultTile
	  		}
	  	}
	  	for j := incZ; j < newLen; j++ {
	  		fN[j] = w.floor[j-incZ]
	  	}
	  	w.floor = fN
	  w.zeroZ += incZ
  }
  if zZ > w.Height() {
	  incZ := zZ-w.Height()
	  newLen := len(w.floor)+incZ
	  fN := make([][]ITile, newLen, newLen)
	  	for j := 0; j < len(w.floor); j++ {
	  		fN[j] = w.floor[j]
	  	}
	  	for j := len(w.floor); j < newLen; j++ {
	  		fN[j] = make([]ITile, width, width)
	  		for i := 0; i < width; i++ {
	  			fN[j][i] = w.defaultTile
	  		}
	  	}
	  	w.floor = fN
	}
}

func (w *World) MinCoords() *Coords {
	return &Coords{-w.zeroX, -w.zeroZ}
}

func (w *World) MaxCoords() *Coords {
	return &Coords{w.Width()-w.zeroX-1, w.Height()-w.zeroZ-1}
}

func (w *World) FindConnect(minCoords, maxCoords *Coords) []*Coords {
	coords := []*Coords{}
	for x := minCoords.X; x <= maxCoords.X; x++ {
		for z := minCoords.Z; z <= maxCoords.Z; z++ {
			if w.GetTile(x, z).Connect() {
				coords = append(coords, &Coords{X: x, Z: z})
			}
		}
	}
	return coords
}
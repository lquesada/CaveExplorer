package scenery

import "fmt"
import "cavernal.com/assets"
import "cavernal.com/world"
import "cavernal.com/model"
import "cavernal.com/lib/g3n/engine/math32"

// --

var concreteColumnModel = &model.NodeSpec{
				Decoder: model.Load(dir, "concretecolumn", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{0, 126, 0},
					Scale: &math32.Vector3{model.MeterScaleFactor, model.MeterScaleFactor*5, model.MeterScaleFactor},
				},
	}

func NewConcreteColumnWith(tile world.ITile) world.ITile {
	nodes := []model.INode{concreteColumnModel.Build()}
	if n := tile.BaseNode(); n != nil {
		nodes = append(nodes, n)
	}
	t := world.NewTile(fmt.Sprintf("concrete column with %s", tile.Name()), model.NewNode(nodes...))
	t.SetWalkThrough(tile.WalkThrough())
	t.SetSeeThrough(tile.SeeThrough())
	t.SetFallThrough(tile.FallThrough())
	return t
}

// --

var concreteTileModel = &model.NodeSpec{
				Decoder: model.Load(dir, "concretetile", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{0, 126, 0},
					Scale: &math32.Vector3{model.MeterScaleFactor, model.MeterScaleFactor, model.MeterScaleFactor},
				},
	}

func NewConcreteTile() world.ITile {
	t := world.NewTile("concrete tile", concreteTileModel.Build())
	t.SetWalkThrough(true)
	t.SetSeeThrough(true)
	return t
}

func NewConcreteTileWith(tile world.ITile) world.ITile {
	nodes := []model.INode{concreteTileModel.Build()}
	if n := tile.BaseNode(); n != nil {
		nodes = append(nodes, n)
	}
	t := world.NewTile(fmt.Sprintf("concrete tile with %s", tile.Name()), model.NewNode(nodes...))
	t.SetWalkThrough(tile.WalkThrough())
	t.SetSeeThrough(tile.SeeThrough())
	t.SetFallThrough(tile.FallThrough())
	return t
}

// --

var concreteWallModel = &model.NodeSpec{
				Decoder: model.Load(dir, "concretewall", assets.Files),
				Transform: &model.Transform{
					Scale: &math32.Vector3{model.MeterScaleFactor, 1, model.MeterScaleFactor},
				},
	}

func NewConcreteWall() world.ITile {
	t := world.NewTile("concrete wall", concreteWallModel.Build())
	return t
}


// --

func NewHole() world.ITile {
	t := world.NewTile("hole", nil)
	t.SetWalkThrough(true)
	t.SetSeeThrough(true)
	t.SetFallThrough(true)
	return t
}

// --

var entranceDoorWithWall = &model.NodeSpec{
				Decoder: model.Load(dir, "entrancedoorwithwall", assets.Files),
				Transform: &model.Transform{
					Scale: &math32.Vector3{model.MeterScaleFactor, 2, model.MeterScaleFactor},
					Rotation: &math32.Vector3{0, math32.Pi, 0},
				},
	}

func NewEntranceDoorWithWall() world.ITile {
	t := world.NewTile("entrance door with wall", entranceDoorWithWall.Build())
	return t
}

// --

var backDoorWithWall = &model.NodeSpec{
				Decoder: model.Load(dir, "doorwithwall", assets.Files),
				Transform: &model.Transform{
					Scale: &math32.Vector3{model.MeterScaleFactor, 2, model.MeterScaleFactor},
					Rotation: &math32.Vector3{0, math32.Pi, 0},
				},
	}

func NewBackDoorWithWall() world.ITile {
	t := world.NewTile("back door with wall", backDoorWithWall.Build())
	return t
}

// --

var nextDoorWithWall = &model.NodeSpec{
				Decoder: model.Load(dir, "doorwithwall", assets.Files),
				Transform: &model.Transform{
					Scale: &math32.Vector3{model.MeterScaleFactor, 2, model.MeterScaleFactor},
				},
	}

func NewNextDoorWithWall() world.ITile {
	t := world.NewTile("next door with wall", nextDoorWithWall.Build())
	return t
}

// --

var caveEntranceDoor = &model.NodeSpec{
				Decoder: model.Load(dir, "caveentrancedoor", assets.Files),
				Transform: &model.Transform{
					Scale: &math32.Vector3{model.MeterScaleFactor, 2, model.MeterScaleFactor},
				},
	}

func NewCaveEntranceDoor() world.ITile {
	t := world.NewTile("cave entrance door", caveEntranceDoor.Build())
	return t
}

// --
package game

import (
	"github.com/lquesada/cavernal/world"
	"github.com/lquesada/cavernal/entity"
	"github.com/lquesada/cavernal/model"

	"github.com/lquesada/cavernal/lib/g3n/engine/math32"
	"math/rand"
	"time"
)

type State struct{
	// Config
	WaitBeforePickupable float32
	FallSpeedUp float32

	Player entity.IPlayer
	Enemies []entity.IEnemy
	Items []entity.IItem
	Attrezzo []entity.IEntity
	Attacks []entity.IAttack
	World *world.World

	rand *rand.Rand
}

func NewState(waitBeforePickupable float32, fallSpeedUp float32, player entity.IPlayer, world *world.World) *State {
	return &State{
		WaitBeforePickupable: 0.3,
		Player: player,
		Enemies: make([]entity.IEnemy, 0, 200),
		Items:  make([]entity.IItem, 0, 200),
		Attacks:  make([]entity.IAttack, 0, 200),
		Attrezzo:  make([]entity.IEntity, 0, 200),
		World: world,
		rand: rand.New(rand.NewSource(time.Now().Unix())),
		}
}

func (s *State) Characters() []entity.ICharacter {
	list := []entity.ICharacter{s.Player}
	for _, v := range s.Enemies {
		list = append(list, v)
	}
	return list
}

func (s *State) PhysicsEntities() []entity.IEntity {
	list := []entity.IEntity{s.Player}
	for _, v := range s.Enemies {
		list = append(list, v)
	}
	for _, v := range s.Items {
		list = append(list, v)
	}
	for _, v := range s.Attacks {
		list = append(list, v)
	}
	return list
}


func (s *State) Entities() []entity.IEntity {
	list := []entity.IEntity{s.Player}
	for _, v := range s.Enemies {
		list = append(list, v)
	}
	for _, v := range s.Items {
		list = append(list, v)
	}
	for _, v := range s.Attacks {
		list = append(list, v)
	}
	for _, v := range s.Attrezzo {
		list = append(list, v)
	}
	return list
}

func (s *State) Over() bool {
	return s.Player.Destroyed()
}

func (s *State) Update(delta float32, inputCallback func(float32)) (pickedUp []entity.IItem, enemiesKilled []entity.IEnemy) {
	var i int
	pickedUp = []entity.IItem{}
	enemiesKilled = []entity.IEnemy{}

	// ----
	// PREPARATION

	// Pre-tick to store needed State and cleanup.
	for _, v := range s.PhysicsEntities() {
		v.PreTick()
	}


	// ----
	// INPUT

	// Process player input.
	if !s.Player.Destroyed() && inputCallback != nil {
		inputCallback(delta)
	}

	// Process player and enemy thoughts.
	for _, v := range s.Characters() {
		v.Think()
	}

	// Process enemies seeing player.
	for _, v1 := range s.Enemies {
	    if v1.TargetFilter() != nil {
			for _, v2 := range s.Characters() {
		 		if v2.Destroyed() || v2.FallingToVoidPosition() != nil {
 					continue
 				}
				if v1.TargetFilter()(v2) && v1.TargetSeeDistance() >= v1.BorderDistanceTo(v2) {
				    success, _ := s.World.EvaluateRay(v1.Position(), v2.Position(), func(t world.ITile) bool { return t.SeeThrough() })
				    if success {
						v1.SeeTarget(v2)
					}
				}
			}
		}
    }


	// ----
	// PICK UP
    // Process colissions between player and equipables
    i = 0
	for i < len(s.Items) {
		v := s.Items[i]
		if v.Destroyed() {
 			i++
 			continue
 		}
		if !v.IsPickupable() {
 			i++
			continue
		}
		if v.BorderDistanceTo(s.Player) < v.Radius() {
			if s.Player.WantToPickUp(v) {
				s.Player.PickUp(v)
				s.Items = append(s.Items[:i], s.Items[i+1:]...)
				i--
				pickedUp = append(pickedUp, v)
			}
		}
		i++
	}

	// ----
	// PRE-PHYSICS
	for _, v := range s.PhysicsEntities() {
		v.SetOnGround(true)

		// Check for fall
		if v.FallingToVoidPosition() == nil {
			falling, fallDirection := s.World.IsFallingToVoid(v.Position(), v.OuterRadius())
			if falling {
				v.SetFallingToVoidPosition(fallDirection)
			}
		}
		// Fall
		if v.FallingToVoidPosition() != nil {
			v.SetOnGround(false)
			fallDirection := &math32.Vector3{}
			fallDirection.Copy(v.FallingToVoidPosition())
			fallDirection.Sub(v.Position())
			// Already fits in the hole
			if fallDirection.X != 0 || fallDirection.Z != 0 {
				dx := fallDirection.X*delta*(2+s.FallSpeedUp)
				constPush := 0.15 * v.OuterRadius()
			 	if math32.Abs(dx) < constPush && math32.Abs(dx) > 0 && fallDirection.X != 0 {
			 		dx *= constPush/math32.Abs(dx)
			 	}
				if math32.Abs(dx) > math32.Abs(fallDirection.X) {
					dx = fallDirection.X
				}
				dz := fallDirection.Z*delta*(2+s.FallSpeedUp)
			 	if math32.Abs(dz) < constPush && math32.Abs(dz) > 0 && fallDirection.Z != 0 {
			 		dz *= constPush/math32.Abs(dz)
			 	}
				if math32.Abs(dz) > math32.Abs(fallDirection.Z) {
					dz = fallDirection.Z
				}
				v.Position().X += dx
				v.Position().Z += dz
	 			v.Speed().X = 0
 				v.Speed().Z = 0
			}
 		}
 		// Do not fall until Radius in.
	 	{
			fitsHole := true
			deltaX, deltaZ, deltaCount := s.World.DeltaXZForTileCoverage(v.OuterRadius())
			for i := 0; i < deltaCount; i++ {			
				cTileX, cTileZ := s.World.WhereStanding(&math32.Vector3{v.Position().X + deltaX[i], 0, v.Position().Z + deltaZ[i]})
				cTile := s.World.GetTile(cTileX, cTileZ)
				if !cTile.FallThrough() {
					fitsHole = false
				}
			}
			if fitsHole {
	 			v.SetOnGround(false)
			}
		}

		// Kill when falling.
		var hurtFall float32 = -100
		var killFall float32 = -150

		if s.Player == v {
			hurtFall = -40
			killFall = -42
		}
 		if v.Position().Y < hurtFall {
 			v.Damaged(1000*delta)
 		}
 		if v.Position().Y < killFall {
 			v.Destroy()
 		} 		

		// On ground?
		if v.FallingToVoidPosition() == nil && !v.Destroyed() {
			deltaX, deltaZ, deltaCount := s.World.DeltaXZForTileCoverage(v.ClimbRadius())
	 		var onGround bool
			for i := 0; i < deltaCount; i++ {			
	 			cTileX, cTileZ := s.World.WhereStanding(&math32.Vector3{v.Position().X + deltaX[i], 0, v.Position().Z + deltaZ[i]})
				cTile := s.World.GetTile(cTileX, cTileZ)
		 		if v.Position().Y <= cTile.Y()+0.01 && v.Speed().Y <= 0 && !cTile.FallThrough() {
					onGround = true
					v.Speed().Y = 0
					v.Position().Y = cTile.Y()
		 		}
		 	}
		 	if !onGround {
		 		v.SetOnGround(false)
		 	}
		 }

 		// Gravity
		if !v.OnGround() {
			v.Gravity(-s.World.Gravity())
		} else {
			v.Gravity(0)
		}
 	}


	// ----
	// TICK

    // Tick to accelerate, move, turn, attack, etc.
	for _, v := range s.PhysicsEntities() {
		v.Tick(delta)
	}


	// ----
	// COLISIONS

	// Between characters
	for _, v1 := range s.Characters() {
 		if v1.Destroyed() {
 			continue
		}
		for _, v2 := range s.Characters() {
	 		if v2.Destroyed() {
	 			continue
 			}
 			if v1 == v2 {
 				continue
 			}
 			// Not using FormerPosition here, it doesn't matter.
			if entity.CheckColisionBetweenFrames(
				v1.Position().X, v1.Position().Y, v1.Position().Z,
				v1.Position().X, v1.Position().Y, v1.Position().Z,
				[]*entity.Cylinder{entity.SimpleToAbsolute(&entity.SimpleCylinder{Radius: v1.OuterRadius(), Height: v1.Height()})},
				v2.Position().X, v2.Position().Y, v2.Position().Z,
				v2.Position().X, v2.Position().Y, v2.Position().Z,
				[]*entity.Cylinder{entity.SimpleToAbsolute(&entity.SimpleCylinder{Radius: v2.OuterRadius(), Height: v2.Height()})},
			) {
				s.collideCharacter(v1, v2, delta)
				s.collideCharacter(v2, v1, delta)
			}
		}
    }

   	// Between characters and attacks
	for _, a := range s.Attacks {
		if a.Destroyed() {
			continue
		}
		for _, v := range s.Characters() {
	 		if v.Destroyed() {
	 			continue
 			}
 			if a.TargetFilter() == nil || !a.TargetFilter()(v) {
 				continue
 			}
 			// Not using FormerPosition here, it doesn't matter.
			if entity.CheckColisionBetweenFrames(
				a.FormerPosition().X, a.FormerPosition().Y, a.FormerPosition().Z,
				a.Position().X, a.Position().Y, a.Position().Z,
				a.Colision(),
				v.FormerPosition().X, v.FormerPosition().Y, v.FormerPosition().Z,
				v.Position().X, v.Position().Y, v.Position().Z,
				[]*entity.Cylinder{entity.SimpleToAbsolute(&entity.SimpleCylinder{Radius: v.OuterRadius(), Height: v.Height()})},
			) {
				a.Hit(v, delta)
			}
		}
    }


    // ----
    // FIX POSITION
	for _, v := range s.PhysicsEntities() {
 		// Avoid falling back to ground.
		falling, _ := s.World.IsFallingToVoid(v.Position(), v.OuterRadius())
 		if v.FallingToVoidPosition() != nil && !falling {
 			v.Position().X, v.Position().Z = v.FormerPosition().X, v.FormerPosition().Z
 		}

	 	// Walls?
 		s.collideWalls(v, delta)

 		// Climb?
		deltaX, deltaZ, deltaCount := s.World.DeltaXZForTileCoverage(v.ClimbRadius())
		_, isItem := v.(entity.IItem)
		tileX, tileZ := s.World.WhereStanding(v.Position())
		tile := s.World.GetTile(tileX, tileZ)
		if v.Position().Y < tile.Y() && !tile.FallThrough() && v.FallingToVoidPosition() == nil {
 			v.Position().Y = tile.Y()
		}
		for i := 0; i < deltaCount; i++ {			
				cTileX, cTileZ := s.World.WhereStanding(&math32.Vector3{v.Position().X + deltaX[i], 0, v.Position().Z + deltaZ[i]})
			cTile := s.World.GetTile(cTileX, cTileZ)
	 		if v.Position().Y < cTile.Y() && (v.Position().Y + v.ClimbReach() >= cTile.Y() || isItem) && !cTile.FallThrough() && v.FallingToVoidPosition() == nil {
	 			v.Position().Y = cTile.Y()
	 		}
	 	}

 	}


	// ----
	// POST-TICK

    // Post-tick to prepare animation, graphics, etc.
	for _, v := range s.PhysicsEntities() {
		v.PostTick(delta)
	}


	// ----
	// ATTACKS

	// Process attacks.
	for _, v := range s.Characters() {
		if !v.AttackIntention() {
			continue
		}
		for _, a := range v.Attack(delta) {
			s.Attacks = append(s.Attacks, a)
		}
	}


	// ----
	// GENERATES

	// Process generation of objects.
	// Generate objects even when marked for destroy, for items.
	for _, v := range s.Entities() {
		for _, p := range v.Generate() {
			// Items generated falling to the bottom are generated falling to the bottom
			p.SetFallingToVoidPosition(v.FallingToVoidPosition())
			if x, ok := p.(entity.IEnemy); ok {
				s.Enemies = append(s.Enemies, x)
			} else if x, ok := p.(entity.IItem); ok {
				s.Items = append(s.Items, x)
				x.WaitBeforePickupable(s.WaitBeforePickupable)
			} else if x, ok := p.(entity.IAttack); ok {
				s.Attacks = append(s.Attacks, x)
			}
		}
	}


	// ----
	// DESTROY

	// Remove destroyed Enemies
	i = 0 
	for i < len(s.Enemies) {
		v := s.Enemies[i]
		if v.Destroyed() {
			s.Enemies = append(s.Enemies[:i], s.Enemies[i+1:]...)
			i--
			enemiesKilled = append(enemiesKilled, v)
		}
		i++
	}

	// Remove destroyed Attacks
	i = 0 
	for i < len(s.Attacks) {
		v := s.Attacks[i]
		if v.Destroyed() {
			s.Attacks = append(s.Attacks[:i], s.Attacks[i+1:]...)
			i--
		}
		i++
	}

	// Remove destroyed Items
	i = 0 
	for i < len(s.Items) {
		v := s.Items[i]
		if v.Destroyed() {
			s.Items = append(s.Items[:i], s.Items[i+1:]...)
			i--
		}
		i++
	}


	return pickedUp, enemiesKilled
}

func (s *State) Node() model.INode {
	nodes := []model.INode{}
	for _, v := range s.Entities() {
		if n := v.Node(); n != nil {
			nodes = append(nodes,
				n.Transform(
					&model.Transform{
						Position: v.Position(),
						Rotation: &math32.Vector3{0, v.LookAngle(), 0},
					}),
				 )
		}

		tileX, tileZ := s.World.WhereStanding(v.Position())
 		tile := s.World.GetTile(tileX, tileZ)
 		if !tile.FallThrough() && v.Position().Y >= tile.Y() {
 			if sn := v.ShadowNode(); sn != nil {
				deltaX, deltaZ, deltaCount := s.World.DeltaXZForTileCoverage(v.ClimbRadius())
				var shadowY float32
				yFound := false
				for i := 0; i < deltaCount; i++ {			
					cTileX, cTileZ := s.World.WhereStanding(&math32.Vector3{v.Position().X + deltaX[i], 0, v.Position().Z + deltaZ[i]})
					cTile := s.World.GetTile(cTileX, cTileZ)
			 		if !yFound || cTile.Y() > shadowY {
			 			shadowY = cTile.Y()
			 			yFound = true
			 		}
			 	}
			 	nodes = append(nodes,
					sn.Transform(
						&model.Transform{
							Position: &math32.Vector3{0, shadowY + 0.01, 0},
						}),
				 )
 			}
 		}
	}
	nodes = append(nodes, s.World.Node())
    return model.NewNode(nodes...)
}


func (s *State) findColisionsBetweenFrames(formerSourcePosition, sourcePosition *math32.Vector3, cylinder []*entity.Cylinder, filter func(interface{}) bool) []entity.IEntity {
	hits := make([]entity.IEntity, 0, 20)
	for _, v := range s.Entities() {
		if !filter(v) {
			continue
		}
		if entity.CheckColisionBetweenFrames(
			formerSourcePosition.X, formerSourcePosition.Y, formerSourcePosition.Z, sourcePosition.X, sourcePosition.Y, sourcePosition.Z, cylinder,
			v.FormerPosition().X, v.FormerPosition().Y, v.FormerPosition().Z, v.Position().X, v.Position().Y, v.Position().Z, v.Colision()) {
						hits = append(hits, v)
		}
	}
	return hits
}

func (s *State) collideCharacter(e, c entity.ICharacter, delta float32) {
	dist := e.CenterDistanceTo(c)
	if dist == 0 {
		dist = 0.0001
	}
	var fX, fZ float32
	if c.PushForce() > e.PushForce() {
		fX = (e.Position().X-c.Position().X)*math32.Max(2, c.PushForce()-e.PushForce()+1)/2/dist
		fZ = (e.Position().Z-c.Position().Z)*math32.Max(2, c.PushForce()-e.PushForce()+1)/2/dist
	} else {
		fX = (e.Position().X-c.Position().X)*math32.Max(2, c.PushForce()-e.PushForce()+1)/3/dist
		fZ = (e.Position().Z-c.Position().Z)*math32.Max(2, c.PushForce()-e.PushForce()+1)/3/dist
	}
	var forcePushRandom float32 = 0.01
   	fX += rand.Float32() * forcePushRandom * 2 - forcePushRandom
   	fZ += rand.Float32() * forcePushRandom * 2 - forcePushRandom
   	if fX > e.MaxSpeed() {
   		fX = e.MaxSpeed()
   	} else if fX < -e.MaxSpeed() {
   		fX = -e.MaxSpeed()
   	}
   	if fZ > e.MaxSpeed() {
   		fZ = e.MaxSpeed()
   	} else if fZ < -e.MaxSpeed() {
   		fZ = -e.MaxSpeed()
   	}
	e.Position().X += fX*delta
	e.Position().Z += fZ*delta
}

func (s *State) collideWalls(v entity.IEntity, delta float32) {
	distance := v.Position().DistanceTo(v.FormerPosition())
	canWalk := func(t world.ITile) bool {
    		return t.WalkThrough() && t.Y() <= v.Position().Y + v.ClimbReach()
    }
    deltaMove := &math32.Vector3{}
	deltaMove.Add(v.Position())
	deltaMove.Sub(v.FormerPosition())
	deltaMove.Y = 0

    success, successFragment := s.World.EvaluateRayRadius(v.FormerPosition(), v.Position(), v.OuterRadius(), canWalk)
    if !success {
    	// Push back
    	back := &math32.Vector3{}
    	back.Copy(deltaMove)
    	back.MultiplyScalar(1- successFragment/distance)
    	v.Position().Sub(back)
    	if v.FallingToVoidPosition() != nil {
    		return
    	}
    }

	// There is a chance that it's still stuck in the wall because of delta + float precision
	var f float32 = math32.Min(s.World.TileSize()*4, math32.Max(v.MaxSpeed()/2, s.World.TileSize()/4))*delta
	var checkRadius float32 = v.OuterRadius()+s.World.TileSize()/50
	deltaX, deltaZ, deltaCount := s.World.DeltaXZForBorderCoverage(checkRadius)
	deltaX = append(deltaX, -checkRadius, +checkRadius, 0, 0)
	deltaZ = append(deltaZ, 0, 0, -checkRadius, +checkRadius)
	deltaCount += 4
	var fixX, fixZ float32
	for i := 0; i < deltaCount; i++ {
		dx := deltaX[i]
		dz := deltaZ[i]
        if (dx > 0) == (deltaMove.X > 0) && (dz > 0) == (deltaMove.Z > 0) &&
          (dx < 0) == (deltaMove.X < 0) && (dz < 0) == (deltaMove.Z < 0) {
          continue
        }
			tileX, tileZ := s.World.WhereStanding(&math32.Vector3{v.Position().X + dx, 0, v.Position().Z + dz})
		tile := s.World.GetTile(tileX, tileZ)
		if !canWalk(tile) {
			// Stuck!
			if dx < 0 {
				fixX--
			} else if dx > 0 {
				fixX++
			}
			if dz < 0 {
				fixZ--
			} else if dz > 0 {
				fixZ++
			}
		}
	}
	if math32.Abs(fixX) > math32.Abs(fixZ) && fixX != 0 {
		//v.Speed().X = 0
		if fixX < 0 {
			v.Position().X += f
		} else {
			v.Position().X -= f
		}
	} else if math32.Abs(fixZ) > math32.Abs(fixX) && fixZ != 0 {
		if fixZ < 0 {
			v.Position().Z += f
		} else {
			v.Position().Z -= f
		}
	} else if fixX != 0 && fixZ != 0 {
		if fixX < 0 {
			v.Position().X += f
		} else {
			v.Position().X -= f
		}
		if fixZ < 0 {
			v.Position().Z += f
		} else {
			v.Position().Z -= f
		}
	}
	// Stuck front towards a corner.
	// Try and move close to 45 degrees front-left and front-right, if success, do it.
	if !success && v.OnGround() {
		remaining := distance-successFragment

		angle := math32.Atan2(deltaMove.Z, deltaMove.X)

		rayLength := s.World.TileSize()/5
		leftAngle := angle-math32.Pi/3.95
		lRay := &math32.Vector3{v.Position().X + math32.Cos(leftAngle)*rayLength, 0, v.Position().Z + math32.Sin(leftAngle)*rayLength}
		successLeft, _ := s.World.EvaluateRayRadius(v.Position(), lRay, v.OuterRadius(), canWalk)

		rightAngle := angle+math32.Pi/3.95
		rRay := &math32.Vector3{v.Position().X + math32.Cos(rightAngle)*rayLength, 0, v.Position().Z + math32.Sin(rightAngle)*rayLength}
		successRight, _ := s.World.EvaluateRayRadius(v.Position(), rRay, v.OuterRadius(), canWalk)

		lDest := &math32.Vector3{v.Position().X + math32.Cos(leftAngle)*remaining, 0, v.Position().Z + math32.Sin(leftAngle)*remaining}
		rDest := &math32.Vector3{v.Position().X + math32.Cos(rightAngle)*remaining, 0, v.Position().Z + math32.Sin(rightAngle)*remaining}

		if successLeft && successRight {
			if rand.Float32() < 0.5 {
				v.Position().Copy(lDest)
			} else {
				v.Position().Copy(rDest)
			}
		} else if successLeft {
			v.Position().Copy(lDest)
		} else if successRight {
			v.Position().Copy(rDest)
		}
    }

}

package enemies

import 	"cavernal.com/assets/humanoid/customization"
import 	"cavernal.com/assets/humanoid/human"
import 	"cavernal.com/assets/humanoid/skeleton"
import 	"cavernal.com/assets/humanoid/orc"
import 	"cavernal.com/assets/humanoid/zombie"
import 	"cavernal.com/entity/humanoid"
import 	"cavernal.com/entity"
import 	"cavernal.com/model"
import  "cavernal.com/lib/g3n/engine/math32"

type SimpleEnemy struct{ 
	humanoid.Humanoid

        // Config
        attackDistance float32
        targetSeeDistance float32
        targetStopWhenCloserThan float32

        // Runtime
        target  entity.IEntity

}

func (e *SimpleEnemy) Evil() {
	// Marker for enemy
}

func (e *SimpleEnemy) Target() entity.IEntity {
	return e.target
}

func (e *SimpleEnemy) SeeTarget(t entity.IEntity) {
	if t == e.target {
		return
	}
	if e.target == nil || e.target.Destroyed() || e.InnerEntity().BorderDistanceTo(t) < e.InnerEntity().BorderDistanceTo(e.target) {
		e.target = t
	}
}

func (e *SimpleEnemy) TargetSeeDistance() float32 {
	return e.targetSeeDistance
}

func (e *SimpleEnemy) Think() {
	// Update move intention
	if e.target != nil {
		if e.target.Destroyed() || e.target.FallingToVoidPosition() != nil {
			e.target = nil
		}
	}
	if e.target != nil && e.InnerEntity().BorderDistanceTo(e.target) >= e.targetStopWhenCloserThan {
		e.MoveIntention().Copy(e.target.Position())
		e.MoveIntention().Sub(e.Position())
		e.MoveIntention().Y = 0
		e.MoveIntention().Normalize()
	}

	// Update look intention
	if e.target != nil {
		x := e.target.Position().X-e.Position().X
		z := e.target.Position().Z-e.Position().Z
		e.SetLookIntention(math32.Atan2(x, z))
	} else if e.MoveIntention().Length() != 0 {
		e.SetLookIntention(math32.Atan2(e.MoveIntention().X, e.MoveIntention().Z))
	}

	if e.target != nil && e.InnerEntity().BorderDistanceTo(e.target) <= e.attackDistance {
		e.WantToAttack()
	}

	e.Humanoid.Think()
}


func NewSimpleHumanEnemy(name string) *SimpleEnemy {
	i := humanoid.NewInventory(humanoid.StandardSlots, 40)
	
	h := human.New(name, i, map[entity.DrawSlotId]model.INode{
			humanoid.Face: customization.FaceScaredModel.Build(),
			humanoid.ProtrudeHeadFrontOrHangShort: human.BlackBeardFrontModel.Build(),
			humanoid.ProtrudeHeadFrontOrHangLong: human.BlackBeardLongModel.Build(),
			humanoid.ProtrudeHeadSidesShort: human.BlackBeardSideModel.Build(),
		})
	h.SetTargetFilter(entity.FilterGood)
	h.SetHealth(1)
	h.SetMaxHealth(1)
	
	return &SimpleEnemy{
		Humanoid: *h,
		targetSeeDistance: 6,
		targetStopWhenCloserThan: 0.5,
		attackDistance: 1,
	}
}


func NewSimpleSkeletonEnemy(name string) *SimpleEnemy {
	i := humanoid.NewInventory(humanoid.StandardSlots, 40)

	h := skeleton.New(name, i, nil)
	h.SetTargetFilter(entity.FilterGood)
	h.SetHealth(1)
	h.SetMaxHealth(1)
	
	return &SimpleEnemy{
		Humanoid: *h,
		targetSeeDistance: 6,
		targetStopWhenCloserThan: 0.5,
		attackDistance: 1,
	}
}


func NewSimpleZombieEnemy(name string) *SimpleEnemy {
	i := humanoid.NewInventory(humanoid.StandardSlots, 40)

	h := zombie.New(name, i, map[entity.DrawSlotId]model.INode{
			humanoid.Face: customization.FaceScaredModel.Build(),
		})
	h.SetTargetFilter(entity.FilterGood)
	h.SetHealth(1)
	h.SetMaxHealth(1)
	
	return &SimpleEnemy{
		Humanoid: *h,
		targetSeeDistance: 6,
		targetStopWhenCloserThan: 0.5,
		attackDistance: 1,
	}
}

func NewSimpleOrcEnemy(name string) *SimpleEnemy {
	i := humanoid.NewInventory(humanoid.StandardSlots, 40)

	h := orc.New(name, i, map[entity.DrawSlotId]model.INode{
			humanoid.Face: customization.FaceScaredModel.Build(),
		})
	h.SetTargetFilter(entity.FilterGood)
	h.SetHealth(1)
	h.SetMaxHealth(1)
	
	return &SimpleEnemy{
		Humanoid: *h,
		targetSeeDistance: 6,
		targetStopWhenCloserThan: 0.5,
		attackDistance: 1,
	}
}

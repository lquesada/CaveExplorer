package humanoid

import "cavernal.com/entity"
import "cavernal.com/model"
import "cavernal.com/lib/g3n/engine/math32"

var (

	Walk = "Walk"
	AttackRight = "AttackRight"
	AttackLeft = "AttackLeft"

	LegLeftPosition = math32.NewVector3(13, 44, 0)
	LegRightPosition = math32.NewVector3(-13, 44, 0)
	ArmLeftPosition = math32.NewVector3(17, 75.5, 0)
	ArmRightPosition = math32.NewVector3(-17, 75.5, 0)
	HandLeftPosition = math32.NewVector3(8.5, 0, 30.5)
	HandRightPosition = math32.NewVector3(-8.5, 0, 30.5)
	LegLeftTransformPosition = (&math32.Vector3{}).Copy(LegLeftPosition).MultiplyScalar(model.GlobalScaleFactor)
	LegRightTransformPosition = (&math32.Vector3{}).Copy(LegRightPosition).MultiplyScalar(model.GlobalScaleFactor)
	ArmLeftTransformPosition = (&math32.Vector3{}).Copy(ArmLeftPosition).MultiplyScalar(model.GlobalScaleFactor)
	ArmRightTransformPosition = (&math32.Vector3{}).Copy(ArmRightPosition).MultiplyScalar(model.GlobalScaleFactor)
	HandLeftTransformPosition = (&math32.Vector3{}).Copy(HandLeftPosition).MultiplyScalar(model.GlobalScaleFactor)
	HandRightTransformPosition = (&math32.Vector3{}).Copy(HandRightPosition).MultiplyScalar(model.GlobalScaleFactor)
)

type Humanoid struct{
	entity.Character

	inventory *Inventory

	baseNodes map[entity.DrawSlotId]model.INode
	sequences map[string]*model.Sequence
	transforms map[entity.DrawSlotId]*model.Transform

	currentAttack string

	jumpPower float32
	jumpTimeLimit *model.Sequence
	repeatedJump bool

	attackFinishedCallback func()
}

func NewHumanoid(name string, i *Inventory, baseNodes map[entity.DrawSlotId]model.INode, sequences map[string]*model.Sequence, transforms map[entity.DrawSlotId]*model.Transform) *Humanoid {
	h := &Humanoid{
		Character: *entity.NewCharacter(entity.NewEntity(name)),
		inventory: i,
		baseNodes: baseNodes,
		sequences: sequences,
		transforms: transforms,
		currentAttack: AttackRight,
		jumpTimeLimit: model.NewTimer(0),
	}

	var baseSpeed float32 = 4

	h.SetMaxSpeed(baseSpeed)
	h.SetFriction(baseSpeed * 5)
	h.SetAcceleration(baseSpeed * 30)
	h.SetMinSpeed(baseSpeed * 0.1)
	h.SetRotationSpeed(baseSpeed * 5)
	h.SetMaxHealth(100)
	h.SetHealth(100)
	h.SetDamageable(true)
	h.SetRadius(0.28)
	h.SetHeight(1.8)
	h.SetHandHeight(0.6)
	h.SetWalkAnimationSpeed(22)
	h.SetWalkAnimationWhenStopping(22)
	h.SetAttackAnimationSpeed(37)
	h.SetPushForce(10)
	h.SetPushFactor(1)
	h.SetOuterRadius(0.6)
	h.SetClimbRadius(0.05)
	h.SetClimbReach(0.25)
	h.SetDontPickUpUntilFarDistance(0.15)
	h.SetJumpTimeLimit(0.10)
	h.SetJumpPower(2)
	h.SetJumpMoveFactor(0.1)

	return h
}

func (h *Humanoid) SetJumpTimeLimit(seconds float32) {
	cur := h.jumpTimeLimit.Counter()
	if !h.jumpTimeLimit.Running() {
		cur = 0
	}
	h.jumpTimeLimit = model.NewTimer(seconds)
	h.jumpTimeLimit.Add(cur)
}

func (h *Humanoid) SetJumpPower(v float32) {
	h.jumpPower = v
}

func (h *Humanoid) PreTick() {
	h.Character.PreTick()
	if h.inventory != nil {
		h.inventory.PreTick()
		for _, v := range h.inventory.AllItems() {
			v.Position().Copy(h.Character.InnerEntity().Position())
		}
	}
	h.repeatedJump = false
}

func (h *Humanoid) Tick(delta float32) {
	h.Character.Tick(delta)
	if h.inventory != nil {
		h.inventory.Tick(delta)
	}
}

func (h *Humanoid) PostTick(delta float32) {
	h.Character.PostTick(delta)
	if h.inventory != nil {
		h.inventory.PostTick(delta)
	}

	// Must be in PostTick so we've got the updated formerPosition vs. position
	displacement := entity.Distance2D(h.Character.InnerEntity().FormerPosition(), h.Character.InnerEntity().Position())/delta
	walkSequence := h.sequences[Walk]
	if displacement > 0.01 && h.Character.MoveIntention().Length() > 0 {
		if !walkSequence.Running() || walkSequence.Stopping() {
			walkSequence.Start()
		}
		walkSequence.Add(h.Character.WalkAnimationSpeed()*displacement)
	} else if walkSequence.Running() {
		walkSequence.Stop()
		walkSequence.Add(h.Character.WalkAnimationWhenStopping()*delta)
	}

	if h.sequences[h.currentAttack].Running() {
		h.sequences[h.currentAttack].Add(h.Character.AttackAnimationSpeed()*100*delta)
	} else if h.attackFinishedCallback != nil {
		h.attackFinishedCallback()
	}
	if h.OnGround() {
		h.SetJumping(false)
	}
	h.jumpTimeLimit.Add(delta)
}

func (h *Humanoid) Inventory() entity.IInventory {
	return h.inventory
}

func (h *Humanoid) WantToAttack() bool {
	return h.manageAttack(false)
}

func (h *Humanoid) WantToRepeatAttack() bool {
	return h.manageAttack(true)
}

func (h *Humanoid) manageAttack(repeat bool) bool {
	if !h.sequences[h.currentAttack].Running() {
		attack, ok := h.Inventory().Attack(repeat).(*attackInfo)
		if !ok {
			attack = nil
		}
		if attack != nil {
			if attack.rightHand {
				h.currentAttack = AttackRight
			} else {
				h.currentAttack = AttackLeft
			}
			h.attackFinishedCallback = attack.attackFinishedCallback
			h.sequences[h.currentAttack].Start()
			h.SetAttacking(true)
			return true
		}
	}
	return false
}
func (e *Humanoid) Attack(delta float32) []entity.IAttack {
	// TODO adjust for non melees
	if e.EquippedWeapon() != nil {
		return e.EquippedWeapon().Attack(delta, 1, e.Position(), e.LookIntention(), e.Position().Y, e.HandHeight(), e.Height(), e.TargetFilter())
	}
	// TODO default weapon
	return nil
}

func (e *Humanoid) Jump() {
	if !e.OnGround() {
		return 
	}
	e.jumpTimeLimit.Reset()
	e.jumpTimeLimit.Start()
	//if !e.jumpTimeLimit.Running() {
		//return
	//}
	e.SetJumping(true)
	e.InnerEntity().Speed().Y += e.jumpPower
	e.repeatedJump = true
	e.SetOnGround(false)
}

func (e *Humanoid) RepeatJump() {
	if !e.IsJumping() {
		return
	}
	if !e.jumpTimeLimit.Running() {
	    return
	}
	e.repeatedJump = true
	e.InnerEntity().Speed().Y += e.jumpPower
}

func (h *Humanoid) AttackIntention() bool {
	poll := h.sequences[h.currentAttack].PollActions()
	for _, p := range poll {
		if p.Action == HumanoidAttackAction {
			h.SetAttacking(false)
			return true
		}
	}
	return false
}

func (e *Humanoid) EquippedWeapon() entity.IEquipable {
	return e.inventory.EquippedWeapon()
}

func (h *Humanoid) Node() model.INode {
	covered := h.CoveredByEquipables()
	nodes := []model.INode{}
    for k, n := range h.baseNodes {
		if covered[k] {
			continue
		}
		switch k {
		case ArmLeft:
			n = n.Transform(
					h.transforms[ArmLeft],
					&model.Transform{
						Position: ArmLeftTransformPosition,
					})
		case ArmRight:
			n = n.Transform(
					h.transforms[ArmRight],
					&model.Transform{
						Position: ArmRightTransformPosition,
					})
		case LegLeft:
			n = n.Transform(
					h.transforms[LegLeft],
					&model.Transform{
						Position: LegLeftTransformPosition,
					})
		case LegRight:
			n = n.Transform(
					h.transforms[LegRight],
					&model.Transform{
						Position: LegRightTransformPosition,
					})
		}
		if n != nil {
        	nodes = append(nodes, n)
        }
    }
    for slot, e := range h.inventory.EquipmentSlots() {
    	if e == nil {
    		continue
    	}
    	if e.IsCountable() && e.Count() == 0 {
    		continue
    	}
		if slot == AlternateLeftSlot ||  slot == AlternateRightSlot {
			// TODO draw alternates on the back
			continue
		}

    	n := h.equipableNode(slot, e, covered)
    	if n != nil {
	        nodes = append(nodes, n)
    	}
	}    	
    for slot, e := range h.Character.StyleEquipables() {
    	if e == nil {
    		continue
    	}
    	n := h.equipableNode(slot, e, covered)
    	if n != nil {
	        nodes = append(nodes, n)
    	}
	}    	
	return model.NewNode(nodes...)
}

func (h *Humanoid) equipableNode(slot entity.SlotId, e entity.IEquipable, covered map[entity.DrawSlotId]bool) model.INode {
	if slot == RingRightSlot && covered[ProtrudeHandRight] {
		return nil
	}
	if slot == RingLeftSlot && covered[ProtrudeHandLeft] {
		return nil
	}
	if slot == AmuletSlot && covered[ProtrudeNeck] {
		return nil
	}
	nodes := []model.INode{}
	if en := e.EquippedNode(); en != nil {
		if slot == HandRightSlot {
				nodes = append(nodes, en.Transform(
					h.transforms[HandRight],
					&model.Transform{
						Position: HandRightTransformPosition,
					},
					h.transforms[ArmRight],
					&model.Transform{
						Position: ArmRightTransformPosition,
					},
					))
		} else if slot == HandLeftSlot {
				nodes = append(nodes, en.Transform(
					h.transforms[HandLeft],
					&model.Transform{
						Position: HandLeftTransformPosition,
					},
					h.transforms[ArmLeft],
					&model.Transform{
						Position: ArmLeftTransformPosition,
					},
					))
		} else {
			nodes = append(nodes, e.EquippedNode())
		}
	}
	if y, ok := e.(IHumanoidEquipable); ok {
    	if n := y.ArmRightNode(); n != nil && slot != RingLeftSlot {
    		nodes = append(nodes, n.Transform(
					h.transforms[ArmRight],
					&model.Transform{
						Position: ArmRightTransformPosition,
					},
    			))
    	}
    	if n := y.ArmLeftNode(); n != nil && slot != RingRightSlot {
    		nodes = append(nodes, n.Transform(
					h.transforms[ArmLeft],
					&model.Transform{
						Position: ArmLeftTransformPosition,
					},
    			))
    	}
    	if n := y.LegRightNode(); n != nil {
    		nodes = append(nodes, n.Transform(
					h.transforms[LegRight],
					&model.Transform{
						Position: LegRightTransformPosition,
					},
    			))
    	}
    	if n := y.LegLeftNode(); n != nil {
    		nodes = append(nodes, n.Transform(
					h.transforms[LegLeft],
					&model.Transform{
						Position: LegLeftTransformPosition,
					},
    			))
    	}
    }
    return model.NewNode(nodes...)
}

func (e *Humanoid) CoveredByEquipables() map[entity.DrawSlotId]bool {
	res := map[entity.DrawSlotId]bool{}
	for _, v := range e.Inventory().EquipmentSlots() {
		if v != nil {
			for _, c := range v.CoverDrawSlots() {
				res[c] = true
			}
		}
	}
	return res
}

func (e *Humanoid) Generate() []entity.IEntity {
	list := e.Character.Generate()
	if e.Destroyed() {
		e.Inventory().DropAllLoot()
	}
	for _, v := range e.Inventory().Generate() {
		v.Position().Copy(e.Character.InnerEntity().Position())
		e.Character.AddDontPickUpUntilFar(v)
		list = append(list, v)
	}
	return list
}

func (e *Humanoid) Damaged(damage float32) {
	// TODO fix because not all equipables are equipped
	var defense float32
	if e.inventory != nil {
		for _, q := range e.inventory.EquipmentSlots() {
			if q != nil {
				defense += float32(q.DefenseValue())
			}
		}
	}
	// when def = 0% damage -- get 100% of damage
	// when def = 100% damage -- get 50% of damage
	// when def >= 200% damage -- get 0% (actually 5%) of damage
	// ALWAYS gets a minimum of 5% of original damage
	d := damage * (1 - defense/damage)*0.5+0.5
	d = math32.Max(damage/20, d)
	e.Entity.Damaged(d)
}

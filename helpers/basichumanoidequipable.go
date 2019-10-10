package helpers

import "cavernal.com/model"
import "cavernal.com/entity"
import "cavernal.com/entity/humanoid"
import "cavernal.com/lib/g3n/engine/math32"

type BasicHumanoidEquipable struct {
	BasicEquipable

	cloneSpecification *BasicHumanoidEquipableSpecification

	armRightNode model.INode
	armLeftNode model.INode
	legRightNode model.INode
	legLeftNode model.INode
}

func (e *BasicHumanoidEquipable) ArmRightNode() model.INode {
	return e.armRightNode
}

func (e *BasicHumanoidEquipable) ArmLeftNode() model.INode {
	return e.armLeftNode
}

func (e *BasicHumanoidEquipable) LegRightNode() model.INode {
	return e.legRightNode
}

func (e *BasicHumanoidEquipable) LegLeftNode() model.INode {
	return e.legLeftNode
}

func (e *BasicHumanoidEquipable) Clone() entity.IItem {
  return e.cloneSpecification.New()
}


type BasicHumanoidEquipableSpecification struct {
	BasicEquipableSpecification

	ArmRightNode model.INode
	ArmLeftNode model.INode
	LegRightNode model.INode
	LegLeftNode model.INode
	CoverDrawSlots []entity.DrawSlotId
}

func (s *BasicHumanoidEquipableSpecification) New() *BasicHumanoidEquipable {
	e := s.BasicEquipableSpecification.New()
	e.SetCoverDrawSlots(append(e.CoverDrawSlots(), s.CoverDrawSlots...))
	return &BasicHumanoidEquipable{
		BasicEquipable: *e,
		armRightNode: s.ArmRightNode,
		armLeftNode: s.ArmLeftNode,
		legRightNode: s.LegRightNode,
		legLeftNode: s.LegLeftNode,
		cloneSpecification: s,
	}
}

// --

func NewBasicArmor(category entity.ItemCategory, name string, defense int, main, armLeft, armRight *model.NodeSpec) *BasicHumanoidEquipableSpecification {
	return &BasicHumanoidEquipableSpecification{
		BasicEquipableSpecification: BasicEquipableSpecification{
				Category: category,
				Name: name,
				ItemType: humanoid.Armor,
				BodyType: humanoid.HumanoidBody,
				DefenseValue: defense,
				EquippedNode: main.Build(),
				DroppedNode: model.NewNode(
					main.Build(),
					armLeft.Build().Transform(
						(&model.Transform{
							Position: humanoid.ArmLeftTransformPosition,
							Rotation: &math32.Vector3{math32.Pi/2, 0, 0},
						}).Inverse()),
					armRight.Build().Transform(
						(&model.Transform{
							Position: humanoid.ArmRightTransformPosition,
							Rotation: &math32.Vector3{math32.Pi/2, 0, 0},
						}).Inverse()),
				),
				InventoryNode: model.NewNode(
					main.Build(),
					armLeft.Build().Transform(
						(&model.Transform{
							Position: humanoid.ArmLeftTransformPosition,
							Rotation: &math32.Vector3{math32.Pi/2, 0, 0},
						}).Inverse()),
					armRight.Build().Transform(
						(&model.Transform{
							Position: humanoid.ArmRightTransformPosition,
							Rotation: &math32.Vector3{math32.Pi/2, 0, 0},
						}).Inverse()),
				),
		},
		ArmLeftNode: armLeft.Build(),
		ArmRightNode: armRight.Build(),
		CoverDrawSlots: []entity.DrawSlotId{humanoid.ProtrudeNeck},
		}
}

func NewBasicPants(category entity.ItemCategory, name string, defense int, main, legLeft, legRight *model.NodeSpec) *BasicHumanoidEquipableSpecification {
	return &BasicHumanoidEquipableSpecification{
		BasicEquipableSpecification: BasicEquipableSpecification{
				Category: category,
				Name: name,
				ItemType: humanoid.Pants,
				BodyType: humanoid.HumanoidBody,
				DefenseValue: defense,
				EquippedNode: main.Build(),
				DroppedNode: model.NewNode(
					legLeft.Build(),
					legRight.Build(),
				),
				InventoryNode: model.NewNode(
					legLeft.Build(),
					legRight.Build(),
				),
		},
		LegLeftNode: legLeft.Build(),
		LegRightNode: legRight.Build(),
		}
}

func NewBasicBoots(category entity.ItemCategory, name string, defense int, legLeft, legRight *model.NodeSpec) *BasicHumanoidEquipableSpecification {
	return &BasicHumanoidEquipableSpecification{
		BasicEquipableSpecification: BasicEquipableSpecification{
				Category: category,
				Name: name,
				ItemType: humanoid.Boots,
				BodyType: humanoid.HumanoidBody,
				DefenseValue: defense,
				DroppedNode: model.NewNode(
					legLeft.Build(),
					legRight.Build(),
				),
				InventoryNode: model.NewNode(
					legLeft.Build(),
					legRight.Build(),
				),
		},
		LegLeftNode: legLeft.Build(),
		LegRightNode: legRight.Build(),
		}
}

func NewBasicAmulet(category entity.ItemCategory, name string, defense int, attack int, main *model.NodeSpec) *BasicHumanoidEquipableSpecification {
	return &BasicHumanoidEquipableSpecification{
		BasicEquipableSpecification: BasicEquipableSpecification{
				Category: category,
				Name: name,
				ItemType: humanoid.Amulet,
				BodyType: humanoid.HumanoidBody,
				DefenseValue: defense,
				AttackValue: attack,
				EquippedNode: main.Build(),
				DroppedNode: main.Build(),
				InventoryNode: main.Build(),
		},
		}

}

func NewBasicBackpack(category entity.ItemCategory, name string, defense int, attack int, main *model.NodeSpec) *BasicHumanoidEquipableSpecification {
	return &BasicHumanoidEquipableSpecification{
		BasicEquipableSpecification: BasicEquipableSpecification{
				Category: category,
				Name: name,
				ItemType: humanoid.Back,
				BodyType: humanoid.HumanoidBody,
				DefenseValue: defense,
				AttackValue: attack,
				EquippedNode: main.Build(),
				DroppedNode: main.Build(),
				InventoryNode: main.Build(),
		},
		}

}

func NewBasicGloves(category entity.ItemCategory, name string, defense int, armLeft, armRight *model.NodeSpec) *BasicHumanoidEquipableSpecification {
	return &BasicHumanoidEquipableSpecification{
		BasicEquipableSpecification: BasicEquipableSpecification{
				Category: category,
				Name: name,
				ItemType: humanoid.Gloves,
				BodyType: humanoid.HumanoidBody,
				DefenseValue: defense,
				DroppedNode: model.NewNode(
					armLeft.Build(),
					armRight.Build(),
				),
				InventoryNode: model.NewNode(
					armLeft.Build(),
					armRight.Build(),
				),
		},
		ArmLeftNode: armLeft.Build(),
		ArmRightNode: armRight.Build(),
		CoverDrawSlots: []entity.DrawSlotId{humanoid.ProtrudeHandLeft, humanoid.ProtrudeHandRight},
		}
}

func NewBasicHelmet(category entity.ItemCategory, name string, defense int, main *model.NodeSpec) *BasicHumanoidEquipableSpecification {
	return &BasicHumanoidEquipableSpecification{
		BasicEquipableSpecification: BasicEquipableSpecification{
				Category: category,
				Name: name,
				ItemType: humanoid.Helmet,
				BodyType: humanoid.HumanoidBody,
				DefenseValue: defense,
				EquippedNode: main.Build(),
				DroppedNode: main.Build().Transform(&model.Transform{
						Position: &math32.Vector3{0, -0.8, 0},
					},
				),
				InventoryNode: main.Build().Transform(&model.Transform{
						Position: &math32.Vector3{0, -0.9, 0},
					},
				),
		},
		CoverDrawSlots: []entity.DrawSlotId{humanoid.ProtrudeHeadUpOrBackLong, humanoid.ProtrudeHeadSidesShort, humanoid.ProtrudeHeadSidesLong},
	}

}

func NewBasicRing(category entity.ItemCategory, name string, defense, attack int, armLeft, armRight *model.NodeSpec) *BasicHumanoidEquipableSpecification {
	return &BasicHumanoidEquipableSpecification{
		BasicEquipableSpecification: BasicEquipableSpecification{
				Category: category,
				Name: name,
				ItemType: humanoid.Ring,
				BodyType: humanoid.HumanoidBody,
				DefenseValue: defense,
				DroppedNode: model.NewNode(
					armRight.Build(),
				),
				InventoryNode: model.NewNode(
					armRight.Build(),
				),
		},
		ArmLeftNode: armLeft.Build(),
		ArmRightNode: armRight.Build(),
		}
}
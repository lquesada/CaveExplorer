package helpers

import "github.com/lquesada/cavernal/model"
import "github.com/lquesada/cavernal/entity"
import "github.com/lquesada/cavernal/entity/humanoid"

func NewBasicHumanoid(name string, i *humanoid.Inventory, baseNodes, customNodes map[entity.DrawSlotId]model.INode) *humanoid.Humanoid {
	legLeftTransform := model.NewTransform()
	legRightTransform := model.NewTransform()
	armLeftTransform := model.NewTransform()
	armRightTransform := model.NewTransform()
	handLeftTransform := model.NewTransform()
	handRightTransform := model.NewTransform()
	nodes := map[entity.DrawSlotId]model.INode{}
	for k, v := range baseNodes {
		nodes[k] = v
	}
	for k, v := range customNodes {
		nodes[k] = v
	}

	return humanoid.NewHumanoid(name, i, 
		nodes,
		map[string]*model.Sequence{
			humanoid.Walk: humanoid.HumanoidWalkAnimation(legLeftTransform, legRightTransform),
			humanoid.AttackRight: humanoid.HumanoidAttackAnimation(armRightTransform, handRightTransform, true),
			humanoid.AttackLeft: humanoid.HumanoidAttackAnimation(armLeftTransform, handLeftTransform, false),
			},
		map[entity.DrawSlotId]*model.Transform{
			humanoid.LegLeft: legLeftTransform,
			humanoid.LegRight: legRightTransform,
			humanoid.ArmLeft: armLeftTransform,
			humanoid.ArmRight: armRightTransform,
			humanoid.HandLeft: handLeftTransform,
			humanoid.HandRight: handRightTransform,
		})
}

package equipables

import (
"github.com/lquesada/cavernal/entity"
"github.com/lquesada/cavernal/assets"
"github.com/lquesada/cavernal/entity/humanoid"
"github.com/lquesada/cavernal/model"
"github.com/lquesada/cavernal/helpers"
)

// --

var	(
	ironRingArmLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ironring_armleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmLeftPosition,
				},
			}
	ironRingArmRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ironring_armright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmRightPosition,
				},
			}
)

func NewIronRing() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicRing(entity.Common, "iron ring", 0, 0, ironRingArmLeftModel, ironRingArmRightModel).New()
}

// --

var	(
	copperRingArmLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "copperring_armleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmLeftPosition,
				},
			}
	copperRingArmRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "copperring_armright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmRightPosition,
				},
			}
)

func NewCopperRing() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicRing(entity.Common, "copper ring", 0, 0, ironRingArmLeftModel, ironRingArmRightModel).New()
}

// --

var	(
	goldRingArmLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "goldring_armleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmLeftPosition,
				},
			}
	goldRingArmRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "goldring_armright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmRightPosition,
				},
			}
)

func NewGoldRing() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicRing(entity.Rare, "gold ring", 0, 0, ironRingArmLeftModel, ironRingArmRightModel).New()
}

// --

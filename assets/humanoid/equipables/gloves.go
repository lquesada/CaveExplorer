package equipables

import (
"cavernal.com/entity"
"cavernal.com/assets"
"cavernal.com/entity/humanoid"
"cavernal.com/model"
"cavernal.com/helpers"
)

// --

var	(
	leatherGlovesArmLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "leathergloves_armleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmLeftPosition,
				},
			}
	leatherGlovesArmRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "leathergloves_armright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmRightPosition,
				},
			}
)

func NewLeatherGloves() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicGloves(entity.Common, "leather gloves", 1, leatherGlovesArmLeftModel, leatherGlovesArmRightModel).New()
}

// --

var	(
	ironGlovesArmLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "irongloves_armleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmLeftPosition,
				},
			}
	ironGlovesArmRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "irongloves_armright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmRightPosition,
				},
			}
)

func NewIronGloves() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicGloves(entity.Common, "iron gloves", 2, ironGlovesArmLeftModel, ironGlovesArmRightModel).New()
}

// --

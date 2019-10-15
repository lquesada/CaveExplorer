package equipables

import (
"github.com/lquesada/cavernal/entity"
"github.com/lquesada/cavernal/model"
"github.com/lquesada/cavernal/assets"
"github.com/lquesada/cavernal/entity/humanoid"
"github.com/lquesada/cavernal/helpers"
)


var	(
	tShirtMainModel = &model.NodeSpec{
				Decoder: model.Load(dir, "tshirt_main", assets.Files),
	}
	tShirtArmLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "tshirt_armleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmLeftPosition,
				},
			}
	tShirtArmRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "tshirt_armright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmRightPosition,
				},
			}
)

func NewTShirt() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicArmor(entity.Common, "t-shirt", 0, tShirtMainModel, tShirtArmLeftModel, tShirtArmRightModel).New()
}

// --

var	(
	shirtMainModel = &model.NodeSpec{
				Decoder: model.Load(dir, "shirt_main", assets.Files),
	}
	shirtArmLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "shirt_armleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmLeftPosition,
				},
			}
	shirtArmRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "shirt_armright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmRightPosition,
				},
			}
)

func NewShirt() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicArmor(entity.Common, "shirt", 0, shirtMainModel, shirtArmLeftModel, shirtArmRightModel).New()
}

// --

var	(
	leatherArmorMainModel = &model.NodeSpec{
				Decoder: model.Load(dir, "leatherarmor_main", assets.Files),
	}
	leatherArmorArmLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "leatherarmor_armleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmLeftPosition,
				},
			}
	leatherArmorArmRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "leatherarmor_armright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmRightPosition,
				},
			}
)

func NewLeatherArmor() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicArmor(entity.Common, "leather armor", 1, leatherArmorMainModel, leatherArmorArmLeftModel, leatherArmorArmRightModel).New()
}

// --

var	(
	ironArmorMainModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ironarmor_main", assets.Files),
	}
	ironArmorArmLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ironarmor_armleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmLeftPosition,
				},
			}
	ironArmorArmRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ironarmor_armright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmRightPosition,
				},
			}
)

func NewIronArmor() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicArmor(entity.Common, "iron armor", 2, ironArmorMainModel, ironArmorArmLeftModel, ironArmorArmRightModel).New()
}

// --

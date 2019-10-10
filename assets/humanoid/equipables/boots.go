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
	shoesWithSocksLegLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "shoeswithsocks_legleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegLeftPosition,
				},
			}
	shoesWithSocksLegRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "shoeswithsocks_legright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegRightPosition,
				},
			}
)

func NewShoesWithSocks() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicBoots(entity.Common, "shoes with socks", 0, shoesWithSocksLegLeftModel, shoesWithSocksLegRightModel).New()
}

// --

var	(
	leatherShoesLegLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "leathershoes_legleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegLeftPosition,
				},
			}
	leatherShoesLegRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "leathershoes_legright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegRightPosition,
				},
			}
)

func NewLeatherShoes() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicBoots(entity.Common, "leather shoes", 1, leatherShoesLegLeftModel, leatherShoesLegRightModel).New()
}

// --

var	(
	leatherBootsLegLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "leatherboots_legleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegLeftPosition,
				},
			}
	leatherBootsLegRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "leatherboots_legright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegRightPosition,
				},
			}
)

func NewLeatherBoots() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicBoots(entity.Common, "leather boots", 1, leatherBootsLegLeftModel, leatherBootsLegRightModel).New()
}

// --

var	(
	ironBootsLegLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ironboots_legleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegLeftPosition,
				},
			}
	ironBootsLegRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ironboots_legright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegRightPosition,
				},
			}
)

func NewIronBoots() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicBoots(entity.Common, "iron boots", 2, ironBootsLegLeftModel, ironBootsLegRightModel).New()
}

// --
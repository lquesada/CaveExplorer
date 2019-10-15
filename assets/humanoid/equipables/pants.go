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
	shortJeansMainModel = &model.NodeSpec{
				Decoder: model.Load(dir, "shortjeans_main", assets.Files),
	}
	shortJeansLegLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "shortjeans_legleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegLeftPosition,
				},
			}
	shortJeansLegRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "shortjeans_legright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegRightPosition,
				},
			}
)

func NewShortJeans() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicPants(entity.Common, "short jeans", 0, shortJeansMainModel, shortJeansLegLeftModel, shortJeansLegRightModel).New()
}

// --

var	(
	khakiTrousersMainModel = &model.NodeSpec{
				Decoder: model.Load(dir, "khakitrousers_main", assets.Files),
	}
	khakiTrousersLegLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "khakitrousers_legleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegLeftPosition,
				},
			}
	khakiTrousersLegRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "khakitrousers_legright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegRightPosition,
				},
			}
)

func NewKhakiTrousers() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicPants(entity.Common, "khaki trousers", 0, khakiTrousersMainModel, khakiTrousersLegLeftModel, khakiTrousersLegRightModel).New()
}

// --

var	(
	leatherPantsMainModel = &model.NodeSpec{
				Decoder: model.Load(dir, "leatherpants_main", assets.Files),
	}
	leatherPantsLegLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "leatherpants_legleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegLeftPosition,
				},
			}
	leatherPantsLegRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "leatherpants_legright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegRightPosition,
				},
			}
)

func NewLeatherPants() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicPants(entity.Common, "leather pants", 1, leatherPantsMainModel, leatherPantsLegLeftModel, leatherPantsLegRightModel).New()
}

// --
var	(
	ironPantsMainModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ironpants_main", assets.Files),
	}
	ironPantsLegLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ironpants_legleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegLeftPosition,
				},
			}
	ironPantsLegRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ironpants_legright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegRightPosition,
				},
			}
)

func NewIronPants() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicPants(entity.Common, "iron pants", 2, ironPantsMainModel, ironPantsLegLeftModel, ironPantsLegRightModel).New()
}

// --

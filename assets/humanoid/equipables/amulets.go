package equipables

import (
"github.com/lquesada/cavernal/entity"
"github.com/lquesada/cavernal/assets"
"github.com/lquesada/cavernal/helpers"
"github.com/lquesada/cavernal/model"
)

// --

var crystalAmuletModel = &model.NodeSpec{
				Decoder: model.Load(dir, "crystalamulet", assets.Files),
	}

func NewCrystalAmulet() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicAmulet(entity.Artifact, "crystal amulet", 1, 0, crystalAmuletModel).New()
}

// --

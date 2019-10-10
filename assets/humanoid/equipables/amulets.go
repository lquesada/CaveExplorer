package equipables

import (
"cavernal.com/entity"
"cavernal.com/assets"
"cavernal.com/helpers"
"cavernal.com/model"
)

// --

var crystalAmuletModel = &model.NodeSpec{
				Decoder: model.Load(dir, "crystalamulet", assets.Files),
	}

func NewCrystalAmulet() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicAmulet(entity.Artifact, "crystal amulet", 1, 0, crystalAmuletModel).New()
}

// --

package equipables

import (
"github.com/lquesada/cavernal/entity"
"github.com/lquesada/cavernal/assets"
"github.com/lquesada/cavernal/model"
"github.com/lquesada/cavernal/helpers"
)

// --

var ragBackpackModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ragbackpack", assets.Files),
	}

func NewRagBackpack() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicBackpack(entity.Common, "rag backpack", 0, 0, ragBackpackModel).New()
}

// --

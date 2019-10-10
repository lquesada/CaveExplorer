package equipables

import (
"cavernal.com/entity"
"cavernal.com/assets"
"cavernal.com/model"
"cavernal.com/helpers"
)

// --

var ragBackpackModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ragbackpack", assets.Files),
	}

func NewRagBackpack() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicBackpack(entity.Common, "rag backpack", 0, 0, ragBackpackModel).New()
}

// --

package equipables

import (
"github.com/lquesada/cavernal/entity"
"github.com/lquesada/cavernal/model"
"github.com/lquesada/cavernal/assets"
"github.com/lquesada/cavernal/helpers"
)

// --

var earHoodModel = &model.NodeSpec{
				Decoder: model.Load(dir, "earhood", assets.Files),
			}

func NewEarHood() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicHelmet(entity.Common, "ear hood", 0, earHoodModel).New()
}

// --

var leatherHelmetModel = &model.NodeSpec{
				Decoder: model.Load(dir, "leatherhelmet", assets.Files),
	}

func NewLeatherHelmet() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicHelmet(entity.Common, "leather helmet", 1, leatherHelmetModel).New()
}

// --

var ironHelmetModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ironhelmet", assets.Files),
	}

func NewIronHelmet() *helpers.BasicHumanoidEquipable {
	return helpers.NewBasicHelmet(entity.Common, "iron helmet", 2, ironHelmetModel).New()
}

// --

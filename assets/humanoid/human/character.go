package human

import "github.com/lquesada/cavernal/assets"
import "github.com/lquesada/cavernal/model"
import 	"github.com/lquesada/cavernal/entity/humanoid"
import "github.com/lquesada/cavernal/entity"
import "github.com/lquesada/cavernal/helpers"

type pathMarker struct{} // Needed to find directory
var (
	dir = model.DirOf(pathMarker{})
	legLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "legleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegLeftPosition,
				},
			}
	legRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "legright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.LegRightPosition,
				},
			}
	armLeftModel = &model.NodeSpec{
				Decoder: model.Load(dir, "armleft", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmLeftPosition,
				},
			}
	armRightModel = &model.NodeSpec{
				Decoder: model.Load(dir, "armright", assets.Files),
				Transform: &model.Transform{
					Position: humanoid.ArmRightPosition,
				},
			}
	headModel = &model.NodeSpec{
				Decoder: model.Load(dir, "head", assets.Files),
			}
	torsoModel = &model.NodeSpec{
				Decoder: model.Load(dir, "torso", assets.Files),
			}
	BlondeHairMediumModel = &model.NodeSpec{
		Decoder: model.Load(dir, "blondehairmedium", assets.Files),
	}
	BlondeHairLongModel = &model.NodeSpec{
		Decoder: model.Load(dir, "blondehairlong", assets.Files),
	}
	BlondeBeardFrontModel = &model.NodeSpec{
		Decoder: model.Load(dir, "blondebeardfront", assets.Files),
	}
	BlondeBeardLongModel = &model.NodeSpec{
		Decoder: model.Load(dir, "blondebeardlong", assets.Files),
	}
	BlondeBeardSideModel = &model.NodeSpec{
		Decoder: model.Load(dir, "blondebeardside", assets.Files),
	}
	BlackHairMediumModel = &model.NodeSpec{
		Decoder: model.Load(dir, "blackhairmedium", assets.Files),
	}
	BlackHairLongModel = &model.NodeSpec{
		Decoder: model.Load(dir, "blackhairlong", assets.Files),
	}
	BlackBeardFrontModel = &model.NodeSpec{
		Decoder: model.Load(dir, "blackbeardfront", assets.Files),
	}
	BlackBeardLongModel = &model.NodeSpec{
		Decoder: model.Load(dir, "blackbeardlong", assets.Files),
	}
	BlackBeardSideModel = &model.NodeSpec{
		Decoder: model.Load(dir, "blackbeardside", assets.Files),
	}
)

func New(name string, i *humanoid.Inventory, customNodes map[entity.DrawSlotId]model.INode) *humanoid.Humanoid {
	return helpers.NewBasicHumanoid(name, i, map[entity.DrawSlotId]model.INode{
			humanoid.Head: headModel.Build(),
			humanoid.Torso: torsoModel.Build(),
			humanoid.LegLeft: legLeftModel.Build(),
			humanoid.LegRight: legRightModel.Build(),
			humanoid.ArmLeft: armLeftModel.Build(),
			humanoid.ArmRight: armRightModel.Build(),
		}, customNodes)
}

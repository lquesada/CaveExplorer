package skeleton

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

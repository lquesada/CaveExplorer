package humanoid

import (
	"github.com/lquesada/cavernal/entity"
	"github.com/lquesada/cavernal/model"
)

type IHumanoidEquipable interface {
	entity.IEquipable

	ArmRightNode() model.INode
	ArmLeftNode() model.INode
	LegRightNode() model.INode
	LegLeftNode() model.INode
}
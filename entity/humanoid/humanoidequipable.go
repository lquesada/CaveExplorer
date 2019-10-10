package humanoid

import (
	"cavernal.com/entity"
	"cavernal.com/model"
)

type IHumanoidEquipable interface {
	entity.IEquipable

	ArmRightNode() model.INode
	ArmLeftNode() model.INode
	LegRightNode() model.INode
	LegLeftNode() model.INode
}
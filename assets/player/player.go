package player

import 	"cavernal.com/assets/humanoid/human"
import 	"cavernal.com/assets/humanoid/customization"
import 	"cavernal.com/entity/humanoid"
import 	"cavernal.com/model"
import "cavernal.com/entity"

type HumanPlayer struct{ 
	humanoid.Humanoid
}

func New(name string, useSlotsCount, itemCount int) *HumanPlayer {
	i := humanoid.NewInventory(humanoid.StandardSlots, itemCount)

	h := human.New(name, i, map[entity.DrawSlotId]model.INode{
			humanoid.Face: customization.FaceScaredModel.Build(),
			humanoid.ProtrudeHeadUpOrBackShort: human.BlondeHairMediumModel.Build(),
			humanoid.ProtrudeHeadUpOrBackLong: human.BlondeHairLongModel.Build(),
		})
	h.SetTargetFilter(entity.FilterEvil)
	h.SetPushForce(9)

	return &HumanPlayer{Humanoid: *h}
}

func (e *HumanPlayer) Good() {
	// Marker for player
}

func (e *HumanPlayer) Think() {
	e.Humanoid.Think()
}

func (e *HumanPlayer) WantToPickUp(i entity.IItem) bool {
	return e.Humanoid.WantToPickUp(i) && e.Inventory().CanPickUp(i)
}

func (e *HumanPlayer) PickUp(i entity.IItem) {
	e.Inventory().PickUp(i)
}

func (e *HumanPlayer) Damaged(damage float32) {
	e.Humanoid.Damaged(damage)
}

package humanoid

import "fmt"
import "cavernal.com/model"
import "cavernal.com/lib/g3n/engine/math32"

var (
	HumanoidAttackAction model.ActionId = "HumanoidAttackAction"

	// For tests
	n = model.NewTransform()

	legStanding = &model.Transform{
		Position: math32.NewVector3(0, 0, 0),
		Rotation: math32.NewVector3(0, 0, 0),
	}
	legForward = &model.Transform{
		Position: math32.NewVector3(0, 0, 0.03),
		Rotation: math32.NewVector3(-0.3, 0, 0),
	}
	legBackward = &model.Transform{
		Position: math32.NewVector3(0, 0, -0.03),
		Rotation: math32.NewVector3(0.3, 0, 0),
	}

	HumanoidWalkAnimation = func(legLeft, legRight *model.Transform) *model.Sequence {
		s, err := model.New(
			400,
			2000./400,
			true,
			false,
			400,
			[]*model.Step{
				model.NewStep(0, true, true, nil, nil,
					map[*model.Transform]*model.TransformStep{
						legLeft: &model.TransformStep{Set: legStanding, Mutate: legBackward},
						legRight: &model.TransformStep{Set: legStanding, Mutate: legForward},
					},
				),
				model.NewStep(100, false, false, nil, nil,
					map[*model.Transform]*model.TransformStep{
						legLeft: &model.TransformStep{Set: legBackward, Mutate: legStanding},
						legRight: &model.TransformStep{Set: legForward, Mutate: legStanding},
					},
				),
				model.NewStep(200, true, true, nil, nil,
					map[*model.Transform]*model.TransformStep{
						legLeft: &model.TransformStep{Set: legStanding, Mutate: legForward},
						legRight: &model.TransformStep{Set: legStanding, Mutate: legBackward},
					},
				),
				model.NewStep(300, false, false, nil, nil,
					map[*model.Transform]*model.TransformStep{
						legLeft: &model.TransformStep{Set: legForward, Mutate: legStanding},
						legRight: &model.TransformStep{Set: legBackward, Mutate: legStanding},
					},
				),
			})
			if err != nil {
				panic(fmt.Sprintf("Critical error. Animation is broken: %v", err))
			}
			return s
		}
	HumanoidWalkAnimationTest = HumanoidWalkAnimation(n, n)

	armFront = &model.Transform{
		Rotation: math32.NewVector3(0, 0, 0),
	}
	armUp = &model.Transform{
		Rotation: math32.NewVector3(-math32.Pi/2, 0, 0),
	}
	armDown = &model.Transform{
		Rotation: math32.NewVector3(-math32.Pi/2, 0, 0),
	}
	armRightDownFrontCenter = &model.Transform{
		Rotation: math32.NewVector3(0.6, 0.4, 0),
	}
	armRightDownFrontCenterHit = &model.Transform{
		Rotation: math32.NewVector3(0.3, 0.2, 0),
	}
	armLeftDownFrontCenter = &model.Transform{
		Rotation: math32.NewVector3(0.6, -0.4, 0),
	}
	armLeftDownFrontCenterHit = &model.Transform{
		Rotation: math32.NewVector3(0.3, -0.2, 0),
	}
	handUp = &model.Transform{
		Rotation: math32.NewVector3(0, 0, 0),
	}
	handMid = &model.Transform{
		Rotation: math32.NewVector3(math32.Pi/4, 0, 0),
	}
	handHit = &model.Transform{
		Rotation: math32.NewVector3(math32.Pi/3, 0, 0),
	}
	handFront = &model.Transform{
		Rotation: math32.NewVector3(math32.Pi/2, 0, 0),
	}

	HumanoidAttackAnimation = func(arm, hand *model.Transform, right bool) *model.Sequence {
		var armDownFrontCenter, armDownFrontCenterHit *model.Transform
		if right {
			armDownFrontCenter = armRightDownFrontCenter
			armDownFrontCenterHit = armRightDownFrontCenterHit
		} else {
			armDownFrontCenter = armLeftDownFrontCenter
			armDownFrontCenterHit = armLeftDownFrontCenterHit
		}
		s, err := model.New(
			120,
			4000./400,
			false,
			true,
			2,
			[]*model.Step{
				model.NewStep(0, true, true, nil, nil,
					map[*model.Transform]*model.TransformStep{
						arm: &model.TransformStep{Set: armFront, Mutate: armUp},
						hand: &model.TransformStep{Set: handUp, Mutate: handMid},
					},
				),
				model.NewStep(15, false, false, nil, nil,
					map[*model.Transform]*model.TransformStep{
						arm: &model.TransformStep{Set: armUp, Mutate: armDownFrontCenterHit},
						hand: &model.TransformStep{Set: handMid, Mutate: handHit},
					},
				),
				model.NewStep(32, false, false, nil, []model.ActionId{HumanoidAttackAction},
					map[*model.Transform]*model.TransformStep{
						arm: &model.TransformStep{Set: armDownFrontCenterHit, Mutate: armDownFrontCenter},
						hand: &model.TransformStep{Set: handHit, Mutate: handFront},
					},
				),
				model.NewStep(60, false, false, nil, nil,
					map[*model.Transform]*model.TransformStep{
						arm: &model.TransformStep{Set: armDownFrontCenter, Mutate: armDown},
						hand: &model.TransformStep{Set: handFront},
					},
				),
				model.NewStep(80, false, false, nil, nil,
					map[*model.Transform]*model.TransformStep{
						arm: &model.TransformStep{Set: armDown, Mutate: armFront},
						hand: &model.TransformStep{Set: handFront, Mutate: handUp},
					},
				),
			})
			if err != nil {
				panic(fmt.Sprintf("Critical error. Animation is broken: %v", err))
			}
			return s
		}
	HumanoidAttackAnimationTest = HumanoidAttackAnimation(n, n, true)
)
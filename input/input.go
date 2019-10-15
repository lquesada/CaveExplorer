package input

import (
	"github.com/lquesada/cavernal/lib/g3n/engine/util/application"
	"github.com/lquesada/cavernal/lib/g3n/engine/window"
	"github.com/lquesada/cavernal/lib/g3n/engine/gui"
	"github.com/lquesada/cavernal/lib/g3n/engine/math32"
 "github.com/lquesada/cavernal/entity"
 "github.com/lquesada/cavernal/model"
)

type Input struct{
	panel *gui.Panel
	player entity.IPlayer
	wantToAttack bool
	jump bool
	wantToAttackCurrent bool
	wantToAttackAngle float32
	waitingForAngle bool
	mousePressed bool
	padAttack bool
	mouseAngle float32
	isometricAngle float32
	heightCenterFactor float32
	worldRotateAngle float32
	shiftAngle float32
	hasShiftAngle bool
    isoSin float32
    isoCos float32
    sequenceKeepLooking *model.Sequence
}

func New(panel *gui.Panel, player entity.IPlayer, isometricAngle, heightCenterFactor, minimumWaitAngleAttack float32, worldRotateAngle float32) *Input {
	return &Input{
		panel: panel,
		player: player,
		isometricAngle: isometricAngle,
		heightCenterFactor: heightCenterFactor,
		worldRotateAngle: worldRotateAngle,
    	isoSin: math32.Sin(isometricAngle),
    	isoCos: math32.Cos(isometricAngle),
    	sequenceKeepLooking: model.NewTimer(minimumWaitAngleAttack),
	}
}
func (i *Input) Subscribe() {
	    app := application.Get()
		i.panel.SubscribeID(window.OnMouseDown, i, i.mouseDown)
		i.panel.SubscribeID(window.OnMouseUp, i, i.mouseUp)
		i.panel.SubscribeID(window.OnCursor, i, i.mouseCursor)
		app.Window().SubscribeID(window.OnKeyDown, i, i.keyDown)
}

func (i *Input) Unsubscribe() {
	    app := application.Get()
		i.panel.UnsubscribeID(window.OnMouseDown, i)
		i.panel.UnsubscribeID(window.OnMouseUp, i)
		i.panel.UnsubscribeID(window.OnCursor, i)
		app.Window().UnsubscribeID(window.OnKeyDown, i)
}

func (i* Input) mouseDown(name string, ev interface{}) {
	e := ev.(*window.MouseEvent)
	if e.Button == window.MouseButton1 {
		i.wantToAttack = true
		// This assumes GUI is at the right side.
		i.wantToAttackAngle = i.mousePositionToAngle(e.Xpos, e.Ypos)
		i.mousePressed = true
	}
}

func (i* Input) mouseUp(name string, ev interface{}) {
	e := ev.(*window.MouseEvent)
	if e.Button == window.MouseButton1 {
		i.mousePressed = false
	}
}

func (i* Input) mouseCursor(name string, ev interface{}) {
	e := ev.(*window.CursorEvent)
	i.mouseAngle = i.mousePositionToAngle(e.Xpos, e.Ypos)
}

func (i *Input) mousePositionToAngle(x, y float32) float32 {
	tx := (x - i.panel.Width()/2) / i.isoCos
	ty := (y - (i.panel.Height()/2)*i.heightCenterFactor) / i.isoSin
	return entity.NormalizeAngle(-math32.Atan2(ty, tx)-i.worldRotateAngle+math32.Pi/2)

}

func (i* Input) keyDown(name string, ev interface{}) {
		e := ev.(*window.KeyEvent)
		// L to attack in current direction
        if e.Keycode == window.KeyL {
			i.wantToAttackCurrent = true
			return
		}
		if e.Keycode == window.KeySpace {
			i.jump = true
		}
		if (e.Keycode == window.KeyKP8 ||
			e.Keycode == window.KeyKP9 ||
			e.Keycode == window.KeyKP6 ||
			e.Keycode == window.KeyKP3 ||
			e.Keycode == window.KeyKP2 ||
			e.Keycode == window.KeyKP1 ||
			e.Keycode == window.KeyKP4 ||
			e.Keycode == window.KeyKP7 ||
			e.Keycode == window.KeyU ||
			e.Keycode == window.KeyK ||
			e.Keycode == window.KeyJ ||
			e.Keycode == window.KeyH) {
			i.padAttack = true
		}
}

func (i* Input) Apply(delta float32) {
    app := application.Get()

	shift := app.KeyState().Pressed(window.KeyLeftShift) || app.KeyState().Pressed(window.KeyRightShift)

	// Movement is treated as continuous
	var x, z float32
    if app.KeyState().Pressed(window.KeyUp) || app.KeyState().Pressed(window.KeyW) {
            z -= 1
    }
    if app.KeyState().Pressed(window.KeyRight) || app.KeyState().Pressed(window.KeyD) {
            x += 1
    }
    if app.KeyState().Pressed(window.KeyDown) || app.KeyState().Pressed(window.KeyS) {
            z += 1
    }
    if app.KeyState().Pressed(window.KeyLeft) || app.KeyState().Pressed(window.KeyA) {
            x -= 1
    }
    i.player.MoveIntention().Set(0, 0, 0)
  	angle := entity.NormalizeAngle(-math32.Atan2(z, x)-i.worldRotateAngle+math32.Pi/2)
    if x != 0 || z != 0 {
		i.player.MoveIntention().X = math32.Sin(angle)
		i.player.MoveIntention().Z = math32.Cos(angle)
		i.player.MoveIntention().Normalize()
	}
	if !i.hasShiftAngle && shift {
		if x == 0 && z == 0 {
			angle = i.player.LookIntention()
		}
		i.shiftAngle = angle
		i.hasShiftAngle = true
	}
	if !shift {
		i.hasShiftAngle = false
	}

	// Look intention
	if i.sequenceKeepLooking.Running() {
		i.sequenceKeepLooking.Add(delta)
	}

	// Shift keeps look intention
	if !shift {
		if x != 0 || z != 0 && !i.player.IsAttacking() {
			i.player.SetLookIntention(angle)
		}
		angleDeltaIntention := math32.Abs(entity.NormalizeAngle(i.player.LookAngle()-i.wantToAttackAngle))
		tooFarAngle := angleDeltaIntention > 0.0001
		if !tooFarAngle || !i.player.IsAttacking() {
			i.waitingForAngle = false
		}
		if i.player.IsAttacking() || (i.waitingForAngle && tooFarAngle) || i.sequenceKeepLooking.Running() {
			i.player.SetLookIntention(i.wantToAttackAngle)
		}
	} else {
		i.player.SetLookIntention(i.shiftAngle)
	}

	// Attack (one-off)
	// Input for one-off attack (i.e. KeyDown)
	var attack bool

	if i.wantToAttackCurrent {
		attack = true
		i.wantToAttackCurrent = false
		if shift {
			i.wantToAttackAngle = i.shiftAngle
		} else {
			i.wantToAttackAngle = angle
		}
	}

	// When using pad
	var pX, pZ float32
	// Keypad to attack in 8 angles
	if app.KeyState().Pressed(window.KeyKP8) {
			pZ -= 1
	}
	if app.KeyState().Pressed(window.KeyKP9) {
			pZ -= 1
			pX += 1
		}
	if app.KeyState().Pressed(window.KeyKP6) {
		pX += 1
	}
	if app.KeyState().Pressed(window.KeyKP3) {
		pX += 1
		pZ += 1
	}
	if app.KeyState().Pressed(window.KeyKP2) {
		pZ += 1
	}
	if app.KeyState().Pressed(window.KeyKP1) {
		pZ += 1
		pX -= 1
	}
	if app.KeyState().Pressed(window.KeyKP4) {
		pX -= 1
	}
	if app.KeyState().Pressed(window.KeyKP7) {
		pX -= 1
		pZ -= 1
	}
	// UHJK to attack in 4 angles
	if app.KeyState().Pressed(window.KeyU) {
		pZ -= 1
	}
	if app.KeyState().Pressed(window.KeyK) {
		pX += 1
	}
	if app.KeyState().Pressed(window.KeyJ) {
		pZ += 1
	}
	if app.KeyState().Pressed(window.KeyH) {
		pX -= 1			
	}
	if i.padAttack {
		i.padAttack = false
  	    if pX != 0 || pZ != 0 {
		  	i.wantToAttackAngle = entity.NormalizeAngle(-math32.Atan2(pZ, pX)-i.worldRotateAngle+math32.Pi/2)
		  	i.shiftAngle = i.wantToAttackAngle
		  	i.wantToAttack = true
		  	attack = true
		}
	}

	if i.wantToAttack {
		attack = true
		i.wantToAttack = false
	}

	// Attack (repeat)
	var repeatAttack bool
	if app.KeyState().Pressed(window.KeyL) {
		repeatAttack = true
		if shift {
			i.wantToAttackAngle = i.shiftAngle
		} else {
			i.wantToAttackAngle = angle
		}
	}
	if pX != 0 || pZ != 0 {
	  	i.wantToAttackAngle = entity.NormalizeAngle(-math32.Atan2(pZ, pX)-i.worldRotateAngle+math32.Pi/2)
	  	repeatAttack = true
	}
	if i.mousePressed {
		repeatAttack = true
		i.wantToAttackAngle = i.mouseAngle
	}

	if attack {
		if i.player.WantToAttack() {
			i.waitingForAngle = true
			i.sequenceKeepLooking.Reset()
			i.sequenceKeepLooking.Start()
		}
	} else if repeatAttack {
		if i.player.WantToRepeatAttack() {
			i.waitingForAngle = true
			i.sequenceKeepLooking.Reset()
			i.sequenceKeepLooking.Start()
		}
	}

	// Jump
	if i.jump {
		i.player.Jump()
		i.jump = false
	} else if app.KeyState().Pressed(window.KeySpace) {
		i.player.RepeatJump()
	}
}

package main

import (
	"flag"
	"fmt"
	"time"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"github.com/lquesada/cavernal/assets/enemies"
	"github.com/lquesada/cavernal/assets/scenery"
	"github.com/lquesada/cavernal/assets/player"
	"github.com/lquesada/cavernal/helpers"
	"github.com/lquesada/cavernal/input"
	"github.com/lquesada/cavernal/hud"
	"github.com/lquesada/cavernal/assets/equipables"
	"github.com/lquesada/cavernal/entity/humanoid"
	humanoidequipables "github.com/lquesada/cavernal/assets/humanoid/equipables"
	"github.com/lquesada/cavernal/entity"
	"github.com/lquesada/cavernal/world"
	"github.com/lquesada/cavernal/lib/g3n/engine/core"
	"github.com/lquesada/cavernal/lib/g3n/engine/math32"
	"github.com/lquesada/cavernal/game"
	"github.com/lquesada/cavernal/lib/g3n/engine/light"
	"github.com/lquesada/cavernal/lib/g3n/engine/util/application"
	"github.com/lquesada/cavernal/lib/g3n/engine/window"
)

var pprofPort = flag.Int("pprof_port", 0, "open port for pprof")

func main() {
    flag.Parse()
    if *pprofPort != 0 {
		go func() { http.ListenAndServe(fmt.Sprintf("localhost:%d", *pprofPort), nil) }()
		// Yes, no graceful degradation when doing profiling. Whatever.
		// go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
		// go tool pprof http://localhost:6060/debug/pprof/heap
	    // go tool pprof http://localhost:6060/debug/pprof/allocs
	}
	app, _ := application.Create(application.Options{
		Title:  "Cavernal",
		Width:  800,
		Height: 600,
	})
	app.Window().UnsubscribeID(window.OnWindowSize, app)

	app.CameraPersp().SetPosition(0, 10, 10)
	app.CameraPersp().LookAt(&math32.Vector3{0, 0, 0})
	app.Orbit().Enabled = false

	c := &RogueController{
		shadows: 0.3,
		isometricAngle: math32.DegToRad(30),
		heightCenterFactor: 0.975,
		minimumWaitAngleAttack: 0.5,
		worldRotateAngle: -math32.Pi/4,
		tilesFarToEnterDoor: 0.7,
		initialLevel: 0,
		finalLevel: 8,
		tileSize: 2.2,
		gravity: 40,
		waitBeforePickupable: 0.3,
		fallSpeedUp: 1,

		artifacts: []entity.IItem{
			equipables.NewDemonbloodSword(),
			equipables.NewFinnSword(),
			equipables.NewGrassSword(),
		},

		rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
	c.Restart()

	c.Register(app)
	app.Run()

	c.Dispose()
}

// Game controller
type RogueController struct{
	app *application.Application

	// Config
	shadows float32
	isometricAngle float32
	heightCenterFactor float32
	minimumWaitAngleAttack float32
	worldRotateAngle float32
	tilesFarToEnterDoor float32
	initialLevel int
	finalLevel int
	tileSize float32
	gravity float32
	waitBeforePickupable float32
	fallSpeedUp float32
	artifacts []entity.IItem

	backMarkers map[int]entity.IEntity
	nextMarkers map[int]entity.IEntity
	levels map[int]*game.State
	level int
	wentNext bool
	wentBack bool
	farEnoughFromDoor bool
	rand *rand.Rand

	currentBackMarker entity.IEntity
	currentNextMarker entity.IEntity
	currentLevel *game.State
	currentArtifact entity.IItem

	player entity.IPlayer
	hud hud.IHUD
    input *input.Input

	brightness       float32

	retryStep        int
	retryCounter     float32
	anyKeyPressedOrMouseClicked  bool
	wentNextStep			 int
	wentNextCounter		float32
	wentBackStep			 int
	wentBackCounter		float32
	lightWinCounter float32
	ambientLight     *light.Ambient
	directionalLight *light.Directional

    clickedAngle float32
    clicked bool
}

func (c *RogueController) Register(app *application.Application) {
	c.app = app
	c.app.Window().SubscribeID(window.OnKeyDown, nil, c.retryOnAction)
	c.app.Window().SubscribeID(window.OnMouseDown, nil, c.retryOnAction)
	c.app.SubscribeID(application.OnBeforeRender, nil, c.updateScene)
}

func (c *RogueController) Unregister() {
  if c.app == nil {
  	return
  }
  c.app.Window().UnsubscribeID(window.OnKeyDown, nil)
  c.app.Window().UnsubscribeID(window.OnMouseDown, nil)
  c.app.UnsubscribeID(application.OnBeforeRender, nil)
  c.app = nil
}

func (c *RogueController) Dispose() {
  c.Unregister()
  if c.input != nil {
    c.input.Unsubscribe()
  }
  if c.hud != nil {
    c.hud.Dispose()
  }
}

func (c *RogueController) Restart() {
	c.brightness = 1

	if c.ambientLight != nil {
		c.ambientLight.Dispose()
	}
	c.ambientLight = light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0)
	if c.directionalLight != nil {
		c.directionalLight.Dispose()
	}
	c.directionalLight = light.NewDirectional(math32.NewColor("white"), 0)
	c.directionalLight.SetPosition(1, 1, 1)

	player := player.New("The Human", 10, 40)
	player.Inventory().Equip(player.Inventory().PickUp(equipables.NewIronSword()), humanoid.HandRightSlot)
	//player.Inventory().Equip(player.Inventory().PickUp(equipables.NewWoodenShield()), humanoid.HandLeftSlot)
	/*player.Inventory().PickUp(equipables.NewArrow())
	player.Inventory().PickUp(equipables.NewArrow())
	player.Inventory().PickUp(equipables.NewArrow())
	player.Inventory().PickUp(equipables.NewArrow())
	player.Inventory().PickUp(equipables.NewArrow())
	player.Inventory().PickUp(equipables.NewArrow())
	player.Inventory().PickUp(equipables.NewShortBow())*/
		player.Inventory().AutoEquip(player.Inventory().PickUp(humanoidequipables.NewShortJeans()))
	player.Inventory().AutoEquip(player.Inventory().PickUp(humanoidequipables.NewShoesWithSocks()))
	player.Inventory().AutoEquip(player.Inventory().PickUp(humanoidequipables.NewTShirt()))
	player.Inventory().AutoEquip(player.Inventory().PickUp(humanoidequipables.NewEarHood()))
	player.Inventory().AutoEquip(player.Inventory().PickUp(humanoidequipables.NewRagBackpack()))
	player.SetMaxHealth(10000)
	player.SetHealth(10000)
	c.player = player

	c.hud = humanoid.NewHUD(c.player)
	c.input = input.New(c.hud.Panel(), c.player, c.isometricAngle, c.heightCenterFactor, c.minimumWaitAngleAttack, c.worldRotateAngle)
	c.input.Subscribe()

	c.level = c.initialLevel
	c.setupWorld()

	c.wentNext = false
	c.wentBack = false

	c.retryStep = 0
	c.retryCounter = 0
	c.anyKeyPressedOrMouseClicked = false
	c.wentNextStep = 0
	c.wentNextCounter = 0
	c.wentBackStep = 0
	c.wentBackCounter = 0
	c.lightWinCounter = 0
    c.clickedAngle = 0
    c.clicked = false

	c.generateLevel(c.level)
	// Hack to set player correctly in the first level
	c.level--
	c.nextLevel()
}

func (c *RogueController) setupWorld() {
	c.backMarkers = map[int]entity.IEntity{}
	c.nextMarkers = map[int]entity.IEntity{}
	c.levels = map[int]*game.State{}
	c.currentArtifact = nil
	if len(c.artifacts) > 0 {
		picker := int(c.rand.Float32()*float32(len(c.artifacts)))
		c.currentArtifact = c.artifacts[picker]
		c.artifacts = append(c.artifacts[:picker], c.artifacts[picker+1:]...)
	}
}

func (c *RogueController) updateScene(_ string, _ interface{}) {
	c.app.Scene().DisposeChildren(true)
	c.app.Scene().Dispose()

	delta := c.app.FrameDeltaSeconds()

	// Manage game over
	if c.over() || c.retryStep != 0 {
		switch c.retryStep {
		case 0:
			c.retryStep = 1
			c.retryCounter = 0
		case 1:
			if c.retryCounter < 1 {
				c.retryCounter += delta
				c.brightness = math32.Min(1, math32.Max(0, 1-c.retryCounter))
			} else {
				c.retryStep = 2
				c.anyKeyPressedOrMouseClicked = false
			}
		case 2:
			if c.anyKeyPressedOrMouseClicked {
				c.retryStep = 3
				c.Restart()
			}
		case 3:
			if c.retryCounter < 2 {
				c.retryCounter += delta
				c.brightness = math32.Min(1, math32.Max(0, c.retryCounter-1))
			} else {
				c.retryStep = 0
				c.brightness = 1
			}
		}
	} else if c.wentNext || c.wentNextStep != 0 {
		switch c.wentNextStep {
		case 0:
			c.wentNextStep = 1
			c.wentNextCounter = 0
		case 1:
			if c.wentNextCounter < 0.5 {
				c.wentNextCounter += delta
				c.brightness = math32.Min(1, math32.Max(0, 1-(c.wentNextCounter*2)))
			} else {
				c.nextLevel()
				c.wentNextStep = 2
			}
		case 2:
			if c.wentNextCounter < 1 {
				c.wentNextCounter += delta
				c.brightness = math32.Min(1, math32.Max(0, (c.wentNextCounter*2)-1))
			} else {
				c.wentNextStep = 3
			}
		case 3:
			c.brightness = 1
			c.wentNextStep = 0
		}
	} else if c.wentBack || c.wentBackStep != 0 {
		switch c.wentBackStep {
		case 0:
			c.wentBackStep = 1
			c.wentBackCounter = 0
		case 1:
			if c.wentBackCounter < 0.5 {
				c.wentBackCounter += delta
				c.brightness = math32.Min(1, math32.Max(0, 1-(c.wentBackCounter*2)))
			} else {
				if c.level == 0 {
					// Going back from 0 = exit
					c.app.Quit()
				} else {
					c.previousLevel()
				}
				c.wentBackStep = 2
			}
		case 2:
			if c.wentBackCounter < 1 {
				c.wentBackCounter += delta
				c.brightness = math32.Min(1, math32.Max(0, (c.wentBackCounter*2)-1))
			} else {
				c.wentBackStep = 3
			}
		case 3:
			c.brightness = 1
			c.wentBackStep = 0
		}
	} else {
		c.Update(delta)
	}

	c.app.Gl().ClearColor(0.0, 0.0, 0.0, 1)

	// Default to drawing the whole scene.
    width, height := c.app.Window().FramebufferSize()
    c.app.Gl().Viewport(0, 0, int32(width), int32(height))
    aspect := float32(int32(width)) / float32(height)
    c.app.Camera().SetAspect(aspect)

	c.hud.PreRender()

	scene := core.NewNode()

    // Increase brightness if won
	if c.player.Inventory().WhereHas(c.currentArtifact) != -1 {
		c.lightWinCounter += delta
		if c.lightWinCounter > 1 {
			c.lightWinCounter = 1
		}
	} else {
		c.lightWinCounter -= delta
		if c.lightWinCounter < 0 {
			c.lightWinCounter = 0
		}

	}
	c.ambientLight = light.NewAmbient(&math32.Color{1.0, 1.0+c.lightWinCounter*0.5, 1.0+c.lightWinCounter}, 0)

	c.ambientLight.SetIntensity((1 - c.shadows) * c.brightness)
	scene.Add(c.ambientLight)
	c.directionalLight.SetIntensity(c.shadows*c.brightness)
	scene.Add(c.directionalLight)

	gameNode := c.currentLevel.Node().G3NNode()

	gameNode.SetPosition(
		-(c.player.Position().X*math32.Cos(c.worldRotateAngle)+c.player.Position().Z*math32.Sin(c.worldRotateAngle)),
		0,
		-(c.player.Position().Z*math32.Cos(c.worldRotateAngle)-c.player.Position().X*math32.Sin(c.worldRotateAngle)))
	gameNode.SetRotation(0, c.worldRotateAngle, 0)
	scene.Add(gameNode)

	c.app.SetScene(scene)   
}


func (c *RogueController) retryOnAction(evname string, ev interface{}) {
	if evname == window.OnKeyDown || evname == window.OnMouseDown {
		c.anyKeyPressedOrMouseClicked = true
	}
}

func (c *RogueController) hasLevel(level int) bool {
	_, ok := c.levels[level]
	return ok
}

func (c *RogueController) generateLevel(level int) {
	// Walkable
	o := func() world.ITile { return scenery.NewConcreteTile() }
	// Special
	u := func() world.ITile { t := scenery.NewConcreteColumnWith(scenery.NewConcreteTile()); return t }
	// Connector
	x := func() world.ITile { t := scenery.NewConcreteWall(); t.SetConnect(true); return t }
	// Hole
	h := func() world.ITile {
		t := scenery.NewHole()
		// TODO remove this bandaid
		t.SetWalkThrough(false)
		t.SetFallThrough(false)
		return t
	}
	// Wall
	w := func() world.ITile { return scenery.NewConcreteWall() }
	// With height
	fu := func(y float32) func() world.ITile { return func() world.ITile { t := u(); t.SetY(y*0.15); return t } } 
	fo := func(y float32) func() world.ITile { return func() world.ITile { t := o(); t.SetY(y*0.15); return t } } 
	n := func() world.ITile { return nil }
	marker := func() entity.IEntity { return (&helpers.BasicEntitySpecification{
		Name: "marker",
		Radius: 0.5,
		ShadowOpacity: 0,
		Height: 1.15,
		}).New()
	}

	var s *game.State
	var nextMarker, backMarker entity.IEntity

	if level == c.finalLevel { // FINAL LEVEL
		wGen := helpers.NewWorldGenerator(
			nil,
			&world.Room{
				Floor: [][]func() world.ITile{
					{w, n, x, n, w},
					{w, o, o, o, w},
					{w, o, o, o, w},
					{w, o, o, o, w},
					{w, w, w, w, w},
				},
				InterestX: 2,
				InterestZ: 4,
			},
			&world.Room{
				Floor: [][]func() world.ITile{
					{w, w, w, w, w, w, w, w, w},
					{w, o, o, h, h, h, o, o, w},
					{w, o, h, h, fo(2), h, h, o, w},
					{w, h, h, fo(2), fo(3), fo(2), h, h, w},
					{w, h, fo(2), fo(3), fu(4), fo(3), fo(2), h, w},
					{w, h, h, fo(2), fo(3), fo(2), h, h, w},
					{w, o, h, h, fo(2), h, h, o, w},
					{w, o, o, h, fo(1), h, o, o, w},
					{w, o, o, o, o, o, o, o, w},
					{w, o, o, o, o, o, o, o, w},
					{w, w, w, o, x, o, w, w, w},
				},
				InterestX: 4,
				InterestZ: 4,
			},
			helpers.MatchAnythingComparator,
			o,
		)
		currentLevel, start, end, _, _ := wGen.Generate(c.rand.Int(), c.tileSize, c.gravity, 0, 0, false)
		s = game.NewState(c.waitBeforePickupable, c.fallSpeedUp, c.player, currentLevel)
		currentLevel.SetTile(start.Coords.X+start.Room.InterestX, start.Coords.Z+start.Room.InterestZ, scenery.NewBackDoorWithWall())

		backMarker = marker()
		wGen.PlaceImportantEntityInRoom(s, start, backMarker)
		backMarker.Position().Z -= currentLevel.TileSize()/2

		if c.currentArtifact != nil {
			wGen.PlaceImportantEntityInRoom(s, end, c.currentArtifact)
		}
	/*} else if level == -1 { // OUT OF THE CAVE LEVEL
		wGen := helpers.NewWorldGenerator(
			nil,
			&world.Room{
				Floor: [][]func() world.ITile{
					{n, n, n, n, n, x, n, n, n, n, n},
					{w, o, o, o, o, o, o, o, o, o, w},
					{w, w, w, w, w, w, w, w, w, w, w},
				},
				InterestX: 5,
				InterestZ: 2,
			},
			&world.Room{
				Floor: [][]func() world.ITile{
					{w, w, w, w, w, w, w, w, w, w, w},
					{w, o, o, o, o, o, o, o, o, o, w},
					{w, o, o, o, o, o, o, o, o, o, w},
					{w, o, o, o, o, o, o, o, o, o, w},
					{w, o, o, o, o, x, o, o, o, o, w},
				},
				InterestX: 5,
				InterestZ: 0,
			},
			helpers.MatchAnythingComparator,
			o,
		)
		currentLevel, start, end, _, _ := wGen.Generate(c.rand.Int(), c.tileSize, c.gravity, 0, 0, false)
		s = game.NewState(c.waitBeforePickupable, c.fallSpeedUp, c.player, currentLevel)
		currentLevel.SetTile(end.Coords.X+end.Room.InterestX, end.Coords.Z+end.Room.InterestZ, scenery.NewCaveEntranceDoor())

		// Only for starting position
		backMarker = marker()
		wGen.PlaceImportantEntityInRoom(s, start, backMarker)
		backMarker.Position().Z -= currentLevel.TileSize()/2

		nextMarker = marker()
		wGen.PlaceImportantEntityInRoom(s, end, nextMarker)
		nextMarker.Position().Z += currentLevel.TileSize()/2
	*/
    } else { // STANDARD LEVEL
		wGen := helpers.NewWorldGenerator(
			helpers.NewBasicWallRoomSet(o, o, w, h, x),
			&world.Room{
				Floor: [][]func() world.ITile{
					{h, w, w, w, x, w, w, w, h},
					{h, w, o, o, o, o, o, w, h},
					{h, w, o, o, o, o, o, w, h},
					{h, w, o, o, o, o, o, w, h},
					{h, w, w, w, w, w, w, w, h},
					{h, h, h, h, h, h, h, h, h},
					{h, h, h, h, h, h, h, h, h},
				},
				InterestX: 4,
				InterestZ: 4,
			},
			&world.Room{
				Floor: [][]func() world.ITile{
					{h, h, h},
					{w, x, w},
				},
				InterestX: 1,
				InterestZ: 1,
			},
			helpers.NonObstacleOverlapComparator,
			o, // connector
		)
		criticalPathCount := 1+int(math32.Round(c.rand.Float32()*float32(level*4)))
		extraRoomCount := 2+int(math32.Round(c.rand.Float32()*float32(level*3)))
		currentLevel, start, end, criticalPath, extra := wGen.Generate(c.rand.Int(), c.tileSize, c.gravity, criticalPathCount, extraRoomCount, true)

		// Override entry/exit tiles
		var backTile world.ITile
		if level == 0 {
			backTile = scenery.NewEntranceDoorWithWall()
		} else {
			backTile = scenery.NewBackDoorWithWall()
		}
		currentLevel.SetTile(start.Coords.X+start.Room.InterestX, start.Coords.Z+start.Room.InterestZ, backTile)
		currentLevel.SetTile(end.Coords.X+end.Room.InterestX, end.Coords.Z+end.Room.InterestZ, scenery.NewNextDoorWithWall())

		s = game.NewState(c.waitBeforePickupable, c.fallSpeedUp, c.player, currentLevel)
		entities :=	[]*helpers.EntitySet{}
		// Can do if level == ... { } but can also generalize using the min and max (e.g. level+1 means spawn also in level 0)
		entities = append(entities,
			helpers.NewEntitySet(2*level, 5*(level+1), func() entity.IEntity {
				e := enemies.NewSimpleSkeletonEnemy("skeleton")
				r := c.rand.Float32()

				{
				var item entity.IItem
				switch {
				case r < 0.05:
					item = equipables.NewCrystalStaff()
				case r < 0.20:
					item = equipables.NewRapier()
				default:
					item = equipables.NewIronSword()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}

				{
				var item entity.IItem
				switch {
				case r < 0.05:
					item = equipables.NewWoodenShield()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.18:
					item = humanoidequipables.NewLeatherHelmet()
				} 
				if item != nil {
					e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				return e
			}),
			helpers.NewEntitySet(0, 7*level, func() entity.IEntity {
				e := enemies.NewSimpleOrcEnemy("orc")
				r := c.rand.Float32()
				{
				var item entity.IItem
				switch {
				case r < 0.05:
					item = equipables.NewLongSword()
				default:
					item = equipables.NewIronSword()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.05:
					item = equipables.NewIronShield()
				case r < 0.30:
					item = equipables.NewWoodenShield()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.05:
					item = humanoidequipables.NewIronHelmet()
				case r < 0.35:
					item = humanoidequipables.NewLeatherHelmet()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.65:
					item = humanoidequipables.NewLeatherArmor()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.65:
					item = humanoidequipables.NewLeatherPants()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.65:
					item = humanoidequipables.NewLeatherBoots()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				return e
			}),
			helpers.NewEntitySet(0, 4*(level-1), func() entity.IEntity {
				e := enemies.NewSimpleOrcEnemy("orc berserker")
				r := c.rand.Float32()
				{
				var item entity.IItem
				switch {
				case r < 0.05:
					item = equipables.NewDemonbloodSword()
				default:
					item = equipables.NewLongSword()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.05:
					item = equipables.NewSpikeShield()
				case r < 0.30:
					item = equipables.NewReinforcedShield()
				default:
					item = equipables.NewIronShield()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.60:
					item = humanoidequipables.NewIronHelmet()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.95:
					item = humanoidequipables.NewIronArmor()
				default:
					item = humanoidequipables.NewLeatherArmor()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.95:
					item = humanoidequipables.NewIronPants()
				default:
					item = humanoidequipables.NewLeatherPants()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.60:
					item = humanoidequipables.NewIronBoots()
				default:
					item = humanoidequipables.NewLeatherBoots()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.30:
					item = humanoidequipables.NewIronGloves()
				default:
					item = humanoidequipables.NewLeatherGloves()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				return e
			}),
			helpers.NewEntitySet(0, 4*(level-2), func() entity.IEntity {
				e := enemies.NewSimpleSkeletonEnemy("skeleton king")
				r := c.rand.Float32()

				{
				var item entity.IItem
				switch {
				case r < 0.10:
					item = equipables.NewFinnSword()
				default:
					item = equipables.NewCrystalStaff()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}

				{
				var item entity.IItem
				switch {
				case r < 0.10:
					item = humanoidequipables.NewCrystalAmulet()
				case r < 0.25:
					item = humanoidequipables.NewGoldRing()
				case r < 0.50:
					item = humanoidequipables.NewCopperRing()
				default:
					item = humanoidequipables.NewIronRing()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.18:
					item = humanoidequipables.NewLeatherHelmet()
				} 
				if item != nil {
					e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				return e
			}),
			helpers.NewEntitySet(0, 4*level, func() entity.IEntity {
				e := enemies.NewSimpleZombieEnemy("zombie")
				r := c.rand.Float32()
				{
				var item entity.IItem

				switch {
				case r < 0.05:
					item = equipables.NewGrassSword()
				case r < 0.20:
					item = equipables.NewScarletSword()
				case r < 0.45:
					item = equipables.NewRootSword()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.80:
					item = humanoidequipables.NewShirt()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.98:
					item = humanoidequipables.NewKhakiTrousers()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				{
				var item entity.IItem
				switch {
				case r < 0.55:
					item = humanoidequipables.NewLeatherShoes()
				} 
				if item != nil {
				e.Inventory().AutoEquip(e.Inventory().PickUp(item))
				}
				}
				return e
			}),

			// helpers.NewEntitySetCountable(5, 5, 0, 1, func() entity.IEntity { return consumables.NewBanana() }),
		)

		wGen.Populate(s, append(criticalPath, extra...), entities)

		backMarker = marker()
		wGen.PlaceImportantEntityInRoom(s, start, backMarker)
		backMarker.Position().Z -= currentLevel.TileSize()/2

		nextMarker = marker()
		wGen.PlaceImportantEntityInRoom(s, end, nextMarker)
		nextMarker.Position().Z += currentLevel.TileSize()/2
	} 

	c.nextMarkers[level] = nextMarker
	c.backMarkers[level] = backMarker
	c.levels[level] = s
	if c.level == level {
		c.jumpToLevel(level)
	}
}

func (c *RogueController) nextLevel() {
	level := c.level+1
	if !c.hasLevel(level) {
		c.generateLevel(level)
	}
	c.jumpToLevel(level)

    c.player.Position().Copy(c.currentBackMarker.Position())
    c.player.Position().Z -= c.player.OuterRadius()*2
    c.player.SetLookAngle(+math32.Pi)
    c.player.Speed().X = 0
    c.player.Speed().Z = 0
}

func (c *RogueController) previousLevel() {
	level := c.level-1
	if !c.hasLevel(level) {
		c.generateLevel(level)
	}
	c.jumpToLevel(level)

    c.player.Position().Copy(c.currentNextMarker.Position())
    c.player.Position().Z += c.player.OuterRadius()*2
    c.player.SetLookAngle(0)
    c.player.Speed().X = 0
    c.player.Speed().Z = 0
}

func (c *RogueController) jumpToLevel(level int) {
	c.level = level
	c.currentLevel = c.levels[level]
	c.currentNextMarker = c.nextMarkers[level]
	c.currentBackMarker = c.backMarkers[level]
	c.wentBack = false
	c.wentNext = false
	c.farEnoughFromDoor = false
}

func (c *RogueController) over() bool {
	return c.player.Destroyed()
}

func (c *RogueController) Update(delta float32) {
	if !c.farEnoughFromDoor {
		if (c.currentNextMarker == nil || c.currentNextMarker.BorderDistanceTo(c.player) > c.currentLevel.World.TileSize()*c.tilesFarToEnterDoor) &&
		   (c.currentBackMarker == nil || c.currentBackMarker.BorderDistanceTo(c.player) > c.currentLevel.World.TileSize()*c.tilesFarToEnterDoor) {
			c.farEnoughFromDoor = true
		}
	}
	c.wentBack = false
	c.wentNext = false
	if c.currentNextMarker != nil && c.farEnoughFromDoor && c.currentNextMarker.BorderDistanceTo(c.player) < 0 {
		c.wentNext = true
	}
	if c.currentBackMarker != nil && c.farEnoughFromDoor && c.currentBackMarker.BorderDistanceTo(c.player) < 0 {
		c.wentBack = true
	}
	c.levels[c.level].Update(delta, c.input.Apply)
}

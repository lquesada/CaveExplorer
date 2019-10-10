package humanoid

import "fmt"
import "reflect"
import "cavernal.com/lib/g3n/engine/util/application"
import "cavernal.com/lib/g3n/engine/window"
import "cavernal.com/lib/g3n/engine/gui"
import "cavernal.com/lib/g3n/engine/math32"
import "cavernal.com/lib/g3n/engine/renderer"
import "cavernal.com/lib/g3n/engine/texture"
import "cavernal.com/lib/g3n/engine/material"
import "cavernal.com/lib/g3n/engine/camera"
import "cavernal.com/lib/g3n/engine/core"
import "cavernal.com/lib/g3n/engine/graphic"
import 	"cavernal.com/lib/g3n/engine/light"
import "cavernal.com/entity"

type HUD struct {
	gamePanel *gamePanel
	SideBar *SideBar
	slots map[entity.SlotId]bool

	player entity.IPlayer
	guiScale float32
	width int
	height int
	itemSlotsCount int
}

func NewHUD(player entity.IPlayer) *HUD {
	app := application.Get()

	app.Gui().RemoveAll(true)
	app.Gui().SetLayout(gui.NewDockLayout())
	gamePanel := newgamePanel()
	app.Gui().Add(gamePanel)
  	app.SetPanel3D(gamePanel)

   	g := &HUD{
   		gamePanel: gamePanel,
   		player: player,
   	}

    mainPanel := gui.NewPanel(0, 0)
    mainPanel.SetColor4(&math32.Color4{0.8, 0.8, 0.8, 1})
    mainPanel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockRight})
    mainPanel.SetLayout(gui.NewDockLayout())

	g.SideBar = &SideBar{
		items: make([]*itemBox, 0),
		mainPanel: mainPanel,
	}

   	app.Gui().Add(g.SideBar.Panel())
   	app.Window().SubscribeID(window.OnKeyDown, g, g.KeyboardInput)

   	return g
}

func (g *HUD) Panel() *gui.Panel {
	return g.gamePanel.Panel
}

func (g *HUD) KeyboardInput(evname string, ev interface{}) {
    kev := ev.(*window.KeyEvent)
    switch kev.Keycode {
    	case window.KeyE:
    		g.player.Inventory().Flip()
    }
}

func (g *HUD) Dispose() {
	app := application.Get()
	app.Window().UnsubscribeID(window.OnKeyDown, g)
	g.SideBar.Dispose()
	g.gamePanel.DisposeChildren(true)
	g.gamePanel.Dispose()
}

func (g *HUD) GatherInput() {
	s := g.SideBar
	app := application.Get()
	i := s.inventory
	shift := app.KeyState().Pressed(window.KeyLeftShift) || app.KeyState().Pressed(window.KeyRightShift)
	if g.gamePanel.PollSwitch() {
		item := i.HoldItem()
		if item != nil {
			// Drop this item.
			// Whether it was already equipped.
			whereSlot := i.WhereEquipped(item)
			if whereSlot != "" {
				i.DropEquipped(whereSlot)
			}
			// Or it was in items, then equip.
			whereItem := i.WhereHas(item)
			if whereItem != -1 {
				i.DropItemSlot(whereItem)
			} 
		}
	}
	for k, ib := range s.slots {
		item := i.Slots()[k]
		if ib.PollSwitch() {
			// A tap on a 0-item removes it.
			if item != nil && item.IsCountable() && item.Count() == 0 {
				i.RemoveEquipped(k)
			}
			if shift && item != nil {
				if k == AlternateLeftSlot || k == AlternateRightSlot {
					// Shift + click on alternate slots tries to move to hand.
					if k == AlternateLeftSlot {
						if !i.ShuffleEquipment(k, HandLeftSlot) {
							i.ShuffleEquipment(k, HandRightSlot)
						}
					} else if k == AlternateRightSlot {
						if !i.ShuffleEquipment(k, HandRightSlot) {
							i.ShuffleEquipment(k, HandLeftSlot)
						}
					}
				}  else {
					i.AutoUnequip(k)
				}
				continue
			}
			if i.HoldItem() == nil {
				// Hold this item.
				item := i.Slots()[k]
				if item != nil {
					i.SetHoldItem(item)
					i.SetHoldAmount(item.Count())
				}
				continue
			} else {
				// Try to equip it.
				item := i.HoldItem()
				if i.CanEquipItem(item, k) {
					// Whether it was already equipped, then shuffleequipment.
					whereSlot := i.WhereEquipped(item)
					if whereSlot != "" {
						i.ShuffleEquipment(i.WhereEquipped(item), k)
						i.SetHoldItem(nil)
						continue
					}
					// Or it was in items, then equip.
					whereItem := i.WhereHas(item)
					if whereItem != -1 {
						i.Equip(whereItem, k)
						i.SetHoldItem(nil)
						continue
					} 
				} else {
					// Revert drag&drop.
					i.SetHoldItem(nil)
					continue
				}
			}
			fmt.Println("TODO UNADRESSED SLOT", k)
		}
		if ib.PollDrop() {
			// A tap on a 0-item removes it.
			if item != nil && item.IsCountable() && item.Count() == 0 {
				i.RemoveEquipped(k)
			}
			if item != nil {
				i.DropEquipped(k)
			}
		}
		if ib.PollUse() {
			// A tap on a 0-item removes it.
			if item != nil && item.IsCountable() && item.Count() == 0 {
				i.RemoveEquipped(k)
			}
			// TODO use items
		}
	}
	for k, ib := range s.items {
		item := i.Items()[k]
		if ib.PollSwitch() {
			if shift {
				i.AutoEquip(k)
				continue
			}
			if i.HoldItem() == nil {
				// Hold this item.
				if item != nil {
					i.SetHoldItem(item)
					i.SetHoldAmount(item.Count())
				}
				continue
			} else {
				// Try to move it.
				// Whether it was already equipped, then unequip.
				whereSlot := i.WhereEquipped(i.HoldItem())
				if whereSlot != "" {
					i.Unequip(whereSlot, k)
					i.SetHoldItem(nil)
					continue
				}
				// Or it was in items, then try and shuffle
				whereItem := i.WhereHas(i.HoldItem())
				if whereItem != -1 {
					i.ShuffleItem(whereItem, k)
					i.SetHoldItem(nil)
					continue
				} 
			}
			fmt.Println("TODO UNADRESSED ITEM", k)
		}
		if ib.PollDrop() {
			if item != nil {
				i.DropItemSlot(k)
			}
		}
		if ib.PollUse() {
			// TODO use items
		}
	}
}

func (g *HUD) PreRender() {
	app := application.Get()

	// Draw the whole screen with no scene (i.e. the HUD)
    // Get framebuffer size and sets the viewport accordingly
    width, height := app.Window().FramebufferSize()
	if int32(width) < int32(g.rightPanelWidth()) {
		return
	}

    var minHUDScale float32 = 1.5
    var maxHUDScale float32 = float32(height)/g.minimumHeight()
	var autoHUDScale float32 = float32(-height+width)/g.basePanelWidth()
    var guiScale float32 = math32.Min(maxHUDScale, math32.Max(minHUDScale, autoHUDScale))

    newSlots := map[entity.SlotId]bool{}
	for k, _ := range g.player.Inventory().Slots() {
		newSlots[k] = true
	}
	itemSlotsCount := len(g.player.Inventory().Items())
    if g.width != width || g.height != height || g.guiScale != guiScale || !reflect.DeepEqual(newSlots, g.slots) || itemSlotsCount != g.itemSlotsCount {
	    g.width = width
	    g.height = height
	    g.guiScale = guiScale
	    g.slots = newSlots
        g.itemSlotsCount = itemSlotsCount
    	g.update()
    }
    g.SideBar.SetInventory(g.player.Inventory())

    app.Gl().Viewport(0, 0, int32(width), int32(height))
    // Sets the HUD root panel size to the size of the framebuffer
    if app.Gui() != nil {
    	app.Gui().SetSize(float32(width), float32(height))
    }

	app.SetScene(core.NewNode())
    app.Renderer().Render(app.Camera())

    // Gather and process mouse input from all itemboxes.
    g.GatherInput()

    // Draw the inventory on top of the HUD
    g.SideBar.Render()

	// Draw the scene but not the HUD area
    // Sets perspective camera aspect ratio
    aspect := float32(int32(width) - int32(g.rightPanelWidth())) / float32(height)
    app.Camera().SetAspect(aspect)

    app.Gl().Viewport(0, 0, int32(width)-int32(g.rightPanelWidth()), int32(height))
}

func (g *HUD) itemInventorySize() float32 {
	return math32.Floor(64*g.guiScale)
}

func (g *HUD) minimumHeight() float32 {
	return 1000 // without guiScale, for autoHUDScale
}
func (g *HUD) basePanelWidth() float32 {
	return 5*64 // without guiScale, for autoHUDScale
}

func (g *HUD) rightPanelWidth() float32 {
	return 5*g.itemInventorySize()
}

func (g *HUD) equippedBodyPanelWidth() float32 {
	return 3*g.itemInventorySize()
}

func (g *HUD) topPanelHeight() float32 {
	return 200*g.guiScale
}

func (g *HUD) update() {
	s := g.SideBar
	s.mainPanel.RemoveAll(true)
    s.mainPanel.SetWidth(g.rightPanelWidth())
	s.slots = make(map[entity.SlotId]*itemBox, len(g.slots))
	s.items = make([]*itemBox, g.itemSlotsCount)

	{
		equippedPanel := gui.NewPanel(0, 0)
    	equippedPanel.SetColor4(&math32.Color4{0.5, 0.5, 0.5, 1})
    	equippedPanel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockTop})
    	equippedPanel.SetHeight(5*g.itemInventorySize())
    	equippedPanel.SetLayout(gui.NewDockLayout())
	    s.mainPanel.Add(equippedPanel)
	    {
			equippedLeftPanel := gui.NewPanel(0, 0)
	    	equippedLeftPanel.SetColor4(&math32.Color4{0.3, 0.3, 0.3, 1})
	    	equippedLeftPanel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockCenter})
	    	equippedLeftPanel.SetLayout(gui.NewDockLayout())
		    equippedPanel.Add(equippedLeftPanel)
			{
				equippedWeaponPanel := gui.NewPanel(0, 0)
		    	equippedWeaponPanel.SetColor4(&math32.Color4{0.2, 0.2, 0.2, 1})
		    	equippedWeaponPanel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockCenter})
	    		equippedWeaponPanel.SetLayout(gui.NewGridLayout(1))
			    equippedLeftPanel.Add(equippedWeaponPanel)
				{
					equippedWeapon1Panel := gui.NewPanel(0, 0)
		    		equippedWeapon1Panel.SetColor4(&math32.Color4{0.1, 0.1, 0.1, 1})
		    		equippedWeapon1Panel.SetWidth(g.itemInventorySize()*2)
		    		equippedWeapon1Panel.SetHeight(g.itemInventorySize())
					equippedWeapon1Panel.SetLayout(gui.NewGridLayout(2))
			    	equippedWeaponPanel.Add(equippedWeapon1Panel)
		    		{
			    		s.slots[HandLeftSlot] = newEquippedBox(g.itemInventorySize(),)
		    			equippedWeapon1Panel.Add(s.slots[HandLeftSlot])
		    		} ; {
			    		s.slots[HandRightSlot] = newEquippedBox(g.itemInventorySize())
		    			equippedWeapon1Panel.Add(s.slots[HandRightSlot])
		    		}					    	
				} ; {
					equippedWeapon2Panel := gui.NewPanel(0, 0)
		    		equippedWeapon2Panel.SetColor4(&math32.Color4{0.1, 0.1, 0.1, 1})
		    		equippedWeapon2Panel.SetWidth(g.itemInventorySize()*2)
		    		equippedWeapon2Panel.SetHeight(g.itemInventorySize())
					equippedWeapon2Panel.SetLayout(gui.NewGridLayout(2))
			    	equippedWeaponPanel.Add(equippedWeapon2Panel)
		    		{
			    		s.slots[AlternateLeftSlot] = newEquippedBox(g.itemInventorySize())
		    			equippedWeapon2Panel.Add(s.slots[AlternateLeftSlot])
		    		} ; {
			    		s.slots[AlternateRightSlot] = newEquippedBox(g.itemInventorySize())
		    			equippedWeapon2Panel.Add(s.slots[AlternateRightSlot])
		    		}					    	
				}
			} ; {
				equippedStatsPanel := gui.NewPanel(0, 0)
		    	equippedStatsPanel.SetColor4(&math32.Color4{0.2, 0.2, 0.2, 1})
		    	equippedStatsPanel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockBottom})
    			equippedStatsPanel.SetWidth(g.itemInventorySize())
    			equippedStatsPanel.SetHeight(g.itemInventorySize())
			    equippedLeftPanel.Add(equippedStatsPanel)
			    {
			    		s.holdItemSlot = newHoldBox(g.itemInventorySize())
		    			equippedStatsPanel.Add(s.holdItemSlot)
			    }
		    }					
		} ; {
			equippedBodyPanel := gui.NewPanel(0, 0)
	    	equippedBodyPanel.SetColor4(&math32.Color4{0.3, 0.3, 0.3, 1})
	    	equippedBodyPanel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockRight})
	    	equippedBodyPanel.SetWidth(g.equippedBodyPanelWidth())
	    	equippedBodyPanel.SetLayout(gui.NewGridLayout(3))
		    equippedPanel.Add(equippedBodyPanel)
		    {
		    	s.slots[AmuletSlot] = newEquippedBox(g.itemInventorySize())
			    equippedBodyPanel.Add(s.slots[AmuletSlot])
		    } ;	{
		    	s.slots[HelmetSlot] = newEquippedBox(g.itemInventorySize())
		    	equippedBodyPanel.Add(s.slots[HelmetSlot])
		    } ; {
		    	s.slots[BackSlot] = newEquippedBox(g.itemInventorySize())
		    	equippedBodyPanel.Add(s.slots[BackSlot])
		    } ; {
		    	s.slots[RingLeftSlot] = newEquippedBox(g.itemInventorySize())
		    	equippedBodyPanel.Add(s.slots[RingLeftSlot])
		    } ; {
		    	s.slots[ArmorSlot] = newEquippedBox(g.itemInventorySize())
		    	equippedBodyPanel.Add(s.slots[ArmorSlot])
		    } ; {
		    	s.slots[RingRightSlot] = newEquippedBox(g.itemInventorySize())
		    	equippedBodyPanel.Add(s.slots[RingRightSlot])
		    } ; {
		    	s.slots[GlovesSlot] = newEquippedBox(g.itemInventorySize())
		    	equippedBodyPanel.Add(s.slots[GlovesSlot])
		    } ; {
		    	s.slots[PantsSlot] = newEquippedBox(g.itemInventorySize())
		    	equippedBodyPanel.Add(s.slots[PantsSlot])
		    } ; {
		    	s.slots[BootsSlot] = newEquippedBox(g.itemInventorySize())
		    	equippedBodyPanel.Add(s.slots[BootsSlot])
		    } ;
		} ; {
			equippedConsumablePanel := gui.NewPanel(0, 0)
	    	equippedConsumablePanel.SetColor4(&math32.Color4{0.7, 0.7, 0.7, 1})
	    	equippedConsumablePanel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockBottom})
	    	equippedConsumablePanel.SetHeight(2*g.itemInventorySize())
	    	equippedConsumablePanel.SetLayout(gui.NewGridLayout(5))
		    equippedPanel.Add(equippedConsumablePanel)
		    for _, k := range ConsumableSlotsSorted {
	    		s.slots[k] = newEquippedBox(g.itemInventorySize())
	    		equippedConsumablePanel.Add(s.slots[k])
		    }
		}
	} ; {
	    flexPanel := gui.NewPanel(0, 0)
	    flexPanel.SetColor4(&math32.Color4{0.4, 0.4, 0.4, 1})
	    flexPanel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockCenter})
	    flexPanel.SetHeight(g.topPanelHeight())
	    flexPanel.SetLayout(gui.NewDockLayout())
	    s.mainPanel.Add(flexPanel)
	    {
		    infoPanel := gui.NewPanel(0, 0)
		    infoPanel.SetColor4(&math32.Color4{0.4, 0.4, 0.4, 1})
		    infoPanel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockTop})
		    infoPanel.SetHeight(g.topPanelHeight())
		    infoPanel.SetLayout(gui.NewDockLayout())
		    flexPanel.Add(infoPanel)
		} ; {
			inventoryPanel := gui.NewPanel(0, 0)
	    	inventoryPanel.SetColor4(&math32.Color4{0.6, 0.6, 0.6, 1})
			inventoryPanel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockCenter})
			inventoryPanel.SetLayout(gui.NewDockLayout())
			flexPanel.Add(inventoryPanel)
			{
		 	   carryPanel := gui.NewPanel(0, 0)
			    carryPanel.SetColor4(&math32.Color4{0.2, 0.2, 0.2, 1})
		    	carryPanel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockCenter})
		    	carryPanel.SetLayout(gui.NewGridLayout(5))
		    	inventoryPanel.Add(carryPanel)
			    for i := 0; i < g.itemSlotsCount; i++ {
		    		s.items[i] = newItemBox(g.itemInventorySize())
		    		carryPanel.Add(s.items[i])
			    }
			}
		}
	}
}

type gamePanel struct {
	*gui.Panel

	switched bool
}

func newgamePanel() *gamePanel {
   	panel := gui.NewPanel(0, 0)
   	panel.SetRenderable(false)
   	panel.SetColor4(&math32.Color4{1, 1, 1, 1})
   	panel.SetLayoutParams(&gui.DockLayoutParams{Edge: gui.DockCenter})
   	gamePanel := &gamePanel{
   		Panel: panel,
   	}
	gamePanel.SubscribeID(window.OnMouseDown, gamePanel, func(name string, ev interface{}) { gamePanel.switched = true })
   	return gamePanel
}

func (g *gamePanel) Dispose() {
	g.UnsubscribeID(window.OnMouseDown, g)
	g.Panel.DisposeChildren(true)
	g.Panel.Dispose()
}

func (g *gamePanel) PollSwitch() bool {
	if g.switched {
		g.switched = false
		return true
	}
	return false
}

type SideBar struct {
	mainPanel	*gui.Panel
	slots map[entity.SlotId]*itemBox
	items []*itemBox
	holdItemSlot *itemBox

	inventory entity.IInventory
}

func (s *SideBar) Dispose() {
	for _, k := range s.slots {
		k.DisposeChildren(true)
		k.Dispose()
	}
	for _, k := range s.items {
		k.DisposeChildren(true)
		k.Dispose()
	}
	s.mainPanel.DisposeChildren(true)
	s.mainPanel.Dispose()
	s.holdItemSlot.Dispose()
}

func (s *SideBar) Panel() *gui.Panel {
	return s.mainPanel
}

func (s *SideBar) SetInventory(inventory entity.IInventory) {
	s.inventory = inventory
}

func (s *SideBar) Render() {
	for k, _ := range s.slots {
		s.drawItemFromSlot(k)
	}
	for i := 0; i < len(s.inventory.Items()); i++ {
		s.drawItemFromItems(s.items[i], i)
	}
	if s.inventory.HoldItem() != nil {
		s.holdItemSlot.drawItem(s.inventory.HoldItem(), s.inventory.HoldItem(), s.inventory.HoldAmount())
	}
}

func (s *SideBar) drawItemFromSlot(slot entity.SlotId) {
	ib := s.slots[slot]
	if ib == nil {
		return
	}
	if v, ok := s.inventory.Slots()[slot]; ok {
		ib.drawItem(v, s.inventory.HoldItem(), s.inventory.HoldAmount())
	} else {
		ib.Panel.SetColor(&math32.Color{0.5, 0.5, 0.5})
	}
}

func (s *SideBar) drawItemFromItems(ib *itemBox, id int) {
	ib.drawItem(s.inventory.Items()[id], s.inventory.HoldItem(), s.inventory.HoldAmount())
}

type itemBox struct {
	*gui.Panel
	renderer *renderer.Renderer
	camera *camera.Perspective
	axis *graphic.AxisHelper
	ambientLight     *light.Ambient
	directionalLight *light.Directional
	scene *core.Node

	equippedSlot bool
	holdSlot bool

	switched bool
	switchedPart bool
	used bool
	dropped bool
	mouseOver bool
}

func newItemBox(itemInventorySize float32) *itemBox {
   	panel := gui.NewPanel(itemInventorySize, itemInventorySize)

    app := application.Get()
	renderer := renderer.NewRenderer(app.Gl())
	renderer.AddDefaultShaders() // LOG this may return errors
	renderer.SetGuiPanel3D(panel)

	camera := camera.NewPerspective(65, 1, 0.01, 1000)
 	camera.SetPosition(1, 1, 1)
	camera.LookAt(&math32.Vector3{0, 0, 0})
	camera.SetAspect(1)
	camera.Project(&math32.Vector3{100, 1, 0.8})

	ambientLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.5)

	directionalLight := light.NewDirectional(math32.NewColor("white"), 0.5)
	directionalLight.SetPosition(-1, 0, 1)

	ib := &itemBox{
		Panel: panel,
		renderer: renderer,
		camera: camera,
		ambientLight: ambientLight,
		directionalLight: directionalLight,
		axis: graphic.NewAxisHelper(0.05), 
	}

	ib.SubscribeID(window.OnMouseDown, ib, ib.MouseInput)
	ib.SubscribeID(gui.OnCursorEnter, ib, ib.MouseInput)
	ib.SubscribeID(gui.OnCursorLeave, ib, ib.MouseInput)
   	app.Window().SubscribeID(window.OnKeyDown, ib, ib.KeyboardInput)

	return ib
}

func newEquippedBox(itemInventorySize float32) *itemBox {
	i := newItemBox(itemInventorySize)
	i.equippedSlot = true
	return i
}

func newHoldBox(itemInventorySize float32) *itemBox {
	i := newItemBox(itemInventorySize)
	i.holdSlot = true
	return i
}

func (ib *itemBox) MouseInput(name string, ev interface{}) {
	if name == window.OnMouseDown {
	    e := ev.(*window.MouseEvent)
	    switch e.Button {
	    case window.MouseButton1:
			ib.switched = true
	    case window.MouseButton2:
			ib.used = true
	    case window.MouseButton3:
			ib.switchedPart = true
		}
	} else if name == gui.OnCursorEnter {
        	ib.mouseOver = true
    } else if name == gui.OnCursorLeave {
	        ib.mouseOver = false
    }
} 

func (ib *itemBox) KeyboardInput(evname string, ev interface{}) {
  	if !ib.mouseOver {
  		return
  	}
    kev := ev.(*window.KeyEvent)
    switch kev.Keycode {
    	case window.KeyQ:
    		ib.dropped = true
    	case window.KeyE:
    		ib.used = true
    	case window.KeyZ:
    		ib.switched = true
    	case window.KeyX:
    		ib.switchedPart = true
    }
}

func (ib *itemBox) Dispose() {
    app := application.Get()
	ib.UnsubscribeID(window.OnMouseDown, ib)
	ib.UnsubscribeID(gui.OnCursorEnter, ib)
	ib.UnsubscribeID(gui.OnCursorLeave, ib)
	ib.UnsubscribeID(window.OnMouseDown, ib)
	app.Window().UnsubscribeID(window.OnKeyDown, ib)
	ib.Panel.DisposeChildren(true)
	ib.Panel.Dispose()
}

func (ib *itemBox) PollSwitch() bool {
	if ib.switched {
		ib.switched = false
		return true
	}
	return false
}

func (ib *itemBox) PollSwitchPart() bool {
	if ib.switchedPart {
		ib.switchedPart = false
		return true
	}
	return false
}

func (ib *itemBox) PollDrop() bool {
	if ib.dropped {
		ib.dropped = false
		return true
	}
	return false
}

func (ib *itemBox) PollUse() bool {
	if ib.used {
		ib.used = false
		return true
	}
	return false
}

func (ib *itemBox) drawItem(item, holdItem entity.IItem, holdAmount int) {
        app := application.Get()
	    _, height := app.Window().FramebufferSize()
	    if ib.scene != nil {
	    	ib.scene.Dispose()
	    }

        x := int32(ib.Pospix().X)
        y := int32(float32(height)-ib.Pospix().Y-float32(ib.Panel.Height()))
        w := int32(ib.ContentWidth())
        h := int32(ib.ContentHeight())
        app.Gl().Viewport(x, y, w, h)

        scene := core.NewNode()
        // Workaround: If I don't add any object here, the background is not drawn.
        scene.Add(ib.axis)
        scene.Add(ib.ambientLight)
        scene.Add(ib.directionalLight)

        if item != nil {
        	if (item != holdItem || item.IsCountable()) && (!item.IsCountable() || holdItem == nil || item.Count() > holdAmount) || ib.holdSlot {
	            if in := item.InventoryNode(); in != nil {
	            	scene.Add(in.G3NNode())
	            }
       		}
	        if item.IsCountable() {
	    		am := item.Count()
	    	    if item == holdItem {
	        		am = item.Count() - holdAmount
	       		}
	        	if ib.holdSlot {
	        		am = holdAmount
	        	}
	        	if ib.equippedSlot || am > 0 {
		    		// TODO this label looks super-off, fix
		    		amount := createSpriteZ(fmt.Sprintf("%d", am))
		    		amount.SetPosition(1.1, 0, 0)
	        		scene.Add(amount)
	        	}
	        }
	    }
        ib.scene = scene
        ib.renderer.SetScene(scene)
        ib.renderer.Render(ib.camera)
}

func createSpriteZ(text string) *core.Node {
    font := gui.StyleDefault().Font
    font.SetPointSize(48)
    font.SetColor(&math32.Color4{1, 1, 1, 1})
    width, height := font.MeasureText(text)
    img := font.DrawText(text)
    tex := texture.NewTexture2DFromRGBA(img)

    plane_mat := material.NewStandard(math32.NewColor("white"))
    plane_mat.AddTexture(tex)
    plane_mat.SetTransparent(true)

    div := float32(100)
    sprite := graphic.NewSprite(float32(width)/div, float32(height)/div, plane_mat)
    sprite.SetPosition(-float32(width)*0.666/div, -float32(height)/2/div, 0)
    node := core.NewNode()
    node.Add(sprite)
    return node
}

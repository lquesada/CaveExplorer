# CaveExplorer

Copyright (c) 2019, Luis Quesada Torres - https://github.com/lquesada | www.luisquesada.com

CaveExplorer is a voxel-based videogame proof-of-concept, written in Go and on the G3N Go 3D game engine.

# Binary

Check the releases page at https://github.com/lquesada/CaveExplorer/releases for binary Linux releases and videos.

# Features

   - Physics engine
     - Simplified cylinder-based collisions (center, radius, height)
     - Cross-frame collision-checks to avoid objects going through objects
   - Entity system
     - Basic entities: characters (player and enemies), items (equipables, consumables, others), and attacks
     - Humanoid: humanoid equipables, inventory, hud
   - Voxel-based models
     - Includes tooling to process and load voxel models designed with MagicaVoxel (https://ephtracy.github.io), 
     - Articulated and animated characters and equipable items
     - Node system that enables composing models with different transparency and tint values
   - Map generator
     - Start room -> critical path (n rooms) -> end room (+ extra N rooms)
     - Rooms are fully configurable, with connection points to other rooms and logic
   - Contents
     - Humanoid characters: human, orc, skeleton, and zombie
     - ~50 equipables (weapons, shields, armor, boots, helmets, amulets, rings, pants, etc.)
     - Basic story: get to the final level, grab the artifact, go all the way out
     
# Known issues

   - Need to update to newer versions of G3N, paying attention to propagating the custom fixes
   - There is a memory leak triggered when building G3N nodes, which happens thousands of times per frame. Gameplays of 30-60 minutes use up several GBs of RAM. Maybe the newer versions of g3n fix this
   - Attack and defense are not properly calculated (e.g. attack and defense in ammo items or equipped items other than the weapon are not be considered)
   - Missing lots of contents, e.g. other weapon types such as bows, wands; projectiles, etc.
   - The user interface for the inventory is incomplete (and ugly)
   - With some OpenGL setups, when pressing two keys at once (e.g. arrow left+arrow up), one of them takes some delay to go through
   - Many things are not properly tweaked, e.g.
     - Position of objects in inventory or when dropped so they don't appear floating
     - Strengh of enemies when pushing the player
   - Acute lack of documentation and presence of hacks/workarounds in the (very prototypish) code
   - The proof-of-concept game is, from a pure gaming perspective, NOT fun, and quite boring
     - Player is practically immortal
     - Story is bland
     - Enemy, player, and item values (defense, attack, etc.) aren't calibrated
     
# Structure

CaveExplorer is sort of a game engine built on top G3N, and an instantiation of the engine.

-----

G3N offers all the OpenGL logic.

Please note that the G3N library included at lib/g3n is a modified version of (an outdated version of) the G3N Go 3D game engine (https://github.com/g3n/engine):
   - Enable Viewport to work in panels with y > 0. there was no metric to be passed with "y"
   - Allow reading textures from []bytes so as to load models from arrays encoded in the binary
   - Allow translucent materials.

Please refer to lib/customfixes for the diffs.

-----

The following directories compose the CaveExplorer sort-of-game-engine:

   - ./entity: Defines and implements basic entity logic (e.g. characters, attacks, collisions, and items)
     - ./entity/humanoid: Defines and implement basic humanoid-specific logic (e.g. humanoid inventory, humanoid animations)
   - ./game: Implements the game state, basic game state logic (attacks hurting characters, picking up objects, etc.), and physics (collisions with walls, collisions with other entities, etc.)
   - ./helpers: Implements helpers to be used during asset and game instance definition
   - ./hud: Defines the interface for a HUD
   - ./input: Implements the input library for mouse and keyboard (including moving, angles, attacking, repeat attacks, jumping, etc.)
   - ./lib/\*: Contains modified libraries
     - ./lib/g3n: Contains the modified G3N 3D Go Game Engine library
   - ./model: Defines and implements model and node classes to work with MagicaVoxel nodes, compose them, and alter transparency and tint.
   - ./world: Defines and implements tiles, rooms, the world, and a world generator.

-----

The following files and directories compose an instantiation of the CaveExplorer sort-of-game-engine:

   - ./main.go: Contains all game logic for an instance, including how to generate the world, the player, and the enemies, what happens when moving a level forward or backward, what to put in the last level, when to exit the game, etc. Relies heavily on helpers and uses the assets listed below.
   - ./assets: Contains definitions and implementations for all entities in the game
     - ./assets/consumables: Consumable entities (unused, incomplete)
     - ./assets/enemies: Instantiation of enemies (orcs, skeletons, zombies) with proper values so they can be used in the game
     - ./assets/equipables: Objects equipable by characters in general (anything that fits in a hand/paw regardless of the size of the character)
     - ./assets/humanoid: Assets related to humanoid characters
     - ./assets/humanoid/customization: Customization features for humanoids (basically faces)
     - ./assets/humanoid/equipables: Humanoid-equipable items (armors, pants, helmets, boots, anything that fits a humanoid with a specific size)
     - ./assets/humanoid/human: Human-looking humanoid... Can use all humanoid equipables (armor, etc.) and equipables (weapons, shields, etc.)
     - ./assets/humanoid/orc: Orc-looking humanoid
     - ./assets/humanoid/skeleton: Skeleton-looking humanoid
     - ./assets/humanoid/zombie: Zombie-looking humanoid
     - ./assets/miscitems: Miscellaneous items (unused, incomplete)
     - ./assets/player: Instantiation of a human with proper values and plumbed to the input library so it can be used in the game
     - ./assets/scenery: Tiles, walls, doors, etc.
   - ./sources: Contains the MagicaVoxel source files and exported files used by the assets

# How to build or run from the code
First, download and install all dependencies:

```
$ go get github.com/lquesada/cavernal

```

First regenerate all models from MagicaVoxel source and exported files. Please refer to regenmodels.sh or https://github.com/lquesada/voxel-3d-models/tree/master/tools/voxobjrename for details on this step.

```
$ bash regenmodels.sh
```

Then compile and run:

```
$ go run main.go
```

or

```
$ go build main.go
$ ./main
```

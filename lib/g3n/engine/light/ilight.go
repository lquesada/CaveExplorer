// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"cavernal.com/lib/g3n/engine/core"
	"cavernal.com/lib/g3n/engine/gls"
)

// ILight is the interface that must be implemented for all light types.
type ILight interface {
	RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int)
}

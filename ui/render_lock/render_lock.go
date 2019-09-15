package render_lock

import (
	"github.com/gizak/termui/v3"
	"sync"
)

type RenderLockStruct struct {
	mu sync.Mutex
}

var Lock = RenderLockStruct{}

func RenderLock(items ...termui.Drawable) {
    Lock.mu.Lock()
	termui.Render(items...)
    Lock.mu.Unlock()
}

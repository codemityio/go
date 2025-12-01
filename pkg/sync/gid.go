package sync

import (
	"bytes"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	fallbackIDCounter atomic.Int64 //nolint:gochecknoglobals // it is a global counter
	fallbackIDs       sync.Map     //nolint:gochecknoglobals // map[goroutineID string] int64
)

// GID returns a goroutine ID or a fallback unique ID per goroutine.
func GID() int64 {
	var buf [64]byte

	stack := runtime.Stack(buf[:], false)

	// try to parse "goroutine 123 ["
	line := string(buf[:stack])

	if strings.HasPrefix(line, "goroutine ") {
		idField := strings.Fields(line)[1]

		if id, err := strconv.ParseInt(idField, 10, 64); err == nil {
			return id
		}
	}

	// fallback: assign unique ID per goroutine using stack signature
	stackSig := string(bytes.TrimSpace(buf[:stack]))

	if val, ok := fallbackIDs.Load(stackSig); ok {
		res, _ := val.(int64)

		return res
	}

	id := fallbackIDCounter.Add(1)

	fallbackIDs.Store(stackSig, id)

	return id
}

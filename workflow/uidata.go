package workflow

const (
	xPadding = 20
	yPadding = 20
	xGap     = 245
	yGap     = 125
)

// uidata represents all uidata of a workflow.
type uidata = map[string]uidatum

// uidatum represents the position of an object.
type uidatum struct {
	XPos int64 `plist:"xpos,omitempty"`
	YPos int64 `plist:"ypos,omitempty"`
}

func (i *Info) buildUIData() {
	i.UIData = make(uidata)

	// A map of object depths to object UIDs at that depth
	depthMap := make(map[int64][]string)

	for _, obj := range i.Objects {
		uid := obj["uid"].(string)
		depth := findDepth(uid, i.Connections)
		depthMap[depth] = append(depthMap[depth], uid)
	}

	for depth, uids := range depthMap {
		for idx, uid := range uids {
			xpos := int64(xPadding + depth*xGap)
			ypos := int64(yPadding + idx*yGap)

			i.UIData[uid] = uidatum{
				XPos: xpos,
				YPos: ypos,
			}
		}
	}
}

func findDepth(uid string, conns map[string][]Connection) int64 {
	pointers := make([]string, 0)

	for connUID, conns := range conns {
		for _, conn := range conns {
			if conn.To == uid {
				pointers = append(pointers, connUID)
			}
		}
	}

	depth := int64(0)

	for _, pointer := range pointers {
		pointerDepth := findDepth(pointer, conns)
		if pointerDepth >= depth {
			depth = pointerDepth + 1
		}
	}

	return depth
}

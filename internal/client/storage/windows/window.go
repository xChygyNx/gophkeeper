package windows

type Window float32

const (
	WindowSwitcherWidth  Window = 200
	WindowSwitcherHeight Window = 80
	WindowAuthWidth      Window = 500
	WindowAuthHeight     Window = 100
	WindowMainWidth      Window = 1400
	WindowMainHeight     Window = 300
)

func (w Window) Size() float32 {
	return float32(w)
}

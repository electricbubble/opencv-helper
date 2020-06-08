package opencv_helper

type DebugMode int

const (
	// DmOff no output
	DmOff DebugMode = iota
	// DmEachMatch output matched and mismatched values
	DmEachMatch
	// DmNotMatch output only values that do not match
	DmNotMatch
)

var debug = DmOff

const dmOutputMsg = `[DEBUG] The current value is '%.4f', the expected value is '%.4f'`

func Debug(dm DebugMode) {
	debug = dm
}

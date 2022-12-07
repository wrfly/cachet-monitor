package cachet

type State int

const (
	StateUnknown State = iota
	StateOperational
	StatePerformanceIssues
	StatePartialOutage
	StateMajorOutage
)

func (s State) String() string {
	switch s {
	case StateOperational:
		return "Operational"
	case StatePerformanceIssues:
		return "PerformanceIssues"
	case StatePartialOutage:
		return "PartialOutage"
	case StateMajorOutage:
		return "MajorOutage"

	default:
		return "unknown"
	}
}

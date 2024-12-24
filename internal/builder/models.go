package builder

type (
	Workflow struct {
		Name        string
		Manager     Config
		Transitions []Transition
	}
	Config struct {
		Expr    string
		Func    string
		Type    string
		Package string
	}
	Transition struct {
		Name   string
		Dist   string
		From   string
		Guards []Config
	}
)

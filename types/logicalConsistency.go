package types

type LogicalConsistency struct {
	Checked                  bool
	ContradictionsExaminable bool
	Range                    NoneHighRange
}

// TODO: implement conversion/parsing functions

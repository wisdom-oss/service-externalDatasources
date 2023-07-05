package types

type LogicalConsistency struct {
	Checked                  bool
	ContradictionsExaminable bool
	Range                    NoneHighRange
}

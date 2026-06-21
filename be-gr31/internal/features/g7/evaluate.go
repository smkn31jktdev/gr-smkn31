package g7

import (
	g7model "be-gr31/internal/model/g7"
	"be-gr31/internal/features/g7/content/evaluate"
)

// EvalResult untuk backward compat
type EvalResult = evaluate.EvalResult

// EvalReport untuk backward compat
type EvalReport = evaluate.EvalReport

// EvaluateJurnals delegates ke evaluate.Jurnals
func EvaluateJurnals(jurnals []g7model.G7) EvalReport {
	return evaluate.Jurnals(jurnals)
}

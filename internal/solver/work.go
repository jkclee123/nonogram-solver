package solver

import "nonogram-solver/types"

// WorkType enumerates the kinds of work handled by workers
type WorkType int

const (
	WorkOverlap WorkType = iota
	WorkCrossReference
)

// WorkItem represents a unit of work in the queue
type WorkItem struct {
	Type   WorkType
	LineID types.LineID
	Color  int // color id for targeted operations
}

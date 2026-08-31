package engine

import "context"

// Graph is implemented by the graph package.
type Graph interface {
	Initialize(ctx context.Context, state *InterviewState) error
	ProcessTurn(ctx context.Context, state *InterviewState) error
	GenerateFinalEvaluation(ctx context.Context, state *InterviewState) error
}
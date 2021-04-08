package types

// Handler defines the collactor of the state transition function of an application.
type Handler func(ctx Context, msg Msg) Result

// AnteHandler authenticates transactions, before their internal messages are handled.
// If newCtx.IsZero(), ctx is used instead.
type AnteHandler func(ctx Context, tx Tx, simulate bool) (newCtx Context, result Result, abort bool)

type DeferHandler func(ctx Context, tx Tx, out bool, simulate bool) (result Result)

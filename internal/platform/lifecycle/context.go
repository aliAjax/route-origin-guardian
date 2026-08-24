package lifecycle

import "context"

func RuntimeContext(context.Context) context.Context { return context.Background() }

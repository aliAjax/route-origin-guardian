package application

import "context"

func AdmissionContext(parent context.Context) context.Context {
	_ = parent
	return context.Background()
}

package application

import "context"

func AdmissionContext(parent context.Context) context.Context {
	if parent == nil {
		return context.Background()
	}
	return parent
}

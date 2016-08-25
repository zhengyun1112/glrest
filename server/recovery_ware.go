package server

import (
	"net/http"
	"runtime"
	"fmt"
	"github.com/zhengyun1112/glrest/logger"
)

// RecoveryWare is the recovery middleware which can cover the panic situation.
type RecoveryWare struct {
	printStack bool
	stackAll   bool
	stackSize  int
}

// ServeHTTP implements the Middleware interface, just recover from the panic. Would provide information on the web page
// if in debug mode.
func (m *RecoveryWare) ServeHttp(w http.ResponseWriter, r *http.Request, pathParams Params, next RawHandler) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			stack := make([]byte, m.stackSize)
			stack = stack[:runtime.Stack(stack, m.stackAll)]
			logger.Error("PANIC: %s\n%s", err, stack)
			if m.printStack {
				fmt.Fprintf(w, "PANIC: %s\n%s", err, stack)
			}
		}
	}()

	next(w, r, pathParams)
}

// NewRecoveryWare returns a new recovery middleware. Would log the full stack if enable the printStack.
func NewRecoveryWare(printStack, stackAll bool) PreMiddleware {
	return &RecoveryWare{
		printStack: printStack,
		stackAll:   stackAll,
		stackSize:  1024 * 8,
	}
}

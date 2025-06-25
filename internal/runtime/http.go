package runtime

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MichelLacerda/nox/internal/token"
)

var mux = http.NewServeMux()

func NewHttpModule() *MapInstance {
	return NewMapInstance(map[string]any{
		"route": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(i *Interpreter, args []any) any {
				route, ok1 := args[0].(string)
				if !ok1 {
					i.Runtime.ReportRuntimeError(nil, "First argument to http.route must be a string (path)")
					return nil
				}

				switch handler := args[1].(type) {
				case *Function:
					mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
						safeHttpHandlerCall(i, w, r, handler)
					})

				case *Instance:
					mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
						methodFn := handler.Get(&token.Token{Lexeme: strings.ToLower(r.Method)})

						if methodFn == nil {
							http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
							return
						}

						if callable, ok := methodFn.(Callable); ok {
							safeHttpHandlerCall(i, w, r, callable)
						} else {
							http.Error(w, "Handler method is not callable", http.StatusInternalServerError)
						}
					})

				default:
					i.Runtime.ReportRuntimeError(nil, "Second argument to http.route must be a function or class instance")
				}

				return nil
			},
		},

		"serve": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(i *Interpreter, args []any) any {
				port, ok := args[0].(float64)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "http.serve expects a port number")
					return nil
				}

				server := &http.Server{
					Addr:    fmt.Sprintf(":%d", int(port)),
					Handler: mux,
				}

				fmt.Printf("Starting HTTP server on port %d\n", int(port))
				err := server.ListenAndServe()
				if err != nil {
					i.Runtime.ReportRuntimeError(nil, "Server error: "+err.Error())
				}
				return nil
			},
		},
	})
}

func safeHttpHandlerCall(i *Interpreter, w http.ResponseWriter, r *http.Request, fn Callable) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(w, "Internal Server Error: ", r)
		}
	}()

	reqMap := NewMapInstance(map[string]any{
		"method": r.Method,
		"url":    r.URL.String(),
		"query": NewMapInstance(func() map[string]any {
			query := make(map[string]any)
			for k, v := range r.URL.Query() {
				query[k] = strings.Join(v, ",")
			}
			return query
		}()),
		"params": NewMapInstance(func() map[string]any {
			params := make(map[string]any)
			for k, v := range r.URL.Query() {
				params[k] = strings.Join(v, ",")
			}
			return params
		}()),
		"headers": NewMapInstance(func() map[string]any {
			headers := make(map[string]any)
			for k, v := range r.Header {
				headers[k] = strings.Join(v, ",")
			}
			return headers
		}()),
		"body": func() string {
			if r.Body == nil {
				return ""
			}
			defer r.Body.Close()
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				i.Runtime.ReportRuntimeError(nil, "Failed to read request body: "+err.Error())
				return ""
			}
			return string(bodyBytes)
		}(),
	})

	respMap := NewMapInstance(map[string]any{
		"write": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(_ *Interpreter, args []any) any {
				if len(args) == 1 {
					fmt.Fprint(w, args[0])
				}
				return nil
			},
		},
		"set_header": &BuiltinFunction{
			ArityValue: 2,
			CallFunc: func(_ *Interpreter, args []any) any {
				key, ok1 := args[0].(string)
				value, ok2 := args[1].(string)
				if ok1 && ok2 {
					w.Header().Set(key, value)
				}
				return nil
			},
		},
		"set_status": &BuiltinFunction{
			ArityValue: 1,
			CallFunc: func(_ *Interpreter, args []any) any {
				status, ok := args[0].(float64)
				if !ok {
					i.Runtime.ReportRuntimeError(nil, "http.set_status expects a status code")
					return nil
				}
				w.WriteHeader(int(status))
				return nil
			},
		},
	})

	result := fn.Call(i, []any{reqMap, respMap})
	if str, ok := result.(string); ok {
		fmt.Fprint(w, str)
	}
}

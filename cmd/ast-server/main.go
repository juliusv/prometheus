package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/prometheus/prometheus/promql"
)

func main() {
	listenAddr := flag.String("listen-addr", ":8000", "Web API listen address.")
	http.HandleFunc("/parse", func(w http.ResponseWriter, r *http.Request) {
		expr, err := promql.ParseExpr(r.FormValue("expr"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing expression: %v", err.Error()), http.StatusBadRequest)
			return
		}
		buf, err := json.Marshal(expr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling AST: %v", err.Error()), http.StatusBadRequest)
			return
		}
		w.Write(buf)
	})
	http.ListenAndServe(*listenAddr, nil)
}

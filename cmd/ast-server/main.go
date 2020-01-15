package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/prometheus/prometheus/promql"
	prom_httputil "github.com/prometheus/prometheus/util/httputil"
)

func main() {
	listenAddr := flag.String("listen-addr", ":8000", "Web API listen address.")
	promURL := flag.String("prometheus-url", "http://demo.robustperception.io:9090/", "Prometheus server to forward API queries to.")
	flag.Parse()

	prometheusURL, err := url.Parse(*promURL)
	if err != nil {
		log.Fatalln("Error parsing Prometheus proxy URL:", err)
	}

	http.HandleFunc("/parse", func(w http.ResponseWriter, r *http.Request) {
		expr, err := promql.ParseExpr(r.FormValue("expr"))
		if err != nil {
			errJSON, err := json.Marshal(map[string]string{"type": "error", "message": fmt.Sprintf("Expression incomplete or buggy: %v", err)})
			if err != nil {
				http.Error(w, fmt.Sprintf("Error marshaling error JSON: %v", err), http.StatusInternalServerError)
				return
			}
			http.Error(w, string(errJSON), http.StatusBadRequest)
			return
		}
		regex, err := regexp.Compile("^(?:.*)$")
		if err != nil {
			panic(err)
		}
		prom_httputil.SetCORS(w, regex, r)
		buf, err := json.Marshal(translateAST(expr))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling AST: %v", err), http.StatusBadRequest)
			return
		}
		w.Write(buf)
	})
	http.Handle("/api/v1/", httputil.NewSingleHostReverseProxy(prometheusURL))
	http.ListenAndServe(*listenAddr, nil)
}

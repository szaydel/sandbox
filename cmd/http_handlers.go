package main

import (
	"fmt"
	"net/http"
	"strings"
)

func prometheusMetricsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := metricsReport.All()
	if err != nil {
		handleErr(err, false)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError,
		)
		return
	}
	for _, l := range data {
		fmt.Fprint(w, l)
	}
	
}

func allInfoHandler(w http.ResponseWriter, r *http.Request) {
	data, err := metricsReport.ToJSON()
	if err != nil {
		handleErr(err, false)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError,
		)
		return
	}
	fmt.Fprint(w, string(data))
}

func roleInfoHandler(w http.ResponseWriter, r *http.Request) {
	role := strings.TrimPrefix(r.URL.Path, "/info/")
	data, err := metricsReport.RoleToJSON(role)
	if err != nil {
		handleErr(
			fmt.Errorf(
				"failed getting metrics for %s with: %s", role, err),
			false,
		)
		if err == errNoInfoForRole {
			http.NotFound(w, r)
		} else {
			http.Error(w,
				http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError,
			)
		}
		return
	}
	fmt.Fprint(w, string(data))
}

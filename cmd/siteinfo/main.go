//
// main.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/render"
	"github.com/weirdgiraffe/watdatcloud"
	"github.com/weirdgiraffe/watdatcloud/aws"
	"github.com/weirdgiraffe/watdatcloud/azure"
	"github.com/weirdgiraffe/watdatcloud/gcp"
)

var addr = flag.String("addr", "127.0.0.1:8080", "address to listen on")

func main() {
	flag.Parse()

	l := watdatcloud.NewRangeLookuper(
		aws.NewAWS(),
		azure.NewAzure(),
		gcp.NewGCP(),
	)
	err := l.UpdateRanges()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Route("/{addr}", func(r chi.Router) {
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			addr, err := net.LookupHost(chi.URLParam(r, "addr"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			ret := []watdatcloud.Range{}
			for i := range addr {
				res, err := l.Lookup(addr[i])
				if err != nil {
					if err.(*watdatcloud.RangeNotFound) != nil {
						continue
					}
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				ret = append(ret, res)
			}
			err = json.NewEncoder(w).Encode(ret)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		})
	})
	log.Printf("run on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, r))
}

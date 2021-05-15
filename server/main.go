package main

import "github.com/KMSIGN/nsd-hack/server/app/api"

func main() {
	api.Configure()
	if err := api.ListenAndServe(); err != nil {
		panic(err)
	}
}

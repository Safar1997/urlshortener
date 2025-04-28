package main

import (
	"net/http"
	"testing"
)

func Test_postHandler(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		r      *http.Request
		urlMap map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postHandler(tt.args.w, tt.args.r, tt.args.urlMap)
		})
	}
}

func Test_getHandler(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		r      *http.Request
		urlMap map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getHandler(tt.args.w, tt.args.r, tt.args.urlMap)
		})
	}
}

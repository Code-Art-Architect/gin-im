package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type R struct {
	Code  int    `json:"code,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Data  any    `json:"data,omitempty"`
	Rows  any    `json:"rows,omitempty"`
	Total any    `json:"total,omitempty"`
}

func Response(w http.ResponseWriter, code int, msg string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	r := R{
		Code: code,
		Rows: data,
		Msg:  msg,
	}
	ans, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
	_, _ = w.Write(ans)
}

func Fail(w http.ResponseWriter, msg string) {
	Response(w, 500, msg, nil)
}

func OK(w http.ResponseWriter, msg string) {
	Response(w, 200, msg, nil)
}

func OKList(w http.ResponseWriter, data any, total any) {
	List(w, 200, data, total)
}

func List(w http.ResponseWriter, code int, data any, total any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	r := R{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	ans, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
	_, _ = w.Write(ans)
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

type App struct {
	Router 		*mux.Router
	Middlewares *Middleware
}

type shortenReq struct {
	URL string `json:"url" validate:"nonzero"`
}

type shortResp struct {
	Shortlink string `json:"shortlink"`
}

func (a *App) Initialize() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Router = mux.NewRouter()
	a.Middlewares = &Middleware{}
	a.initalizeRouters()
}

func (a *App) initalizeRouters() {
	//a.Router.HandleFunc("/api/shorten", a.createShortlink).Methods("POST")
	//a.Router.HandleFunc("/api/info", a.getShortlinkInfo).Methods("GET")
	//a.Router.HandleFunc("/{shortlink:[a-zA-Z0-9]{1,11}}", a.redirect).Methods("GET")
	m := alice.New(a.Middlewares.LoggingHandler,a.Middlewares.recoverHandler)
	a.Router.Handle("/api/shorten",m.ThenFunc(a.createShortlink)).Methods("POST")
	a.Router.Handle("/api/info",m.ThenFunc(a.getShortlinkInfo)).Methods("GET")
	a.Router.Handle("/{shortlink:[a-zA-Z0-9]{1,11}}",m.ThenFunc(a.redirect)).Methods("GET")
}

func (a *App) createShortlink(w http.ResponseWriter, r *http.Request)  {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ResponseError(w, StatusError{http.StatusBadRequest,fmt.Errorf("parse params failed %v",req)})
		return
	}
	if err := validator.Validate(req); err != nil {
		ResponseError(w, StatusError{http.StatusBadRequest,fmt.Errorf("validate params failed %v",req)})
		return
	}
	defer r.Body.Close()

	record, err := getRecordByUrl(req.URL)
	if err != nil {
		ResponseError(w,err)
		return
	} else if record.Id != 0 {
		ResponseJson(w,http.StatusOK,shortResp{Shortlink:record.Key})
		return
	}

	lastId, err := getRecordLastId()
	if err != nil {
		ResponseError(w,err)
		return
	}
	newId := lastId + 1
	key := Base62encode(int(newId))
	var link *Link = new(Link)
	link.Key = key
	link.Url = req.URL
	_, er := addRecord(link)
	if er != nil  {
		ResponseError(w,er)
		return
	}
	ResponseJson(w,http.StatusOK,shortResp{Shortlink:key})
}

func (a *App) getShortlinkInfo(w http.ResponseWriter, r *http.Request)  {
	vals := r.URL.Query()
	s := vals.Get("shortlink")
	record, err := getRecordByKey(s)
	if err != nil {
		ResponseError(w,err)
		return
	}
	ResponseJson(w, http.StatusOK,record)
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	key := vars["shortlink"]
	record, err := getRecordByKey(key)
	if err != nil {
		ResponseError(w,err)
		return
	}
	http.Redirect(w,r,record.Url,http.StatusTemporaryRedirect) 

}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr,a.Router))
}

func ResponseError(w http.ResponseWriter, err error)  {
	switch e := err.(type) {
	case Error:
		log.Printf("http %d - %s",e.Status(),e.Error())
		ResponseJson(w,e.Status(),e.Error())
	default:
		ResponseJson(w,http.StatusInternalServerError,http.StatusText(http.StatusInternalServerError))
	}
}

func ResponseJson(w http.ResponseWriter, code int, payload interface{})  {
	resp,_ := json.Marshal(payload)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
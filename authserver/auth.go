package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ErrDumpingRequest   = errors.New("error dumping new request %v")
	ErrUnmarshalRequest = errors.New("error unmarshal request %v")
	ErrProvider         = errors.New("error getting new provider %v")
	ErrInvalidState     = errors.New("error, got invalid state got %s expected %s")
	ErrIncorrectCode    = errors.New("error invalid code, %v")
	ErrInvalidToken     = errors.New("error invalid token %v")
	ErrExtractingClaims = errors.New("could not extract claims %v")
)

var sampleState = "state-test"

type server struct {
	clientId     string
	clientSecret string
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oConfig      *oauth2.Config
}

type authRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Email        string `json:"email"`
}

type claims struct {
	Email    string   `json:"email"`
	Verified bool     `json:"verified"`
	Groups   []string `json:"groups"`
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {

	if _, err := w.Write([]byte("welcome to demo service")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

}
func (s *server) login(w http.ResponseWriter, r *http.Request) {
	var loginReq *authRequest

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf(ErrDumpingRequest.Error(), err), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(reqBody, &loginReq); err != nil {
		http.Error(w, fmt.Sprintf(ErrUnmarshalRequest.Error(), err), http.StatusBadGateway)
		return
	}

	s.clientId = loginReq.ClientID
	s.clientSecret = loginReq.ClientSecret

	provider, err := oidc.NewProvider(r.Context(), "http://127.0.0.1:5556/dex")
	if err != nil {
		http.Error(w, fmt.Sprintf(ErrProvider.Error(), err), http.StatusInternalServerError)
		return
	}

	oConfig := &oauth2.Config{
		ClientID:     loginReq.ClientID,
		ClientSecret: loginReq.ClientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	s.oConfig = oConfig

	authURL := oConfig.AuthCodeURL(sampleState)

	http.Redirect(w, r, authURL, http.StatusSeeOther)

}

func (s *server) callback(w http.ResponseWriter, r *http.Request) {
	var claims claims
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	if state != sampleState {
		http.Error(w, fmt.Sprintf(ErrInvalidState.Error(), state, sampleState), http.StatusUnauthorized)
		return
	}

	oauthToken, err := s.oConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf(ErrIncorrectCode.Error(), err), http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok {
		http.Error(w, errors.New("invalid conversion for id_token").Error(), http.StatusUnauthorized)
		return
	}

	idToken, err := s.verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
		return
	}

	if err := idToken.Claims(&claims); err != nil{
		http.Error(w, ErrExtractingClaims.Error(), http.StatusUnauthorized)
		return
	}

}

func main() {

	s := server{}
	router := mux.NewRouter()
	router.HandleFunc("/", s.index)
	router.HandleFunc("/login", s.login)

	httpServer := &http.Server{
		Handler: router,
		Addr:    ":5555",
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

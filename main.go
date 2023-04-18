package main

import (
	"OPNsenseProxyAPI/opnsense"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type syncAliasesRequest struct {
	Host    string   `json:"host"`
	Aliases []string `json:"aliases"`
}

var apiKey string
var apiSecret string
var address string
var domainName string

func main() {
	apiKey = os.Getenv("API_KEY")
	apiSecret = os.Getenv("API_SECRET")
	address = os.Getenv("OPNSENSE_ADDRESS")
	domainName = os.Getenv("DOMAIN_NAME")
	if apiKey == "" {
		log.Fatalf("API_KEY not set")
	}
	if apiSecret == "" {
		log.Fatalf("API_SECRET not set")
	}
	if address == "" {
		log.Fatalf("OPNSENSE_ADDRESS not set")
	}
	if domainName == "" {
		log.Fatalf("DOMAIN_NAME not set")
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))
	r.Post("/sync", handleSyncAliasesRequest)
	log.Infof("Running API on port 9657")
	http.ListenAndServe(":9657", r)
}

func handleSyncAliasesRequest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request syncAliasesRequest
	err := decoder.Decode(&request)
	if err != nil {
		log.Errorf("Error while decoding sync request: %v", err)
	}
	opnsenseClient := opnsense.NewClient(address, apiKey, apiSecret)
	// check if host exists
	exists, err := opnsenseClient.DoesHostOverrideExist(request.Host)
	if err != nil {
		log.Errorf("Error while checking if host override exists: %v", err)
	}
	if !exists {
		hostIP, err := getIPAddress(r)
		if err != nil {
			log.Errorf("Error while extracting host IP: %v", err)
		}
		hostname := strings.Replace(request.Host, fmt.Sprintf(".%v", domainName), "", -1)
		log.Infof("%v does not exist. Creating host override with hostname (%v), domain (%v) and IP (%v)", request.Host, hostname, domainName, hostIP)
		hostOverride := opnsense.NewHostOverride(hostname, domainName, hostIP)
		_, err = opnsenseClient.CreateHostOverride(hostOverride)
		if err != nil {
			log.Fatalf("Error while creating host override: %v", err)
		}
	}
	// sync aliases
	_, err = opnsenseClient.SyncAliases(request.Host, request.Aliases, domainName)
	if err != nil {
		log.Fatalf("Error while syncing alias overrides: %v", err)
	}
}

func getIPAddress(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return ip, nil
}

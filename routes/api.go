package routes

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/proc"
	"golang.org/x/net/context"
	"net/http"
	"time"
	"strings"
)

type response struct {
	Code   string `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
	Proc   proc.Process `json:"process,omitempty"`
}

type responseError struct {
	Code   string `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
	Cause  string `json:"cause,omitempty"`
}

type ListRequest struct {
	Code string 		`json:"code,omitempty"`
	Daemon []proc.Daemon	`json:"daemon,omitempty"`
	Docker []proc.Docker	`json:"docker,omitempty"`
	Service []proc.Service	`json:"service,omitempty"`
}

type ShooterRequest struct {
	Name string
	ServiceType int
	KillExecutor bool
}

func handleShooter(ctx context.Context, w http.ResponseWriter, r *http.Request) *HttpError {
	var timeStart = time.Now()
	log.Debugf("routes.api.handleShooter - START '%v'", timeStart)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var t ShooterRequest
	err := decoder.Decode(&t)
	if err != nil {
		log.Warnf("routes.api.handleShooter - '%v'", err.Error())
		json.NewEncoder(w).Encode(responseError{Code: "400", Status: "Invalid params", Cause: err.Error()})
		return nil
	}
	if strings.EqualFold(t.Name, "") {
		log.Warnf("routes.api.handleShooter - Name field is mandatory and should not be empty")
		json.NewEncoder(w).Encode(responseError{Code: "400", Status: "Invalid params", Cause: "Name field is mandatory and should not be empty"})
		return nil
	}

	var p = proc.Manager{}
	p.ConfigManager()
	p.LoadProcesses()
	proc, err := p.Shooter(t.Name, t.ServiceType, t.KillExecutor)
	switch {
		case proc != nil:
			json.NewEncoder(w).Encode(response{Code: "200", Status: "Succes", Proc: proc})
		case err == nil:
			log.Warning("routes.api.handleShooter - 'TASK NOT FOUND'")
			json.NewEncoder(w).Encode(responseError{Code: "404", Status: "Empty result", Cause: "TASK NOT FOUND"})
		default:
			log.Errorf("routes.api.handleShooter - '%v'", err.Error())
			json.NewEncoder(w).Encode(responseError{Code: "500", Status: "Server error", Cause: err.Error()})
	}

	log.Debugf("routes.api.handleShooter - FINISH - '%v'", time.Since(timeStart))
	return nil
}

func handleList(ctx context.Context, w http.ResponseWriter, r *http.Request) *HttpError {
	var timeStart = time.Now()
	log.Debugf("routes.api.handleList - START - '%v'", timeStart)

	var p = proc.Manager{}
	p.ConfigManager()

	var err error
	if err = p.LoadProcesses(); err == nil{
		json.NewEncoder(w).Encode(ListRequest{Code: "200", Daemon: p.Daemons, Docker: p.Dockers, Service: p.Services})
		return nil
	}

	json.NewEncoder(w).Encode(responseError{Code: "500", Status: "Server error", Cause: err.Error()})
	log.Debugf("routes.api.handleList - FINISH - '%v'", time.Since(timeStart))
	return nil
}

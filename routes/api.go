package routes

import (
	"encoding/json"
	"golang.org/x/net/context"
	"net/http"
	"github.com/stratio/valkiria/proc"
	"strconv"
	log "github.com/Sirupsen/logrus"
	"time"
)

type response struct {
	Status string `json:"status,omitempty"`
	Code string `json:"code,omitempty"`
}

type responseError struct {
	Status string `json:"status,omitempty"`
	Cause string  `json:"cause,omitempty"`
	Code string `json:"code,omitempty"`
}

func handleChaos(ctx context.Context, w http.ResponseWriter, r *http.Request) *HttpError {

	log.Debug("routes.api.handleDaemon")

	daemon, err := strconv.Atoi(r.URL.Query()["daemon"][0])
	service, err := strconv.Atoi(r.URL.Query()["service"][0])
	docker, err := strconv.Atoi(r.URL.Query()["docker"][0])

	if err == nil {
		if ! proc.IsSessionLock() {
			proc := proc.Processes{}
			proc.LoadProcesses()
			err = proc.Chaos(daemon, service, docker)
			if err != nil {
				log.Debug("routes.api.handleDaemon - ERROR:" + err.Error())
				json.NewEncoder(w).Encode(response{Status: "error: " + err.Error()})
			} else {

				json.NewEncoder(w).Encode(response{Code: "200", Status: "succes"})
			}
		} else {
			json.NewEncoder(w).Encode(response{Code: "200", Status: "session locked"})
		}
	} else {
		log.Debugf("routes.api.handleDaemon - ERROR: '%v'", err.Error())
	}
	return nil
}

func handleShooter(ctx context.Context, w http.ResponseWriter, r *http.Request) *HttpError {
	var timeStart = time.Now()
	log.Debugf("routes.api.handleShooter - START '%v'", timeStart)
	//01. Validate params
	name := r.URL.Query().Get("name")
	typeService, err := strconv.Atoi(r.URL.Query().Get("typeService"))
	killExecutor, err := strconv.ParseBool(r.URL.Query().Get("killExecutor"))
	if err != nil {
		log.Infof("routes.api.handleShooter - ERROR: '%v'", err.Error())
		json.NewEncoder(w).Encode(responseError{Cause: "Invalid params.", Code: "400", Status: "Error: " + err.Error()})
	} else {
		//02. Load process
		var p = proc.Processes{}
		p.LoadProcesses()
		b, err := p.Shooter(name, typeService, killExecutor)
		if b {
			json.NewEncoder(w).Encode(response{Code: "200", Status: "Succes"})
		} else {
			if err == nil {
				log.Info("routes.api.handleShooter - WARNNING: 'Service not found'")
				json.NewEncoder(w).Encode(response{Code: "404", Status: "Service not found"})
			} else {
				log.Infof("routes.api.handleShooter - ERROR: '%v'", err.Error())
				json.NewEncoder(w).Encode(responseError{Cause: "Server error.", Code: "500", Status: "Error: " + err.Error()})
			}
		}
	}
	log.Debugf("routes.api.handleShooter - FINISH - '%v'", time.Since(timeStart))
	return nil
}


func handleList(ctx context.Context, w http.ResponseWriter, r *http.Request) *HttpError {
	var timeStart = time.Now()
	log.Debugf("routes.api.handleList - START - '%v'", timeStart)
	//01. Validate params

	//02. Load process
	var p = proc.Processes{}
	err := p.LoadProcesses()

	//03. Paint processes in log
	if err == nil {
		for _, d := range p.Daemons {
			log.Infof("routes.api.handleList daemon - '%v' '%v' '%v' -", d.Name, d.Pid, d.Path)
		}
		for _, s := range p.Services {
			log.Infof("routes.api.handleList service - '%v' '%v' '%v' '%v' '%v' -", s.TaskName, s.Pid, s.Name, s.Executor, s.Ppid)
		}
		for _, do := range p.Dockers {
			log.Infof("routes.api.handleList docker - '%v' '%v' '%v' -", do.TaskName, do.Image, do.Name)
		}
		json.NewEncoder(w).Encode(response{Code: "200", Status: "succes"})
	} else {
		json.NewEncoder(w).Encode(responseError{Cause: "Server error.", Code: "500", Status: "Error: " + err.Error()})
	}
	log.Debugf("routes.api.handleList - FINISH - '%v'", time.Since(timeStart))
	return nil
}
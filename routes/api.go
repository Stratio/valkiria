package routes

import (
	"encoding/json"
	//log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
	"strings"
	"github.com/Stratio/valkiria/manager"
	"github.com/Stratio/valkiria/plugin"
)

const(
	code200 = "200"
	code400 = "400"
	invalidParams = "Invalid params"
	mandatoryField = "Name field is mandatory and should not be empty"
)

type DefaultErrorResponse struct{
	Code string 		`json:"code,omitempty"`
	Cause string 		`json:"cause,omitempty"`
	Status string 		`json:"status,omitempty"`
}

type KillResponse struct {
	Code string 		`json:"code,omitempty"`
	Process []plugin.Process`json:"killProcess,omitempty"`
	Errors []string 	`json:"errors,omitempty"`
}

type ListResponse struct {
	Code string 		`json:"code,omitempty"`
	Process []plugin.Process`json:"findProcess,omitempty"`
	Errors []string 		`json:"errors,omitempty"`
}

type KillRequest struct {
	Name string
	Properties string
}

func handleShooter(ctx context.Context, w http.ResponseWriter, r *http.Request) *HttpError {
	//.
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var sr KillRequest
	err := decoder.Decode(&sr)
	if err != nil {
		json.NewEncoder(w).Encode(DefaultErrorResponse{Code: code400, Status: invalidParams, Cause: err.Error()})
		return nil
	}
	if strings.EqualFold(sr.Name, "") {
		json.NewEncoder(w).Encode(DefaultErrorResponse{Code: code400, Status: invalidParams, Cause: mandatoryField})
		return nil
	}
	//.
	var p = manager.NewManager()
	procs, errs := p.Shooter(sr.Name, sr.Properties)
	//.
	json.NewEncoder(w).Encode(KillResponse{Code: code200, Process: procs, Errors: parseErrors(errs)})

	return nil
}

func handleList(ctx context.Context, w http.ResponseWriter, r *http.Request) *HttpError {
	//.
	var manager = manager.NewManager()
	//.
	procs, errs := manager.Read()
	//.
	json.NewEncoder(w).Encode(ListResponse{Code: code200, Process: procs, Errors: parseErrors(errs)})
	return nil
}

func parseErrors(errs []error)(ausErrors []string){
	for _, e := range errs{
		if e != nil{
			ausErrors = append(ausErrors, e.Error())
		}
	}
	return
}




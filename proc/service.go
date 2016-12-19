package proc

import (
	log "github.com/Sirupsen/logrus"
	procinfo "github.com/c9s/goprocinfo/linux"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"time"
	"regexp"
)

const (
	procDirectory = "/proc/"
	statusFile    = "/status"
	abc           = "abcdefghijklmnopqrstuvwxyz"
	mesos         = "mesos"
)

type Service struct {
	Pid            uint64
	Name           string
	TaskName       string
	Ppid           int64
	Executor       bool
	ChaosTimeStamp int64
}

func (s *Service) Kill() (err error) {
	log.Debug("proc.service.Kill")
	s.ChaosTimeStamp = time.Now().UTC().UnixNano()
	log.Infof("proc.service.Kill - '%v' '%v' '%v' '%v' '%v'", s.Pid, s.Name, s.Ppid, s.TaskName, s.ChaosTimeStamp)
	err = syscall.Kill(int(s.Pid), 9)
	if err != nil {
		log.Infof("proc.service.Kill - ERROR: '%v'", err.Error())
	}
	return
}

func ReadAllChildProcess(daemons []Daemon, daemonList []string, blackList []string)(aux []Service, err error){

	if files, err := ioutil.ReadDir(procDirectory); err == nil {
		for _, file := range files {
			if !strings.ContainsAny(file.Name(), abc) {
				status, err := procinfo.ReadProcessStatus(procDirectory + file.Name() + statusFile)
				if err == nil {
					link, _ := os.Readlink(procDirectory + file.Name() + "/cwd")
					validatePath, _ := regexp.Compile("^/var/lib/mesos/slave/slaves/.*/frameworks/.*/executors/.*/runs/.*")
					if !isInBlackList(status.Name, blackList) && validatePath.MatchString(link) {
						splitTaskName := strings.Split(link, "/")
						taskName := splitTaskName[10]
						s := Service{Pid: status.Pid, Name: status.Name, Ppid: status.PPid, TaskName: taskName}
						d, _ := ReadAllDaemons(daemonList)
						if status.PPid == 1 || (len(d) > 0 && status.PPid == int64(d[0].Pid)){
								s.Executor = true
						}
						aux = append(aux, s)
						log.Debugf("proc.service.ReadAllServices - append - '%v' '%v' '%v' '%v'", taskName, status.Pid, status.Name, status.PPid)

					}
				}
			}

		}
	}
	return
}

func isInBlackList(name string, blackListServices []string) (res bool) {
	// log.Debug("proc.service.isInBlackList")
	for _, blackService := range blackListServices {
		if strings.Compare(name, blackService) == 0 {
			res = true
		}
	}
	// log.Debugf("proc.service.isInBlackList - '%v' blackList '%v'", name, res)
	return
}

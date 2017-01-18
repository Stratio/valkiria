package mesos

import(
	procinfo "github.com/c9s/goprocinfo/linux"
	"io/ioutil"
	"os"
	"strings"
	"regexp"
	"github.com/Stratio/valkiria/plugin"
	"github.com/Stratio/valkiria/proc"
	//log "github.com/Sirupsen/logrus"
)

const (
	procDirectory = "/proc/"
	statusFile    = "/status"
	abc           = "abcdefghijklmnopqrstuvwxyz"
)

func (m *MesosConfig) GetServices() (func ()([]plugin.Process, error)){
	return func ()([]plugin.Process, error){
		return ReadAllChildProcess(m.DaemonListForChildServices, m.BlackListServices)
	}
}

func ReadAllChildProcess(daemonList []string, blackList []string)(aux []plugin.Process, err error){
	if files, err := ioutil.ReadDir(procDirectory); err == nil {
		for _, file := range files {
			if !strings.ContainsAny(file.Name(), abc) {
				status, err := procinfo.ReadProcessStatus(procDirectory + file.Name() + statusFile)
				if err == nil {
					link, _ := os.Readlink(procDirectory + file.Name() + "/cwd")
					validatePath, _ := regexp.Compile("^/var/lib/mesos/slave/slaves/.*/frameworks/.*/executors/.*/runs/.*")
					if !isInBlackList(status.Name, blackList) && validatePath.MatchString(link) {
						splitTaskName := strings.Split(link, "/")
						frameWorkId := splitTaskName[8]
						taskName := splitTaskName[10]
						s := proc.Service{Pid: status.Pid, Name: status.Name, Ppid: status.PPid, KillName: taskName, FrameWorkName: frameWorkId}
						d, _ := ReadAllDaemons(daemonList)
						if status.PPid == 1 || (len(d) > 0 && status.PPid == int64(d[0].(*proc.Daemon).Pid)){
							s.Executor = true
						}
						aux = append(aux, &s)

					}
				}
			}

		}
	}
	return
}

func isInBlackList(name string, blackListServices []string) (res bool) {
	for _, blackService := range blackListServices {
		if strings.Compare(name, blackService) == 0 {
			res = true
		}
	}
	return
}


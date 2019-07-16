package topz

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/shirou/gopsutil/process"
)

func handleError(res http.ResponseWriter, err error) {
	res.WriteHeader(http.StatusInternalServerError)
	res.Write([]byte(err.Error()))
}

type ProcInfo struct {
	PID 			int32
	MemoryPercent	float32
	MemoryInfo		*process.MemoryInfoStat
	CPUPercent		float64
	Command			string
}

func HandleTopz(res http.ResponseWriter, req *http.Request) {
	// list of PIDs of processes running now
	pids, err := process.Pids()

	if err != nil {
		handleError(res, err)
	}

	// here we'll keep data for each PID
	procs := make([]*ProcInfo, len(pids))

	// wait for a collection of goroutines to finish
	wg := sync.WaitGroup{}
	wg.Add(len(pids))

	for ix := range pids {
		proc, err := process.NewProcess(pids[ix])
		if err != nil {
			handleError(res, err)
			return
		}

		// gather data about each process
		go func(i int) {
			var err error

			p := &ProcInfo{}
			p.PID = pids[i]

			if p.Command, err = proc.Cmdline(); err != nil {
				log.Printf("Error getting Command line: %v", err)
			}
			if p.MemoryInfo, err = proc.MemoryInfo(); err != nil {
				log.Printf("Error getting memory info: %v", err)
			}
			if p.MemoryPercent, err = proc.MemoryPercent(); err != nil {
				log.Printf("Error getting Memory: %v", err)
			}
			if p.CPUPercent, err = proc.Percent(100 * time.Millisecond); err != nil {
				log.Printf("Error getting CPU: %v", err)
			}
			if len(p.Command) > 0 {
				procs[i] = p
			}

			wg.Done()
		}(ix)
	}

	wg.Wait()

	res.WriteHeader(http.StatusOK)
	w := tabwriter.NewWriter(res, 0, 0, 1, ' ', 0)

	for ix := range procs {
		proc := procs[ix]
		if proc == nil {
			continue
		}
		fmt.Fprintf(w, "%d\t%g\t%g\t%s\n", proc.PID, proc.CPUPercent, proc.MemoryPercent, proc.Command)
	}

	w.Flush()
}

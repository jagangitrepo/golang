package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

type RequestType int

const (
	RequestType_Scan = iota + 1
	RequestType_Quit
	RequestType_Print
)

type Request struct {
	request_type RequestType
	data         interface{}
}

type Task struct {
	id        int
	worker_id int
	request   Request
}

type Worker struct {
	name  string
	taskQ chan Task
	done  chan bool
}

func (w *Worker) Push(task Task) {
	w.taskQ <- task
}

func (w *Worker) ProcessQ() {
	exit_loop := false
	for !exit_loop {
		select {
		case task := <-w.taskQ:
			{
				log.Printf("Processing from workerid: %s task: %+v", w.name, task.request)
			}
		case <-w.done:
			log.Println("Got close signal workerid: ", w.name)
			exit_loop = true
			break
		}
	}
	log.Println("Exited workerid: ", w.name)
}

type Workers struct {
	workers       map[int]Worker
	taskQ         chan Task
	no_of_workers int
	wg            *sync.WaitGroup
}

func (w *Workers) Init() {
	for i := 1; i <= w.no_of_workers; i++ {
		worker := Worker{
			fmt.Sprintf("Worker%d", i),
			make(chan Task),
			make(chan bool),
		}
		w.workers[i] = worker
		go worker.ProcessQ()
		// w.wg.Add(i)
	}
}

func (w *Workers) Push(tasks ...Task) {
	for _, task := range tasks {
		wrkr := w.workers[task.worker_id]
		wrkr.Push(task)
	}
}

func SubmitWorkerTask(wrkr *Workers) {
	ticker := time.NewTicker(time.Millisecond * 10)
	timer := time.NewTimer(time.Second * 3)
	defer ticker.Stop()
	defer timer.Stop()

	for {
		select {
		case <-ticker.C:
			wrkr.Push(Task{
				rand.Intn(1000),
				(rand.Intn(3) % wrkr.no_of_workers) + 1,
				Request{
					RequestType_Scan, "Scan Request",
				},
			})
		case exit_time := <-timer.C:
			_ = exit_time
			wrkr.CloseAllWorkers()
			return
		}
	}
}

func (w *Workers) WaitandListen() {
	w.wg.Wait()
}

func (w *Workers) CloseAllWorkers() {
	for _, wrk := range w.workers {
		wrk.done <- true
		// w.wg.Done()
	}
}

func RunWorkers() {
	wrkrs := Workers{
		make(map[int]Worker),
		make(chan Task),
		3,
		&sync.WaitGroup{},
	}
	wrkrs.Init()

	go SubmitWorkerTask(&wrkrs)
	wrkrs.WaitandListen()

	// time.Sleep(time.Second * 6)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	// Singal Handler
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println("Received Signal: ", sig)
		wrkrs.CloseAllWorkers()
		time.Sleep(time.Millisecond * 500)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
func main() {

	flag.Parse()
	if *cpuprofile != "" {
		log.SetOutput(ioutil.Discard)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		RunWorkers()
		pprof.StopCPUProfile() // Stop profiling here
	} else {
		RunWorkers()
	}
}

// Output:
// 2024/10/11 20:07:07 Processing from workerid: Worker2 task: {request_type:1 data:Scan Request}
// 2024/10/11 20:07:07 Processing from workerid: Worker1 task: {request_type:1 data:Scan Request}
// 2024/10/11 20:07:07 Processing from workerid: Worker3 task: {request_type:1 data:Scan Request}
// 2024/10/11 20:07:07 Processing from workerid: Worker1 task: {request_type:1 data:Scan Request}

// Received Signal:  interrupt
// 2024/10/11 20:07:07 Got close signal workerid:  Worker1
// 2024/10/11 20:07:07 Exited workerid:  Worker1
// 2024/10/11 20:07:07 Got close signal workerid:  Worker3
// 2024/10/11 20:07:07 Exited workerid:  Worker3
// 2024/10/11 20:07:07 Got close signal workerid:  Worker2
// 2024/10/11 20:07:07 Exited workerid:  Worker2
// exiting

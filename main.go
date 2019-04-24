package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan sched)             // broadcast channel

type sched struct {
	Time            string
	Gomaxprocs      int
	Idleprocs       int
	Threads         int
	Spinningthreads int
	Idlethreads     int
	Runqueue        int
	Gcwaiting       int
	Nmidlelocked    int
	Stopwait        int
	Sysmonwait      int
	Ps              []p
	Ms              []m
	Gs              []g
}

type p struct {
	ID          int
	Status      int
	Schedtick   int
	Syscalltick int
	M           int
	Runqsize    int
	Gfreecnt    int
}

type m struct {
	ID         int
	P          int
	Curg       int
	Mallocing  int
	Throwing   int
	Preemptoff string
	Locks      int
	Dying      int
	Spinning   bool
	Blocked    bool
	Lockedg    int
}

type g struct {
	ID            int
	Status        int
	StatusMessage string
	M             int
	Lockedm       int
}

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	schedChan := make(chan sched)

	go parseStdin(schedChan)

	http.Handle("/", http.FileServer(http.Dir("./frontend/dist")))
	http.HandleFunc("/ws", makeHandleConnections(schedChan))

	log.Fatal(http.ListenAndServe("localhost:6061", nil))
}

func makeHandleConnections(schedChan chan sched) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		// Make sure we close the connection when the function returns
		defer ws.Close()

		// Register our new client
		clients[ws] = true

		for s := range schedChan {
			for client := range clients {
				err := client.WriteJSON(s)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}

	}

}

func parseStdin(schedChan chan sched) {
	scanner := bufio.NewScanner(os.Stdin)

	var schedule sched

	for scanner.Scan() {
		split := strings.Split(strings.TrimLeft(scanner.Text(), " "), ":")

		firstsched := false

		switch string(split[0][0]) {
		case "S":
			if firstsched != true {
				schedChan <- schedule
				fmt.Printf("SCHED %v\n", schedule)
				firstsched = true
			}
			var svals []int
			schedTime := strings.Split(split[0], " ")[1]

			schedVals := strings.Split(strings.TrimLeft(split[1], " "), " ")

			for _, val := range schedVals {
				i, err := strconv.Atoi(strings.Split(val, "=")[1])
				if err != nil {
					log.Fatal(err)
				}
				svals = append(svals, i)
			}

			schedule = sched{
				Time:            schedTime,
				Gomaxprocs:      svals[0],
				Idleprocs:       svals[1],
				Threads:         svals[2],
				Spinningthreads: svals[3],
				Idlethreads:     svals[4],
				Runqueue:        svals[5],
				Gcwaiting:       svals[6],
				Nmidlelocked:    svals[7],
				Stopwait:        svals[8],
				Sysmonwait:      svals[9],
				Ps:              []p{},
				Ms:              []m{},
				Gs:              []g{},
			}
		case "P":

			var pvals []int
			id, err := strconv.Atoi(split[0][1:])
			if err != nil {
				log.Fatal(err)
			}

			vals := strings.Split(strings.TrimLeft(split[1], " "), " ")

			for _, val := range vals {
				i, err := strconv.Atoi(strings.Split(val, "=")[1])
				if err != nil {
					log.Fatal(err)
				}
				pvals = append(pvals, i)
			}

			process := p{
				ID:          id,
				Status:      pvals[0],
				Schedtick:   pvals[1],
				Syscalltick: pvals[2],
				M:           pvals[3],
				Runqsize:    pvals[4],
				Gfreecnt:    pvals[5],
			}

			schedule.Ps = append(schedule.Ps, process)
		case "M":
			var mvals []int
			var spinning bool
			var blocked bool
			var preemptoff string
			id, err := strconv.Atoi(split[0][1:])
			if err != nil {
				log.Fatal(err)
			}

			vals := strings.Split(strings.TrimLeft(split[1], " "), " ")

			for j, val := range vals {
				switch j {
				case 4:
					preemptoff = strings.Split(val, "=")[1]
				case 7:
					spinning, err = strconv.ParseBool(strings.Split(val, "=")[1])
					if err != nil {
						log.Fatal(err)
					}
				case 8:
					blocked, err = strconv.ParseBool(strings.Split(val, "=")[1])
					if err != nil {
						log.Fatal(err)
					}
				default:
					i, err := strconv.Atoi(strings.Split(val, "=")[1])
					if err != nil {
						log.Fatal(err)
					}
					mvals = append(mvals, i)
				}
			}

			thread := m{
				ID:         id,
				P:          mvals[0],
				Curg:       mvals[1],
				Mallocing:  mvals[2],
				Throwing:   mvals[3],
				Preemptoff: preemptoff,
				Locks:      mvals[4],
				Dying:      mvals[5],
				Spinning:   spinning,
				Blocked:    blocked,
				Lockedg:    mvals[6],
			}

			schedule.Ms = append(schedule.Ms, thread)
		case "G":

			var gvals []int
			var statusMessage string

			id, err := strconv.Atoi(split[0][1:])
			if err != nil {
				log.Fatal(err)
			}

			s := strings.TrimLeft(split[1], " ")

			var re = regexp.MustCompile(`\(.*\)`)
			statusMessage = strings.TrimSuffix(strings.TrimPrefix(re.FindString(s), "("), ")")
			s = re.ReplaceAllString(s, "")

			vals := strings.Split(s, " ")

			for _, val := range vals {

				i, err := strconv.Atoi(strings.Split(val, "=")[1])
				if err != nil {
					log.Fatal(err)
				}
				gvals = append(gvals, i)
			}

			goroutine := g{
				ID:            id,
				Status:        gvals[0],
				StatusMessage: statusMessage,
				M:             gvals[1],
				Lockedm:       gvals[2],
			}

			schedule.Gs = append(schedule.Gs, goroutine)
		}
		//fmt.Printf("%v\n", schedule) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

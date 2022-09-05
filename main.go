package main

import (
	"log"
	"math"
	"net/rpc"
	"os"
	"time"
)

const (
	SAMPLES_SIZE = 10000
	ORI          = "A"
	DEST         = "E"
)

func main() {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatal("undefined HOST")
	}

	client, err := rpc.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	benchmark(client, ORI, DEST)
}

func benchmark(client *rpc.Client, ori string, dest string) {
	var samples []time.Duration
	for i := 0; i < SAMPLES_SIZE; i++ {
		log.Printf("sending request to find the shortest path between %s and %s", ori, dest)

		args := ShortestPathArgs{
			Ori:  ori,
			Dest: dest,
		}
		reply := ShortestPathReply{}
		start := time.Now()
		client.Call("Graph.ShortestPath", args, &reply)
		rtt := time.Since(start)

		samples = append(samples, rtt)
		log.Printf("shortest path received %v", reply.Path)
	}

	var mean float64
	for _, sample := range samples {
		mean += float64(sample)
	}
	mean = mean / float64(len(samples))

	var sd float64
	for _, sample := range samples {
		sd += math.Pow((float64(sample) - mean), 2)
	}
	sd = math.Sqrt(sd / float64(len(samples)))

	log.Printf("average RTT is %.2f (+- %.2f)", mean, sd)
}

type ShortestPathArgs struct {
	Ori  string
	Dest string
}

type ShortestPathReply struct {
	Path []string
}

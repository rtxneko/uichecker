package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"sync"
)

func main() {
	filename := flag.String("f", "ip.txt", "Filename for list of IPs")
	threads := flag.Int("t", 100, "Number of threads")
	outfile := flag.String("o", "hits.txt", "name of output file")
	remoteAddress := flag.String("a", "", "dokodemo-door remote address")
	remotePort := flag.String("p", "", "dokodemo-door remote port")
	_, err := os.Create(*outfile)
	if err != nil {
		log.Panicln("Could not create file for write: ", err)
	}
	of, err := os.OpenFile(*outfile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Panicln("Could not open file for write: ", err)
	}
	wg := sync.WaitGroup{}
	flag.Parse()
	var semaphore = make(chan struct{}, *threads)
	f, err := os.Open(*filename)
	if err != nil {
		log.Fatalln("Could not open IP list: ", err)
	}
	input := bufio.NewScanner(f)
	for input.Scan() {
		wg.Add(1)
		go brute(input.Text(), semaphore, &wg, of, *remoteAddress, *remotePort)
	}
	wg.Wait()
	err = of.Sync()
	if err != nil {
		log.Panicln(err)
	}
	of.Close()
}

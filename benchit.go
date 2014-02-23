// Benchit - run benchmarks and pipe output to a file

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"
)

var (
	description = flag.String("descr", "Benchmark", "description of the benchmark.")
)

func runBenchMark(s string, d string) {
	os.Chdir(s)
	out, _ := exec.Command("bash", "-c", "go test -run=Xxx -bench=.").Output()
	r, _ := regexp.Compile("Benchmark")
	if r.Match(out) {
		t := time.Now()
		dir := "benchmarks"
		os.MkdirAll(dir, 0755)
		filename := dir + "/" + d + t.Format("20060102150405") + ".txt"
		fmt.Printf("%s", out)
		_ = ioutil.WriteFile(filename, out, 0600)
	}
	if s != "." {
		os.Chdir("..")
	}
}

func main() {
	flag.Parse()
	cmd_string := `ls -al | \
                       awk '$1 ~ /^d/ && $NF ~ /(^[.]$|^[a-zA-Z0-9]+$)/ \
                       {print $NF}'`
	cmd := exec.Command("bash", "-c", cmd_string)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		runBenchMark(scanner.Text(), *description)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	cmd.Wait()
}

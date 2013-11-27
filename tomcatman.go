/*
TODO:
  Check if instance is running from pidfile
  Enable multiple strings of commands
  Figure out how to change the absolute path to the config file
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var start *string
var stop *string
var restart *string
var path = map[string]string{}

func main() {
	configure()
	parseArgs()
	run()
}

// configure reads the config file to determine which flags should be available.
func configure() {
	file, err := os.Open("/path/to/tomcatman.conf")
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			fmt.Fprintln(os.Stderr, "Error: please ensure tomcatman.conf exists and the path is configured correctly.")
			os.Exit(1)
		}
		log.Fatal(err)
	}
	var data []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = strings.Split(scanner.Text(), "=")
		// Assign path["instance"] = "/path/to/instance"
		path[data[0]] = data[1]
	}
}

// parseArgs retrieves and verifies the command line arguments.
func parseArgs() {
	// Set up flags
	start = flag.String("start", "instance", "start the given Tomcat instance")
	stop = flag.String("stop", "instance", "stop the given Tomcat instance")
	restart = flag.String("restart", "instance", "restart the given Tomcat instance")
	flag.Parse()

	if *start == "instance" && *stop == "instance" && *restart == "instance" {
		fmt.Fprintln(os.Stderr, "Error processing args. Run tomcatman -h to see valid options.")
	}
}

// run chooses the correct branch of execution by calling execute on the proper instance.
func run() {
	if *start != "instance" {
		if _, ok := path[*start]; ok {
			execute("start", path[*start])
		} else {
			log.Fatal("Error: instance name \"", *start, "\" not found in tomcatman.conf.")
		}
	}
	if *stop != "instance" {
		if _, ok := path[*stop]; ok {
			execute("stop", path[*stop])
		} else {
			log.Fatal("Error: instance name \"", *stop, "\" not found in tomcatman.conf.")
		}
	}
	if *restart != "instance" {
		if _, ok := path[*restart]; ok {
			execute("restart", path[*restart])
		} else {
			log.Fatal("Error: instance name \"", *start, "\" not found in tomcatman.conf.")
		}
	}
}

// execute starts/stops/restarts the given target.
func execute(command, path string) {
	var scriptName string

	switch command {
	case "start":
		os.Setenv("CATALINA_BASE", path)
		scriptName = "./startup.sh"
	case "stop":
		os.Setenv("CATALINA_BASE", path)
		scriptName = "./shutdown.sh"
	case "restart":
		execute("stop", path)
		time.Sleep(5 * time.Second) //TODO: Remove this magic number
		execute("start", path)
		return
	default:
		log.Fatal("unknown command")
	}

	if err := os.Chdir(os.Getenv("CATALINA_HOME") + "/bin"); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(scriptName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	cmd.Wait()
}

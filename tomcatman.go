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
			fmt.Fprintln(os.Stderr, "Error: please ensure tomcatman.conf exists and is in the same directory as this program.")
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
	if *stop != "instance" {
		if _, ok := path[*stop]; ok {
			execute("stop")
		} else {
			fmt.Fprintln(os.Stderr, "Error:", *start, "not found in tomcatman.conf.")
			os.Exit(1)
		}
	}

	if *start != "instance" {
		if _, ok := path[*start]; ok {
			execute("start")
		} else {
			fmt.Fprintln(os.Stderr, "Error:", *start, "not found in tomcatman.conf.")
			os.Exit(1)
		}
	}

	if *restart != "instance" {
		if _, ok := path[*restart]; ok {
			execute("restart")
		} else {
			fmt.Fprintln(os.Stderr, "Error:", *start, "not found in tomcatman.conf.")
			os.Exit(1)
		}
	}
}

// execute starts/stops/restarts the given target.
func execute(command string) {
	switch command {
	case "start":
		os.Setenv("CATALINA_BASE", path[*start])
		if err := os.Chdir(os.Getenv("CATALINA_HOME") + "/bin"); err != nil {
			log.Fatal(err)
		}
		cmd := exec.Command("./startup.sh")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err = cmd.Start(); err != nil {
			log.Fatal(err)
		}
		go io.Copy(os.Stdout, stdout)
		go io.Copy(os.Stderr, stderr)
		cmd.Wait()
	case "stop":
		os.Setenv("CATALINA_BASE", path[*stop])
		if err := os.Chdir(os.Getenv("CATALINA_HOME") + "/bin"); err != nil {
			log.Fatal(err)
		}
		cmd := exec.Command("./shutdown.sh")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err = cmd.Start(); err != nil {
			log.Fatal(err)
		}
		go io.Copy(os.Stdout, stdout)
		go io.Copy(os.Stderr, stderr)
		cmd.Wait()
	case "restart":
		// Stop
		os.Setenv("CATALINA_BASE", path[*restart])
		if err := os.Chdir(os.Getenv("CATALINA_HOME") + "/bin"); err != nil {
			log.Fatal(err)
		}
		cmd := exec.Command("./shutdown.sh")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err = cmd.Start(); err != nil {
			log.Fatal(err)
		}
		go io.Copy(os.Stdout, stdout)
		go io.Copy(os.Stderr, stderr)
		cmd.Wait()

		time.Sleep(5 * time.Second)

		// Start
		cmd = exec.Command("./startup.sh")
		stdout, err = cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		stderr, err = cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err = cmd.Start(); err != nil {
			log.Fatal(err)
		}
		go io.Copy(os.Stdout, stdout)
		go io.Copy(os.Stderr, stderr)
		cmd.Wait()
	}
}

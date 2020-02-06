package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/rs/xid"
)

var (
	count   int
	format  string
	tpltxt  string
	verbose bool
)

func init() {
	flag.IntVar(&count, "n", 1, "Number of xIDs to generate when called with no other arguments.")
	flag.StringVar(&format, "f", "string", "One of string, inspect, time, timestamp, raw, machine, pid, counter or template.")
	flag.StringVar(&tpltxt, "t", "", "The Go template used to format the output.")
	flag.BoolVar(&verbose, "v", false, "Turn on verbose mode.")
}

func main() {
	flag.Parse()
	args := flag.Args()

	var print func(xid.ID)
	switch format {
	case "string":
		print = printString
	case "inspect":
		print = printInspect
	case "time":
		print = printTime
	case "timestamp":
		print = printTimestamp
	case "machine":
		print = printMachine
	case "pid":
		print = printPid
	case "counter":
		print = printCounter
	case "raw":
		print = printRaw
	case "template":
		print = printTemplate
	default:
		fmt.Println("Bad formatting function:", format)
		os.Exit(1)
	}

	if len(args) == 0 {
		for i := 0; i < count; i++ {
			args = append(args, xid.New().String())
		}
	}

	var ids []xid.ID
	for _, arg := range args {
		id, err := xid.FromString(arg)
		if err != nil {
			fmt.Printf("Error when parsing %q: %s\n\n", arg, err)
			flag.PrintDefaults()
			os.Exit(1)
		}
		ids = append(ids, id)
	}

	for _, id := range ids {
		if verbose {
			fmt.Printf("%s: ", id)
		}
		print(id)
	}
}

func printString(id xid.ID) {
	fmt.Println(id.String())
}

func printInspect(id xid.ID) {
	const inspectFormat = `
REPRESENTATION:

  String: %v
     Raw: %v

COMPONENTS:

       Time: %v
  Timestamp: %v
    Machine: %v(%v)
        Pid: %v
    Counter: %v

`
	fmt.Printf(inspectFormat,
		id.String(),
		strings.ToUpper(hex.EncodeToString(id.Bytes())),
		id.Time(),
		id.Time().Unix(),
		string(id.Machine()), strings.ToUpper(hex.EncodeToString(id.Machine())),
		id.Pid(),
		id.Counter(),
	)
}

func printTime(id xid.ID) {
	fmt.Println(id.Time())
}
func printTimestamp(id xid.ID) {
	fmt.Println(id.Time().Unix())
}

func printMachine(id xid.ID) {
	os.Stdout.Write(id.Machine())
	os.Stdout.WriteString("\n")
}

func printPid(id xid.ID) {
	fmt.Println(id.Pid())
}
func printCounter(id xid.ID) {
	fmt.Println(id.Counter())
}

func printRaw(id xid.ID) {
	os.Stdout.Write(id.Bytes())
	os.Stdout.WriteString("\n")
}

func printTemplate(id xid.ID) {
	b := &bytes.Buffer{}
	t := template.Must(template.New("").Parse(tpltxt))
	t.Execute(b, struct {
		String    string
		Raw       string
		Time      time.Time
		Timestamp int64
		Machine   string
		Pid       uint16
		Counter   int32
	}{
		String:    id.String(),
		Raw:       strings.ToUpper(hex.EncodeToString(id.Bytes())),
		Time:      id.Time(),
		Timestamp: id.Time().Unix(),
		Machine:   string(id.Machine()),
		Pid:       id.Pid(),
		Counter:   id.Counter(),
	})
	b.WriteByte('\n')
	io.Copy(os.Stdout, b)
}

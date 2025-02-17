package output

import (
	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/shiena/ansicolor"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Writer is an interface which writes output to somewhere for nuclei events.
type Writer interface {
	// Close closes the output writer interface
	Close()
	// Write writes the event to file and/or screen.
	Write(*ResultEvent) error
	WriteSuccess(*ResultSuccess) error
	// Request logs a request in the trace log
	Request(ip, port, requestType string, err error)
}

type Info struct {
	Banner  string
	Service string
	Cert    string
}

type ResultEvent struct {
	WorkingEvent interface{} `json:"WorkingEvent"`
	Info         *Info       `json:"info,inline"`
	Time         time.Time   `json:"time"`
	Target       string      `json:"Target"`
}

type ResultSuccess struct {
	StepIP  string    `json:"Step-ip"`
	Country string    `json:"country"`
	Time    time.Time `json:"time"`
	Target  string    `json:"target"`
	Domain  string    `json:"domain"`
	Ping    int64     `json:"ping"`
	Status  bool
}

type StandardWriter struct {
	w           io.Writer
	json        bool
	outputFile  map[string]*fileWriter
	outputMutex *sync.Mutex
	traceFile   *fileWriter
	traceMutex  *sync.Mutex
}

var decolorizerRegex = regexp.MustCompile(`\x1B\[[0-9;]*[a-zA-Z]`)

func NewStandardWriter(nocolor, json bool, file string, traceFile string, cdnname []string) (*StandardWriter, error) {
	w := ansicolor.NewAnsiColorWriter(os.Stdout)
	if nocolor {
		color.NoColor = true
	}
	var outfilemap map[string]*fileWriter = make(map[string]*fileWriter)
	if file != "" {
		if len(cdnname) > 0 {
			output, err := newFileOutputWriter("./rst.txt")
			if err != nil {
				return nil, errors.Wrap(err, "could not create output file")
			}
			outfilemap["default"] = output
			var dfile = file
			for i := range cdnname {
				file = strings.Replace(dfile, ".txt", "", 1)
				file = file + "-" + cdnname[i] + ".txt"
				output, err := newFileOutputWriter(file)
				if err != nil {
					return nil, errors.Wrap(err, "could not create output file")
				}
				outfilemap[cdnname[i]] = output
			}
		} else {
			output, err := newFileOutputWriter(file)
			if err != nil {
				return nil, errors.Wrap(err, "could not create output file")
			}
			outfilemap["default"] = output
		}
	}
	var traceOutput *fileWriter
	if traceFile != "" {
		output, err := newFileOutputWriter(traceFile)
		if err != nil {
			return nil, errors.Wrap(err, "could not create output file")
		}
		traceOutput = output
	}
	writer := &StandardWriter{
		w:           w,
		json:        json,
		outputFile:  outfilemap,
		outputMutex: &sync.Mutex{},
		traceFile:   traceOutput,
		traceMutex:  &sync.Mutex{},
	}
	return writer, nil
}

// Write writes the event to file and/or screen.
func (w *StandardWriter) Write(event *ResultEvent) error {
	if event == nil {
		return nil
	}
	event.Time = time.Now()

	var data []byte
	var err error
	if w.json {
		data, err = w.formatJSON(event)
	} else {
		data = w.formatScreen(event)
	}

	if err != nil {
		return errors.Wrap(err, "could not format output")
	}
	if len(data) == 0 {
		return nil
	}

	w.outputMutex.Lock()
	defer w.outputMutex.Unlock()

	_, _ = w.w.Write(data)
	_, _ = w.w.Write([]byte("\n"))

	if w.outputFile != nil {
		if !w.json {
			data = decolorizerRegex.ReplaceAll(data, []byte(""))
		}
		if writeErr := w.outputFile["default"].Write(data); writeErr != nil {
			return errors.Wrap(err, "could not write to output")
		}
	}

	return nil
}

// Write writes the event to file and/or screen.
func (w *StandardWriter) WriteSuccess(event *ResultSuccess) error {
	if event == nil {
		return nil
	}
	event.Time = time.Now()

	var data []byte
	var err error
	if w.json {
		data, err = w.formatSuccessJSON(event)
	} else {
		data = w.formatSuccessScreen(event)
	}

	if err != nil {
		return errors.Wrap(err, "could not format output")
	}
	if len(data) == 0 {
		return nil
	}

	w.outputMutex.Lock()
	defer w.outputMutex.Unlock()

	_, _ = w.w.Write(data)
	_, _ = w.w.Write([]byte("\n"))

	if w.outputFile != nil {
		if !w.json {
			data = decolorizerRegex.ReplaceAll(data, []byte(""))
		}
		if writeErr := w.outputFile[strings.ToLower(event.Domain)].Write(data); writeErr != nil {
			return errors.Wrap(err, "could not write to output")
		}
	}

	return nil
}

func (w *StandardWriter) Close() {
	if w.outputFile != nil {
		for i := range w.outputFile {
			w.outputFile[i].Close()
		}
	}
	if w.traceFile != nil {
		w.traceFile.Close()
	}
}

type JSONTraceRequest struct {
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Error string `json:"error"`
	Type  string `json:"type"`
}

// Request writes a log the requests trace log
func (w *StandardWriter) Request(ip, port, requestType string, err error) {
	if w.traceFile == nil {
		return
	}
	request := &JSONTraceRequest{
		IP:   ip,
		Port: port,
		Type: requestType,
	}
	if err != nil {
		request.Error = err.Error()
	} else {
		request.Error = "none"
	}

	data, err := jsoniter.Marshal(request)
	if err != nil {
		return
	}
	w.traceMutex.Lock()
	_ = w.traceFile.Write(data)
	w.traceMutex.Unlock()
}

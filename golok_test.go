package golok

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/septianw/log15"
)

func CaptureOutput(f func()) string {
	var old = os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	outC := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old
	out := <-outC

	return out
}

func TestGetconfig(t *testing.T) {
	var conf logConfig
	var l, r uintptr

	t.Logf("conf: %+v", conf)

	// Test default value
	conf = Getconfig()
	if conf.logfile != "/var/log/bara.log" {
		t.Errorf("Expected %s, have %s\n", "/var/log/bara.log", conf.logfile)
	}
	if conf.level != log15.LvlWarn {
		t.Errorf("Expected %s, have %s\n", log15.LvlWarn, conf.level)
	}

	/*
		Testing log level OFF
	*/
	Loglevel = "off"
	conf = Getconfig()
	if conf.level != OFF {
		t.Errorf("Expected %s, have %s\n", OFF, conf.level)
	}

	/*
		Testing file log
	*/
	Logfile = "/tmp/golok.log"
	_ = os.Remove(Logfile)
	conf = Getconfig()
	if conf.logfile != Logfile {
		t.Errorf("Expected %s, have %s\n", Logfile, conf.logfile)
	}
	_, err := os.Stat(Logfile)

	if os.IsNotExist(err) {
		t.Errorf("File %s expected to be present, have %+v instead\n", Logfile, err)
	}

	/*
		Testing log level
	*/
	Loglevel = "crit"
	conf = Getconfig()
	if conf.level != log15.LvlCrit {
		t.Errorf("Expected %s, have %s\n", log15.LvlCrit, conf.level)
	}

	Loglevel = "error"
	conf = Getconfig()
	if conf.level != log15.LvlError {
		t.Errorf("Expected %s, have %s\n", log15.LvlError, conf.level)
	}

	Loglevel = "warn"
	conf = Getconfig()
	if conf.level != log15.LvlWarn {
		t.Errorf("Expected %s, have %s\n", log15.LvlWarn, conf.level)
	}

	Loglevel = "info"
	conf = Getconfig()
	if conf.level != log15.LvlInfo {
		t.Errorf("Expected %s, have %s\n", log15.LvlInfo, conf.level)
	}

	Loglevel = "debug"
	conf = Getconfig()
	if conf.level != log15.LvlDebug {
		t.Logf("type of log15.LvlDebug is : %+v", reflect.TypeOf(log15.LvlDebug).Kind())
		t.Logf("type of conf.level is : %+v", reflect.TypeOf(conf.level).Kind())
		t.Errorf("Expected %s, have %s\n", log15.LvlDebug, conf.level)
	}

	Loglevel = "all"
	conf = Getconfig()
	if conf.level != log15.LvlDebug {
		t.Errorf("Expected %s, have %s\n", log15.LvlDebug, conf.level)
	}

	/*
		Testing file format
	*/
	// refer to this : http://stackoverflow.com/questions/9643205/how-do-i-compare-two-functions-for-pointer-equality-in-the-latest-go-weekly
	Logfileformat = "human"
	conf = Getconfig()
	l = reflect.ValueOf(log15.TerminalFormat()).Pointer()
	r = reflect.ValueOf(conf.logfileformat).Pointer()
	if l != r {
		t.Errorf("Expected %s, have %s\n", l, r)
	}

	Logfileformat = "both"
	conf = Getconfig()
	l = reflect.ValueOf(log15.LogfmtFormat()).Pointer()
	r = reflect.ValueOf(conf.logfileformat).Pointer()
	if l != r {
		t.Errorf("Expected %s, have %s\n", l, r)
	}

	Logfileformat = "machine"
	conf = Getconfig()
	l = reflect.ValueOf(log15.JsonFormat()).Pointer()
	r = reflect.ValueOf(conf.logfileformat).Pointer()
	if l != r {
		t.Errorf("Expected %s, have %s\n", l, r)
	}

	/*
		Testing screen format
	*/
	Logscreenformat = "human"
	conf = Getconfig()
	l = reflect.ValueOf(log15.TerminalFormat()).Pointer()
	r = reflect.ValueOf(conf.logscreenformat).Pointer()
	if l != r {
		t.Errorf("Expected %s, have %s\n", l, r)
	}

	Logscreenformat = "both"
	conf = Getconfig()
	l = reflect.ValueOf(log15.LogfmtFormat()).Pointer()
	r = reflect.ValueOf(conf.logscreenformat).Pointer()
	if l != r {
		t.Errorf("Expected %s, have %s\n", l, r)
	}

	Logscreenformat = "machine"
	conf = Getconfig()
	l = reflect.ValueOf(log15.JsonFormat()).Pointer()
	r = reflect.ValueOf(conf.logscreenformat).Pointer()
	if l != r {
		t.Errorf("Expected %s, have %s\n", l, r)
	}

}

func TestWritelog(t *testing.T) {
	Logfile = "/tmp/golok.log"
	Loglevel = "all"
	Logscreenformat = "human"

	var infooutput = CaptureOutput(func() {
		Writelog("info", "info test")
	})

	var warnoutput = CaptureOutput(func() {
		Writelog("warn", "warning test")
	})

	var eroroutput = CaptureOutput(func() {
		Writelog("error", "error test")
	})

	var critoutput = CaptureOutput(func() {
		Writelog("crit", "critical test")
	})

	var dbugoutput = CaptureOutput(func() {
		Writelog("debug", "debug test")
	})

	var oddtest = CaptureOutput(func() {
		Writelog("debug", "Test odd parameter", "This key")
	})
	t.Log(oddtest)

	var niltest = CaptureOutput(func() {
		Writelog("debug", "Test odd parameter", nil, "value", "", 0.45)
	})
	t.Log(niltest)

	var offtest = CaptureOutput(func() {
		Loglevel = "off"
		Writelog("debug", "Off test")
	})

	t.Log(infooutput, warnoutput, eroroutput, critoutput, dbugoutput)
	if !strings.Contains(infooutput, "info test") && !strings.Contains(infooutput, "INFO") {
		t.Errorf("expected %s, have %s", "INFO[<date>|<hh:mm:ss>] info test", warnoutput)
	}

	if !strings.Contains(warnoutput, "warning test") && !strings.Contains(warnoutput, "WARN") {
		t.Errorf("expected %s, have %s", "WARN[<date>|<hh:mm:ss>] warning test", warnoutput)
	}

	if !strings.Contains(eroroutput, "error test") && !strings.Contains(eroroutput, "EROR") {
		t.Errorf("expected %s, have %s", "EROR[<date>|<hh:mm:ss>] error test", warnoutput)
	}

	if !strings.Contains(critoutput, "critical test") && !strings.Contains(critoutput, "CRIT") {
		t.Errorf("expected %s, have %s", "CRIT[<date>|<hh:mm:ss>] critical test", warnoutput)
	}

	if !strings.Contains(dbugoutput, "debug test") && !strings.Contains(dbugoutput, "DBUG") {
		t.Errorf("expected %s, have %s", "DBUG[<date>|<hh:mm:ss>] debug test", warnoutput)
	}

	if !strings.Contains(oddtest, "=nil") {
		t.Errorf("expected %s, have %s", "This key=nil", oddtest)
	}

	if !strings.Contains(niltest, `="value again"`) {
		t.Errorf("expected %s, have %s", `nil=value nil="value again"`, niltest)
	}

	if offtest != "" {
		t.Errorf("Expecting empty log, but resulting %+v", offtest)
	}

}

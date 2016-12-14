package golok

import (
	"bytes"
	//	"fmt"
	"io"
	//	"log"
	"os"
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
	Loglevel = "all"
	Logfile = "/tmp/golok.log"

	conf = Getconfig()
	t.Logf("conf: %+v", conf)

	if conf.logfile != Logfile {
		t.Errorf("Expected %s, have %s\n", Logfile, conf.logfile)
	}

	if conf.level != log15.LvlDebug {
		t.Errorf("Expected %s, have %s\n", log15.LvlDebug, conf.level)
	}
}

func TestWritelog(t *testing.T) {
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

	t.Log(infooutput, warnoutput, eroroutput, critoutput, dbugoutput)
	if !strings.Contains(infooutput, "info test") && !strings.Contains(warnoutput, "INFO") {
		t.Errorf("expected %s, have %s", "INFO[<date>|<hh:mm:ss>] info test", warnoutput)
	}

	if !strings.Contains(warnoutput, "warning test") && !strings.Contains(warnoutput, "WARN") {
		t.Errorf("expected %s, have %s", "WARN[<date>|<hh:mm:ss>] warning test", warnoutput)
	}

	if !strings.Contains(eroroutput, "error test") && !strings.Contains(warnoutput, "EROR") {
		t.Errorf("expected %s, have %s", "EROR[<date>|<hh:mm:ss>] error test", warnoutput)
	}

	if !strings.Contains(critoutput, "critical test") && !strings.Contains(warnoutput, "CRIT") {
		t.Errorf("expected %s, have %s", "CRIT[<date>|<hh:mm:ss>] critical test", warnoutput)
	}

	if !strings.Contains(dbugoutput, "debug test") && !strings.Contains(warnoutput, "DBUG") {
		t.Errorf("expected %s, have %s", "DBUG[<date>|<hh:mm:ss>] debug test", warnoutput)
	}

}

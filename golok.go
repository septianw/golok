package golok

import (
	//	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/septianw/log15"
)

var Loglevel string
var Logfile string
var Logfileformat string
var Logscreenformat string

const OFF = -1
const CRIT = log15.LvlCrit
const ERROR = log15.LvlError
const WARN = log15.LvlWarn
const INFO = log15.LvlInfo
const DEBUG = log15.LvlDebug
const ALL = DEBUG

type logConfig struct {
	level           log15.Lvl
	logfile         string
	logfileformat   log15.Format
	logscreenformat log15.Format
}

/*
  Get human config and convert it to machine config
*/
func Getconfig() logConfig {
	var config logConfig

	//	Critical : 0
	//	Error : 1
	//	Warn : 2
	//	Info : 3
	//	Debug : 4

	switch Loglevel {
	case "off":
		config.level = OFF
	case "crit":
		config.level = CRIT
	case "error":
		config.level = ERROR
	case "warn":
		config.level = WARN
	case "info":
		config.level = INFO
	case "debug":
		config.level = DEBUG
	case "all":
		config.level = ALL
	default:
		config.level = WARN
	}

	if strings.Compare(Logfile, "") == 0 {
		config.logfile = "/var/log/bara.log"
	} else {
		config.logfile = Logfile
	}

	if strings.Compare(Logfileformat, "") == 0 {
		config.logfileformat = log15.JsonFormat()
	} else {
		switch Logfileformat {
		case "human":
			config.logfileformat = log15.TerminalFormat()
		case "both":
			config.logfileformat = log15.LogfmtFormat()
		case "machine":
			config.logfileformat = log15.JsonFormat()
		}
	}

	if strings.Compare(Logscreenformat, "") == 0 {
		config.logscreenformat = log15.TerminalFormat()
	} else {
		switch Logscreenformat {
		case "human":
			config.logscreenformat = log15.TerminalFormat()
		case "both":
			config.logscreenformat = log15.LogfmtFormat()
		case "machine":
			config.logscreenformat = log15.JsonFormat()
		}
	}

	_, err := os.Stat(config.logfile)

	if os.IsNotExist(err) {
		ioutil.WriteFile(config.logfile, []byte(""), 0664)
	}

	return config
}

/*
 Just write to log, and don't care about config.
 get config and write log.
*/
func Writelog(tipe string, message string, ctx ...interface{}) {
	var config logConfig = Getconfig()

	var Log = log15.New()

	if config.level == OFF {
		Log.SetHandler(log15.DiscardHandler())
	} else {
		Log.SetHandler(log15.MultiHandler(
			log15.LvlFilterHandler(config.level, log15.StreamHandler(os.Stdout, config.logscreenformat)),
			log15.LvlFilterHandler(config.level, log15.Must.FileHandler(config.logfile, config.logfileformat)),
			log15.LvlFilterHandler(log15.LvlError, log15.StreamHandler(os.Stderr, log15.TerminalFormat()))))
	}

	if len(ctx)%2 != 0 {
		ctx = append(ctx, nil)
	}

	for i, ct := range ctx {
		if (i%2 == 0) && (ct == nil) {
			ctx[i] = "nil"
		} else if (i%2 == 0) && (reflect.TypeOf(ct).Kind() == reflect.Float64) {
			ctx[i] = strconv.Itoa(int(ct.(float64)))
		}
	}

	switch tipe {
	case "info":
		Log.Info(message, ctx)
	case "warn":
		Log.Warn(message, ctx)
	case "error":
		Log.Error(message, ctx)
	case "crit":
		Log.Crit(message, ctx)
	case "debug":
		Log.Debug(message, ctx)
	}
}

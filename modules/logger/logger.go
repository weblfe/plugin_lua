package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weblfe/plugin_lua/core"
	"github.com/yuin/gopher-lua"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type (
	LuaFunctionTable struct {
		logger *logrus.Logger
	}
)

var (
	defaultLogger = NewLogger()
	Funcs         = defaultLogger.methods()
)

func NewLogger() *LuaFunctionTable {
	var logger = new(LuaFunctionTable)
	return logger.init()
}

// Create function create(file string,level string,mode number) logger
func Create(L *lua.LState) int {
	var args = core.GetArgs(L)
	if len(args) <= 0 {
		return 0
	}
	var (
		logger = NewLogger()
		table  = L.NewTypeMetatable(Name)
	)
	switch len(args) {
	case 1:
		var (
			v       = args[0]
			str, ok = v.(string)
		)
		if ok {
			logger.setOut(str)
		}
	case 2:
		var (
			v       = args[0]
			str, ok = v.(string)
		)
		if ok {
			logger.setOut(str)
		}
		var v2 = args[1]
		if str, ok = v2.(string); ok {
			logger.setLevel(str)
		}
	case 3:
		v := args[0]
		if str, ok := v.(string); ok {
			if n, ok2 := args[2].(lua.LNumber); ok2 {
				var m = uint32(n)
				logger.setOut(str, os.FileMode(m))
			} else {
				logger.setOut(str, os.ModePerm)
			}
		}
		v2 := args[1]
		if str, ok := v2.(string); ok {
			logger.setLevel(str)
		}
	}
	for k, fn := range logger.methods() {
		table.RawSet(lua.LString(k), L.NewFunction(fn))
	}
	L.Push(table)
	return 1
}

func (l *LuaFunctionTable) methods() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"create":     Create,
		"logInfo":    l.logInfo,
		"logInfoLn":  l.logInfoLn,
		"logError":   l.logError,
		"logErrorLn": l.logErrorLn,
		"logDebug":   l.logDebug,
		"logDebugLn": l.logDebugLn,
		"logWarnLn":  l.logWarnLn,
		"logWarn":    l.logWarn,
		"logTrace":   l.logTrace,
		"logTraceLn": l.logTraceLn,
		"setLevel":   l.logSetLevel,
		"getLevel":   l.logGetLevel,
	}
}

func (l *LuaFunctionTable) logSetLevel(L *lua.LState) int {
	var args = core.GetArgs(L)
	if len(args) <= 0 {
		return 0
	}
	if len(args) > 0 {
		if level, ok := args[0].(string); ok {
			l.setLevel(level)
		}
	}
	return 1
}

func (l *LuaFunctionTable) logGetLevel(L *lua.LState) int {
	var level = ""
	switch l.logger.Level {
	case logrus.DebugLevel:
		level = "debug"
	case logrus.WarnLevel:
		level = "warn"
	case logrus.InfoLevel:
		level = "info"
	case logrus.ErrorLevel:
		level = "error"
	case logrus.TraceLevel:
		level = "trace"
	}
	L.Push(lua.LString(level))
	return 1
}

func (l *LuaFunctionTable) setOut(out string, mod ...os.FileMode) *LuaFunctionTable {
	mod = append(mod, os.ModePerm)
	if out == "" {
		return l
	}
	var file, err = filepath.Abs(out)
	if err != nil {
		fmt.Println(err.Error())
		return l
	}
	if _, e := os.Stat(file); e != nil {
		if !os.IsNotExist(e) {
			return l
		}
		_ = os.MkdirAll(filepath.Dir(file), mod[0])
	}
	if fd, err2 := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, mod[0]); err2 == nil {
		l.logger.Out = fd
		runtime.SetFinalizer(l, (*LuaFunctionTable).destroy)
	}
	return l
}

func (l *LuaFunctionTable) setLevel(level string) *LuaFunctionTable {
	switch strings.ToLower(level) {
	case "info":
		l.logger.SetLevel(logrus.InfoLevel)
	case "error":
		l.logger.SetLevel(logrus.ErrorLevel)
	case "debug":
		l.logger.SetLevel(logrus.DebugLevel)
	case "warn":
		l.logger.SetLevel(logrus.WarnLevel)
	case "trace":
		l.logger.SetLevel(logrus.TraceLevel)
	}
	return l
}

func (l *LuaFunctionTable) init() *LuaFunctionTable {
	l.logger = logrus.New()
	return l
}

func (l *LuaFunctionTable) logInfoLn(L *lua.LState) int {
	var args = core.GetArgs(L)
	if len(args) <= 0 {
		return 0
	}
	l.logger.Infoln(args...)
	return 1
}

func (l *LuaFunctionTable) logInfo(state *lua.LState) int {
	var args = core.GetArgs(state)
	if len(args) <= 0 {
		return 0
	}
	l.logger.Infoln(args...)
	return 1
}

func (l *LuaFunctionTable) logTrace(state *lua.LState) int {
	var args = core.GetArgs(state)
	if len(args) <= 0 {
		return 0
	}
	l.logger.Trace(args...)
	return 1
}

func (l *LuaFunctionTable) logWarn(state *lua.LState) int {
	var args = core.GetArgs(state)
	if len(args) <= 0 {
		return 0
	}
	l.logger.Warn(args...)
	return 1
}

func (l *LuaFunctionTable) logWarnLn(state *lua.LState) int {
	var args = core.GetArgs(state)
	if len(args) <= 0 {
		return 0
	}
	l.logger.Warnln(args...)
	return 1
}

func (l *LuaFunctionTable) logTraceLn(state *lua.LState) int {
	var args = core.GetArgs(state)
	if len(args) <= 0 {
		return 0
	}
	l.logger.Traceln(args...)
	return 1
}

func (l *LuaFunctionTable) logDebug(L *lua.LState) int {
	var args = core.GetArgs(L)
	if len(args) <= 0 {
		return 0
	}
	l.logger.Debug(args...)
	return 1
}

func (l *LuaFunctionTable) logDebugLn(L *lua.LState) int {
	var args = core.GetArgs(L)
	if len(args) <= 0 {
		return 0
	}
	l.logger.Debugln(args...)
	return 1
}

func (l *LuaFunctionTable) logErrorLn(L *lua.LState) int {
	var args = core.GetArgs(L)
	if len(args) <= 0 {
		return 0
	}
	l.logger.Errorln(args...)
	return 1
}

func (l *LuaFunctionTable) logError(L *lua.LState) int {
	var args = core.GetArgs(L)
	if len(args) <= 0 {
		return 0
	}
	l.logger.Error(args...)
	return 1
}

func (l *LuaFunctionTable) destroy() {
	if l.logger.Out == nil {
		return
	}
	runtime.SetFinalizer(l, nil)
	if l.logger.Out != os.Stderr && l.logger.Out != os.Stdin {
		if closer, ok := l.logger.Out.(io.WriteCloser); ok {
			_ = closer.Close()
		}
		l.logger.Out = nil
	}
}

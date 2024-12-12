package log

import (
	"encoding/json"
	"fmt"
	"maps"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// log level

type LogLevel int

const (
	LLVerbose LogLevel = iota
	LLDebug
	LLInfo
	LLWarn
	LLError
	LLCritical
)

func GetLogLevel(level string) LogLevel {
	if strings.ToLower(level) == "critical" {
		return LLCritical
	}
	if strings.ToLower(level) == "error" {
		return LLError
	}
	if strings.ToLower(level) == "warn" || strings.ToLower(level) == "warning" {
		return LLWarn
	}
	if strings.ToLower(level) == "info" {
		return LLInfo
	}
	if strings.ToLower(level) == "debug" {
		return LLDebug
	}
	return LLVerbose
}

// connection protocol

type ConnProtocol int

const (
	Tcp ConnProtocol = iota
	Udp
)

func getNetwork(protocol ConnProtocol) string {
	if protocol == Udp {
		return "udp"
	}
	return "tcp"
}

// common keys

type Keys struct {
	Timestamp string
	Message   string
	Severity  string
}

type RLog struct {
	Address    string
	Protocol   ConnProtocol
	Timeout    time.Duration
	Keys       Keys
	CommonData map[string]interface{}
	LogLevel   LogLevel
	TimeFormat string
	IsUtc      bool
}

func NewRLog(host string, port int, logLevel, timeFormat string, isUtc bool) *RLog {
	return NewRLogExt(host, port, logLevel, timeFormat, isUtc, Tcp, 3)
}

func NewRLogExt(host string, port int, logLevel, timeFormat string, isUtc bool, protocol ConnProtocol, timeout time.Duration) *RLog {
	address := ""
	if host != "" && port != 0 {
		address = host + ":" + strconv.Itoa(port)
	}

	keys := Keys{
		Message:   "message",
		Timestamp: "timestamp",
		Severity:  "severity",
	}

	tf := timeFormat
	if tf == "" {
		tf = time.DateTime
	}

	if os.Getenv("GU_REMOTE_LOG_DEBUG") == "true" {
		fmt.Println("Configured rlog. Address:", address, ", Keys:", keys, ", Timeformat:", tf)
	}

	rlog := &RLog{
		Address:    address,
		Protocol:   protocol,
		Timeout:    timeout,
		Keys:       keys,
		LogLevel:   GetLogLevel(logLevel),
		TimeFormat: tf,
		IsUtc:      isUtc,
	}
	return rlog
}

func (l *RLog) UpdateCommonData(data map[string]interface{}) {
	l.CommonData = data
}

func (l *RLog) UpdateKeys(message, severity, timestamp string) {
	if message != "" {
		l.Keys.Message = message
	}
	if severity != "" {
		l.Keys.Severity = severity
	}
	if timestamp != "" {
		l.Keys.Timestamp = timestamp
	}
}

func (l *RLog) Verbose(logs ...any) {
	if l.LogLevel > LLVerbose {
		return
	}
	msg := getMsg(logs...)
	Verbose(msg)
	l.sendJsonWithSeverity(msg, nil, "VERBOSE")
}

func (l *RLog) VerboseD(data map[string]interface{}, logs ...any) {
	if l.LogLevel > LLVerbose {
		return
	}
	msg := getMsg(logs...)
	Verbose(msg)
	l.sendJsonWithSeverity(msg, data, "VERBOSE")
}

func (l *RLog) Debug(logs ...any) {
	if l.LogLevel > LLDebug {
		return
	}

	msg := getMsg(logs...)
	Debug(msg)
	l.sendJsonWithSeverity(msg, nil, "DEBUG")
}

func (l *RLog) DebugD(data map[string]interface{}, logs ...any) {
	if l.LogLevel > LLDebug {
		return
	}
	msg := getMsg(logs...)
	Debug(msg)
	l.sendJsonWithSeverity(msg, data, "DEBUG")
}

func (l *RLog) Info(logs ...any) {
	if l.LogLevel > LLInfo {
		return
	}
	msg := getMsg(logs...)
	Info(msg)
	l.sendJsonWithSeverity(msg, nil, "INFO")
}

func (l *RLog) InfoD(data map[string]interface{}, logs ...any) {
	if l.LogLevel > LLInfo {
		return
	}
	msg := getMsg(logs...)
	Info(msg)
	l.sendJsonWithSeverity(msg, data, "INFO")
}

func (l *RLog) Warn(logs ...any) {
	if l.LogLevel > LLWarn {
		return
	}
	msg := getMsg(logs...)
	Warn(msg)
	l.sendJsonWithSeverity(msg, nil, "WARN")
}

func (l *RLog) WarnD(data map[string]interface{}, logs ...any) {
	if l.LogLevel > LLWarn {
		return
	}
	msg := getMsg(logs...)
	Warn(msg)
	l.sendJsonWithSeverity(msg, data, "WARN")
}

func (l *RLog) Error(logs ...any) {
	if l.LogLevel > LLError {
		return
	}
	msg := getMsg(logs...)
	Error(msg)
	l.sendJsonWithSeverity(msg, nil, "ERROR")
}

func (l *RLog) ErrorD(data map[string]interface{}, logs ...any) {
	if l.LogLevel > LLError {
		return
	}
	msg := getMsg(logs...)
	Error(msg)
	l.sendJsonWithSeverity(msg, data, "ERROR")
}

func (l *RLog) Critical(logs ...any) {
	if l.LogLevel > LLCritical {
		return
	}
	msg := getMsg(logs...)
	Critical(msg)
	l.sendJsonWithSeverity(msg, nil, "CRITICAL")
}

func (l *RLog) CriticalD(data map[string]interface{}, logs ...any) {
	if l.LogLevel > LLCritical {
		return
	}
	msg := getMsg(logs...)
	Critical(msg)
	l.sendJsonWithSeverity(msg, data, "CRITICAL")
}

func (l *RLog) LogString(logs ...any) {
	msg := getMsg(logs...)
	l.sendString(msg)
	Plain(msg)
}

// private

func (l *RLog) connect() net.Conn {
	network := getNetwork(l.Protocol)
	conn, err := net.Dial(network, l.Address)
	if err != nil {
		// Err(err, "dial to network", network, "address", l.Address)
		return nil
	}

	// update deadlines
	deadline := time.Now().Add(l.Timeout * time.Second)
	conn.SetDeadline(deadline)
	conn.SetWriteDeadline(deadline)
	conn.SetReadDeadline(deadline)

	return conn
}

func (l *RLog) sendJsonWithSeverity(msg string, add map[string]interface{}, severity string) {
	if l.Address == "" {
		if os.Getenv("GU_REMOTE_LOG_DEBUG") == "true" {
			fmt.Println("No address for remote logging.")
		}
		return
	}

	// create map
	data := map[string]interface{}{}
	data[l.Keys.Severity] = severity
	data[l.Keys.Message] = msg

	// timestamp
	var ts time.Time
	if l.IsUtc {
		ts = time.Now().UTC()
	} else {
		ts = time.Now()
	}
	if os.Getenv("GU_REMOTE_LOG_DEBUG") == "true" {
		fmt.Println("ts:", ts)
	}
	data[l.Keys.Timestamp] = ts.Format(l.TimeFormat)

	// copy additional
	if add != nil {
		maps.Copy(data, add)
	}
	if l.CommonData != nil {
		maps.Copy(data, l.CommonData)
	}

	go l.sendJson(data)
}

func (l *RLog) sendJson(data map[string]interface{}) bool {
	if l.Address == "" {
		if os.Getenv("GU_REMOTE_LOG_DEBUG") == "true" {
			fmt.Println("No address for remote logging.")
		}
		return false
	}

	conn := l.connect()
	if conn == nil {
		if os.Getenv("GU_REMOTE_LOG_DEBUG") == "true" {
			fmt.Println("No connection for remote logging.")
		}
		return false
	}
	defer conn.Close()

	payload, err := json.Marshal(data)
	if err != nil {
		Err(err, "cannot marshal data", data)
		return false
	}

	if os.Getenv("GU_REMOTE_LOG_DEBUG") == "true" {
		fmt.Println(string(payload))
	}

	_, err = conn.Write(payload)
	if err != nil {
		Err(err, "cannot write payload")
		return false
	}
	return true
}

func (l *RLog) sendString(data string) bool {
	if l.Address == "" {
		return false
	}

	if os.Getenv("GU_REMOTE_LOG_DEBUG") == "true" {
		fmt.Println(data)
	}

	conn := l.connect()
	if conn == nil {
		return false
	}
	defer conn.Close()

	_, err := conn.Write([]byte(data))
	if err != nil {
		Err(err, "cannot write string")
		return false
	}
	return true
}

func getMsg(logs ...any) string {
	msg := fmt.Sprintln(logs...)
	return strings.TrimSuffix(msg, "\n")
}

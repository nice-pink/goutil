package log

import (
	"encoding/json"
	"maps"
	"net"
	"strconv"
	"time"
)

const (
	SEVERITY_KEY = "Severity"
)

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

type RLog struct {
	Address      string
	Protocol     ConnProtocol
	Timeout      time.Duration
	TimestampKey string
	MessageKey   string
	CommonData   map[string]interface{}
}

func NewRLog(host string, port int, protocol ConnProtocol, timeout time.Duration) *RLog {
	address := host + ":" + strconv.Itoa(port)
	rlog := &RLog{
		Address:      address,
		Protocol:     protocol,
		Timeout:      timeout,
		MessageKey:   "message",
		TimestampKey: "timestamp",
	}
	return rlog
}

func (l *RLog) UpdateCommonData(data map[string]interface{}) {
	l.CommonData = data
}

func (l *RLog) UpdateKeys(message, timestamp string) {
	if message != "" {
		l.MessageKey = message
	}
	if timestamp != "" {
		l.TimestampKey = timestamp
	}
}

func (l *RLog) Verbose(msg string, data map[string]interface{}) {
	l.sendJsonWithSeverity(msg, data, "VERBOSE")
	Verbose(data[l.MessageKey])
}

func (l *RLog) Info(msg string, data map[string]interface{}) {
	l.sendJsonWithSeverity(msg, data, "INFO")
	Info(data[l.MessageKey])
}

func (l *RLog) Debug(msg string, data map[string]interface{}) {
	l.sendJsonWithSeverity(msg, data, "DEBUG")
	Debug(data[l.MessageKey])
}

func (l *RLog) Warn(msg string, data map[string]interface{}) {
	l.sendJsonWithSeverity(msg, data, "WARN")
	Warn(data[l.MessageKey])
}

func (l *RLog) Error(msg string, data map[string]interface{}) {
	l.sendJsonWithSeverity(msg, data, "ERROR")
	Error(data[l.MessageKey])
}

func (l *RLog) Critical(msg string, data map[string]interface{}) {
	l.sendJsonWithSeverity(msg, data, "CRITICAL")
	Critical(data[l.MessageKey])
}

func (l *RLog) LogString(msg string) {
	l.sendString(msg)
	Plain(msg)
}

// private

func (l *RLog) connect() net.Conn {
	network := getNetwork(l.Protocol)
	conn, err := net.Dial(network, l.Address)
	if err != nil {
		Err(err, "dial to network", network, "address", l.Address)
		return nil
	}

	// update deadlines
	deadline := time.Now().Add(l.Timeout)
	conn.SetDeadline(deadline)
	conn.SetWriteDeadline(deadline)
	conn.SetReadDeadline(deadline)

	return conn
}

func (l *RLog) sendJsonWithSeverity(msg string, add map[string]interface{}, severity string) bool {
	// create map
	data := map[string]interface{}{}
	data[SEVERITY_KEY] = severity
	data[l.TimestampKey] = time.Now().Format(time.DateTime)
	data[l.MessageKey] = msg

	// copy additional
	maps.Copy(data, add)
	maps.Copy(data, l.CommonData)

	return l.sendJson(data)
}

func (l *RLog) sendJson(data map[string]interface{}) bool {
	conn := l.connect()
	if conn == nil {
		return false
	}
	defer conn.Close()

	payload, err := json.Marshal(data)
	if err != nil {
		Err(err, "cannot marshal data", data)
		return false
	}

	_, err = conn.Write(payload)
	if err != nil {
		Err(err, "cannot write payload")
		return false
	}
	return true
}

func (l *RLog) sendString(data string) bool {
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

package log

import (
	"encoding/json"
	"maps"
	"net"
	"strconv"
	"time"
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
}

func NewRLog(host string, port int) *RLog {
	return NewRLogExt(host, port, Tcp, 3)
}

func NewRLogExt(host string, port int, protocol ConnProtocol, timeout time.Duration) *RLog {
	address := host + ":" + strconv.Itoa(port)

	keys := Keys{
		Message:   "message",
		Timestamp: "timestamp",
		Severity:  "severity",
	}

	rlog := &RLog{
		Address:  address,
		Protocol: protocol,
		Timeout:  timeout,
		Keys:     keys,
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

func (l *RLog) Verbose(msg string, data map[string]interface{}) {
	l.sendJsonWithSeverity(msg, data, "VERBOSE")
	Verbose(data[l.Keys.Message])
}

func (l *RLog) Info(msg string, data map[string]interface{}) {
	l.sendJsonWithSeverity(msg, data, "INFO")
	Info(data[l.Keys.Message])
}

func (l *RLog) Debug(msg string, data map[string]interface{}) {
	l.sendJsonWithSeverity(msg, data, "DEBUG")
	Debug(data[l.Keys.Message])
}

func (l *RLog) Warn(msg string, data map[string]interface{}) {
	l.sendJsonWithSeverity(msg, data, "WARN")
	Warn(data[l.Keys.Message])
}

func (l *RLog) Error(msg string, data map[string]interface{}) {
	l.sendJsonWithSeverity(msg, data, "ERROR")
	Error(data[l.Keys.Message])
}

func (l *RLog) Critical(msg string, data map[string]interface{}) {
	l.sendJsonWithSeverity(msg, data, "CRITICAL")
	Critical(data[l.Keys.Message])
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
	if l.Address == "" {
		return false
	}

	// create map
	data := map[string]interface{}{}
	data[l.Keys.Severity] = severity
	data[l.Keys.Timestamp] = time.Now().Format(time.DateTime)
	data[l.Keys.Message] = msg

	// copy additional
	maps.Copy(data, add)
	maps.Copy(data, l.CommonData)

	return l.sendJson(data)
}

func (l *RLog) sendJson(data map[string]interface{}) bool {
	if l.Address == "" {
		return false
	}

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
	if l.Address == "" {
		return false
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

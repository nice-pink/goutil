package log

import (
	"encoding/json"
	"net"
	"strconv"
	"time"

	"github.com/google/martian/v3/log"
	"github.com/nice-pink/goutil/pkg/log"
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
	Timeout      time.Time
	TimestampKey string
	MessageKey   string
	CommonData   map[string]interface{}
}

func NewRLog(host string, port int, protocol ConnProtocol, timeout time.Time) *RLog {
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

func (l *RLog) Verbose(data map[string]interface{}) {
	l.sendJsonWithSeverity(data, "VERBOSE")
	log.Verbose(data[l.MessageKey])
}

func (l *RLog) Info(data map[string]interface{}) {
	l.sendJsonWithSeverity(data, "INFO")
	log.Info(data[l.MessageKey])
}

func (l *RLog) Debug(data map[string]interface{}) {
	l.sendJsonWithSeverity(data, "DEBUG")
	log.Debug(data[l.MessageKey])
}

func (l *RLog) Warn(data map[string]interface{}) {
	l.sendJsonWithSeverity(data, "WARN")
	log.Warn(data[l.MessageKey])
}

func (l *RLog) Error(data map[string]interface{}) {
	l.sendJsonWithSeverity(data, "ERROR")
	log.Error(data[l.MessageKey])
}

func (l *RLog) Critical(data map[string]interface{}) {
	l.sendJsonWithSeverity(data, "CRITICAL")
	log.Critical(data[l.MessageKey])
}

func (l *RLog) LogString(msg string) {
	l.sendString(msg)
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
	conn.SetDeadline(l.Timeout)
	conn.SetWriteDeadline(l.Timeout)
	conn.SetReadDeadline(l.Timeout)

	return conn
}

func (l *RLog) sendJsonWithSeverity(data map[string]interface{}, severity string) bool {
	data[SEVERITY_KEY] = severity
	success := l.sendJson(data)
	delete(data, SEVERITY_KEY)
	return success
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

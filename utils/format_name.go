package utils

import "strings"

var commonInitialisms = map[string]string{
	"acl":   "ACL",
	"api":   "API",
	"ascii": "ASCII",
	"cpu":   "CPU",
	"css":   "CSS",
	"dns":   "DNS",
	"eof":   "EOF",
	"guid":  "GUID",
	"html":  "HTML",
	"http":  "HTTP",
	"https": "HTTPS",
	"id":    "ID",
	"ip":    "IP",
	"json":  "JSON",
	"lhs":   "LHS",
	"qps":   "QPS",
	"ram":   "RAM",
	"rhs":   "RHS",
	"rpc":   "RPC",
	"sla":   "SLA",
	"smtp":  "SMTP",
	"sql":   "SQL",
	"ssh":   "SSH",
	"tcp":   "TCP",
	"tls":   "TLS",
	"ttl":   "TTL",
	"udp":   "UDP",
	"ui":    "UI",
	"uid":   "UID",
	"uuid":  "UUID",
	"uri":   "URI",
	"url":   "URL",
	"utf8":  "UTF8",
	"vm":    "VM",
	"xml":   "XML",
	"xmpp":  "XMPP",
	"xsrf":  "XSRF",
	"xss":   "XSS",
}

// Camelize camelize name
func Camelize(name string) string {
	temp := []byte(name)
	for i := 0; i < len(temp); i++ {
		switch temp[i] {
		case '-':
			temp[i] = ' '
		case '_':
			temp[i] = ' '
		}
	}
	ss := strings.Split(string(temp), " ")
	for i, s := range ss {
		if r, ok := commonInitialisms[s]; ok {
			ss[i] = r
		}
	}
	name = strings.Join(ss, " ")
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// CamelizeDownFirst camelize and down first name
func CamelizeDownFirst(name string) string {
	name = Camelize(name)
	return strings.ToLower(name[:1]) + name[1:]
}

// DownFirst down first name
func DownFirst(name string) string {
	return strings.ToLower(name[:1]) + name[1:]
}

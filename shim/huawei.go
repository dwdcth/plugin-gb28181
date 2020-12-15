package shim

import (
	"fmt"
	"regexp"
	"strings"
	"github.com/Monibuca/plugin-gb28181/sip"
)

var recordRe = regexp.MustCompile(`(?m)<DeviceID>(\d+)<\/DeviceID>`)
var cmdRe = regexp.MustCompile(`(?m)<CmdType>(\w+)<\/CmdType>`)
var snRe = regexp.MustCompile(`(?m)<SN>(\w+)<\/SN>`)

func getCmd(body string) string {
	match := cmdRe.FindAllStringSubmatch(body, -1)
	if len(match) > 0 {
		return match[0][1]
	}
	return ""
}

func getDeviceIdFromResp(body string) string {
	match := recordRe.FindAllStringSubmatch(body, -1)
	if len(match) > 0 {
		return match[0][1]
	}
	return ""
}

func getDeviceIdFromHost(host string) string {
	return strings.Split(host, "@")[0]
}

func getSn(body string) string {
	match := snRe.FindAllStringSubmatch(body, -1)
	if len(match) > 0 {
		return match[0][1]
	}
	return ""
}

// todo SN
func GetTidFromResponse(msg *sip.Message, filterCmd string) string {
	body := msg.Body
	method := string(msg.GetMethod())
	if getCmd(body) == filterCmd {
		return fmt.Sprintf("%s_%s_%s_%s", method, filterCmd, getDeviceIdFromResp(body), getSn(body))
	}
	return ""
}

func GetTidFromHost(msg *sip.Message, filterCmd string) string {
	host := msg.To.Uri.Host()
	body := msg.Body
	method := string(msg.GetMethod())
	if getCmd(body) == filterCmd {
		return fmt.Sprintf("%s_%s_%s_%s", method, filterCmd, getDeviceIdFromHost(host), getSn(body))
	}
	return ""
}

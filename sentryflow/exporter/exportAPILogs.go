// SPDX-License-Identifier: Apache-2.0

package exporter

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/5gsec/SentryFlow/protobuf"
)

// == //

// apiLogStreamInform structure
type apiLogStreamInform struct {
	Hostname  string
	IPAddress string

	stream protobuf.SentryFlow_GetAPILogServer
	error  chan error
}

// apiLogStreamInformV2 structure
type apiLogStreamInformV2 struct {
	Hostname  string
	IPAddress string

	stream protobuf.SentryFlow_GetAPILogV2Server
	error  chan error
}

// GetAPILog Function (for gRPC)
func (exs *ExpService) GetAPILog(info *protobuf.ClientInfo, stream protobuf.SentryFlow_GetAPILogServer) error {
	log.Printf("[Exporter] Client %s (%s) connected (GetAPILog)", info.HostName, info.IPAddress)

	currExporter := &apiLogStreamInform{
		Hostname:  info.HostName,
		IPAddress: info.IPAddress,
		stream:    stream,
	}

	ExpH.exporterLock.Lock()
	ExpH.apiLogExporters = append(ExpH.apiLogExporters, currExporter)
	ExpH.exporterLock.Unlock()

	return <-currExporter.error
}

// GetAPILogV2 Function (for gRPC)
func (exs *ExpService) GetAPILogV2(info *protobuf.ClientInfo, stream protobuf.SentryFlow_GetAPILogV2Server) error {
	log.Printf("[Exporter] Client %s (%s) connected (GetAPILogV2)", info.HostName, info.IPAddress)

	currExporter := &apiLogStreamInformV2{
		Hostname:  info.HostName,
		IPAddress: info.IPAddress,
		stream:    stream,
	}

	ExpH.exporterLock.Lock()
	ExpH.apiLogExportersV2 = append(ExpH.apiLogExportersV2, currExporter)
	ExpH.exporterLock.Unlock()

	return <-currExporter.error
}

// SendAPILogs Function
func (exp *ExpHandler) SendAPILogs(apiLog *protobuf.APILog) error {
	failed := 0
	total := len(exp.apiLogExporters)

	for _, exporter := range exp.apiLogExporters {
		log.Print("Sending api log!!!!")
		log.Printf("Sending api log right here!!!!! %+v\n", apiLog)
		if err := exporter.stream.Send(apiLog); err != nil {
			log.Printf("[Exporter] Failed to export an API log to %s (%s): %v", exporter.Hostname, exporter.IPAddress, err)
			failed++
		}
	}

	if failed != 0 {
		msg := fmt.Sprintf("[Exporter] Failed to export API logs properly (%d/%d failed)", failed, total)
		return errors.New(msg)
	}

	return nil
}

// SendAPILogsV2 Function
func (exp *ExpHandler) SendAPILogsV2(apiLog *protobuf.APILogV2) error {
	failed := 0
	total := len(exp.apiLogExportersV2)

	for _, exporter := range exp.apiLogExportersV2 {
		if err := exporter.stream.Send(apiLog); err != nil {
			log.Printf("[Exporter] Failed to export an API log(V2) to %s (%s): %v", exporter.Hostname, exporter.IPAddress, err)
			failed++
		}
	}

	if failed != 0 {
		msg := fmt.Sprintf("[Exporter] Failed to export API logs(V2) properly (%d/%d failed)", failed, total)
		return errors.New(msg)
	}

	return nil
}

// == //

// InsertAPILog Function
func InsertAPILog(apiLog interface{}) {
	switch data := apiLog.(type) {
	case *protobuf.APILog:
		ExpH.exporterAPILogs <- data
		// Make a string with labels
		var labelString []string
		for k, v := range data.SrcLabel {
			labelString = append(labelString, fmt.Sprintf("%s:%s", k, v))
		}
		sort.Strings(labelString)

		// Update Stats per namespace and per labels
		UpdateStats(data.SrcNamespace, strings.Join(labelString, ","), data.GetPath())
	case *protobuf.APILogV2:
		ExpH.exporterAPILogsV2 <- data
	}
}

// == //

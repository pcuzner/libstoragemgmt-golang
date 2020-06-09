// SPDX-License-Identifier: 0BSD

package libstoragemgmt

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"

	errors "github.com/libstorage/libstoragemgmt-golang/errors"
)

// TmoSetCb used to register timeout value for plugin
type TmoSetCb func(timeout uint32) error

// TmoGetCb used to register timeout value for plugin
type TmoGetCb func() uint32

// CapabilitiesCb returns what the plugin is capable of
type CapabilitiesCb func(system *System) (*Capabilities, error)

// JobStatusCb callback returns the job status for the specified job
// The returned array should be [JobStatusType, uint8, result]
type JobStatusCb func(jobId string) ([]interface{}, error)

// JobFreeCb callback for freeing job resources
type JobFreeCb func(jobID string) error

// PoolsCb callback for pools
type PoolsCb func(search ...string) ([]Pool, error)

// PluginRegisterCb callback to register needed information
type PluginRegisterCb func(p *PluginRegister) error

// PluginUnregisterCb callback to allow plugin to cleanup resources
type PluginUnregisterCb func() error

//SystemsCb callback to retrieve systems
type SystemsCb func() ([]System, error)

//VolumesCb callback to retrieve volumes
type VolumesCb func(search ...string) ([]Volume, error)

// RequiredCallbacks are the callbacks that plugins must implement
type RequiredCallbacks struct {
	TimeOutSet       TmoSetCb
	TimeOutGet       TmoGetCb
	JobStatus        JobStatusCb
	JobFree          JobFreeCb
	Capabilities     CapabilitiesCb
	Systems          SystemsCb
	Pools            PoolsCb
	PluginRegister   PluginRegisterCb
	PluginUnregister PluginUnregisterCb
}

// CallBacks callbacks for plugin to implement
type CallBacks struct {
	Required RequiredCallbacks
	Volumes  VolumesCb
}

type handler func(p *Plugin, params json.RawMessage) (interface{}, error)

// Plugin represents plugin
type Plugin struct {
	tp        transPort
	cb        *CallBacks
	callTable map[string]handler
}

// PluginRegister data passed to PluginRegister callback
type PluginRegister struct {
	URI      string
	Password string
	Timeout  uint32
	Flags    uint32
}

// PluginInit initializes the plugin with the specified callbacks
func PluginInit(callbacks *CallBacks, cmdLineArgs []string) (*Plugin, error) {
	if len(cmdLineArgs) == 2 {
		fd, err := strconv.ParseInt(cmdLineArgs[1], 10, 32)
		if err != nil {
			return nil, err
		}

		// Only information I could find which pretains to how to do this.
		// https://play.golang.org/p/0uEcuPk291
		f := os.NewFile(uintptr(fd), "client")
		s, err := net.FileConn(f)
		if err != nil {
			return nil, err
		}

		tp := transPort{uds: s, debug: false}
		return &Plugin{tp: tp, cb: callbacks, callTable: buildTable(callbacks)}, nil
	}
	return nil, &errors.LsmError{
		Code:    errors.LibBug,
		Message: fmt.Sprintf("Plugin called with invalid args: %s\n", cmdLineArgs)}
}

func noSupport(tp *transPort, method string) {
	tp.sendError(&errors.LsmError{
		Code: errors.NoSupport,
		Message: fmt.Sprintf(
			"method %s not supported", method)})
}

// Run the plugin, looping processing requests and sending responses.
func (p *Plugin) Run() {
	for {
		request, err := p.tp.readRequest()
		if err != nil {
			if lsmError, ok := err.(*errors.LsmError); ok == true {

				if lsmError.Code != errors.TransPortComunication {
					p.tp.sendError(lsmError)
					fmt.Printf("Returned error %+v\n", lsmError)
					continue
				} else {
					fmt.Printf("Communication error: exiting! %s\n", lsmError)
				}
				return
			}
			fmt.Printf("Unexpected error, exiting! %s\n", err)
			return
		}

		var response interface{}
		if f, ok := p.callTable[request.Method]; ok == true && f != nil {
			fmt.Printf("Executing %s(%s)\n", request.Method, string(request.Params))
			response, err = f(p, request.Params)
			if err != nil {
				p.tp.sendError(err)
			} else {
				p.tp.sendResponse(response)
			}

			// Need to shut down the connection.
			if request.Method == "plugin_unregister" {
				p.tp.close()
				return
			}
		} else {
			noSupport(&p.tp, request.Method)
		}
	}
}

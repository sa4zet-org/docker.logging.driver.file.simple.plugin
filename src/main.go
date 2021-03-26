package main

import (
	"encoding/json"
	dockerPluginSdk "github.com/docker/go-plugins-helpers/sdk"
	dockerDaemonLogger "github.com/moby/moby/daemon/logger"
	"net/http"
)

func main() {
	//https://docs.docker.com/engine/extend/plugins_logging/
	pluginHandler := dockerPluginSdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)
	fileLoggingDriver := newFileLoggingDriver()
	pluginHandler.HandleFunc("/LogDriver.StartLogging", startLoggingHandler(fileLoggingDriver))
	pluginHandler.HandleFunc("/LogDriver.StopLogging", stopLoggingHandler(fileLoggingDriver))
	pluginHandler.HandleFunc("/LogDriver.Capabilities", capabilitiesHandler())
	if err := pluginHandler.ServeUnix("fileLoggingDriver", 0); err != nil {
		panic(err)
	}
}

type StartLoggingRequest struct {
	File string
	Info dockerDaemonLogger.Info
}

type StopLoggingRequest struct {
	File string
}

type CapabilitiesResponse struct {
	Err string
	Cap dockerDaemonLogger.Capability
}

type PluginResponse struct {
	Err string
}

func respond(err error, responseWriter http.ResponseWriter) {
	var pluginResponse PluginResponse
	if err != nil {
		pluginResponse.Err = err.Error()
	}
	_ = json.NewEncoder(responseWriter).Encode(&pluginResponse)
}

func startLoggingHandler(fileLoggingDriver *FileLoggingDriver) func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		var startLoggingRequest StartLoggingRequest
		if err := json.NewDecoder(httpRequest.Body).Decode(&startLoggingRequest); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}
		err := fileLoggingDriver.StartLogging(startLoggingRequest.File, startLoggingRequest.Info)
		respond(err, responseWriter)
	}
}

func stopLoggingHandler(fileLoggingDriver *FileLoggingDriver) func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		var stopLoggingRequest StopLoggingRequest
		if err := json.NewDecoder(httpRequest.Body).Decode(&stopLoggingRequest); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}
		err := fileLoggingDriver.StopLogging(stopLoggingRequest.File)
		respond(err, responseWriter)
	}
}

func capabilitiesHandler() func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		_ = json.NewEncoder(responseWriter).Encode(&CapabilitiesResponse{
			Cap: dockerDaemonLogger.Capability{ReadLogs: false},
		})
	}
}

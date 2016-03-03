package runner

import (
	"flag"
	"log"
	"os"
)

var (
	registryAddress = flag.String(
		"registry-address",
		"",
		"Address of the service registry. Alternatively set the REGISTRY_ADDRESS environment variable")
	logFile = flag.String(
		"log-file",
		"",
		"Logfile of the observed application.\nAlternatively set the LOG_FILE environment variable")
	externalHostPort = flag.String(
		"external-host-port",
		"",
		"Port that is exposed as the external endpoint (e.g. on the docker host)Alternatively set the EXTERNAL_HOST_PORT environment variable")
	externalHostAddress = flag.String(
		"external-host-ip",
		"",
		"Exposed IP address of the host that runs this service/container. Alternatively set the EXTERNAL_HOST_IP environment variable")
	serviceLabel = flag.String(
		"service-label",
		"",
		"Label that identifies this service on the registry. Alternatively set the SERVICE_LABEL environment variable")
	serviceCommand = flag.String(
		"run",
		"",
		"Executable of the service. Pass with quotation marks if additional arguments are given.")
)

var (
	config struct {
		registryAddress     string
		logFile             string
		externalHostPort    string
		externalHostAddress string
		serviceLabel        string
		serviceID           string
		serviceCommand      string
	}
)

func loadConfig() {
	flag.Parse()

	config.registryAddress = takeOrElse(
		*registryAddress, os.Getenv("REGISTRY_ADDRESS"))
	config.logFile = takeOrElse(
		*logFile, os.Getenv("LOG_FILE"))
	config.externalHostPort = takeOrElse(
		*externalHostPort, os.Getenv("EXTERNAL_HOST_PORT"))
	config.externalHostAddress = takeOrElse(
		*externalHostAddress, os.Getenv("EXTERNAL_HOST_ADDRESS"))
	config.serviceLabel = takeOrElse(
		*serviceLabel, os.Getenv("SERVICE_LABEL"))
}

// takeOrElse returns the first argument if not empty, otherwise the second
func takeOrElse(this string, that string) string {
	if this != "" {
		return this
	}
	if that != "" {
		return that
	}
	log.Fatal("Missing argument. Run with -help for details")
	return ""
}

package runner

import (
	"flag"
	"os"
)

var (
	registryAddress = flag.String(
		"registry-address",
		"",
		"Address of the service registry. Alternatively set the REGISTRY_ADDRESS environment variable")
	logFilePath = flag.String(
		"log-file",
		"",
		"Logfile of the observed application. Alternatively set the LOG_FILE_PATH environment variable")
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
	serviceArguments = flag.String(
		"run",
		"",
		"Executable of the service. Pass with quotation marks if additional arguments are given.")
)

var (
	loggingHubAddress string
)

func Main() {
	flag.Parse()

	*registryAddress = takeOrElse(
		*registryAddress, os.Getenv("REGISTRY_ADDRESS"))
	*logFilePath = takeOrElse(
		*logFilePath, os.Getenv("LOG_FILE_PATH"))
	*externalHostPort = takeOrElse(
		*externalHostPort, os.Getenv("EXTERNAL_HOST_PORT"))
	*externalHostAddress = takeOrElse(
		*externalHostAddress, os.Getenv("EXTERNAL_HOST_ADDRESS"))
	*serviceLabel = takeOrElse(
		*serviceLabe, os.Getenv("SERVICE_LABEL"))

	registerService()
	startHeartBeat()
	startLogging()

	startService

}

// takeOrElse returns the first argument if not empty, otherwise the second
func takeOrElse(this string, that string) string {
	if this != "" {
		return this
	}
	if that != "" {
		return that
	}
	return ""
}

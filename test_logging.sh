#/bin/bash
export EXTERNAL_HOST_ADDRESS=127.0.0.1
export EXTERNAL_HOST_PORT=20000
export LOG_FILE=hawai-logging.log
export SERVICE_COMMAND="/home/daniel/projects/src/gitlab.com/danck/hawai-logginghub/hawai-logginghub"
export SERVICE_LABEL=logging
export REGISTRY_ADDRESS='http://192.168.29.140:32000/service'

./hawai-runner

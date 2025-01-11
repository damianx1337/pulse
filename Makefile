IMG_NAME=pulse
#VTAG=$(shell date +"%y%m%d-%H%M")
#VTAG=$(shell date +"%y%m%d")-$(shell git branch --show-current)
VTAG=$(shell date +"%y%m%d")
CT_NAME=${IMG_NAME}_ct
MAINTAINER=$(shell whoami)
CT_FILE=Containerfile_localtesting


build:
	podman build -f ${CT_FILE} -t ${MAINTAINER}/${IMG_NAME}:${VTAG}
	podman image prune -f
	podman builder prune -f
	podman system prune -f

run:
	podman run -it --rm --name ${CT_NAME} ${MAINTAINER}/${IMG_NAME}:${VTAG}

view-logs:
	podman logs ${CT_NAME}

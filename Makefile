BUILD_DIR = build

SRC_FOLDER = 
APP_NAME = money-clock
APP_ID = "com.github.sebastianrau.money-clock"

GIT_VERSION_TAG=$(shell git describe --tags --abbrev=0)
GIT_VERSION_TAG_FULL=$(shell git describe --tags --abbrev=2)
GIT_BUILD=$(shell git rev-parse --short HEAD)
GIT_DATE=$(shell git log -1 --date=format:"%Y/%m/%d" --format="%ad" )

ifeq ($(GIT_VERSION_TAG),)
	GIT_VERSION_TAG = 0.0.0
	GIT_VERSION_TAG_FULL = 0.0.0-000000
endif

APP_TAGS = "build=${GIT_BUILD}","date=${GIT_DATE}","tag=${GIT_VERSION_TAG_FULL}"




version:
	sed -ie "s/Version = \"*.*.*\"/Version = \"${GIT_VERSION_TAG}\"/" FyneApp.toml

app: app.windows64 app.darwin app.darwinArm app.linux64

app.windows64: version
	cd ${BUILD_DIR} && GOARCH=amd64 fyne package -os windows -icon logo.png --src ../${SRC_FOLDER} --appVersion ${GIT_VERSION_TAG} --release --tags ${APP_TAGS} --appID ${APP_ID} --name ${APP_NAME}
# TODO add zip of package

app.darwin: version	
	cd ${BUILD_DIR} && GOARCH=amd64 fyne package -os darwin -icon logo.png --src ../${SRC_FOLDER} --appVersion ${GIT_VERSION_TAG} --release --tags ${APP_TAGS} --appID ${APP_ID} --name ${APP_NAME}

app.darwinArm: version
	cd ${BUILD_DIR} && GOARCH=arm64 fyne package -os darwin -icon logo.png --src ../${SRC_FOLDER} --appVersion ${GIT_VERSION_TAG} --release --tags ${APP_TAGS} --appID ${APP_ID} --name ${APP_NAME}.arm


app.linux64: version
	cd ${BUILD_DIR} && GOARCH=amd64 fyne package -os darwin -icon logo.png --src ../${SRC_FOLDER} --appVersion ${GIT_VERSION_TAG} --release --tags ${APP_TAGS} --appID ${APP_ID} --name ${APP_NAME}	


app.linuxArm: version
	cd ${BUILD_DIR} && GOARCH=arm64 fyne package -os darwin -icon logo.png --src ../${SRC_FOLDER} --appVersion ${GIT_VERSION_TAG} --release --tags ${APP_TAGS} --appID ${APP_ID} --name ${APP_NAME}	


cli-lint:
	golangci-lint run *.go

clean:
	-rm -f ${BUILD_DIR}/*
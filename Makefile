BUILD_DIR = build

SRC_FOLDER = 
APP_NAME = money-clock
APP_ID = "com.github.sebastianrau.money-clock"

GIT_VERSION_TAG=$(shell git describe --tags --abbrev=0)
GIT_VERSION_TAG_FULL=$(shell git describe --tags --abbrev=2)
GIT_BUILD=$(shell git rev-parse --short HEAD)
GIT_DATE=$(shell git log -1 --date=format:"%Y/%m/%d" --format="%ad" )

APP_TAGS = "build=${GIT_BUILD}","date=${GIT_DATE}","tag=${GIT_VERSION_TAG_FULL}"

version:
	sed -ie "s/Version = \"*.*.*\"/Version = \"${GIT_VERSION_TAG}\"/" FyneApp.toml

app: app.windows64 app.darwin app.darwinArm app.linux64

app.windows64: version
	cd ${BUILD_DIR} && GOARCH=amd64 fyne package -os windows -icon ../../logo.png --src ../${SRC_FOLDER} --appVersion ${GIT_VERSION_TAG} --release --tags ${APP_TAGS} --appID ${APP_ID} --name ${APP_NAME}
# TODO add zip of package
app.darwin: version	
	cd ${BUILD_DIR} && GOARCH=amd64 fyne package -os darwin -icon ../../logo.png --src ../${SRC_FOLDER} --appVersion ${GIT_VERSION_TAG} --release --tags ${APP_TAGS} --appID ${APP_ID} --name ${APP_NAME}
	cd ${BUILD_DIR} && zip -vr ${APP_NAME}.app.zip  ${APP_NAME}.app

app.darwinArm: version
	cd ${BUILD_DIR} && GOARCH=arm64 fyne package -os darwin -icon ../../logo.png --src ../${SRC_FOLDER} --appVersion ${GIT_VERSION_TAG} --release --tags ${APP_TAGS} --appID ${APP_ID} --name ${APP_NAME}.arm
	cd ${BUILD_DIR} && zip -vr ${APP_NAME}.app.arm.zip ${APP_NAME}.arm.app

app.linux64: version
	cd ${BUILD_DIR} && GOARCH=amd64 fyne package -os darwin -icon ../../logo.png --src ../${SRC_FOLDER} --appVersion ${GIT_VERSION_TAG} --release --tags ${APP_TAGS} --appID ${APP_ID} --name ${APP_NAME}	
# TODO add zip of package

app.linuxArm: version
	cd ${BUILD_DIR} && GOARCH=arm64 fyne package -os darwin -icon ../../logo.png --src ../${SRC_FOLDER} --appVersion ${GIT_VERSION_TAG} --release --tags ${APP_TAGS} --appID ${APP_ID} --name ${APP_NAME}	
# TODO add zip of package

streamdeck.icons:
	cd streamdeck/ && zip -vr 'Monitor Control Icons.streamDeckIconPack' com.github.sebastianraufocusrite-mackie-control.sdIconPack/ -x "*.DS_Store"
	
cli-lint:
	golangci-lint run  cmd/... pkg/...

clean:
	-rm -f ${BUILD_DIR}/*
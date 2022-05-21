# make ARCH="amd64"
# make ARCH="386"
BASEDIR=$(DESTDIR)/usr/bin
ARCH=amd64
VERSION=1.3.0

default:
	@#echo "Generating binary data..."
	@#go-bindata public/... templates/...
	@echo "Building..."
	env GOOS=linux GOARCH=${ARCH} go build -ldflags="-s -w -X main.RunMode=production" speedtest
	@#echo "Compressing..."
	@#upx -q pyxsoftUI

.PHONY: install
upload: default
	#rsync -a -P --rsync-path="sudo rsync" -e "ssh -i ~/.ssh/cdn-global.pem" powerwaf admin@3.217.61.177:/opt/ui/
	rsync -a -P -e "ssh -i ~/.ssh/cdn-global.pem" speedtest root@www.webspeed.tools:/opt/speedtest/
	@echo "Restarting pwafui..."
	ssh -i ~/.ssh/cdn-global.pem root@www.webspeed.tools "systemctl restart speedtest"

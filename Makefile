all: build

build:
	@go build -o mercimek

release:
	@goxc
	@rmdir debian/

deploy: release
	@scp release/0.1/mercimek*.deb mercimek:
	@ssh do 'sudo dpkg -i mercimek*.deb'
	@ssh do 'sudo systemctl restart mercimek'

.PHONY: all build vet test release deploy

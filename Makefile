all: build

build:
	@go build -o mercimek

release:
	@goxc
	@rmdir debian/

deploy: release
	@scp release/0.1/mercimek*.deb mercimek:
	@ssh ilber 'sudo dpkg -i mercimek*.deb'
	@ssh ilber 'sudo service mercimek restart'

.PHONY: all build vet test release deploy

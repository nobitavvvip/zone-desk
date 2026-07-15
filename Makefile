.PHONY: build dev package clean

build:
	./scripts/build.sh

dev:
	./scripts/dev.sh

clean:
	rm -rf target/ frontend/dist/
	rm -rf frontend/node_modules/

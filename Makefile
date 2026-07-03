.PHONY: build frontend backend dev package install upgrade clean

build: frontend backend

frontend:
	cd frontend && npm install && npm run build

backend:
	cd backend && go build -o ../build/zonedesk ./cmd/server

dev:
	./scripts/dev.sh

package:
	./scripts/package.sh

install:
	./scripts/install.sh

upgrade:
	./scripts/upgrade.sh

clean:
	rm -rf build/ frontend/dist/
	rm -rf frontend/node_modules/

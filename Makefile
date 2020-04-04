dev: prepare
	cd backend && go run . -tlscert server.crt -tlskey server.key

prepare:
	cd frontend && yarn run build
	rm -rf backend/dist
	cp -r frontend/dist backend/dist
	cp backend/manifest.webmanifest backend/dist/
	cp backend/doit.png backend/dist/
	cp backend/loader.gif backend/dist/
	cp backend/service-worker.js backend/dist

dist: prepare
	cd backend && go build -o doit
	cp -r backend/doit backend/dist backend/doit.png backend/manifest.webmanifest backend/service-worker.js doit/
	tar cvfz doit.tar.gz doit

clean:
	rm -f doit.tar.gz
	rm -rf doit


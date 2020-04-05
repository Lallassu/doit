dev: prepare
	cd backend && go run . -tlscert server.crt -tlskey server.key -hostport 0.0.0.0:8443

prepare:
	cd frontend && yarn install && yarn run build
	rm -rf backend/dist
	cp -r frontend/dist backend/dist
	cp backend/manifest.webmanifest backend/dist/
	cp backend/doit.png backend/dist/
	cp backend/loader.gif backend/dist/
	cp backend/service-worker.js backend/dist

dist: prepare
	mkdir doit
	cd backend && go build -o doit
	cp -r backend/doit backend/dist backend/doit.png backend/manifest.webmanifest backend/service-worker.js doit/
	tar cvfz doit.tar.gz doit
	rm -rf doit

clean:
	rm -f doit.tar.gz
	rm -rf doit


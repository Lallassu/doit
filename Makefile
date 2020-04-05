dev: prepare
	cd backend && go run . -tlscert server.crt -tlskey server.key -hostport 0.0.0.0:8443

prepare:
	cd frontend && yarn install && yarn run build
	mkdir -p doit
	cp -r frontend/dist doit/dist

dist: prepare
	mkdir -p doit
	cd backend && go build -o doit
	cp -r backend/doit doit/
	tar cvfz doit.tar.gz doit
	rm -rf doit

clean:
	rm -f doit.tar.gz
	rm -rf doit


dev: prepare
	go run . -tlscert cert.crt -tlskey cert.key

prepare:
	cd ../frontend && npm run build
	rm -rf ../backend/dist
	cp -r ../frontend/dist ../backend/dist
	cp manifest.webmanifest ./dist
	cp doit.png ./dist
	cp loader.gif ./dist
	cp service-worker.js ./dist

dist: prepare
	go build -o doit
	tar cvfz doit.tar.gz doit dist doit.png manifest.webmanifest service-worker.js


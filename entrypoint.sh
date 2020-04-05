#!/bin/sh
if ! touch /data/.verify_access; then
  log "ERROR: /data doesn't seem to be writable. Please make sure attached directory is writable by uid=$(id -u)"
  exit 2
fi
# Remove test-file
rm /data/.verify_access || true

# If we don't have a TLSCERT we generate a simple self-signed one
if [ ! -f $DOIT_TLSCERT ]; then
  echo "INFO: $DOIT_TLSCERT not found, generating self signed certificate"
  openssl req -new -newkey rsa:2048 -days 3650 -nodes -x509 -keyout $DOIT_TLSKEY -out $DOIT_TLSCERT -subj "/C=US/ST=Oregon/L=Portland/O=Company Name/OU=Org/CN=www.example.com"
fi

# Run
/doit/doit -hostport $DOIT_HOST:$DOIT_PORT -mailfrom $DOIT_MAILFROM -database $DOIT_DATABASE -tlscert $DOIT_TLSCERT -tlskey $DOIT_TLSKEY



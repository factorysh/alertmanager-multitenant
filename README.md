# Alertmanager Multitenant

[![Build Status](https://drone.bearstech.com/api/badges/factorysh/alertmanager-multitenant/status.svg)](https://drone.bearstech.com/factorysh/alertmanager-multitenant)

## Description

A middleware package for validate and authorized an alert posting with the alertmanager api, using JWT.
Reject the request if the JWT is bad or if the "project" label is different in the JWT and the body. Authorized just the alert posting

## Demo

There is a main and a docker-compose file in the `demo` directory. It run alertmanager and mailhog for testing the alerts posting.
The main build a proxy (called `demo` on the compose file) : send a request to the proxy and check on the alertmanager (:9093) and mailhog (:8025) interface

	make demo

Know you can send a http request to the proxy

	curl -i
	-H "JWT: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o" \
	-H "Content-Type: application/json" \
	-d '[{"labels": {"project": "testing_project"}}]' \
	MY_DOCKER_IP:9000/api/v2/alerts

The JWT for testing was generated with https://jwt.io/, with the secret : `secret` ... hardcoded in the main

You can test the authorization with errors

	curl -i
	-H "JWT: bad_jwt" \
	-H "Content-Type: application/json" \
	-d '[{"labels": {"project": "testing_project"}}]' \
	MY_DOCKER_IP:9000/api/v2/alerts

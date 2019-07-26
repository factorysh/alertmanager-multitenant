#! /bin/sh

cd demo && make demo
curl -i -H "JWT: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o" \
	-H "Content-Type: application/json" \
	-d '[{"labels": {"project": "test"}}]' \
	0.0.0.0:9000/api/v2/alerts | grep 200

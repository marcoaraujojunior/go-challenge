#Go Challange

[![Build Status](https://travis-ci.org/marcoaraujojunior/go-challenge.svg?branch=master)](https://travis-ci.org/marcoaraujojunior/go-challenge)
[![Coverage Status](https://coveralls.io/repos/github/marcoaraujojunior/go-challenge/badge.svg)](https://coveralls.io/github/marcoaraujojunior/go-challenge)

Micro app built in golang to storage simple invoice data

##Usage

docker-compose up -d


##To Create New Invoice

curl -H "Content-Type: application/json" -H 'Authorization:Basic dXNlcm5hbWU6cGFzc3dvcmQ=' -X POST --data '{"referencemonth": 9, "referenceyear": 2016, "document": "1234", "description": "Test Invoice 1234",  "amount": 56.78 }' http://localhost:8888/v1/invoice


##To Update Invoice

curl -H "Content-Type: application/json" -H 'Authorization:Basic dXNlcm5hbWU6cGFzc3dvcmQ=' -X PUT --data '{"description": "Test Update Invoice 1234"}' http://localhost:8888/v1/invoice/1234

##To Delete Invoice (Soft Delete)

curl -H 'Authorization:Basic dXNlcm5hbWU6cGFzc3dvcmQ=' -X DELETE http://localhost:8888/v1/invoice/1234

##To list all invoices

curl -H 'Authorization:Basic dXNlcm5hbWU6cGFzc3dvcmQ=' http://localhost:8888/v1/invoices

##To list invoices with filter and sort

http://localhost:8888/v1/invoices?sort=-id,created_at&month=9&year=2016&document=1234

### Filter

You can filter this fields:

month
year
document
Ex: month=9&year=2016&document=1234

### Sort

To sort descending use -
Ex: -id, -created_at

To sort ascending use just fields:
Ex: id, created_at

License
----

MIT

# prosig-code-challenge
Prosigliere backend code challenge

## Requirements
* Go 1.25.2
* Docker

## Instructions

To run using docker execute the following commands in the root folder:
```bash
docker build -t app .
docker run -p 8080:8080 app:latest
```

To run locally run the following command in the root folder:
```bash
./local.sh
```

To run unit tests run the following command in the root folder:
```bash
./unit_tests.sh
```

## Author
* Matias Kopp (koppmatias97@gmail.com)
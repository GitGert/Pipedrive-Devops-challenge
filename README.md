# Pipedrive-Devops-challenge

All of the tasks requirements can be found [in this pdf](Software_Engineering_Intern_in_DevOps_test_task.pdf)

## Description of every finished part

### PART 1  - API

I wrote a simple http server using [gorilla mux](https://github.com/gorilla/mux) with 3 endpoints (GET /deals, POST /deals, PUT /deals).

For making api calls I stored API_TOKEN and COMPANY_DOMAIN in .env file where they are being read into constants.go

GET /deals - sends a request the Pipedrive endpoint with a maximum of 20 deals to be returned.

POST /deals - expects a body to be included in the request, the body will be parsed and forwarded to the Pipedrive "Add a deal" endpoint

PUT /deals - expects a body and a dealId to be in the request, both of them will be forwarded to the Pipedrive "Update a deal" endpoint


### PART 2 - Instrumentation

LOG EVERYTHING -
I decided to log all request and the status codes that my server gives (note that a request to my server can yield a 200 response while the Pipedrive API gave an error. The reason
I decided to do it this way is because I see myself as a middle man, and I am not validating the data that is being sent to my server for possible issues with Pipedrive API. My 
server does validate fields and the data type in those fields but not the correctness of the request.)

all the logs have a timestamp and are colored red or green to find problems at a glance.

GET /metrics - sends a request to all other 3 endpoints and returns the response time/request duration of each one.

### PART 3 - CI

I used [Github actions](.github/workflows/test_on_commit.yml) to run a golang linter and tests.


### PART 4 - CD

I wrote another [GitHub action](.github/workflows/master_merged.yml) to run only when a pull request is merged to master.

### PART 5 -Reproducibility

Created this reamde and [dockerfile](Dockerfile) for easy reproducibility

## HOW TO RUN

Before you can run anything You will need to have a Pipedrive account with at least one deal.
you will also need to add a .env file to the project root with your pipedrive COMPANY_DOMAIN and API_TOKEN
the file should look like this:
```bash
COMPANY_DOMAIN=fakecompanyName1
API_TOKEN=fakeAPIToken
```

There are two options on how to run this web server. First and recommended on is with docker the second one is doing it locally.

### how to run - DOCKER

You will need to have [docker](https://www.docker.com/) installed on your machine. I recommend using [docker desktop](https://www.docker.com/products/docker-desktop/)

to build the docker container run:
```bash
docker build -t pipedrive_app:latest .
```

Now you can either go to your docker desktop app and run the container there.

Or use the command line:
```bash
docker build -t pipedrive_app:latest .
docker run -p 8080:8080 pipedrive_app
```

### how to run - LOCAL

You will need to have [Golang](https://go.dev/) installed on your machine.

To install all dependencies run:
```bash
go mod download
```

To run the server:
```bash
go run main.go
```
### easy testing
After running the server using either of the methods you can test the endpoints using a tool of your choice or:

inside of bash_test_requests are 4 curl request for testing each of the endpoints.
you can call them from root like this:
```bash
sh bash_test_requests/<filename>.sh
```
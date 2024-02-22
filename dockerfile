FROM golang:1.21.6-alpine

#Set the working directory to /app
WORKDIR /app

#Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

#Download and install any required Go dependencies
RUN go mod Download

#Copy the entire source code to the working directory
COPY . .

#Build the Go application
RUN go build -o main . 

#Expose the port specified by the PORT environment variable
EXPOSE 3000

#Set the entry point of the container to the executable
CMD ["./main"]

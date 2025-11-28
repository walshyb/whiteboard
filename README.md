# Real-time Whiteboard

My own Muro / Excalidraw clone!

This just started off as an idea for a real-time whiteboard app, purely for learning purpose because I was recently asked to design this actually in a systems design interview. I wasn't expected to know the workings of web sockets and was able to cobble together a solution, but realized I had HUGE gaps in my knowledge that I definitely need to fill.

I originally thought "oh yeah, this can just be a quick little project where I just throw down a canvas, hook up a little websocket, and badda bing!" But as I continued down the path of learning, I kept finding more things that I should learn before progressing, or kept adding my own scope creep as I got more ideas.

All in all, this project (thus far) has helped me learn more about: Golang, websockets, HTML canvas, MongoDB (and nosql), Protobufs (AHH I love protobufs; idk how i ever built apis without server/client type safety??)

## System Requirements

- Node ^25.2.x
- [Protoc](https://protobuf.dev/installation/) ^25.2
- Go ^1.25.x
- make ^3.x
- Redis
- MongoDB

## Getting Started

### clone the repo

```bash
git clone git@github.com:walshyb/whiteboard.git && cd whiteboard
```

## Using Docker:

### Build the docker images and run the mongo, redis, and server containers:

```bash
# Note: Default behavior binds to local ports, so ensure that ports 8080, 27017, and 6379 are free
docker compose up --build -d
```

### Run the frontend

```bash
# Install Client deps
cd client && npm install
cd .. and make
cd client && npm run dev
```

## Full Local Development:

### 2. Download and build dependencies

```bash
# Install Client deps
cd client
npm install
cd ..

# Install Go deps
go mod download

# Build protobuf files
make
```

### 3. Run the services

- Run Redis (system-dependent)
- Run MongoDB (system-dependent)

```bash
# Run the backend
go run .

# And in a new terminal shell in current workspace,
# Run the client:
cd client
npm run dev
```

_Note_: Currently the frontend and backend run separately.

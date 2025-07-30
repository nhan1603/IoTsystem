# IoTSystem

Simulation of IoT data ingestion system

### Backend Setup

- Install go 1.23

- Install `make` command

- cd `/api`

- `go mod vendor`

- Setup Docker and Docker Compose, and make sure it is running

- `cd ..`

- Setup project by running `make setup`

- Seed the necessary data into the database: `make pg-drop` then `make pg-migrate`

- Start server with `make run`

- The API URLs for checking back-end availability is: `http://localhost:3001/api/public/ping`

### Frontend Setup

- Install nodejs

- cd `/web`

- `npm install`

- `npm start`

- The Web URL: `http://localhost:3000`

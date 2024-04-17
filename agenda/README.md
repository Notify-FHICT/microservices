# Agenda Service

The Agenda Service is a microservice designed to manage events, notes, and tags for users. It provides functionalities to create, read, update, and delete events, as well as link notes and tags to events.

## Features

- Create, read, update, and delete events
- Link and unlink notes and tags to events
- Prometheus metrics for monitoring request duration and processed requests

## Usage

### Running with Docker

You can run the Agenda Service using Docker. Simply pull the Docker image from the GitHub Container Registry:
``docker pull ghcr.io/notify-fhict/microservices/agenda:main``

### API Endpoints

The service exposes the following API endpoints:

- `/create`: Create a new event
- `/read/{id}`: Read an event by ID
- `/update`: Update an existing event
- `/delete`: Delete an event
- `/link_noteID`: Link or unlink a note ID to an event
- `/link_tagID`: Link or unlink a tag ID to an event
- `/dashboard/{userID}`: Read all events for a specific user

### Environment Variables

The Agenda Service requires the following environment variables to be set:

- `MONGO_URI`: MongoDB connection URI
- `RABBITMQ_URI`: RabbitMQ connection URI

(These values are currently hardcoded)

## Getting Started

1. Set up MongoDB and RabbitMQ instances.
2. Set the required environment variables.
3. Run the Agenda Service using Docker.
4. Access the API endpoints to manage events and notes.

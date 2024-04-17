# Agenda Service

The Agenda Service is a microservice designed to manage notes and events. It provides functionalities for creating, reading, updating, and deleting notes, as well as linking notes to events.

## Usage

### Running Locally

To run the Agenda Service locally, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/Notify-FHICT/microservices/notes.git
   ```
2. Navigate to the repository directory:

   ```bash
   cd notes
   ```
3. Build the Docker image:

   ```bash
   docker build -t agenda-service .
   ```
4. Run the Docker container:

   ```bash
   docker run -p 3000:3000 agenda-service
   ```

### Running from Docker Registry

Alternatively, you can pull the pre-built Docker image from the GitHub Container Registry:

```bash
docker pull ghcr.io/notify-fhict/microservices/notes:main
```

Then run the container:

```bash
docker run -p 3000:3000 ghcr.io/notify-fhict/microservices/notes:main
```

### API Endpoints

- **Create Note**: `POST /create`
- **Read Note**: `GET /read/:id`
- **Update Note**: `PUT /update`
- **Delete Note**: `DELETE /delete`
- **Link Tag ID**: `PUT /link_tagID`
- **Link Event**: `PUT /link_event`
- **Update Content**: `PUT /update_content`

### Environment Variables

- `MONGODB_URI`: MongoDB connection URI
- `PORT`: Port on which the server listens (default: 3000)

(These values are currently hardcoded)

## Dependencies

The Agenda Service relies on the following technologies:

- [Go](https://golang.org/): Programming language used for backend development.
- [MongoDB](https://www.mongodb.com/): NoSQL database for storing notes.
- [Docker](https://www.docker.com/): Containerization platform for easy deployment.
- [GitHub Container Registry](https://github.com/features/packages): Hosts pre-built Docker images for easy distribution.

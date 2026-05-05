# 🎵 AudioHunt (Shazam Clone)

A full-stack, end-to-end audio recognition application inspired by Shazam. AudioHunt allows users to record audio snippets through a beautifully designed web interface and matches them against a database of audio fingerprints using custom Digital Signal Processing (DSP) algorithms.

## ✨ Features

- **Audio Recognition Engine:** Custom Go-based DSP backend that extracts peaks and generates unique audio fingerprints to match recorded samples against known tracks.
- **Premium User Interface:** A modern, immersive glassmorphic UI with dark mode, smooth micro-animations, and dynamic visual feedback during audio recording.
- **Microphone Integration:** Seamless in-browser audio recording using `RecordRTC`.
- **Scalable Architecture:** Clean separation of concerns with a Go backend REST API, a React frontend, and a PostgreSQL database.

## 🛠 Tech Stack

### Frontend
- **Framework:** React 19 + TypeScript + Vite
- **Styling:** Custom CSS with glassmorphic and premium dark-mode aesthetics
- **Audio Capture:** RecordRTC for browser microphone access
- **HTTP Client:** Axios

### Backend
- **Language:** Go (Golang) 1.25
- **Routing:** Gorilla Mux
- **Audio Processing:** Custom DSP logic using `madelynnblue/go-dsp` and `go-audio/wav` for peak extraction and audio fingerprinting.
- **Database:** PostgreSQL (via Docker)

## 📂 Project Structure

```
.
├── backend/                  # Go backend application
│   ├── cmd/                  # Entry points (main.go)
│   ├── internal/             # Application code (API routing, DSP logic, song management)
│   ├── migrations/           # Database setup and schema files
│   ├── uploads/              # Temporary storage for uploaded audio chunks
│   ├── .env                  # Backend environment variables
│   └── go.mod                # Go module dependencies
├── frontend/                 # React frontend application
│   ├── public/               # Static assets
│   ├── src/                  # React components (AudioRecorder, MatchResults, etc.)
│   ├── package.json          # Node dependencies
│   └── vite.config.ts        # Vite configuration
├── docker-compose.yml        # Docker composition for the PostgreSQL database
└── README.md                 # Project documentation
```

## 🚀 Getting Started

### Prerequisites

- [Go 1.25+](https://golang.org/doc/install)
- [Node.js](https://nodejs.org/) & npm
- [Docker Desktop](https://www.docker.com/products/docker-desktop) (for running PostgreSQL)

### 1. Start the Database
The backend relies on PostgreSQL. A `docker-compose.yml` file is provided to spin up the database easily with the correct schema migrations.

```bash
# From the root directory
docker-compose up -d
```
*Note: The database runs on port `5434` to avoid conflicts with local Postgres installations.*

### 2. Run the Go Backend

The backend server processes audio chunks and communicates with the database.

```bash
cd backend
go mod download
go run cmd/main.go
```
*The backend server will typically start on `http://localhost:8080`.*

### 3. Run the React Frontend

Open a new terminal window to start the Vite development server.

```bash
cd frontend
npm install
npm run dev
```
*The frontend will be accessible at `http://localhost:5173`. Open this in your browser to interact with the application!*

## 🧠 How it Works (Under the Hood)

1. **Recording:** The user clicks the record button on the frontend. `RecordRTC` captures the audio via the browser's MediaRecorder API and converts it into a valid audio format (like `.wav`).
2. **Transmission:** The chunk is sent to the backend `/api/recognize` endpoint via a `multipart/form-data` POST request.
3. **Signal Processing:** The Go backend reads the audio file, applies an FFT (Fast Fourier Transform), and finds constellation peaks to generate a unique fingerprint for the sample.
4. **Matching:** The generated fingerprint hashes are matched against the pre-computed hashes of songs in the PostgreSQL database.
5. **Results:** The backend returns the best-matching song with its confidence score, which the frontend displays to the user in a beautiful result card.

## 📝 License
This project is for educational and portfolio purposes.

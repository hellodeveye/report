# Feishu Report Assistant

This project is a web application designed to help users automatically generate draft weekly and monthly reports by pulling and summarizing their daily reports from Feishu.

## Project Structure

The project is organized into three main directories:

-   `/backend`: The Go application that handles business logic, communicates with the Feishu Open Platform, and serves the API.
-   `/frontend`: The Vue.js single-page application that provides the user interface.
-   `/deployment`: Contains Docker, Docker Compose, and Kubernetes configuration files for deployment.

## Tech Stack

-   **Backend**: Go, Gorilla/Mux, Redis
-   **Frontend**: Vue.js, Vite, Tailwind CSS, Tiptap Editor
-   **Deployment**: Docker, Docker Compose, Kubernetes

## Getting Started

### Prerequisites

-   Docker and Docker Compose
-   `kubectl` for Kubernetes deployment
-   A Feishu App with `App ID` and `App Secret`

### Running with Docker Compose (Recommended for Local Development)

1.  **Configure Feishu App Credentials:**
    Create a `.env` file in the root directory and add your Feishu app details:

    ```env
    FEISHU_APP_ID=your_app_id
    FEISHU_APP_SECRET=your_app_secret
    ```

2.  **Build and Run:**

    ```bash
    docker-compose up --build
    ```

3.  **Access the application:**
    -   Frontend: `http://localhost:5173`
    -   Backend API: `http://localhost:8080`

### Running with Kubernetes

Navigate to the `deployment` directory and apply the configurations:

```bash
# Deploy Redis
kubectl apply -f redis.yaml

# Deploy Backend
kubectl apply -f backend.yaml

# Deploy Frontend
kubectl apply -f frontend.yaml
```

Make sure to configure secrets for the Feishu App credentials within your Kubernetes cluster. 
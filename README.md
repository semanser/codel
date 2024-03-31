<img src="./.github/demo.png" />
<div align="center">Fully autonomous AI Agent that can perform complicated tasks and projects using terminal, browser, and editor.</div>
</br>

**Discord: https://discord.gg/uMaGSHNjzc**

# Features
- 🔓 Secure. Everything is running in a sandboxed Docker environment.
- 🤖 Autonomous. Automatically detects the next step and performs it.
- 🔍 Built-in browser. Fetches latest information from the web (tutorials, docs, etc.) if needed.
- 📙 Built-in text editor. View all the modified files right in your browser.
- 🧠 All the history commands and outputs are saved in the PostgreSQL database.
- 📦 Automatic Docker-image picker based on the user task.
- 🤳 Self-hosted
- 💅 Modern UI

# How to run
## Prerequisites
- golang
- nodejs
- docker
- postgresql

## Environment variables
First, run `cp ./backend/.env.example ./backend/.env && cp ./frontend/.env.example ./frontend/.env.local` to generate the env files for both backend and frontend.

### Backend
Edit the `.env` file in `backend` folder
- `OPEN_AI_KEY` - OpenAI API key
- `DATABASE_URL` - PostgreSQL database URL (eg. `postgres://user:password@localhost:5432/database`)
- `DOCKER_HOST` - Docker SDK API (eg. `DOCKER_HOST=unix:///Users/<my-user>/Library/Containers/com.docker.docker/Data/docker.raw.sock`) [more info](https://stackoverflow.com/a/62757128/5922857)

Optional:
- `PORT` - Port to run the server (default: `8080`)
- `OPEN_AI_MODEL` - OpenAI model (default: `gpt-4-0125-preview`). The list of supported OpenAI models can be found [here](https://pkg.go.dev/github.com/sashabaranov/go-openai#pkg-constants).
### Frontend
Edit the `.env.local` file in `frontend` folder
- `VITE_API_URL` - Backend API URL. *Omit* the URL scheme (e.g., `localhost:8080` *NOT* `http://localhost:8080`).

## Steps
### Backend
Run the command(s) in `backend` folder
- Run `go run .` to start the server

>The first run can be a long wait because the dependencies and the docker images need to be download to setup the backend environment.
When you see output below, the server has started successfully:
```
<your-date> <your-time> connect to http://localhost:<your-port>/playground for GraphQL playground
```
### Frontend
Run the command(s) in `frontend` folder
- Run `yarn` to install the dependencies
- Run `yarn dev` to run the web app

### Enjoy
- Open your browser and visit the web app URL to enjoy!

# Roadmap
- [x] Agent API
- [x] Frontend
- [x] Backend API + PostgreSQL integration
- [x] Docker runner
- [x] Terminal output streaming
- [ ] Editor output
- [ ] SWE-bench
- [ ] Better way to run it (eg a single docker command)

See more detailed roadmap [here](https://github.com/semanser/codel/milestones).



# Credits
This project wouldn't be possible without:
- https://arxiv.org/abs/2308.00352
- https://arxiv.org/abs/2403.08299
- https://www.cognition-labs.com/introducing-devin
- https://github.com/go-rod/rod
- https://github.com/semanser/JsonGenius

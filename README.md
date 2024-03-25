<div align="center">
  <img src="./.github/logo.png" width="200" />
</div>

<div align="center">Fully autonomous AI Agent that can perform complicated tasks and projects using terminal, browser, and editor.</div>

<img src="./.github/demo.png" />

Discord: https://discord.gg/uMaGSHNjzc

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

## Environment variables
Add to `.env` file in the `backend` folder.

### Backend
- `OPEN_AI_KEY` - OpenAI API key
- `DATABASE_URL` - PostgreSQL database URL (eg. `postgres://user:password@localhost:5432/database`)
- `DOCKER_HOST` - Docker SDK API (eg. `DOCKER_HOST=unix:///Users/<my-user>/Library/Containers/com.docker.docker/Data/docker.raw.sock`) [more info](https://stackoverflow.com/a/62757128/5922857)
### Frontend
Frontend env variables can be set by creating a `.env.local` file in the `frontend` folder.
- `VITE_API_URL` - Backend API URL (eg. `localhost:8080`)

## Steps
- Run `go run .` in `backend` folder
- Run `yarn dev` in `frontend` folder
- Open your browser and enjoy!

# Roadmap
- [x] Agent API
- [x] Frontend
- [x] Backend API + PostgreSQL integration
- [x] Docker runner
- [x] Terminal output streaming
- [ ] Browser output streaming (in progress)
- [ ] Editor output
- [ ] SWE-bench
- [ ] Better way to run it (eg a single docker command)



# Credits
This project wouldn't be possible without:
- https://arxiv.org/abs/2308.00352
- https://arxiv.org/abs/2403.08299
- https://www.cognition-labs.com/introducing-devin
- https://github.com/semanser/JsonGenius

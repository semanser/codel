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

# Usage
The simplest way to run Codel is to use a pre-built Docker image. You can find the latest image on the [Github Container Registry](https://github.com/semanser/codel/pkgs/container/codel).


> [!IMPORTANT]
> Don't forget to set the required environment variables.

```bash
docker run -d \
  -e OPEN_AI_KEY=<your_open_ai_key> \
  -p 3000:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  ghcr.io/semanser/codel:latest
```

Alternatively, you can create a .env file and run the Docker image with the following command:
```bash
docker run -d \
  --env-file .env \
  -p 3000:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  ghcr.io/semanser/codel:latest
```

<details>
  <summary>Required environment variables</summary>

    - `OPEN_AI_KEY` - OpenAI API key
</details>

<details>
    <summary>Optional environment variables</summary>

    - `OPEN_AI_MODEL` - OpenAI model (default: gpt-4-0125-preview). The list of supported OpenAI models can be found [here](https://pkg.go.dev/github.com/sashabaranov/go-openai#pkg-constants).
    - `DATABASE_URL` - PostgreSQL database URL (eg. `postgres://user:password@localhost:5432/database`)
    - `DOCKER_HOST` - Docker SDK API (eg. `DOCKER_HOST=unix:///Users/<my-user>/Library/Containers/com.docker.docker/Data/docker.raw.sock`) [more info](https://stackoverflow.com/a/62757128/5922857)
    - `PORT` - Port to run the server in the Docker container (default: 8080)

    See backend [.env.example](./backend/.env.example) for more details.
</details>

# Development

Check out the [DEVELOPMENT.md](./DEVELOPMENT.md) for more information.

# Roadmap

You can find the project's roadmap [here](https://github.com/semanser/codel/milestones).

# Credits
This project wouldn't be possible without:
- https://arxiv.org/abs/2308.00352
- https://arxiv.org/abs/2403.08299
- https://www.cognition-labs.com/introducing-devin
- https://github.com/go-rod/rod
- https://github.com/semanser/JsonGenius

# cligpt

A simple helper for google gemini that you can ask for a command. Check it and run if want.

## Requirement

A env var for google api or gemini api

```bash
export GEMINI_API_KEY="your-api-key"

# OR

export GOOGLE_API_KEY="your-api-key"
```

## Install

```bash
go install github.com/bhrott/cligpt@0.0.2
```

## Usage

```bash
cligpt list all docker container images
```
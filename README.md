# pisces

[![.github/workflows/ci.yml](https://github.com/mjc-gh/pisces/actions/workflows/ci.yml/badge.svg)](https://github.com/mjc-gh/pisces/actions/workflows/ci.yml)

A tool for analyzing phishing attack sites

## Development

You can build the CLI tool with the following:

```
make build.cli
./build/pisces -h
```

## Usage

```
NAME:
   pisces - A tool for analyzing phishing sites

USAGE:
   pisces [global options] [command [command options]]

VERSION:
   0.0.0

COMMANDS:
   analyze     Analyze and interact one or more URLs for phishing
   collect     Collect HTML and assets for one or more URLs
   screenshot  Screenshot one or more URLs
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Result Output

```json
{
  "action": "fetch",
  "elapsed": 1500000000,
  "url": "https://api.example.com/data",
  "result": {
    "location": "https://example.com/final",
    "redirectLocations": [
      {
        "status_code": 301,
        "location": "https://example.com/redirect1"
      },
      {
        "status_code": 302,
        "location": "https://example.com/final"
      }
    ],
    "body": "<html><body>Final content</body></html>",
    "bodySize": 38,
    "initialBody": "<html><body>Initial content</body></html>",
    "initialBodySize": 41,
    "assetsCount": 2,
    "assets": [
      {
        "url": "https://example.com/style.css",
        "resource_type": "stylesheet",
        "request_headers": {
          "User-Agent": "Mozilla/5.0"
        },
        "response_headers": {
          "Content-Type": "text/css"
        },
        "body": "body { margin: 0; }"
      },
      {
        "url": "https://example.com/script.js",
        "resource_type": "script",
        "request_headers": {
          "Accept": "application/javascript"
        },
        "response_headers": {
          "Content-Type": "application/javascript"
        },
        "body": "console.log('hello');"
      }
    ]
  }
}
```

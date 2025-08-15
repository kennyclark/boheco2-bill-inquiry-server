# BOHECO 2 Bill Inquiry Proxy Server

A lightweight proxy server that enables querying BOHECO 2 electricity bill information. This server handles session management and forwards requests to BOHECO 2's API while maintaining proper headers and authentication.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Requirements

- Download and install [GoLang 1.24+](https://go.dev/doc/install)

## Environement Variables

The following environment variables are required:

| Variable | Description | Default |
| --- | --- | --- |
| BOHECO2_PROXY_SERVER_PORT | The port of the server | 3000 |
| BOHECO2_API_BASE_URL | The base url of BOHECO 2 Online Bill Inquiry API  | https://bill-inquiry-api.onrender.com |

## Running the Server

To run the server, run the following command:
```bash
go run main.go
```

## License

This project is licensed under the **MIT License**.You are free to use, modify, distribute, and sell this software, provided that you include the original copyright notice and license terms.


See [LICENSE](./LICENSE) for full details.
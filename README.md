# go-book-ai

## Overview

`go-book-ai` is a CLI tool for generating book chapters using AI. It leverages language models like OpenAI's GPT to create detailed outlines and content for book chapters.

## Setup

### Prerequisites

- Go 1.20 or later
- OpenAI API Key

### Installation

1. **Clone the repository**:
   ```sh
   git clone https://github.com/yourusername/go-book-ai.git
   cd go-book-ai
   ```

2. **Initialize the Go module**:
   ```sh
   go mod init go-book-ai
   go mod tidy
   ```

3. **Set the OpenAI API Key as an environment variable**:
   ```sh
   export OPENAI_API_KEY=your_openai_api_key
   ```

4. **Build the CLI application**:
   ```sh
   go build ./cmd/bookcli
   ```

## Usage

### Create a New Book

To create a new book with a specific topic:
```sh
./bookcli new "Your Book Topic"
```

### Continue an Existing Book

To continue working on an existing book:
```sh
./bookcli continue "book_id"
```

## Testing

Run the tests to ensure everything is working correctly:
```sh
go test ./...
```

## Error Handling

The application includes robust error handling with retries and exponential backoff for network-related issues and rate limits.

## Contributing

Feel free to submit issues, fork the repository, and send pull requests. Contributions are welcome!

## License

This project is licensed under the MIT License.
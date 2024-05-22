# DocumentSigner

## Overview

DocumentSigner is a robust and secure server solution developed in Go using the Gin web framework. It efficiently manages API requests for a user-centric document management system, emphasizing strong security practices, including HTTPS, to protect sensitive documents during upload, download, and sharing.

## Features

- **Secure Document Handling**: Upload, download, and share documents securely.
- **HTTPS**: Ensures all communications are encrypted.
- **Advanced Authentication**: Robust user authentication and session management.
- **User-Centric**: Designed for seamless document management.

## Getting Started

### Prerequisites

- Go 1.16+
- Gin web framework

### Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/Sci4rr/DocumentSigner.git
   cd DocumentSigner
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Run the server:

   ```sh
   go run server.go
   ```

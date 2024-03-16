# News Search Web Application

This is a simple web application written in Go that allows users to search for news articles using the News API. The application displays search results and supports pagination.

## Features

- Search for news articles using keywords.
- Pagination for navigating through search results.
- Default landing search for when users first access the website.
- Redirects users to a specific page based on their location.

## Setup

1. Clone the repository:

    ```bash
    git clone <repository-url>
    ```

2. Navigate to the project directory:

    ```bash
    cd news-search-web-app
    ```

3. Create a `.env` file in the root directory of the project and provide your News API key:

    ```env
    NEWS_API_KEY=your_news_api_key_here
    PORT=3000
    ```

4. Install dependencies:

    ```bash
    go mod tidy
    ```

5. Run the application:

    ```bash
    go run main.go
    ```

6. Access the application in your web browser at `http://localhost:3000`.

## Usage

- Enter a search query in the search box and press Enter to search for news articles.
- Use pagination controls to navigate through search results.
- The application automatically redirects users to a specific page based on their location.

## Dependencies

- [newsapi](https://github.com/freshman-tech/news-demo-starter-files/news): Go package for interacting with the News API.
- [godotenv](https://github.com/joho/godotenv): Go package for loading environment variables from a .env file.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

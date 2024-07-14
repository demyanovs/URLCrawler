[![Go Reference](https://pkg.go.dev/badge/github.com/demyanovs/urlcrawler.svg)](https://pkg.go.dev/github.com/demyanovs/urlcrawler)

# URLCrawler

A robust and efficient web crawler written in Go (Golang). 
This project aims to provide a powerful tool for crawling and scraping web pages, 
collecting data such as page titles, links, descriptions, keywords, and response codes.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Features

- Multithreaded crawling
- Customizable crawling depth
- Respect for `robots.txt` (URL filtering and crawling delay)
- Configurable delay between requests
- Bulk saving of crawl results
- Export to JSON and CSV files
- and [more](#command-line-options)...


## Installation

### From Binary

You can use pre-built binary for your OS from the release.

### From Source

To install from the source code, you need to have Go installed on your machine. If you don't have Go installed, you can download it from the [official website](https://go.dev/dl/).

1. Clone the repository:

    ```sh
    git clone https://github.com/demyanovs/urlcrawler.git
    cd urlcrawler
    ```

2. Build the project:

    ```sh
    go build -o urlcrawler main.go
    ```

## Usage

This web crawler can be used directly from the command line after installation. It is configured through various flags that allow you to control its behavior.

### Command-Line Options

The following are the primary command-line options available for the web crawler:

- `-u` **(required)**: Specifies the starting URL for the crawler.
- `-depth`: Sets the maximum depth of crawling relative to the starting URL. Default is `0` (infinite).
- `-delay`: Determines the delay between requests in milliseconds to manage load on the server. Default is `1000`.
- `-output`: Specifies the output format for the crawl results. Supported formats are `csv` and `json`. Default is `csv`.
- `output-file`: Specifies the file path to save the crawl results. Default is `results.csv`.
- `limit`: Specifies the maximum number of pages to crawl. Default is `0` (unlimited).
- `timeout`: Specifies the maximum time in milliseconds to wait for a response. Default is `5000`.
- `bulk-size`: Specifies the number of pages to save in each bulk write operation. Default is `30`.
- `q`: quiet mode, suppresses all output except for errors. Default is `false`.
- `ignore-robots`: Ignore robots.txt rules. Default is `false`.
- `queue-len`: Specifies the number of parallel workers to use. Default is `50`.

### Basic Usage

```sh
./urlcrawler -u=https://example.com
```

With depth 2 and limit of 10 URLs:
```sh
./urlcrawler -u=https://example.com -depth=2 -limit=10
```

For the help run: 
```sh
./urlcrawler -h
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](LICENSE.md)
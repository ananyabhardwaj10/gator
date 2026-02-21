# GATOR

GATOR is a CLI tool that can fetch, store, and display blog feeds.

## Requirements

Before using GATOR, make sure you have the following installed:

- Go (1.21 or later recommended)
- PostgreSQL

## Installation

To install GATOR, run:

```bash
go install github.com/ananyabhardwaj10/gator@latest

## Configuration

Before running GATOR, you need to create a configuration file so the program can connect to your database.

Create the following file in your home directory:

~/.gatorconfig.json


Add your PostgreSQL connection string:

```json
{
  "db_url": "postgres://username:password@localhost:5432/database_name?sslmode=disable"
}

Replace the username, password and database_name with their actual values. Make sure the postgres server is running before this.

After installing the program, you can run it using gator <command> [arguments]

## Available Commands: 

gator login <username>

gator register <username>

gator reset

gator users

gator agg <time request>

gator addfeed <feed name> [feed url]

gator feeds 

gator follow <feed url>

gator following

gator unfollow <feed url>

gator browse <limit> (optional)
# Gator

## Dependencies
*You will need both of these installed on your machine to run the program:*
1. Postgres
2. Go

## Installation
To install the Gator CLI, use `go install`:

```bash
go install [github.com/josequiceno2000/gator@latest](https://www.google.com/search?q=https://github.com/josequiceno2000/gator%40latest)

## Configuration Setup

Gator requires a configuration file (`config.yaml`) to store database connection details and other settings. Here's how to set it up:

1.  **Create `config.yaml`:**
    * Create a file named `config.yaml` in the same directory as your `gator` executable.
    * Add the following content to the file, replacing the placeholders with your actual values:

        ```yaml
        db_url: "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
        current_username: ""
        ```

    * **`db_url`:** This is the connection string for your PostgreSQL database. Replace the default values with your database credentials.
    * **`current_username`:** This field will store the username of the currently logged-in user. Initially, it should be empty.

2.  **Database Setup:**
    * Ensure that you have a PostgreSQL database running and that the database specified in `db_url` exists.
    * The database must have the appropriate tables created by running the migrations. If you have not done this already, you must run the goose migrations.

3.  **Run Migrations (if needed):**
    * Navigate to the directory containing your `gator` executable in your terminal.
    * Run the following command to apply the database migrations:

        ```bash
        goose -dir sql/migrations postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
        ```

        * Replace the connection string with your actual `db_url` if it's different.
        * This command will create the necessary tables in your database.

## Running Gator

Once the configuration and database are set up, you can run the Gator CLI.

1.  **Run the `gator` command:**
    * Open your terminal and navigate to the directory containing your `gator` executable.
    * Run the `gator` command followed by the desired command and arguments.

## Example Commands

Here are a few examples of commands you can run with Gator:

* **Register a new user:**

    ```bash
    gator register <username>
    ```

    * Replace `<username>` with the desired username.
    * This command creates a new user account in the database.

* **Log in:**

    ```bash
    gator login <username>
    ```

    * Replace `<username>` with the username of the user you want to log in as.
    * This command sets the `current_username` in the `config.yaml` file.

* **Add a feed:**

    ```bash
    gator addfeed "TechCrunch" "[https://techcrunch.com/feed/](https://techcrunch.com/feed/)"
    ```

    * This command adds a new RSS feed to the database.

* **Follow a feed:**

    ```bash
    gator follow "[https://techcrunch.com/feed/](https://techcrunch.com/feed/)"
    ```

    * This command allows the currently logged-in user to follow a feed.

* **Browse posts:**

    ```bash
    gator browse
    ```

    * This command displays the most recent posts from the feeds the currently logged-in user is following.
    * You can also specify a limit: `gator browse 5`

* **Aggregate feeds:**

    ```bash
    gator agg 1m
    ```

    * This command starts the feed aggregation process, fetching and parsing RSS feeds every minute.
    * You can change the interval (e.g., `1h` for every hour).
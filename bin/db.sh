#!/usr/bin/env bash

# shellcheck source=/dev/null
. ".env"

TTY_CLIENT() {
    docker compose exec mysql mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" "${@}"
}

CMD() {
    docker compose exec -T mysql mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" -e "$1"
}

# Function to execute MongoDB commands
MONGO_CMD() {
    docker compose exec -T mongo mongosh --quiet --username "${MONGO_INITDB_ROOT_USERNAME}" --password "${MONGO_INITDB_ROOT_PASSWORD}" --eval "$1"
}

# Function to execute MongoDB commands with the MongoDB client
MONGO_CLIENT() {
    docker compose exec -T mongo mongosh --username "${MONGO_INITDB_ROOT_USERNAME}" --password "${MONGO_INITDB_ROOT_PASSWORD}" "${@}"
}

# Function to execute PostgreSQL commands
PSQL_CMD() {
    docker compose exec -T postgresql psql -U "${PG_USER}" -c "${@}"
}

# Function to execute PostgreSQL commands with the PostgreSQL client
PSQL_CLIENT() {
    docker compose exec postgresql psql -U "${PG_USER}" "${@}"
}

# Entry function to generate a random password
# Usage: ./cli.sh entry_gen_password
# Description: Generates a random password using openssl and filters out special characters
entry_gen_password() {
    openssl rand -base64 32 | tr -d '+/=' | cut -c -32
}

# Entry function to execute PostgreSQL client commands
# Usage: ./cli.sh entry_pgsql_client <COMMAND>
# Description: Executes the provided PostgreSQL command using the PostgreSQL client
entry_pgsql_client() {
    PSQL_CLIENT "${@}"
}
# Entry function to create a PostgreSQL user
# Usage: ./cli.sh entry_pgsql_create_user <USERNAME>
# Description: Creates a PostgreSQL user with the provided username, optional role, and optional database (default: postgres)
entry_pgsql_create_user() {
    local username=$1
    password=$(entry_gen_password)

    if [[ -z "$username" ]]; then
        echo "Error: No username specified."
        echo "Usage: ./cli.sh entry_pgsql_create_user <USERNAME> [<DATABASE>]"
        exit 1
    fi

    # Check if the user already exists
    existing_user=$(PSQL_CMD "SELECT 1 FROM pg_roles WHERE rolname = '$username';" -t -A)
    if [[ -n "$existing_user" ]]; then
        echo "User '$username' already exists. Skipping creation."
        exit 1
    fi

    PSQL_CMD "CREATE USER $username WITH PASSWORD '$password';" || exit 1

    echo "Creating user: $username"
    echo "Generated password: $password"
    echo "Please save this password as it will not be shown again."

    echo "User '$username' created successfully."
}

# Entry function to create a database in PostgreSQL
# Usage: ./cli.sh create_pgsql_db <OWNER> <DATABASE_NAME>
# Description: Creates a PostgreSQL database with the provided name
entry_pgsql_create_db() {
    local owner=$1
    local database=$2

    if [[ -z "$database" ]]; then
        echo "Error: No database name specified."
        echo "Usage: ./cli.sh create_pgsql_db <DATABASE_NAME>"
        exit 1
    fi

    # Check if the user already exists
    existing_user=$(PSQL_CMD "SELECT 1 FROM pg_roles WHERE rolname = '$owner';" -t -A)
    if [[ -z "$existing_user" ]]; then
        echo "User '$owner' does not exist. Skipping privilege grant."
        exit 1
    fi

    # Check if the database already exists
    existing_db=$(PSQL_CMD "SELECT 1 FROM pg_database WHERE datname='$database';" -t -A)
    if [[ -n "$existing_db" ]]; then
        echo "Database '$database' already exists. Skipping creation."
        return
    fi

    # Create the database with utf8 encoding and Asia/Shanghai timezone
    PSQL_CMD "CREATE DATABASE $database WITH OWNER $owner ENCODING 'UTF8' LC_COLLATE='en_US.utf8' LC_CTYPE='en_US.utf8' TEMPLATE=template0;"
    PSQL_CMD "ALTER DATABASE $database SET timezone TO 'Asia/Shanghai';"
    echo "Database '$database' created successfully with UTF8 encoding and timezone set to Asia/Shanghai, owned by $owner."
}

# Entry function to grant privileges to a user on a database in PostgreSQL
# Usage: ./cli.sh grant_pgsql <DATABASE_NAME> <USERNAME>
# Description: Grants all privileges on the specified database to the specified user
entry_pgsql_grant() {
    local database=$1
    local username=$2

    if [[ -z "$database" ]]; then
        echo "Error: Database name not specified."
        echo "Usage: ./cli.sh grant_pgsql <DATABASE_NAME> <USERNAME>"
        exit 1
    fi

    if [[ -z "$username" ]]; then
        echo "Error: Username not specified."
        echo "Usage: ./cli.sh grant_pgsql <DATABASE_NAME> <USERNAME>"
        exit 1
    fi

    # Check if the user already exists
    existing_user=$(PSQL_CMD "SELECT 1 FROM pg_roles WHERE rolname = '$username';" -t -A)
    if [[ -z "$existing_user" ]]; then
        echo "User '$username' does not exist. Skipping privilege grant."
        exit 1
    fi

    # Check if the database already exists
    existing_db=$(PSQL_CMD "SELECT 1 FROM pg_database WHERE datname='$database';" -t -A)
    if [[ -z "$existing_db" ]]; then
        echo "Database '$database' does not exist. Skipping privilege grant."
        return
    fi

    echo "Granting privileges on database: $database to user: $username"
    PSQL_CMD "GRANT ALL PRIVILEGES ON DATABASE \"$database\" TO \"$username\";"
    if [ $? -eq 0 ]; then
        echo "Privileges successfully granted to user: $username on database: $database"
    else
        echo "Failed to grant privileges. Please check the database and user names, and ensure you have sufficient permissions."
        exit 1
    fi

    echo "Granting privileges on public schema in database: $database to user: $username"
    PSQL_CMD "GRANT USAGE ON SCHEMA public TO \"$username\";"
    PSQL_CMD "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO \"$username\";"
    PSQL_CMD "GRANT SELECT,INSERT,UPDATE,DELETE ON ALL TABLES IN SCHEMA public TO \"$username\";"
    PSQL_CMD "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO \"$username\";"
    PSQL_CMD "GRANT ALL ON SCHEMA public TO \"$username\";"
    echo "Privileges on public schema granted."
    echo "Displaying current privileges for user: $username on database: $database"
    PSQL_CMD "SELECT datname, datcollate, datctype, datacl FROM pg_database WHERE datname = '$database';"
}

# Entry function to delete a PostgreSQL database
# Usage: ./cli.sh entry_pgsql_delete_db <DATABASE_NAME>
# Description: Deletes the specified PostgreSQL database after confirmation
entry_pgsql_delete_db() {
    local database_name="$1"

    if [[ -z "$database_name" ]]; then
        echo "Error: No database name specified."
        echo "Usage: ./cli.sh entry_pgsql_delete_db <DATABASE_NAME>"
        exit 1
    fi

    # Check if the database exists
    existing_db=$(PSQL_CMD "SELECT 1 FROM pg_database WHERE datname = '$database_name';" -t -A)
    if [[ -z "$existing_db" ]]; then
        echo "Database '$database_name' does not exist. Skipping deletion."
        return
    fi

    echo -e "Are you sure you want to delete the database '$database_name'?"
    read -rp "Please type the database name to confirm: " confirm_database
    if [[ "$confirm_database" != "$database_name" ]]; then
        echo "Database deletion cancelled."
        return
    fi

    PSQL_CMD "DROP DATABASE $database_name;"
    echo "Database '$database_name' deleted successfully."
}

# Entry function to delete a PostgreSQL user
# Usage: ./cli.sh entry_pgsql_delete_user <USERNAME>
# Description: Deletes the specified PostgreSQL user after confirmation
entry_pgsql_delete_user() {
    local username="$1"

    if [[ -z "$username" ]]; then
        echo "Error: No username specified."
        echo "Usage: ./cli.sh entry_pgsql_delete_user <USERNAME>"
        exit 1
    fi

    # Check if the user exists
    existing_user=$(PSQL_CMD "SELECT 1 FROM pg_roles WHERE rolname = '$username';" -t -A)
    if [[ -z "$existing_user" ]]; then
        echo "User '$username' does not exist. Skipping deletion."
        return
    fi

    echo -e "Are you sure you want to delete the user '$username'?"
    read -rp "Please type the username to confirm: " confirm_username
    if [[ "$confirm_username" != "$username" ]]; then
        echo "User deletion cancelled."
        return
    fi
    PSQL_CMD "DROP USER $username;" && echo "User '$username' deleted successfully."
}

# Entry function to login to PostgreSQL
# Usage: ./cli.sh entry_pgsql_login <USERNAME> <PASSWORD>
# Description: Logs in to the PostgreSQL database using the provided credentials
entry_pgsql_login() {
    local username="$1"
    local password="$2"

    if [[ -z "$username" || -z "$password" ]]; then
        echo "Error: Username or password not specified."
        echo "Usage: ./cli.sh entry_pgsql_login <USERNAME> <PASSWORD>"
        exit 1
    fi

    PSQL_CLIENT -U "$username" -W "$password"
}

# Entry function to execute MongoDB client commands
# Usage: ./cli.sh entry_mongo_client <COMMAND>
# Description: Executes the provided MongoDB command using the MongoDB client
entry_mongo_client() {
    MONGO_CLIENT "${@}"
}

# Entry function to create a user in MongoDB
# Usage: ./cli.sh mongo_create_user <USERNAME> [<ROLE>] [<DATABASE>]
# Description: Creates a MongoDB user with the provided username, optional role, and optional database (default: admin)
entry_mongo_create_user() {
    local username=$1
    local role=$2
    local database=${3:-admin} # Default to 'admin' if not provided
    password=$(entry_gen_password)

    if [[ -z "$username" ]]; then
        echo "Error: Missing username."
        echo "Usage: ./cli.sh mongo_create_user <USERNAME> [<ROLE>] [<DATABASE>]"
        exit 1
    fi

    # Check if the user already exists in the admin database
    existing_user=$(MONGO_CMD "db.getSiblingDB('admin').getUser('$username')")

    if [[ "$existing_user" != "null" ]]; then
        echo "User '$username' already exists in the admin database. Skipping creation."
        return
    fi

    echo "Creating user: $username"
    echo "Generated password: $password"
    echo "Please save this password as it will not be shown again."

    # Create the user in the admin database
    if [[ -n "$role" && -n "$database" ]]; then
        MONGO_CMD "db.getSiblingDB('admin').createUser({
            user: '$username',
            pwd: '$password',
            roles: [{ role: '$role', db: '$database' }]
        })"
        echo "User '$username' created successfully in the admin database with role '$role' on database '$database'."
    else
        MONGO_CMD "db.getSiblingDB('admin').createUser({
            user: '$username',
            pwd: '$password',
            roles: []
        })"
        echo "User '$username' created successfully in the admin database without any roles."
    fi
}

# Entry function to create a database in MongoDB
# Usage: ./cli.sh mongo_create_db <DATABASE_NAME>
# Description: Creates a MongoDB database with the provided name
entry_mongo_create_db() {
    local database=$1

    if [[ -z "$database" ]]; then
        echo "Error: No database name specified."
        echo "Usage: ./cli.sh mongo_create_db <DATABASE_NAME>"
        exit 1
    fi

    # Check if the database already exists
    existing_db=$(MONGO_CMD "db.getMongo().getDBNames().indexOf('$database') >= 0")

    if [[ "$existing_db" == "true" ]]; then
        echo "Database '$database' already exists. Skipping creation."
        return
    fi

    # Create the database by inserting a document into a collection
    MONGO_CMD "db.getSiblingDB('$database').createCollection('init_collection')"
    MONGO_CMD "db.getSiblingDB('$database').init_collection.insertOne({ initialized: true })"
    echo "Database '$database' created successfully."
}

# Entry function to grant a role to a user in MongoDB
# Usage: ./cli.sh entry_mongo_grant <USERNAME> <ROLE> <DATABASE>
# Description: Grants the specified role on the specified database to the specified user (dbOwner)
entry_mongo_grant() {
    local username=$1
    local role=$2
    local database=$3

    if [[ -z "$username" || -z "$role" || -z "$database" ]]; then
        echo "Error: Missing parameters."
        echo "Usage: ./cli.sh entry_mongo_grant <USERNAME> <ROLE> <DATABASE>"
        exit 1
    fi

    # Check if the user exists in the admin database
    existing_user=$(MONGO_CMD "db.getSiblingDB('admin').getUser('$username')")

    if [[ "$existing_user" == "null" ]]; then
        echo "Error: User '$username' does not exist in the admin database."
        exit 1
    fi

    # Grant the role to the user on the specified database
    MONGO_CMD "db.getSiblingDB('admin').grantRolesToUser('$username', [{ role: '$role', db: '$database' }])"

    echo "Granted role '$role' on database '$database' to user '$username'."
}

# Entry function to delete a database in MongoDB
# Usage: ./cli.sh entry_mongo_delete_db <DATABASE_NAME>
# Description: Deletes the specified MongoDB database
entry_mongo_delete_db() {
    local database=$1

    if [[ -z "$database" ]]; then
        echo "Error: No database name specified."
        echo "Usage: ./cli.sh entry_mongo_delete_db <DATABASE_NAME>"
        exit 1
    fi

    # Check if the database exists
    existing_db=$(MONGO_CMD "db.getMongo().getDBNames().indexOf('$database') >= 0")

    if [[ "$existing_db" != "true" ]]; then
        echo "Database '$database' does not exist. Skipping deletion."
        return
    fi

    # Drop the database
    MONGO_CMD "db.getSiblingDB('$database').dropDatabase()"
    echo "Database '$database' deleted successfully."
}

# Entry function to delete a user in MongoDB
# Usage: ./cli.sh entry_mongo_delete_user <USERNAME>
# Description: Deletes the specified MongoDB user from the admin database
entry_mongo_delete_user() {
    local username=$1

    if [[ -z "$username" ]]; then
        echo "Error: No username specified."
        echo "Usage: ./cli.sh entry_mongo_delete_user <USERNAME>"
        exit 1
    fi

    # Check if the user exists in the admin database
    existing_user=$(MONGO_CMD "db.getSiblingDB('admin').getUser('$username')")

    if [[ "$existing_user" == "null" ]]; then
        echo "Error: User '$username' does not exist in the admin database."
        exit 1
    fi

    # Drop the user
    MONGO_CMD "db.getSiblingDB('admin').dropUser('$username')"
    echo "User '$username' deleted successfully from the admin database."
}

# Entry function to login to the MongoDB database
# Usage: ./cli.sh mongo_login <USERNAME> <PASSWORD>
# Description: Logs in to the MongoDB database using the provided credentials
entry_mongo_login() {
    local username="$1"
    local password="$2"

    if [[ -z "$username" || -z "$password" ]]; then
        echo "Error: Username or password not specified."
        echo "Usage: ./cli.sh mongo_login <USERNAME> <PASSWORD>"
        exit 1
    fi

    docker compose exec -T mongo mongosh --username "$username" --password "$password" --authenticationDatabase "admin"
}

# Entry function to login to the MySQL database
# Usage: ./cli.sh login <USERNAME> <PASSWORD>
# Description: Logs in to the specified MySQL database using the provided credentials
entry_login() {
    local username="$1"
    local password="$2"

    if [[ -z "$username" || -z "$password" ]]; then
        echo "Error: Username or password not specified."
        echo "Usage: ./cli.sh login <USERNAME> <PASSWORD>"
        exit 1
    fi

    TTY_CLIENT -u"$username" -p"$password"
}

# Entry function to create a database in MySQL
# Usage: ./cli.sh create_db <DATABASE_NAME>
# Description: Creates a MySQL database with the provided name
entry_create_db() {
    local database=$1

    if [[ -z "$database" ]]; then
        echo "Error: No database name specified."
        echo "Usage: ./cli.sh create_db <DATABASE_NAME>"
        exit 1
    fi

    # Check if the database already exists
    existing_db=$(CMD "SHOW DATABASES LIKE '$database';")
    if [[ -n "$existing_db" ]]; then
        echo "Database '$database' already exists. Skipping creation."
        return
    fi

    CMD "CREATE DATABASE IF NOT EXISTS $database;"
    echo "Database '$database' created successfully."
}

# Entry function to connect to a MySQL database as root
# Usage: ./cli.sh client [<Options>]
# Description: Connects to the specified MySQL database using the root user
entry_client() {
    TTY_CLIENT "${@}"
}

# Entry function to create a user in MySQL and grant privileges
# Usage: ./cli.sh create_user <USERNAME> <PASSWORD>
# Description: Creates a MySQL user with the provided credentials
entry_create_user() {
    local username=$1
    password=$(entry_gen_password)

    if [[ -z "$username" ]]; then
        echo "Error: No username specified."
        echo "Usage: ./cli.sh create_user <USERNAME>"
        exit 1
    fi

    # Check if the user already exists
    existing_user=$(CMD "SELECT user FROM mysql.user WHERE user = '$username';")
    if [[ -n "$existing_user" ]]; then
        echo "User '$username' already exists. Skipping creation."
        return
    fi

    echo "Creating user: $username"
    echo "Generated password: $password"
    echo "Please save this password as it will not be shown again."

    CMD "CREATE USER IF NOT EXISTS '$username'@'%' IDENTIFIED BY '$password';"
    echo "User '$username' created successfully."
}

# Entry function to delete a user in MySQL
# Usage: ./cli.sh delete_user <USERNAME>
# Description: Deletes a MySQL user with the provided name after confirmation
entry_delete_user() {
    local username=$1

    if [[ -z "$username" ]]; then
        echo "Error: No username specified."
        echo "Usage: ./cli.sh delete_user <USERNAME>"
        exit 1
    fi

    # Check if the user exists
    existing_user=$(CMD "SELECT user FROM mysql.user WHERE user = '$username';")
    if [[ -z "$existing_user" ]]; then
        echo "User '$username' does not exist. Skipping deletion."
        return
    fi

    read -rp "Are you sure you want to delete the user '$username'? Please type the username to confirm: " confirm_username
    if [[ "$confirm_username" != "$username" ]]; then
        echo "User deletion cancelled."
        return
    fi

    CMD "DROP USER '$username'@'%';"
    echo "User '$username' deleted successfully."
}

# Entry function to delete a database in MySQL
# Usage: ./cli.sh delete_db <DATABASE_NAME>
# Description: Deletes a MySQL database with the provided name after confirmation
entry_delete_db() {
    local database=$1

    if [[ -z "$database" ]]; then
        echo "Error: No database name specified."
        echo "Usage: ./cli.sh delete_db <DATABASE_NAME>"
        exit 1
    fi

    # Check if the database exists
    existing_db=$(CMD "SHOW DATABASES LIKE '$database';")
    if [[ -z "$existing_db" ]]; then
        echo "Database '$database' does not exist. Skipping deletion."
        return
    fi

    echo -e "Are you sure you want to delete the database '$database'?"
    read -rp "Please type the database name to confirm: " confirm_database
    if [[ "$confirm_database" != "$database" ]]; then
        echo "Database deletion cancelled."
        return
    fi

    CMD "DROP DATABASE $database;"
    echo "Database '$database' deleted successfully."
}

# Entry function to grant privileges to a user on a database in MySQL
# Usage: ./cli.sh grant <DATABASE_NAME> <USERNAME>
# Description: Grants all privileges on the specified database to the specified user
entry_grant() {
    local database=$1
    local username=$2

    if [[ -z "$database" ]]; then
        echo "Error: Database name not specified."
        echo "Usage: ./cli.sh grant <DATABASE_NAME> <USERNAME>"
        exit 1
    fi

    if [[ -z "$username" ]]; then
        echo "Error: Username not specified."
        echo "Usage: ./cli.sh grant <DATABASE_NAME> <USERNAME>"
        exit 1
    fi

    # Check if the database exists
    existing_db=$(CMD "SHOW DATABASES LIKE '$database';")
    if [[ -z "$existing_db" ]]; then
        echo "Database '$database' does not exist. Skipping privilege grant."
        return
    fi

    # Check if the user exists
    existing_user=$(CMD "SELECT user FROM mysql.user WHERE user = '$username';")
    if [[ -z "$existing_user" ]]; then
        echo "User '$username' does not exist. Skipping privilege grant."
        return
    fi

    echo "Granting privileges on database: $database to user: $username"
    CMD "GRANT ALL PRIVILEGES ON $database.* TO '$username'@'%'; FLUSH PRIVILEGES;"
    echo "Privileges granted. Displaying current privileges for user: $username on database: $database"
    CMD "SHOW GRANTS FOR '$username'@'%';"
}

# Main function example
# Usage: ./cli.sh <COMMAND>
# Description: Main function to execute the default behavior
function main() {
    _usage
}

source "${AWESOME_SHELL_ROOT}/core/usage.sh" && usage "${@}"

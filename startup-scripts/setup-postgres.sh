#!/bin/bash

# Initialize variables with default values
user=""
password=""
database=""

# Parse command-line options and arguments using getopts
while getopts "u:p:d:" opt; do
  case $opt in
    u)
      user="$OPTARG"
      ;;
    p)
      password="$OPTARG"
      ;;
    d)
      database="$OPTARG"
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      exit 1
      ;;
  esac
done

# Check if required options are provided
if [ -z "$user" ] || [ -z "$password" ] || [ -z "$database" ]; then
  echo "Usage: $0 -u <user> -p <password> -d <database>" >&2
  exit 1
fi


# Update the package list
sudo apt-get update
sudo apt-get upgrade -y
sudo apt-get clean

# Install Node to run integration tests
sudo apt install nodejs npm -y
sudo apt-get install locales -y

# Install PostgreSQL and its contrib package (includes additional utilities and extensions)
sudo apt-get install postgresql postgresql-contrib -y

sudo apt remove git -y
sudo apt purge git -y
sudo apt autoremove -y

# Start and enable the PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create a new PostgreSQL user and database (optional)
sudo -u postgres psql -c "ALTER USER $user WITH PASSWORD '$password';"
sudo -u postgres psql -c "CREATE DATABASE $database;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $database TO $user;"

echo "PostgreSQL installation completed."
echo "Check PostgreSQL Version"
psql --version
echo "Check Node Version"
node -v

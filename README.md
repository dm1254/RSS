# RSS
## Installations


* Install Postgres v15 or later

macOS with brew
```vim
brew install postgresql@15
```

Linux
```vim
sudo apt update
sudo apt install postgresql postgresql-contrib
```
* Installing Go

This project requires Go version 1.23 or higher. Follow the instructions below to install Go on your system.

### Option 1: Using Webi (Recommended for Unix-like systems)

[Webi](https://webinstall.dev/golang/) provides a simple command to install Go:

```bash
curl -sS https://webinstall.dev/golang | bash
```
### Option 2: Use package managers

macOS with brew
```vim
brew install go
```
Linux
```vim
sudo apt update
sudo apt install golang-go
```
### Install the gator CLI tool
* Use go install
```vim
go install https://github.com/dm1254/RSS
```
## Build Instructions

### Create a config file
Manually create a .gatorconfig.json file in your home directory 
```vim
touch ~/.gatorconfig.json
```
Add the following content to file
```vim
{
  "db_url": "postgres://your-database-url-goes-here"
}
```
* Replace "db_url": "postgres://your-database-url-goes-here" with the connection string to your PostgreSQL database

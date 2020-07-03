# Oauth2 Authorization Code built with Golang, postgresql, redis 
## with Google Login/Sign Up
## with help from https://godoc.org/gopkg.in/oauth2.v3

# Setting up with Docker and docker-compose
## Prerequisites

  * Download and install [Docker Desktop](https://www.docker.com/get-started) if you don't have it installed.
  * Download and install [Docker Compose](https://docs.docker.com/compose/install) if you don't have it installed.

## Set up Environment Variables
  * Duplicate the file `.env.template` and rename it as `.env`
  * Complete the information needed in the `.env` file. eg. `DB_NAME=hostname` etc. 
  * Please make sure all information is provided.
  * Request for API Keys when needed. 

## Building and Running docker container  
  * Open terminal for Linux/macOS. For Windows open command prompt or Powershell
  * change directory to the project folder. i/e this project folder
  * run the command `docker-compose -f docker-compose.yml build` to build an images using the docker-compose YAML file in the project. This will take a while to finish building. 
  * run the command  `docker-compose -f docker-compose.yml up`. The containers should start running. You can monitor logs whiles it running in your terminal
  * Go to your browser  and visit `127.0.0.1:8080` to access the client app. 
  * run CTRL + C to terminate the containers running. and `docker-compose -f docker-compose.yml down` to remove the images. 

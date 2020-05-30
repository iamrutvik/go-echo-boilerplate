# README #

This README would normally document whatever steps are necessary to get your application up and running.

### What is this repository for? ###

* Quick summary
* Version
* [Learn Markdown](https://bitbucket.org/tutorials/markdowndemo)

### How do I get set up? ###

* Summary of set up
* Configuration
* Dependencies
* Database configuration
* How to run tests
* Deployment instructions

### Contribution guidelines ###

* Writing tests
* Code review
* Other guidelines

### Who do I talk to? ###

* Repo owner or admin
* Other community or team contact

##### Check where  go lint installed - https://github.com/golang/lint
- To find out where golint was installed you can run go list -f {{.Target}} golang.org/x/lint/golint. For golint to be used globally add that directory to the $PATH environment setting.
- Set up the go lint for the project https://github.com/vmware/dispatch/wiki/Configure-GoLand-with-golint
- add things about modd
    - Install modd via `brew install modd`
    - go to the project directory and run `modd`
- add command to run via docker - docket-compose up --build
- add command to mongo db connection string - mongodb://localhost:27017
- add command to [remove from the docker](https://linuxize.com/post/how-to-remove-docker-images-containers-volumes-and-networks/) -  docker system prune --volumes
- command to remove git caches - git rm -rf --cached . 
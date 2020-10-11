# Go Boilerplate #

Dockerized Go Boilerplate with Integrated Prisma. Start your next react project in seconds with the best DX and a focus on performance and best practices.

### Features ###

* ##### Echo Framework #####
  High performance, extensible, minimalist Go web framework.
* ##### Docker Ready #####
  Uses `docker-compose.yml` to make it docker ready. You can utilise this image anywhere including Kubernetes.
* ##### [Prisma](https://www.prisma.io/) #####
  Modern Database access now even for Go, use it to query your database with GraphQL or use it as ORM.
* ##### Swagger Ready #####
  Converts Go annotations to Swagger Documentation 2.0, with Swagger UI. It uses popular package [swag](https://github.com/swaggo/swag).
* ##### API Versioning #####
  Manage your API's version effectively.
* ##### TLS #####
  Automatically install TLS certificates from Let's Encrypt, which you can configure easily.
* ##### HTTP/2 #####
  Inbuilt HTTP/2 support improves speed and provides better user experience.
* ##### Live reload #####
  Install [modd](https://github.com/cortesi/modd) and your app will be auto reload on every save.
* ##### Linter #####
  Pre-configured Go linter with Best Practices. 

### Internal Packages ###

* [Viper](https://github.com/spf13/viper) for configuration
* Used Go Modules for Dependency management.
* [Supervisord](https://github.com/ochinchina/supervisord) for process management, Docker uses it's config file to run compiled executable. 
* [modd](https://github.com/cortesi/modd) for Live Reload
* [Prisma](https://www.prisma.io/) as Database Access Layer
* [JWT Go](https://github.com/dgrijalva/jwt-go) for JWT Authentication
* [Swag](https://github.com/swaggo/swag) for Swagger Documentation
* [Validator](https://github.com/go-playground/validator) for Request validation
 

### Prerequisites ###

* Go lang installed and $GOPATH must be set
* [modd](https://github.com/cortesi/modd) globally installed
* Docker installed if required

### How to start ###

* Clone this repository
* `modd` to start the auto reloaded API server
* `docker-compose up --build` to run auto reload docker service

### Who do I talk to? ###

Send an email to hi@rutvikbhatt.com or open an issue.

###### Few useful command ###### 
- Check where  go lint installed - https://github.com/golang/lint
- To find out where golint was installed you can run go list -f {{.Target}} golang.org/x/lint/golint. For golint to be used globally add that directory to the $PATH environment setting.
- Set up the go lint for the project https://github.com/vmware/dispatch/wiki/Configure-GoLand-with-golint
- things about modd
    - Install modd via `brew install modd`
    - go to the project directory and run `modd`
- add command to run via docker - docket-compose up --build
- add command to mongo db connection string - mongodb://localhost:27017
- add command to [remove from the docker](https://linuxize.com/post/how-to-remove-docker-images-containers-volumes-and-networks/) -  docker system prune --volumes
- command to remove git caches - git rm -rf --cached . 

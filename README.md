# README #

This README would normally document whatever steps are necessary to get your application up and running.

### What is this repository for? ###

* Quick summary
* Version
* [Learn Markdown](https://bitbucket.org/tutorials/markdowndemo)

### How to build and run the app in production? ###
* go to service directory e.g. cd ...\email_action
* docker build -t email_action  -f Dockerfile.production .
* docker run -it -p 8080:8080 email_action

### How to build and run the app in docker locally? (currently setup pipeline for AWS ECR) ###
* go to service directory e.g. cd ...\email_action
* docker build -t email_action  -f Dockerfile.dev .
* docker run -it -p 8080:8080 email_action

### How to build and run the app locally(with auto-reload feature)? ###
* go to service directory e.g. cd ...\email_action
* go mod vendor
* cd front_end_app
* npm install
* npm run serve
* bee run

### How to run frontend page locally ?###
* cd front_end_app
* npm install
* npm run serve

### How to run frontend app and nginx in docker ?###
* cd front_end_app
* docker build -t email_action_app  -f Dockerfile .
* docker run -it -p 80:80 email_action_app

### How to run frontend page in production ?###
* npm run build

### AWS ECR sourced deployment to ECS ###
* code build spec: https://stackoverflow.com/questions/55339872/codepipeline-ecr-source-ecs-deploy-configuration
* openresty with lua to support jwt authentication: https://medium.com/@tumulr/building-an-api-gateway-with-nginx-lua-e3dff45e6e63

### Stripe ##
stripe listen --forward-to http://127.0.0.1:8080/api/v1/webhook
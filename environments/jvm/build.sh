#!/usr/bin/env bash

docker run -it --rm --name my-maven-project -v "$PWD":/usr/src/mymaven -w /usr/src/mymaven maven:3.2-jdk-8 mvn clean install

cp target/*-jar-with-dependencies.jar env.jar

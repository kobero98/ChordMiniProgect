#! /bin/bash
while getopts x: flag
do
    case "${flag}" in
        x) x=${OPTARG};;
    esac
done
rm Docker-compose.yaml
python3 dockerComposeCreator.py $x
docker-compose build
docker-compose up

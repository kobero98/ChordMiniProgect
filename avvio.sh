#! /bin/bash

C=0
while getopts x: flag
do
    case "${flag}" in
        x) x=${OPTARG}
           C=1 
            ;;
    esac
done
if [[ "$C" == "1" ]] 
then
    rm Docker-compose.yaml;
    python3 dockerComposeCreator.py $x;
    docker-compose build;
    docker-compose up;
else
    echo "bash avvio.sh -x [numero_nodi]"
fi

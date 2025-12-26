#!/usr/bin/env bash

FILE="pokemon"
SPECIES_DIR="/home/david/projects/api-data/data/api/v2/pokemon-species/"

if [ -f $FILE ]; then
    rm $FILE
fi

touch $FILE
for dir in /home/david/projects/api-data/data/api/v2/pokemon/*/; do
    NAME="$(jq -r '.name' "${dir}/index.json")"
    FOLDER="$(basename $dir)"
    SPECIES_FILE="${SPECIES_DIR}${FOLDER}/index.json"

    if [ -f $SPECIES_FILE ]; then
        G="$(jq -r '.genera[] | select(.language.name == "en").genus' "${SPECIES_FILE}" | sed 'y/áéí/aei/')"
        G="$(echo $G | sed s/"Pokemon"//)"
        #G="${G// /_}"
        NAME="${NAME},${G}"
    else
        NAME="${NAME},-"
    fi



    echo "${NAME}" >> "${FILE}"
done

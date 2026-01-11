#!/usr/bin/env bash

function join_by {
    local d=${1-} f=${2-}
    if shift 2; then
        printf %s "$f" "${@/#/$d}"
    fi
}

FILE="data/pokemon"
SPECIES_DIR="/home/david/projects/api-data/data/api/v2/pokemon-species/"

if [ -f $FILE ]; then
    rm $FILE
fi

touch $FILE
for dir in /home/david/projects/api-data/data/api/v2/pokemon/*/; do
    DELIMITER=","
    NAME="$(jq -r '.name' "${dir}/index.json")"
    TYPES="$(jq -r '.types[] | .type.name' "${dir}/index.json")"
    TYPES="$(join_by $DELIMITER $TYPES)"

    if ! [[ $TYPES =~ $DELIMITER ]]; then
        TYPES="${TYPES},-"
    fi

    FOLDER="$(basename $dir)"
    SPECIES_FILE="${SPECIES_DIR}${FOLDER}/index.json"

    if [ -f $SPECIES_FILE ]; then
        G="$(jq -r '.genera[] | select(.language.name == "en").genus' "${SPECIES_FILE}" | sed 'y/áéí/aei/')"
        G="$(echo $G | sed s/"Pokemon"//)"
        NAME="${NAME},${G},${TYPES}"
    else
        NAME="${NAME},-,${TYPES}"
    fi

    echo "${NAME}" >> "${FILE}"
done

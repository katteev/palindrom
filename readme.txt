-- hent alle brukere
curl -L -X GET "http://localhost:8080/users"

-- oppdater bruker2
curl -L -X PUT "http://localhost:8080/users/2" -H "Content-Type: application/json" --data-raw "{\"fornavn\": \"oppdatertnavn\",\"etternavn\": \"oppdatert\"}"

-- slett bruker 4
curl -L -X DELETE "http://localhost:8080/users/4"

-- sjekk om bruker 1s fornavn/etternavn er palindrom
curl -L -X GET "http://localhost:8080/users/palindrom/1"

-- hent bruker med id 1
curl -L -X GET "http://localhost:8080/users/1" -H "Content-Type: application/json"

-- lag ny bruker
curl -L -X POST "http://localhost:8080/users" -H "Content-Type: application/json" --data-raw "{\"fornavn\": \"Hannah\",\"etternavn\": \"Anna2\"}"
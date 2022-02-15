#!/bin/bash
url=url
while IFS= read -r thread
do
printf "\n$thread\n"
curl -X POST http://localhost:5050/parse -d "{\"$url\":\"$thread\"}"
done < "threads"
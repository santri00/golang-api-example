##!/bin/bash
#latestHash=$(git log -1 --pretty=format:"%h")
#echo $latestHash
#imgTag=v1.0.0.$latestHash
#docker build -t stripathy/dealer-inventory-api:$imgTag .
#docker run -p 3001:3001 --name dealer-inventory-api -it stripathy/dealer-inventory-api:$imgTag
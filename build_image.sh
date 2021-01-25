#!/bin/bash

imgTag=v1.0.0
docker build -t stripathy/demo-vehicle-api:$imgTag .
docker run -p 3001:3001 --name demo-vehicle-api -it stripathy/demo-vehicle-api:$imgTag
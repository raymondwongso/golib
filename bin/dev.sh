#!/bin/bash

set -euo pipefail

cd dev
docker build -t golib-dev:latest .
cd ..

docker run -it --name golib-dev-local --rm -v $(pwd):/go/src/app golib-dev:latest /bin/bash

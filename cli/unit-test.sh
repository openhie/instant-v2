#!/bin/bash

cd src/core || exit 

go test .

cd ../.. || exit


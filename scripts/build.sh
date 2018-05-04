#!/usr/bin/env bash

dep ensure

echo "Compiling .go files to bin/handlers/ ..."

rm -rf bin/

cd src/handlers/

for folder in */;
  do
    (cd $folder
      for f in *.go;
      do  
        if [ $f == *"_test.go" ] ; then
          echo "-- "$f "Skipped"
          continue;
        fi

        filename="${f%.go}"
    
        if GOOS=linux go build -o "../../../bin/handlers/$filename" ${f}; then
          echo "✓ $filename Compiled"
        else
          echo "✕ Failed to compile $filename!"
          exit 1
        fi
    done)
  done

echo "Done."
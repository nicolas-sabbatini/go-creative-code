#!/usr/bin/env bash

set -e

function generate_html() {
  HTML=$(
    cat <<-END
<!DOCTYPE html>
<!--
Copyright 2018 The Go Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
-->
<html>
  <head>
    <meta charset="utf-8" />
    <title>$1</title>
    <style>
      html {
        width: 100%;
        height: 100%;
        margin: 0px;
        padding: 0px;
      }
      body {
        display: flex;
        justify-content: center;
        align-items: center;
        min-width: 100%;
        min-height: 100%;
        background-color: rgb(0, 0, 0);
        height: 100%;
        margin: 0px;
        padding: 0px;
      }
    </style>
  </head>

  <body>
    <!--
	Add the following polyfill for Microsoft Edge 17/18 support:
	<script src="https://cdn.jsdelivr.net/npm/text-encoding@0.7.0/lib/encoding.min.js"></script>
	(see https://caniuse.com/#feat=textencoder)
	-->
    <script src="wasm_exec.js"></script>
    <script>
      if (!WebAssembly.instantiateStreaming) {
        // polyfill
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
          const source = await (await resp).arrayBuffer();
          return await WebAssembly.instantiate(source, importObject);
        };
      }

      const go = new Go();
      let mod, inst;
      WebAssembly.instantiateStreaming(
        fetch("$1.wasm"),
        go.importObject
      )
        .then((result) => {
          mod = result.module;
          inst = result.instance;
          document.getElementById("runButton").disabled = false;
        })
        .catch((err) => {
          console.error(err);
        });

      async function run() {
        document.getElementById("runButton").remove();
        console.clear();
        await go.run(inst);
        inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
      }
    </script>

    <button onClick="run();" id="runButton" disabled>Run</button>
  </body>
</html>
END
  )
  echo "$HTML"
}

# Run tidy
go mod tidy
# Export build variables
export GOARCH=wasm
export GOOS=js

# Bild all
for project in $(dir cmd); do
  echo "Bilding $project"
  directory="target/$project"
  binary="cmd/$project/$project"
  mkdir -p "$directory"
  cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" "$directory"
  html="$(generate_html "$project")"
  echo "$html" >"$directory/index.html"
  go build -o "$directory/$project.wasm" "$binary.go"
done

echo "Done! :D"

import * as childx from "child_process"; // Node.js child_process module
import fs from "fs"; // Node.js file system module

import resolve from "rollup-plugin-node-resolve";
import commonjs from "rollup-plugin-commonjs";
import {
  terser
} from "rollup-plugin-terser";
import builtins from "builtin-modules";
import globals from "rollup-plugin-node-globals";
import builtins2 from "rollup-plugin-node-builtins";
import json from "rollup-plugin-json";
import replace from "rollup-plugin-replace";

let minify = false;
let productionEnv = false;

run("./src/renderer", "./renderer", minify);
run("./src/main", "./main", minify);

export default [
  // main process (nodejs) config
  {
    input: "./main_go.js",
    output: {
      file: "./main.js",
      format: "cjs",
      sourcemap: false
    },
    external: ["electron", ...builtins],
    plugins: [
      resolve(),
      commonjs({}),
      globals(),
      builtins2(),
      minify ?
      terser({
        ecma: 6
      }) :
      {
        name: "x",
        generateBundle: function () {}
      },
      deleteFile("./main_go.js"),
      deleteFile("./main_go.js.map")
    ]
  },
  // renderer process (browser) config
  {
    input: "./renderer_go.js",
    output: {
      file: "./renderer.js",
      format: "iife",
      sourcemap: false,
      name: "main"
    },
    external: ["electron"],
    plugins: [
      resolve({
        browser: true,
        preferBuiltins: true,
        dedupe: ["react", "react-dom"]
      }),
      replace({
        "process.env.NODE_ENV": JSON.stringify(
          productionEnv ? "production" : "development"
        )
      }),
      json(),
      commonjs({}),
      globals(),
      builtins2(),
      minify ?
      terser({
        ecma: 6
      }) :
      {
        name: "x",
        generateBundle: function () {}
      },
      deleteFile("./renderer_go.js"),
      deleteFile("./renderer_go.js.map")
    ]
  }
];

function run(input, output, minify) {
  let path = "gopherjs build " + input;
  if (minify) {
    path = path + " -m ";
  }

  output = output + "_go";

  path = path + " -o " + output + ".js";

  childx.execSync(path);
}

function deleteFile(filePath) {
  return {
    name: "deleteFile",
    writeBundle: function () {
      fs.unlinkSync(filePath);
    }
  };
}


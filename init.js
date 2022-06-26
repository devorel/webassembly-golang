async function init(WASM_URL) {

    const crypto = require("crypto");
    const fs = require('fs');
    globalThis.require = require;
    globalThis.fs = require("fs");
    globalThis.TextEncoder = require("util").TextEncoder;
    globalThis.TextDecoder = require("util").TextDecoder;

    globalThis.crypto = {
        getRandomValues(b) {
            crypto.randomFillSync(b);
        },
    };
    require("./wasm_exec");
    const go = new Go();

    go.env = Object.assign({TMPDIR: require("os").tmpdir()}, process.env);
    go.exit = process.exit;
    const result = await WebAssembly.instantiate(fs.readFileSync(WASM_URL), go.importObject);

    process.on("exit", (code) => { // Node.js exits if no event handler is pending
        if (code === 0 && !go.exited) {
            // deadlock, make Go print error and stack traces
            go._pendingEvent = {id: 0};
            go._resume();
        }
    });
    // wasm = result.instance;
    // await go.run(result.instance);
    go.run(result.instance);

}


async function run() {
    await init(`./ecdsa.wasm`);
    const a = await global.test("sss");
    console.log('-------a-----');
    console.log(a);

    const b = await global.test("sss");
    console.log('-----b-------');
    console.log(b);
}

run();
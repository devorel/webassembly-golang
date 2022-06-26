async function init(WASM_URL) {
    if (!WebAssembly.instantiateStreaming) { // polyfill
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }

    const go = new Go();
    let mod, inst;
    WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject)
        .then(async (result) => {
            // mod = result.module;
            // inst = result.instance;
            await go.run(result.instance);
            // console.log(inst)
        })
        .catch((err) => {
            console.error(err);
        });

}
(async _ => await init(`./ecdsa.wasm`))();
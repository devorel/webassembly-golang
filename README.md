#how to work with webassembly golang

##nodejs
run
```bash
node init.js
```
```js
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
```

##web browser
run server
```bash
npx serve
```
open /index.html
```html
<script src="./wasm_exec.js"></script>
<script src="./initWeb.js"></script>
<script>
    async function run() {
        const a = await window.test("sss");
        console.log('-------a-----');
        console.log(a);

        const b = await window.test("sss");
        console.log('-----b-------');
        console.log(b);
    }

    setTimeout(run, 1000);
</script>
```

##compile the code 
GOARCH=amd64 GOOS=linux go build -o ecdsa.o ./ecdsa.go

##compile the code for webassembly
GOARCH=wasm GOOS=js go build -o ecdsa.wasm ./jsecdsa.go
WASM_HEADLESS=off GOARCH=wasm GOOS=js go build -o ecdsa.wasm ./jsecdsa.go

####environment variables
export GOROOT=/usr/local/go
export PATH=$PATH:/usr/local/go/bin
export GOPATH=${PWD}
unset $GOPATH
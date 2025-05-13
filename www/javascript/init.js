window.addEventListener("load", function load(event){

    sfomuseum.golang.wasm.fetch("wasm/world_clock_time.wasm").then((rsp) => {
	console.log("OK");
    }).catch((err) => {
	console.error("Failed to load update WASM binary", err);
        return;
    });

    
});


window.addEventListener("load", function load(event){

    var submit_el = document.getElementById("submit");
    
    var timezones;
    
    sfomuseum.golang.wasm.fetch("wasm/world_clock_time.wasm").then((rsp) => {

	world_clock_timezones().then(rsp => {
	    
	    timezones = JSON.parse(rsp);
	    console.log(timezones);

	    submit_el.onclick = function(){
		console.log("CLICK");
		return false;
	    };

	    submit_el.removeAttribute("disabled");
	    
	}).catch((err) => {
	    console.error("Failed to retrieve timezones", err);
	})

	
    }).catch((err) => {
	console.error("Failed to load update WASM binary", err);
        return;
    });

    
});


window.addEventListener("load", function load(event){

    var timezones;

    var feedback_el = document.getElementById("feedback");
    
    var populate_timezones = function(el){

	if (! timezones){
	    return false;
	}

	var count = timezones.length;

	for (var i=0; i < count; i++){

	    var tz = timezones[i];

	    var opt = document.createElement("option");
	    opt.setAttribute("value", tz.name);
	    opt.setAttribute("data-whosonfirst-id", tz["wof:id"]);
	    opt.appendChild(document.createTextNode(tz.label));

	    el.appendChild(opt);
	}
	    
    };

    var render_results = function(results) {

	var results_el = document.getElementById("results");
	results_el.innerHTML = "";
	
	var count = results.length;

	if (! count){
	    return;
	}
	
	var table = document.createElement("table");
	table.setAttribute("class", "table table-hover");

	var thead = document.createElement("thead");
	var tbody = document.createElement("tbody");
	
	var tr = document.createElement("tr");
	
	var tz_header = document.createElement("th");
	tz_header.appendChild(document.createTextNode("Location"));
	tr.appendChild(tz_header);
	
	var date_header = document.createElement("th");
	date_header.appendChild(document.createTextNode("Date"));
	tr.appendChild(date_header);

	var dow_header = document.createElement("th");
	dow_header.appendChild(document.createTextNode("Day"));
	tr.appendChild(dow_header);
	
	var time_header = document.createElement("th");
	time_header.appendChild(document.createTextNode("Time"));
	tr.appendChild(time_header);
	
	thead.appendChild(tr);
	table.appendChild(thead);

	for (var i=0; i < count; i++){

	    var tr = document.createElement("tr");
	    
	    var tz_column = document.createElement("td");
	    tz_column.appendChild(document.createTextNode(results[i].label));
	    tr.appendChild(tz_column);
	    	
	    var date_column = document.createElement("td");
	    date_column.appendChild(document.createTextNode(results[i].date));
	    tr.appendChild(date_column);

	    var dow_column = document.createElement("td");
	    dow_column.appendChild(document.createTextNode(results[i].day_of_week));
	    tr.appendChild(dow_column);
	    
	    var time_column = document.createElement("td");
	    time_column.appendChild(document.createTextNode(results[i].time));
	    tr.appendChild(time_column);
	    
	    tbody.appendChild(tr);
	}

	table.appendChild(tbody);

	var wrapper = document.createElement("div");
	wrapper.setAttribute("class", "table-responsive");
	wrapper.appendChild(table);
	
	results_el.appendChild(wrapper);
    };
    
    var derive_times_input = function(){
	
	const date_el = document.getElementById("date");
	const time_el = document.getElementById("time");	
	const timezone_el = document.getElementById("timezone");
	const other_els = document.getElementsByClassName("other-timezone");

	const date = date_el.value;
	const time = time_el.value;
	const tz = timezone_el.value;
	
	var others = [];
	
	var count_others = other_els.length;
	
	for (var i=0; i < count_others; i++){
	    const tz = other_els[i].value;

	    if (tz != ""){
		others.push(tz);
	    }
	}

	if (!date){
	    return { error: "Missing date" };
	}

	if (! time){
	    return { error: "Missing time" };
	    feedback_el.innerText = "Missing time";
	}

	if (! tz){
	    return { error: "Missing timezone" };
	}
	
	if (others.length == 0){
	    return { error: "No other locations to lookup" }; 
	}

	const dt = date + " " + time;	
	const str_others = others.join(",");
	
	return {
	    datetime: dt,
	    timezone: tz,
	    locations: str_others,
	};
    };

    var derive_times = function(){

	feedback_el.innerHTML = "";
	
	const data = derive_times_input();

	if (data.error) {
	    feedback_el.innerText = data.error;
	    console.error(data.error);
	    return false;
	}
	
	world_clock_time(data.datetime, data.timezone, data.locations).then((rsp) => {
	    const results = JSON.parse(rsp);
	    render_results(results);
	}).catch((err) => {
	    feedback_el.innerText = "Failed to derive times: " + err;
	    console.error("Failed to derive times", err)
	});
	
    };

    var add_other_timezone = function(){

	var others_el = document.getElementById("other-timezones");
	
	const n = others_el.children.length + 1;
	const id = "others-" + n;

	console.log("new id", id);
	
	var wrapper_el = document.createElement("div");
	wrapper_el.setAttribute("id", id);
	
	var select_el = document.createElement("select");
	select_el.setAttribute("class", "form-select other-timezone");
	
	var opt_el = document.createElement("option");
	opt_el.setAttribute("value", "");
	
	select_el.appendChild(opt_el);

	var remove_im = document.createElement("img");
	remove_im.setAttribute("src", "images/remove.svg");
	remove_im.setAttribute("height", 20);
	remove_im.setAttribute("width", 20);	
	remove_im.setAttribute("data-wrapper-id", id);
	
	var remove_el = document.createElement("div");
	remove_el.setAttribute("class", "remove-other");
	remove_el.setAttribute("data-wrapper-id", id);
	remove_el.appendChild(remove_im);

	remove_el.onclick = function(e){

	    var el = e.target;
	    var wrapper_id = el.getAttribute("data-wrapper-id");
	    var wrapper_el = document.getElementById(wrapper_id);

	    if (wrapper_id){
		wrapper_el.remove();
	    }
	    
	    return false;
	};

	wrapper_el.appendChild(remove_el);
	wrapper_el.appendChild(select_el);
	
	others_el.appendChild(wrapper_el);
	
	populate_timezones(select_el);
	return false;
	
    };

    var init = function(){

	var timezone_el = document.getElementById("timezone");    	    
	populate_timezones(timezone_el);
	
	var others = document.getElementsByClassName("other-timezone");
	var others_count = others.length;
	
	for (var i=0; i < others_count; i++){
	    populate_timezones(others[i]);
	}
	
	var add_el = document.getElementById("add-other");
	add_el.style.display = "inline";
	
	add_el.onclick = function(){
	    
	    try {
		add_other_timezone();
	    } catch(err){
		feedback_el.innerText = "Failed to add new location: " + err;		    
		console.error("SAD", err)
	    };
	    
	    return false;
	};
	
	var submit_el = document.getElementById("submit");
	
	submit_el.onclick = function(){
	    
	    try {
		derive_times();
	    } catch (err) {
		console.error("Failed to derive times", err);
	    }
	    
	    return false;
	};
	
	submit_el.removeAttribute("disabled");
	feedback_el.innerHTML = "";
    };

    var setup_offline = function(){

	console.debug("Setup offline support");
	
	const scope = location.pathname;
	
	worldclock.offline.init(scope).then((rsp) => {

	    console.debug("Offline service workers registered for scope " + scope);

	    const purge_func = function(){
		
		worldclock.offline.purge_with_confirmation().then((rsp) => {
		    feedback_el.innerText = "Offline cache has been removed.";
		}).catch((err) => {
		    feedback_el.innerText = "Failed to purge offline cache, " + err;
		});
		
		return false;
	    };
	    
	    var purge_im = document.createElement("img");
	    purge_im.setAttribute("src", "images/purge.svg");
	    purge_im.setAttribute("height", 16);
	    purge_im.setAttribute("width", 16);
	    purge_im.onclick = purge_func;

	    /*
	    var purge_el = document.createElement("span");
	    purge_el.setAttribute("id", "purge");
	    purge_el.appendChild(purge_im);
	    */
	    
	    var header = document.getElementById("header");	
	    header.appendChild(purge_im);
	    
	}).catch((err) => {
	    feedback_el.innerText = "Failed to initialize offline mode, " + err;
	});
	
    };
    
    // Okay, go!
    
    sfomuseum.golang.wasm.fetch("wasm/world_clock_time.wasm").then((rsp) => {

	world_clock_timezones().then(rsp => {

	    try {
		timezones = JSON.parse(rsp);
	    } catch(err) {
		feedback_el.innerText = "Failed to derive timezones list: " + err;
		console.error(err);
		return false
	    }

	    const offline = document.body.hasAttribute("offline");
	    
	    if (offline){

		try {
		    setup_offline();
		} catch (err) {
		    feedback_el.innerText = "Failed to setup offline support: " + err;
		    console.error(err);
		    return false;
		}
	    }

	    try {
		init();
	    } catch(err) {
		feedback_el.innerText = "Failed to initialize application: " + err;
		console.error(err);
		return false;
	    }
	    
	}).catch((err) => {
	    feedback_el.innerText = "Failed to derive timezones list: " + err;	    
	    console.error("Failed to retrieve timezones", err);
	})

	
    }).catch((err) => {
	feedback_el.innerText = "Failed to load application: " + err;	
	console.error("Failed to load update WASM binary", err);
        return;
    });

    
});


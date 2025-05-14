window.addEventListener("load", function load(event){

    var timezones;
    
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
	results_el.appendChild(table);
    };
    
    var derive_times = function(){

	var date_el = document.getElementById("date");
	var timezone_el = document.getElementById("timezone");
	var other_els = document.getElementsByClassName("other-timezone");

	var date = date_el.value;
	var tz = timezone_el.value;
	var others = [];
	
	var count_others = other_els.length;
	
	for (var i=0; i < count_others; i++){
	    others.push(other_els[i].value);
	}
	
	var str_others = others.join(",");

	world_clock_time(date, tz, str_others).then((rsp) => {
	    const results = JSON.parse(rsp);
	    render_results(results);
	}).catch((err) => {
	    console.error("SAD", err)
	});
	
    };

    var add_other_timezone = function(){

	var select_el = document.createElement("select");
	select_el.setAttribute("class", "form-select other-timezone");
	
	var opt_el = document.createElement("option");
	opt_el.setAttribute("value", "");
	
	select_el.appendChild(opt_el);

	var others_el = document.getElementById("other-timezones");   	
	others_el.appendChild(select_el);
	
	populate_timezones(select_el);
	return false;
	
    };
    
    sfomuseum.golang.wasm.fetch("wasm/world_clock_time.wasm").then((rsp) => {

	world_clock_timezones().then(rsp => {

	    try {
		timezones = JSON.parse(rsp);
	    } catch(err) {
		console.error(err);
		return false
	    }

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
		add_other_timezone();
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
	    
	}).catch((err) => {
	    console.error("Failed to retrieve timezones", err);
	})

	
    }).catch((err) => {
	console.error("Failed to load update WASM binary", err);
        return;
    });

    
});


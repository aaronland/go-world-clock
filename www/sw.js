const cache_name = 'world-clock-v0.0.5';

const app_files = [
    // HTML
    "./index.html",

    // Images
    "./images/add.svg",
    "./images/remove.svg",
    "./images/purge.svg",    
    
    // CSS
    "./css/bootstrap.min.css",
    "./css/worldclock.css",    
    
    // Javascript dependencies
    "./javascript/sfomuseum.golang.wasm.bundle.js",

    // Javascript application
    "./javascript/worldclock.offline.js",        
    "./javascript/worldclock.init.js",    

    // WASM
    "./wasm/world_clock_time.wasm",
    
    // Javascript service workers
    "./sw.js"    
];

self.addEventListener("install", (e) => {

    console.log("SW install event", cache_name);

    e.waitUntil((async () => {
	const cache = await caches.open(cache_name);
	console.log("SW cache files", cache_name, app_files);
	await cache.addAll(app_files);
	console.log("SW cache files added", cache_name);
    })());
});

addEventListener("activate", (event) => {
    console.log("SW activate", cache_name);
});

addEventListener("message", (event) => {
    // event is a MessageEvent object
    console.log(`The service worker sent me a message: ${event.data}`);
  });


self.addEventListener('fetch', (e) => {

    // https://developer.mozilla.org/en-US/docs/Web/API/Cache
    
    e.respondWith((async () => {

	console.debug("fetch", cache_name, e.request.url);
	
	const cache = await caches.open(cache_name);
	const r = await cache.match(e.request);
	
	console.debug(`[Service Worker] Fetching resource: ${e.request.url}`);
	
	if (r) {
	    console.debug("return cache", e.request.url);
	    return r;
	}
	
	const response = await fetch(e.request);
	
	console.debug(`[Service Worker] Caching new resource: ${e.request.url}`);
	cache.put(e.request, response.clone());
	
	return response;
    })());
    
});

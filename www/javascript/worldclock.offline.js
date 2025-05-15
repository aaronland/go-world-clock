/* https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps/Offline_Service_workers#service_workers_explained */
/* https://github.com/mdn/pwa-examples/blob/main/js13kpwa/sw.js */
/* https://web.dev/learn/pwa/service-workers/ */
/* https://webkit.org/blog/8090/workers-at-your-service/ */

var worldclock = worldclock || {};

worldclock.offline = (function(){

    var self = {
	
	init: function(scope){

	    return new Promise((resolve, reject) => {
		
		if (! "serviceWorker" in navigator) {
		    reject("Service workers not available.");
		    return
		}
		
		var sw_uri = scope + "sw.js";
		
		var sw_args = {
		    scope: scope,
		};

		console.log("register service worker", sw_uri, sw_args);
		
		navigator.serviceWorker.register(sw_uri, sw_args).then((registration) => {

		    if (navigator.onLine){

			console.debug("update sw registration");
			
			registration.update().then((rsp) => {
			    console.log("sw registration updated");
			}).catch((err) => {
			    console.error("failed to update sw registration", err);
			});
		    }
		    
		    console.log("sw registered", sw_args);		    
		    resolve();
		    
		}).catch((err) => {
		    console.error("Failed to register service worker", err);
		    reject(err);
		});
		
	    });
	    
	},

	purge_with_confirmation: function(){

	    return new Promise((resolve, reject) => {
		
		if (! confirm("Are you sure you want to delete all the application caches? This can not be undone.")){
		    resolve();
		    return;
		}

		if (! navigator.onLine){
		    
		    if (! confirm("Are you really sure? You appear to be offline and deleting the application cache will probably cause offline support to stop working until you are online again.")){
			resolve()
			return;
		    }
		}
		
		self.purge().then((rsp) => {
		    resolve(rsp);
		}).catch((err) => {
		    reject(err);
		});
	    });
	},
	
	purge: function(){

	    return new Promise((resolve, reject) => {
		
		caches.keys().then(function (cachesNames) {
		    
                    console.debug("Delete " + document.defaultView.location.origin + " caches");

                    return Promise.all(cachesNames.map(function (cacheName) {
			
			if (! cacheName.startsWith("worldclock-")){
			    return Promise.resolve();
			}
			
			return caches.delete(cacheName).then(function () {
			    console.debug("Cache with name " + cacheName + " is deleted");
			}); 
                    }))
                
		}).then(function () {
                    console.debug("All " + document.defaultView.location.origin + " caches are deleted");
		    resolve();
		}).catch((err) => {
		    console.error("Failed to remove caches, ",err);
		    reject(err);
		});
		
	    });
	},
    };

    return self;
    
})();

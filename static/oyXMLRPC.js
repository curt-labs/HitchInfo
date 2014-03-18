//
//
//  JavaScript Refactoring for safer, faster, better AJAX.
//
//  Copyright 2005, Pavel Simakov, 
//  http://www.softwaresecretweapons.com
//
//

function oyXMLRPCProvider () {	

	var status = null;
	var url = null;	
	var req = null;
	var msgCount = 0;	
	var inProgress = false;
	var isComplete = false;
		
	var oThis = this;
	
	// checks to see if we have too many messages in log
	var internalCanMsg = function(){
		msgCount++;
		return msgCount < 100;
	}
	
	// adds message to internal log
	var internalOnLog = function(msg){
		if(oThis.onLog && internalCanMsg()) {
			oThis.onLog(msg);
		}
	}
	
	// adds message to internal error handler
	var internalOnError = function(msg){
		if(oThis.onError && internalCanMsg()) {
			oThis.onError(msg);
		}
	}	
	
	// tells us whether we are busy waiting for the response to another requst
	var internalIsBusy = function(){
		return inProgress && !isComplete;
	}	
	
	// internal callback function for the browser; it is called when a state of a request object changes
	var internalRequestComplete = function() {
				
		var STATE_COMPLETED = 4;
		var STATUS_200 = 200;
				
		if (!internalIsBusy()) {
			internalOnError("internalRequestComplete: error - no request submitted");
		}
		
		internalOnLog("internalRequestComplete: readyState " + req.readyState);
		
		if (req.readyState == STATE_COMPLETED) {
			status = req.status;
			inProgress = false;
			isComplete = true;

			internalOnLog("internalRequestComplete: status " + status);
			
			if (status == STATUS_200) {
				internalOnLog("internalRequestComplete: calling callback on content with length " + req.responseText.length + " chars");			
				if(oThis.onComplete) {
					oThis.onComplete(req.responseText, req.responseXML);
				}				 
				internalOnLog("internalRequestComplete: complete on " + new Date());
			} else {
				internalOnError("internalRequestComplete: error - bad status while fetching " + url);
			}
		} else {
			// we need to review other state codes for XMLRPC provider
		}
	}	
	
	// call this function to figure out version of this class
	this.getVersion = function(){
		return "1.0.0";
	}
	
	// call this function to figure out if current browser supports XML HTTP Requests
	this.isSupported = function(){
		var nonEI = window.XMLHttpRequest;
		var onIE = window.ActiveXObject;
		
		if (onIE) {	    		
			onIE = new ActiveXObject("Microsoft.XMLHTTP") != null;
		}
		
		return window.XMLHttpRequest || onIE;
	}
	
	// call this function to find out if more calls are possible and if request has been completely received 
	this.isBusy = function(){
		return internalIsBusy();
	}		

	//  call this function to submit new request
	this.submit = function(_url){	
		if (internalIsBusy()) {
			internalOnError("submit: error - busy processing another request " + _url);			
		}	
		
		msgCount = 0;
		internalOnLog("submit: started on " + new Date() + " for " + _url);
				
		url = _url;	
		status = null;
		inProgress = true;
		isComplete = false;
		
	    if (window.XMLHttpRequest) {
	    
	    	// branch for native XMLHttpRequest object
			internalOnLog("submit: using XMLHttpRequest()");
	    
	        req = new XMLHttpRequest();
	        req.onreadystatechange = internalRequestComplete;
	        req.open("GET", url, true);
	                
        	req.send(null);	
	        	    
	    } else { 
	    	    	
	    	// branch for IE/Windows ActiveX version
	    	if (window.ActiveXObject) {	    		
		        req = new ActiveXObject("Microsoft.XMLHTTP");
		        if (req) {
		        
					internalOnLog("submit: using Microsoft.XMLHTTP");
		        
		            req.onreadystatechange = internalRequestComplete;
		            req.open("GET", url, true);
			    	req.setrequestheader("Pragma","no-cache");
		   	    	req.setrequestheader("Cache-control","no-cache");
		           
		        	req.send();	
		        } else {
					internalOnError("submit: error - unable to create Microsoft.XMLHTTP");
		        }
		    } else {
				internalOnError("submit: error - browser does not support XML HTTP Request");
		    }
	    }
		
		internalOnLog("submit: complete");
	}
	
	// call this function to abort current request
	this.abort = function(){
		internalOnLog("abort: " + url);
		
		if (!internalIsBusy()) {
			internalOnError("abort: error - no request submitted");			
		}
	
		onComplete = null;		
		req.abort();
	}	

	// call this function to find out current url
	this.getUrl = function(){
		return url;
	}
	
	// call this function to find out HTTP status code after response completes
	this.getStatus = function(){
		return status;
	}	
	
	// please can override this; this is function called when fatal error occurs
	this.onError = function(msg){
		
	} 
	
	// user can override this; this function  is called when log message is created	
	this.onLog = function(msg) {
	
	}	
	
	// user can override this;  this function is called when response is received without errors
	this.onComplete = function(responseText, responseXML){
		
	}
	
}
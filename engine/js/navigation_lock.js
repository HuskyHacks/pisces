(function navigationLock(){
	function handler(event) {
		event.preventDefault();
		event.returnValue = true;
	}

	window.addEventListener("beforeunload", handler);
	window.__pisces_beforeUnloadHandler = handler

	return "locked";
}())

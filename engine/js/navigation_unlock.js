(function navigationUnlock() {
	window.removeEventListener("beforeunload", window.__pisces_beforeUnloadHandler);
	window.__pisces_beforeUnloadHandler = undefined;

	return "unlocked";
}())

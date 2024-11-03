function setupReloadSocket(reload = false) {
	const reloadWebsocket = new WebSocket(reloadPath);
	let doReloadNext = reload;
	reloadWebsocket.onopen = function () {
		if (reload === true) {
			window.location.reload();
		} else {
			doReloadNext = true;
		}
	};
	reloadWebsocket.onerror = function onError() {
		setTimeout(() => setupReloadSocket(doReloadNext), 250);
	}
	reloadWebsocket.onclose = function onClose() {
		setTimeout(() => setupReloadSocket(doReloadNext), 250);
	};
}

setupReloadSocket();

document.onload = getTheme();

function setCookie(theme) {
	var d = new Date();
	d.setTime(d.getTime() + (365 * 24 * 60 * 60 * 1000));
	var expires = "expires="+d.toUTCString();
	document.cookie = "theme=" + theme + ";" + expires + ";path=/";
}

function getCookie(cname) {
	var name = cname + "=";
	var decodedCookie = decodeURIComponent(document.cookie);
	var ca = decodedCookie.split(';');
	for(var i = 0; i <ca.length; i++) {
		var c = ca[i];
			while (c.charAt(0) == ' ') {
		c = c.substring(1);
		}

		if (c.indexOf(name) == 0) {
			return c.substring(name.length, c.length);
		}
	}
	return "";
}

function getTheme() {
	if (document.cookie.length == 0) {
		setCookie("light");
	}

	theme = getCookie("theme");
	if (theme == "dark") {
		document.getElementById("theme").href = "/plon/dark.css";
	} else {
		document.getElementById("theme").href = "/plon/light.css";
	}
}

function toggleTheme() {
	theme = getCookie("theme")
	if (theme == "light") {
		setCookie("dark");
	} else {
		setCookie("light");
	}

	location.reload();
}

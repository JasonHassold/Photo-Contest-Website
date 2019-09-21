// Load img, get ratio, and set div background
var img = new Image();
img.src = "/assets/images/tinted-camera-wide-min.png";
var ratio;
var div;
var text = "Enter for a chance to win!";

// Wait for image to load
img.onload = function() {
	// Get ration of image
	ratio = img.naturalHeight / img.naturalWidth;
	console.log(ratio);
	// Set div background image
	document.getElementById("camera").style.backgroundImage = 
		"url('" + img.src + "')";
	document.getElementById("camera").style.backgroundSize = "cover";
	// Get div and adjust
	div = document.getElementById("camera");
	adjust();
	window.addEventListener('resize', adjust);
}

// Updater
function adjust() {
	image();
	padding();
	hTag();
}

// Adjusts div hieght and width using ratio
function image() {
	div.style.width = div.parentElement.clientWidth + "px";
	div.style.height = ratio * div.parentElement.clientWidth + "px";
} 

// Adjust padding above call to action
function padding() {
	if (div.parentElement.clientWidth > 1270) {
		document.getElementById("calltoaction").style.paddingTop = "20%";
	} else if (div.parentElement.clientWidth > 794) {
		document.getElementById("calltoaction").style.paddingTop = "10%";
	} else {
		document.getElementById("calltoaction").style.paddingTop = "0";
	}
}

// Adjust <hx> tag
function hTag() {
	if (div.parentElement.clientWidth <= 535 && 
		document.getElementById("h4").innerHTML == "") {
			document.getElementById("h1").innerHTML = "";
			document.getElementById("h2").innerHTML = "";
			document.getElementById("h4").innerHTML = text;
	} else if (div.parentElement.clientWidth > 535 &&
		div.parentElement.clientWidth <= 692 && 
		document.getElementById("h2").innerHTML == "") {
			document.getElementById("h1").innerHTML = "";
			document.getElementById("h4").innerHTML = "";
			document.getElementById("h2").innerHTML = text;
	} else if (div.parentElement.clientWidth > 692 && 
		document.getElementById("h1").innerHTML == "") {
			document.getElementById("h2").innerHTML = "";
			document.getElementById("h4").innerHTML = "";
			document.getElementById("h1").innerHTML = text;
	}
}
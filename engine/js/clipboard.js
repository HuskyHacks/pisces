(function clipboard() {
	const text = navigator.clipboard.readText();

	navigator.clipboard.writeText("");

	return text;
}())

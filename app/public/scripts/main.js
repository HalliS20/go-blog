// Purpose: This file contains the JavaScript code for the single-page application (SPA) loader.
//
// This code is responsible for loading content from the server and updating the page without a full page reload.

spaLoader()

function spaLoader() {
	document.addEventListener("DOMContentLoaded", () => {
		attachAllListeners()

		window.addEventListener("popstate", () => {
			fetchContent(window.location.pathname)
		}) // Listen for back/forward button clicks and fetch content
	})
}

function fetchContent(url) {
	const contentDiv = document.getElementById("content")
	fetch(url)
		.then((response) => response.text())
		.then((html) => {
			const parser = new DOMParser()
			const doc = parser.parseFromString(html, "text/html")

			// Update the content
			const newContent = doc.getElementById("content")
			contentDiv.innerHTML = newContent.innerHTML

			// Update the head
			const newHead = doc.head
			document.head.innerHTML = newHead.innerHTML

			attachAllListeners() // Reattach listeners to new content
		})
}

function attachButtonListener(button) {
	button.addEventListener(
		"click",
		(e) => {
			e.preventDefault() // Prevent default button behavior
			e.stopPropagation() // Stop event from bubbling up
			const url = button.getAttribute("data-href")
			if (url && url.startsWith("/")) {
				history.pushState(null, "", url)
				fetchContent(url)
			}

			return false // Prevent onclick from firing
		},
		true,
	) // Use capturing to ensure this runs before onclick
}

function attachAllListeners() {
	document
		.querySelectorAll('button[data-href^="/"]')
		.forEach(attachButtonListener)
}


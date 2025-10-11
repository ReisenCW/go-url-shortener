function onSubmit() {
	const urlInput = document.getElementById('urlInput');
	const userIdInput = document.getElementById('userIdInput');
	const shortenedInput = document.getElementById('shortenedUrl');
	const shortenButton = event?.target || document.querySelector('button');

	const longUrl = urlInput.value.trim();
	const userId = userIdInput ? userIdInput.value.trim() || 'anonymous' : 'anonymous';

	if (!longUrl) {
		alert('Please enter a URL to shorten');
		return;
	}

	// disable button to prevent duplicate submissions
	if (shortenButton) shortenButton.disabled = true;
	shortenedInput.value = 'Shortening...';

    console.log('Sending request to shorten URL:', longUrl, 'for user:', userId);

	fetch('http://localhost:9808/create-short-url', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ long_url: longUrl, user_id: userId })
	})
	.then(response => {
		if (!response.ok) throw new Error('Server returned ' + response.status);
		return response.json();
	})
	.then(data => {
		const shortenedInput = document.getElementById('shortenedUrl');
		const copyBtn = document.getElementById('copyBtn');
		if (data.short_url) {
			shortenedInput.value = data.short_url;
			// make input act as a clickable link (open in new tab)
			shortenedInput.style.cursor = 'pointer';
			shortenedInput.addEventListener('click', () => {
				const href = shortenedInput.value;
				if (href) window.open(href, '_blank', 'noopener');
			});
			// show copy button next to input
			if (copyBtn) {
				copyBtn.style.display = 'inline-block';
			}
		} else if (data.error) {
			shortenedInput.value = '';
			alert('Error: ' + data.error);
		} else {
			shortenedInput.value = '';
			alert('Unexpected response from server');
		}
	})
	.catch(err => {
		shortenedInput.value = '';
		alert('Request failed: ' + err.message);
	})
	.finally(() => {
		if (shortenButton) shortenButton.disabled = false;
	});
}

// allow pressing Enter in the URL input to trigger shortening
document.addEventListener('DOMContentLoaded', () => {
	const urlInput = document.getElementById('urlInput');
	urlInput.addEventListener('keydown', (e) => {
		if (e.key === 'Enter') {
			e.preventDefault();
			onSubmit();
		}
	});
	const copyBtn = document.getElementById('copyBtn');
	if (copyBtn) {
		copyBtn.addEventListener('click', () => {
			const shortVal = document.getElementById('shortenedUrl').value;
			if (!shortVal) return;
			navigator.clipboard.writeText(shortVal).then(() => {
				copyBtn.textContent = 'Copied';
				setTimeout(() => { copyBtn.textContent = 'Copy'; }, 1500);
			}).catch(() => {
				alert('Copy failed');
			});
		});
	}
});
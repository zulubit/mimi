export default function mimiJson(element, type = 'application/json') {
  const script = element.querySelector(`script[type="${type}"]`);
  if (script) {
    try {
      return JSON.parse(script.textContent);
    } catch (error) {
      console.error('Invalid JSON in <script>: ', error);
    }
  } else {
    console.warn(`No <script type="${type}"> tag found inside <${element.localName}>.`);
  }
  return {}; // Return an empty object if parsing fails or no script is found
}


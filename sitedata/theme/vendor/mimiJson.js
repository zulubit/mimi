export default function mimiJson(element, attributeName = 'mimi-data') {
  const data = element.getAttribute(attributeName);
  if (data) {
    try {
      return JSON.parse(data);
    } catch (error) {
      console.error(`Invalid JSON in ${attributeName} attribute:`, error);
    }
  } else {
    console.warn(`No ${attributeName} attribute found inside <${element.localName}>.`);
  }
  return {}; // Return an empty object if parsing fails or no attribute is found
}


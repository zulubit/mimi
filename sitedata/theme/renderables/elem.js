import { LitElement, html } from '../vendor/lit.min.js';

class MyElement extends LitElement {
  static properties = {
    data: { type: Object }, // Property to store parsed data
  };

  constructor() {
    super();
    this.data = {}; // Initialize as an empty object
  }

  connectedCallback() {
    super.connectedCallback();
    // Parse the 'mimi-data' attribute as JSON and assign it to 'data'
    const attrValue = this.getAttribute('mimi-data');
    if (attrValue) {
      try {
        this.data = JSON.parse(attrValue);
      } catch (error) {
        console.error('Invalid JSON in mimi-data:', error);
      }
    }
  }

  render() {
    return html`
      <div>
        <h1>${this.data.title || 'No Title'}</h1>
        <p>${this.data.description || 'No Description'}</p>
        <h3>Keywords:</h3>
        <ul>
          ${this.data.meta?.keywords?.map(
            keyword => html`<li>${keyword}</li>`
          )}
        </ul>
      </div>
    `;
  }
}

customElements.define('my-element', MyElement);


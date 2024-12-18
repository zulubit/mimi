import { LitElement, html } from '../vendor/lit.min.js';
import mimiJson from '../vendor/mimiJson.js';

class MyElement extends LitElement {
  static properties = {
    data: { type: Object }, // Data parsed from mimi-data attribute
  };

  constructor() {
    super();
    this.data = {}; 
  }

  connectedCallback() {
    super.connectedCallback();
    // Parse the mimi-data attribute
    this.data = mimiJson(this);
  }

  render() {
    return html`
      <div>
        <h1>${this.data.heading || 'No Heading'}</h1>
        <p>${this.data.subheading || 'No Subheading'}</p>
        <p><strong>CTA:</strong> ${this.data.cta?.text || 'No CTA'} (<a href="${this.data.cta?.link}">Link</a>)</p>
        <p><strong>Background Image:</strong> ${this.data.backgroundImage || 'No Background Image'}</p>
        <slot></slot> <!-- Render child elements here -->
      </div>
    `;
  }
}

customElements.define('my-element', MyElement);


import { LitElement, html, css } from "lit";

class EditHeadingBlock extends LitElement {
  static properties = {
    blockData: { type: Object },
  };

  constructor() {
    super();
    // Default value
    this.blockData = { content: "" };
  }

  static styles = css`
    .edit-heading-block {
      padding: 0.75rem;
      border: 1px solid #ddd;
      border-radius: 4px;
      background-color: #f9f9f9;
    }

    label {
      display: block;
      margin-bottom: 0.5rem;
      font-weight: bold;
    }

    input {
      width: 100%;
      padding: 0.5rem;
      border: 1px solid #ccc;
      border-radius: 4px;
      font-size: 1rem;
    }
  `;

  render() {
    return html`
      <div class="edit-heading-block">
        <label for="heading-content">Heading Text</label>
        <input
          type="text"
          id="heading-content"
          .value=${this.blockData.content || ""}
          @input=${this._handleContentChange}
          placeholder="Enter heading text"
        />
      </div>
    `;
  }

  _handleContentChange(e) {
    const newBlockData = {
      ...this.blockData,
      content: e.target.value,
    };

    this.blockData = newBlockData;

    // Dispatch an event to notify parent components of the change
    this.dispatchEvent(
      new CustomEvent("block-updated", {
        detail: { blockData: this.blockData },
        bubbles: true,
        composed: true,
      }),
    );
  }
}

customElements.define("edit-heading-block", EditHeadingBlock);

export { EditHeadingBlock };

export class DOM {
  static createElement(tagName, attributes = {}, textContent = "") {
    const element = document.createElement(tagName);
    for (const key in attributes) {
      element.setAttribute(key, attributes[key]);
    }
    element.textContent = textContent;
    return element;
  }
  static createOneElement(tagName, attribute = "", contentAttribute = "", text) {
    const element = document.createElement(tagName);
    element.setAttribute('data-value', text);

    element.setAttribute(attribute, contentAttribute);
    return element;
  }

  static createCase(tagName, attribute = "", contentAttribute = "", value) {
    const element = document.createElement(tagName);
    element.setAttribute(attribute, contentAttribute);
    element.dataset.value = value
    return element;
  }

  static getParentAtribut(selector, attributeName) {
    const element = document.querySelector(selector);
    if (element) {
      const parent = element.parentElement;
      return parent.getAttribute(attributeName);
    } else {
      console.error(`Element "${selector}" not found.`);
      return null;
    }
  }

  static getByTag(tagName) {
    return document.getElementsByTagName(tagName);
  }

  static getByClass(className) {
    return document.getElementsByClassName(className);
  }

  static getById(id) {
    return document.getElementById(id);
  }

  static append(parentSelector, child) {
    const parentElement = document.querySelector(parentSelector);
    if (parentElement) {
      parentElement.appendChild(child);
    } else {
      console.error(`Parent element "${parentSelector}" not found.`);
    }
  }

  static setAttribute(selector, attributeName, attributeValue) {
    const elements = document.querySelectorAll(selector);
    elements.forEach((element) => {
      element.setAttribute(attributeName, attributeValue);
    });
  }
  static setAttributeNameTag(nameTag, attributeName, attributeValue) {
    nameTag.setAttribute(attributeName, attributeValue);
  }

  static createOneElement(tagName, attribute = "", contentAttribute = "") {
    const element = document.createElement(tagName);
    element.setAttribute(attribute, contentAttribute);
    return element;
  }
  static append(parentSelector, child) {
    const parentElement = document.querySelector(parentSelector);
    if (parentElement) {
      parentElement.appendChild(child);
    } else {
      console.error(`Parent element "${parentSelector}" not found.`);
    }
  }

  static getAttribute(selector, attributeName) {
    const element = document.querySelector(selector);
    if (element) {
      return element.getAttribute(attributeName);
    } else {
      console.error(`Element "${selector}" not found.`);
      return null;
    }
  }

  static addClass(selector, className) {
    const elements = document.querySelectorAll(selector);
    elements.forEach((element) => {
      element.classList.add(className);
    });
  }

  static removeClass(selector, className) {
    const elements = document.querySelectorAll(selector);
    elements.forEach((element) => {
      element.classList.remove(className);
    });
  }

  static removeElement(selector) {
    const element = document.querySelector(selector);
    if (element && element.parentNode) {
      element.parentNode.removeChild(element);
    }
  }

  static setHTML(parentSelector, child) {
    const parentElement = document.querySelector(parentSelector);
    if (parentElement) {
      parentElement.innerHTML = child;
    } else {
      console.error(`Parent element "${parentSelector}" not found.`);
    }
  }

  static getHTML(parentSelector) {
    const parentElement = document.querySelector(parentSelector);
    if (parentElement) {
      return parentElement.innerHTML;
    } else {
      console.error(`Parent element "${parentSelector}" not found.`);
    }
  }

  static appendByID(parentSelector, child) {
    const parentElement = document.getElementById(parentSelector);
    if (parentElement) {
      parentElement.appendChild(child);
    } else {
      console.error(`Parent element "${parentSelector}" not found.`);
    }
  }

  static addEventListener(selector, event, handler) {
    const elements = document.querySelectorAll(selector);
    elements.forEach((element) => {
      element.addEventListener(event, handler);
    });
  }

  static addEventListenerNode(selector, handler) {
    const elements = document.querySelectorAll(selector);
    elements.forEach((element) => {
      element.addEventListener(event, handler);
    });
  }

  static addEventListenerOne(selector, event, handler) {
    const element = document.getElementById(selector);
    element.addEventListener(event, handler);
  }

  static removeEventListener(selector, event, handler) {
    const elements = document.querySelectorAll(selector);
    elements.forEach((element) => {
      element.removeEventListener(event, handler);
    });
  }

  static getValue(selector) {
    const element = document.querySelector(selector);
    if (element && ["INPUT", "TEXTAREA", "SELECT"].includes(element.tagName)) {
      return element.value;
    } else {
      console.error(`Element "${selector}" not found or not a form field.`);
      return null;
    }
  }
  static getValueByID(selector) {
    const element = document.getElementById(selector);
    if (element && ["INPUT", "TEXTAREA", "SELECT"].includes(element.tagName)) {
      return element.value;
    } else {
      console.error(`Element "${selector}" not found or not a form field.`);
      return null;
    }
  }

  static setValue(selector, value) {
    const element = document.querySelector(selector);
    if (element && ["INPUT", "TEXTAREA", "SELECT"].includes(element.tagName)) {
      element.value = value;
    } else {
      console.error(`Element "${selector}" not found or not a form field.`);
    }
  }

  static getTargetID(e) {
    return e.target.id;
  }
  static getTargetText(e) {
    return e.target.textContent;
  }
  static getTarget(e) {
    return e.target.tagName;
  }
  static getClosestTargetId(e, parent) {
    return e.target.closest(parent).dataset.id;
  }
  static contains(e, content) {
    return e.target.classList.contains(content);
  }

  static insertBefore(node, input) {
    node.parentNode.insertBefore(input, node);
  }
  static removeNode(node) {
    node.remove();
  }
  static querySelectorAll(tag) {
    return document.querySelectorAll(tag);
  }


  static querySelector(tag) {
    return document.querySelector(tag);
  }

     static addOneClass(selector, className) {
      selector.classList.add(className);
  }
  
}





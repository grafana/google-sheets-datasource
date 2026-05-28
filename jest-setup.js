// Jest setup provided by Grafana scaffolding
import './.config/jest-setup';

// Mock canvas getContext('2d') for @grafana/ui's measureText which is not supported in JSDOM
HTMLCanvasElement.prototype.getContext = () => ({
  measureText: (text) => ({ width: text.length * 8 }),
  font: '',
});

// Mock IntersectionObserver which is not available in JSDOM
global.IntersectionObserver = class IntersectionObserver {
  constructor(callback) {
    this.callback = callback;
  }
  observe() {}
  unobserve() {}
  disconnect() {}
};

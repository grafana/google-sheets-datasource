// Jest setup provided by Grafana scaffolding
import './.config/jest-setup';

// Mock IntersectionObserver which is not available in JSDOM
global.IntersectionObserver = class IntersectionObserver {
  constructor(callback) {
    this.callback = callback;
  }
  observe() {}
  unobserve() {}
  disconnect() {}
};

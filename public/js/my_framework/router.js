export class Router {
  constructor() {
    this.routes = {};
    this.currentRoute = '/';
  }
  addRoute(route, callback) {
    this.routes[route] = callback;
  }
  navigate(route) {
    if (this.routes[route]) {
      this.currentRoute = route;
      this.routes[route]();
    } else {
      this.navigate('/');
    }
  }
}
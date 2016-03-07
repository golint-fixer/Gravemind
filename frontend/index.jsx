import React from 'react';
import ReactDOM from 'react-dom';
import { Router, Route, IndexRoute, browserHistory } from 'react-router';

import Application from 'components/Application';
import Home from 'components/Home';
import NotFound from 'components/NotFound';
import About from 'components/About';

// Force load our CSS
require('assets/stylesheets/index.scss');

// Copy static assets
require.context('static', true);

// Render routes and use HTML5 browser history
ReactDOM.render((
  <Router history={browserHistory}>
    <Route path="/" component={Application}>
      <IndexRoute component={Home} />
      <Route path="about" component={About} />
      <Route path="*" component={NotFound}/>
    </Route>
  </Router>
), document.body);

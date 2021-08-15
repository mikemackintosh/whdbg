import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter } from "react-router-dom";

import './App.scss';
import App from './App';
import reportWebVitals from './reportWebVitals';

const render = Component => {
  return ReactDOM.render(
    <BrowserRouter>
      <Component />
    </BrowserRouter>,
    document.getElementById('root')
  );
};

render(App);

if (module.hot) {
  module.hot.accept('./App', () => {
    const NextApp = require('./App').default;
    render(NextApp);
  });
}

reportWebVitals();

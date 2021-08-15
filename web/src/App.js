import React, { useState, useCallback, useMemo, useRef } from 'react';
import Logo from './assets/logo.svg';

import CookieConsent from "react-cookie-consent";
import {
  withRouter,
  Route,
  Switch,
  BrowserRouter as Router,
} from "react-router-dom";

import routes from './routes.js';
import NotFound from './pages/NotFound.js';

function App() {
  return (
    <>
    <header>
        <div className="logo">
          <img className="logo-img" src={Logo} alt="whdbg.dev" />
          <div className="brand">
            <span>Webhook Debugger</span>
          </div>
        </div>
    </header>
    <div className="container-fluid">

      <Switch>
        {routes.map((route, idx) => (
          <Route path={route.path} exact component={route.component} key={idx} />
        ))}
        <Route component={NotFound} />
      </Switch>

      <div className="copyright">
        &copy; 2021 - Mike Mackintosh
      </div>

      <CookieConsent
        location="bottom"
        buttonText="Sounds good!"
        cookieName="postmaster-cc-store1"
        style={{ background: "rgb(32 40 64)" }}
        buttonStyle={{ background: "#335eea", color: "#ffffff", fontSize: "14px", borderRadius: '8px' }}
        expires={7}>
        <b>We use cookies to enhance your experience.<br/></b>
        <span style={{ fontSize: "12px" }}>
          For the record, we use cookies to improve your experience, but we don't collect, store or sell any data.
        </span>
      </CookieConsent>
    </div>
    </>
  );
}

export default App;

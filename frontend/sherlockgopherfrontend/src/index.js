import React from 'react';
import ReactDOM from 'react-dom';
import { Route, BrowserRouter as Router } from 'react-router-dom'
import {createBrowserHistory} from 'history';
import './index.css';
import * as serviceWorker from './serviceWorker';

// Import all Routes
import App from './App'; 
import SitePerformence from './siteperformence.js';
import Controls from './controls.js';
import ScriptingAPI from './scriptingapi.js'
import Nodedetails from './nodedetails.js'
import NodeGraph from './graph.js'
import Screenshotservice from './screenshotservice.js'
import Imagemetadataservice from './imagemetadataservice.js'


const browserHistory = createBrowserHistory();

const router = (
  <Router history={browserHistory}>
    <div>
      <Route exact path="/" component={App}  />
      <Route  path="/Graph" component={NodeGraph} />
      <Route  path="/sitesperformence" component={SitePerformence}  />
      <Route  path="/nodedetails" component={Nodedetails}  />
      <Route  path="/controls" component={Controls}  />
      <Route  path="/screenshots" component={Screenshotservice}  />
      <Route  path="/imagemetadata" component={Imagemetadataservice}  />
      <Route  path="/scriptingapi" component={ScriptingAPI}  />
    </div>
  </Router>
)


ReactDOM.render(
  router,
  document.getElementById('root')
);

serviceWorker.unregister();

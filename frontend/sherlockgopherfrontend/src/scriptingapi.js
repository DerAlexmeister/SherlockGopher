import React from 'react';
import SyntaxHighlighter from 'react-syntax-highlighter';
import { docco } from 'react-syntax-highlighter/dist/esm/styles/hljs';

import SearchBar from './searchbar.js';

import './assets/css/App.css'

export default class ScriptingAPI extends React.Component {
    render() {
        return (
          <div>
            
              <SearchBar></SearchBar>
              <div class="body">
                <p>Scripting API</p> 
                <hr></hr>
                <ul class="list-group" style={{color: "black"}}>
                  <p><h4>General</h4></p>
                  <li class="list-group-item"><a href="#areyouthere">/areyouthere</a></li>
                </ul> 
                <br></br>
                <ul class="list-group" style={{color: "black"}}>
                  <p><h4>Graph</h4></p>
                  <li class="list-group-item"><a href="#graphmeta">/meta</a></li>
                  <li class="list-group-item"><a href="#graphall">/all</a></li>
                  <li class="list-group-item"><a href="#graphalloptimized">/alloptimized</a></li>
                  <li class="list-group-item"><a href="#graphpos">/performenceofsites</a></li>
                  <li class="list-group-item"><a href="#graphdos">/detailsofnode</a></li>
                  <li class="list-group-item"><a href="#graphsearch">/search</a></li>
                </ul>
                <br></br> 
                <ul class="list-group" style={{color: "black"}}>
                  <p><h4>Monitor</h4></p>
                  <li class="list-group-item"><a href="#monitormeta">/meta</a></li>
                </ul> 
                <br></br>
                <ul class="list-group" style={{color: "black"}}>
                  <p><h4>Controller</h4></p>
                  <li class="list-group-item"><a href="#controllerchangestate">/changestate</a></li>
                  <li class="list-group-item"><a href="#controllerstatus">/status</a></li>
                  <li class="list-group-item"><a href="#controllerdropit">/dropit</a></li>
                  </ul>
                <br></br>
                <h3>General</h3><hr></hr>

                <div id="areyouthere" style={{backgroundColor: "#e5e5e5"}}  class="alert alert-dark" role="alert">
                  <h4 class="alert-heading"># /areyouthere</h4>
                  <hr></hr>
                  <p>
                  <ul class="list-group" style={{color: "black"}}>
                    <li class="list-group-item"><b>Address:</b> [Basicaddress]/areyouthere</li>
                    <li class="list-group-item"><b>HTTP/s Methods:</b> <span class="badge badge-pill badge-info">GET</span></li>
                    <li class="list-group-item"><b>Description:</b> Will send you a "yes i am here" as json to respond to the "ping".</li>
                  </ul>
                  </p>
                  <hr></hr>
                  <p class="mb-0">
                    <p>Response: </p>
                    <SyntaxHighlighter language="javascript" style={docco}>
                      {JSON.stringify({
                          "Message": "Yes i am here!"
                      })}
                    </SyntaxHighlighter>
                  </p>
                </div>


                <h3>Graph</h3><hr></hr>
                <div id="graphmeta" style={{backgroundColor: "#e5e5e5"}}  class="alert alert-dark" role="alert">
                  <h4 class="alert-heading"># /graph/v1/meta</h4>
                  <hr></hr>
                  <p>
                  <ul class="list-group" style={{color: "black"}}>
                    <li class="list-group-item"><b>Address:</b> [Basicaddress]/graph/v1/meta</li>
                    <li class="list-group-item"><b>HTTP/s Methods:</b> <span class="badge badge-pill badge-info">GET</span></li>
                    <li class="list-group-item"><b>Description:</b> Will return some meta information of the Neo4J db such as the amount of nodes or number of relationships.</li>
                  </ul>
                  </p>
                  <hr></hr>
                  <p class="mb-0">
                    <p>Response: </p>
                    <SyntaxHighlighter language="javascript" style={docco}>
                      {JSON.stringify([
                          {
                              "amountofimages": 18
                          },
                          {
                              "amountofsheets": 10
                          },
                          {
                              "amountofjs": 20
                          } ])}
                    </SyntaxHighlighter>
                  </p>
                </div>
                
                
                <div id="graphall" style={{backgroundColor: "#e5e5e5"}}  class="alert alert-dark" role="alert">
                  <h4 class="alert-heading"># /graph/v1/all</h4>
                  <hr></hr>
                  <p>
                  <ul class="list-group" style={{color: "black"}}>
                    <li class="list-group-item"><b>Address:</b> [Basicaddress]/graph/v1/all</li>
                    <li class="list-group-item"><b>HTTP/s Methods:</b> <span class="badge badge-pill badge-info">GET</span></li>
                    <li class="list-group-item"><b>Description:</b> Will return json which contains all nodes and their relationships to other nodes.</li>
                  </ul>
                  </p>
                  <hr></hr>
                  <p class="mb-0">
                    <p>Response: </p>
                    <SyntaxHighlighter language="javascript" style={docco}>
                      {JSON.stringify([
                        {
                          "Destination": "www.example.com/6",
                          "DestinationType": "Javascript",
                          "Relationship": "Requires",
                          "Source": "www.example.com/23",
                          "SourceType": "HTML"
                      },
                      {
                          "Destination": "www.example.com/7",
                          "DestinationType": "CSS",
                          "Relationship": "Requires",
                          "Source": "www.example.com/79",
                          "SourceType": "HTML"
                      }
                      ])}
                    </SyntaxHighlighter>
                  </p>
                </div>
                
                
                <div id="graphalloptimized" style={{backgroundColor: "#e5e5e5"}}  class="alert alert-dark" role="alert">
                  <h4 class="alert-heading">/graph/v1/alloptimized</h4>
                  <hr></hr>
                  <p>
                  <ul class="list-group" style={{color: "black"}}>
                    <li class="list-group-item"><b>Address:</b> [Basicaddress]/graph/v1/alloptimized</li>
                    <li class="list-group-item"><b>HTTP/s Methods:</b> <span class="badge badge-pill badge-info">GET</span></li>
                    <li class="list-group-item"><b>Description:</b> 
                    Will return json which contains all nodes and their relationships to other nodes but optimized for the react frontend.
                    </li>
                  </ul>
                  </p>
                  <hr></hr>
                  <p class="mb-0">
                    <p>Response: </p>
                    <SyntaxHighlighter language="javascript" style={docco}>
                      {JSON.stringify({
                          "links": [
                            {
                                "color": "#08206A",
                                "label": "Links",
                                "source": "www.example.com/43",
                                "target": "www.example.com/0"
                            },
                            {
                                "color": "#08206A",
                                "label": "Links",
                                "source": "www.example.com/41",
                                "target": "www.example.com/0"
                            }],
                            "nodes": [
                            {
                                "color": "#7CA9EF",
                                "id": "www.example.com/43"
                            },
                            {
                                "color": "#7CA9EF",
                                "id": "www.example.com/0"
                            },
                            {
                                "color": "#7CA9EF",
                                "id": "www.example.com/41"
                            }] 
                      })}
                    </SyntaxHighlighter>
                  </p>
                </div>
                
                <div id="graphpos" style={{backgroundColor: "#e5e5e5"}}  class="alert alert-dark" role="alert">
                  <h4 class="alert-heading"># /performenceofsites</h4>
                  <hr></hr>
                  <p>
                  <ul class="list-group" style={{color: "black"}}>
                    <li class="list-group-item"><b>Address:</b> [Basicaddress]/graph/v1/performenceofsites</li>
                    <li class="list-group-item"><b>HTTP/s Methods:</b> <span class="badge badge-pill badge-info">GET</span></li>
                    <li class="list-group-item"><b>Description:</b> Will return a list of performence indicators like statuscode and RTT.</li>
                  </ul>
                  </p>
                  <hr></hr>
                  <p class="mb-0">
                    <p>Response: </p>
                    <SyntaxHighlighter language="javascript" style={docco}>
                      {JSON.stringify([
                        {
                          "Address": "www.example.com/0",
                          "ResponseTime": "522",
                          "Status": "409"
                      },
                      {
                          "Address": "www.example.com/1",
                          "ResponseTime": "1034",
                          "Status": "246"
                      }
                      ])}
                    </SyntaxHighlighter>
                  </p>
                </div>
                
                <div id="graphdos" style={{backgroundColor: "#e5e5e5"}} class="alert alert-dark" role="alert">
                  <h4 class="alert-heading"># /detailsofnode</h4>
                  <p>
                  <hr></hr>
                  <ul class="list-group" style={{color: "black"}}>
                    <li class="list-group-item"><b>Address:</b> [Basicaddress]/graph/v1/detailsofnode</li>
                    <li class="list-group-item"><b>HTTP/s Methods:</b> <span class="badge badge-pill badge-warning">POST</span></li>
                    <li class="list-group-item"><b>Description:</b> Will send you a "yes i am here" as json to respond to the "ping".</li>
                  </ul>
                  </p>
                  <hr></hr>
                  <p class="mb-0">
                  <p>Request: </p>
                    <SyntaxHighlighter language="javascript" style={docco}>
                      {JSON.stringify({
                          "url": "www.example.com/imprint"
                      })}
                    </SyntaxHighlighter>
                  <p>Response: </p>
                    <SyntaxHighlighter language="javascript" style={docco}>
                      {JSON.stringify({
                          "www.example.com/imprint": {
                                "Accept-Ranges": "bytes",
                                "Address": "1",
                                "Age": "12",
                                "Allow": "GET, HEAD",
                                "Cache-Control": "max-age=3600",
                                "Connection": "close",
                                "Content-Encoding": "gzip",
                                "Content-Language": "de",
                                "Content-Location": "/imprint",
                                "Content-MD5": "Q2hlY2sgSW50ZWdyaXR5IQ==",
                                "Content-Type": "text/html; charset=utf-8",
                            }
                      })}
                    </SyntaxHighlighter>
                  </p>
                </div>
                
                <div id="graphsearch" style={{backgroundColor: "#e5e5e5"}} class="alert alert-dark" role="alert">
                  <h4 class="alert-heading"># /search</h4>
                  <p>
                  <ul class="list-group" style={{color: "black"}}>
                    <li class="list-group-item"><b>Address:</b> [Basicaddress]/graph/v1/search</li>
                    <li class="list-group-item"><b>HTTP/s Methods:</b> <span class="badge badge-pill badge-warning">POST</span></li>
                    <li class="list-group-item"><b>Description:</b> Will return a message containing information about the request like Fine or an error containing a message.</li>
                  </ul>
                  </p>
                  <hr></hr>
                  <p class="mb-0">
                  <p>Request: </p>
                    <SyntaxHighlighter language="javascript" style={docco}>
                      {JSON.stringify({
                          "url": "www.example.com"
                      })}
                    </SyntaxHighlighter>
                  <p>Response: </p>
                    <SyntaxHighlighter language="javascript" style={docco}>
                      {JSON.stringify({
                          "Status": "Fine"
                      })}
                    </SyntaxHighlighter>
                    </p>
                </div>

                <h3>Monitor</h3>
                <hr></hr>
                <div id="monitormeta" style={{backgroundColor: "#e5e5e5"}} class="alert alert-dark" role="alert">
                  <h4 class="alert-heading"># /meta</h4>
                  <p>
                  <ul class="list-group" style={{color: "black"}}>
                    <li class="list-group-item"><b>Address:</b> [Basicaddress]/monitor/v1/meta</li>
                    <li class="list-group-item"><b>HTTP/s Methods:</b> <span class="badge badge-pill badge-info">GET</span></li>
                    <li class="list-group-item"><b>Description:</b> Will return meta information about the services (Analyser and Crawler)..</li>
                  </ul>
                  </p>
                  <hr></hr>
                  <p class="mb-0">
                    <p>Response: </p>
                    <SyntaxHighlighter language="javascript" style={docco}>
                      {JSON.stringify({
                          "Crawler":{
                            "Website":    45,
                            "Undone":     9,
                            "Processing": 25,
                            "Finished":   10,
                            "Failed":     1,
                          },
                          "Analyser":{
                            "Website":       45,
                            "Undone":        5,
                            "Processing":    25,
                            "CrawlerError":  0,
                            "Saving":        5,
                            "SendToCrawler": 5,
                            "Finished":      5,
                          }
                      })}
                    </SyntaxHighlighter>
                  </p>
                </div>
                <h3>Controller</h3>
                <hr></hr>
                <div id="controllerchangestate" style={{backgroundColor: "#e5e5e5"}} class="alert alert-dark" role="alert">
                  <h4 class="alert-heading"># /changestate</h4>
                  <p>
                  <ul class="list-group" style={{color: "black"}}>
                    <li class="list-group-item"><b>Address:</b> [Basicaddress]/controller/v1/changestate</li>
                    <li class="list-group-item"><b>HTTP/s Methods:</b> <span class="badge badge-pill badge-warning">POST</span></li>
                    <li class="list-group-item"><b>Description:</b> Will change the state of various services.</li>
                    <li class="list-group-item"><b>Operations:</b> Clean, Stop, Pause, Resume</li>
                    <li class="list-group-item"><b>Targets:</b> Crawler, Analyser, All</li>
                  </ul>
                  </p>
                  <hr></hr>
                  <p class="mb-0">
                    <p>Request: </p>
                      <SyntaxHighlighter language="javascript" style={docco}>
                        {JSON.stringify({
                              "operation": "clean",
                              "target": "crawler"
                          })}
                      </SyntaxHighlighter>
                    <p>Response: </p>
                      <SyntaxHighlighter language="javascript" style={docco}>
                        {JSON.stringify({
                            "Status": "Fine"
                        })}
                      </SyntaxHighlighter>
                  </p>
                </div>
                
                <div id="controllerstatus" style={{backgroundColor: "#e5e5e5"}} class="alert alert-dark" role="alert">
                  <h4 class="alert-heading"># /status</h4>
                  <p>
                  <ul class="list-group" style={{color: "black"}}>
                    <li class="list-group-item"><b>Address:</b> [Basicaddress]/controller/v1/status</li>
                    <li class="list-group-item"><b>HTTP/s Methods:</b> <span class="badge badge-pill badge-info">GET</span></li>
                    <li class="list-group-item"><b>Description:</b> Will return the status of the crawler and analyser.</li>
                  </ul>
                  </p>
                  <hr></hr>
                  <p class="mb-0">
                    <p>Response: </p>
                    <SyntaxHighlighter language="javascript" style={docco}>
                      {JSON.stringify(   {
                            "Analyser": "Running",
                            "Crawler": "Paused"
                        })}
                    </SyntaxHighlighter>
                  </p>
                </div>
                
                <div id="controllerdropit" style={{backgroundColor: "#e5e5e5"}} class="alert alert-dark" role="alert">
                  <h4 class="alert-heading"># /dropit</h4>
                  <p>
                  <ul class="list-group" style={{color: "black"}}>
                    <li class="list-group-item"><b>Address:</b> [Basicaddress]/controller/v1/dropit</li>
                    <li class="list-group-item"><b>HTTP/s Methods:</b> <span class="badge badge-pill badge-info">GET</span></li>
                    <li class="list-group-item"><b>Description:</b> Will drop the database.</li>
                  </ul>
                  </p>
                  <hr></hr>
                  <p class="mb-0">
                    <p>Response: </p>
                    <SyntaxHighlighter language="javascript" style={docco}>
                      {JSON.stringify(   {
                          "Message": "Dropped the table."
                      })}
                    </SyntaxHighlighter>
                  </p>
                </div>
          </div>
        </div>
        )
      }
}
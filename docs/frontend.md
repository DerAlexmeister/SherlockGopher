# Docs of the frontend

SherlockGopher provides a seperate frontend/UI which will be delivered by a dedecated webserver.
So it is possible to exchange the frontend while the rest of the service is unaffected. 
As mentioned in the API (see ```api.md```), there are serveral endpoints to access the same data as the frontend does. This enables the user to use own crafted scripts and tools in order to use the data in a different way.

In case the frontend/UI should stay as it is the following will describe in a few sentences what the UI does and how to use it.

## Technical Aspects

```
    - Webserver: Go-Gin-Gonic
    - Frontend: REACT (Standard)
    - Rendering: Client-Side-Rendering
```

The default UI provides a REACT-App which can be accessed via the browser. The webserver will deliver the UI to the user via the browser (to see more about the system architecture ```architecture.md```). 

``` Notice:  Javascript must be activated```

Tested Browsers:
```
    - Firefox
    - Chromium
    - Brave
```

## SherlockGophers UI

### Home - /


<p align="center"><img alt="Sherlock Gopher" src="https://github.com/ob-algdatii-20ss/SherlockGopher/blob/develop/assets/frontend/Bildschirmfoto%20vom%202020-05-02%2017-03-35.png" width="950" height="350"></p>

Home will be the entry point to the UI. The user can
enter a address, which will be crawled and from there on navigate to a different sites by clicking ``` Menu``` in the top right corner.

<hr>

### Sidebar navigation

<p align="center"><img alt="Sherlock Gopher" src="https://github.com/ob-algdatii-20ss/SherlockGopher/blob/develop/assets/frontend/sidebar.png" width="650" height="350"></p>

Incase you clicked on the menu-button the side navigation bar will open up. There are some options to click on. They are all description in the following.

<hr>

### Graph - /Graph

The graph can be shown in two ways. As 2D-Directed-Graph or as 3D-Cluster. The user is able to switch between both 2D and 3D with little switch on the left side below the meta information.

#### 2D-Directed-Graph
<p align="center"><img alt="Sherlock Gopher" src="https://github.com/ob-algdatii-20ss/SherlockGopher/blob/develop/assets/frontend/Bildschirmfoto%20vom%202020-05-03%2023-47-25.png" width="750" height="350"></p>

#### 3D-Cluster

<p align="center"><img alt="Sherlock Gopher" src="https://github.com/ob-algdatii-20ss/SherlockGopher/blob/develop/assets/frontend/Bildschirmfoto%20vom%202020-05-02%2017-15-12.png" width="750" height="350"></p>

<hr>

### PoS (PerformanceOfSites) - /sitesperformence

<p align="center"><img alt="Sherlock Gopher" src="https://github.com/ob-algdatii-20ss/SherlockGopher/blob/develop/assets/frontend/Bildschirmfoto%20vom%202020-05-02%2017-03-13.png" width="650" height="350"></p>

On this page, the user can see the "performance" of each page. E.g the RTT so how long it took to request the site and get a response. Also the Response Code will be shown.

<hr>

### Controls - /controls

#### Statuschart 

<p align="center"><img alt="Sherlock Gopher" src="https://github.com/ob-algdatii-20ss/SherlockGopher/blob/develop/assets/frontend/piegraphcontrols.png" width="850" height="350"></p>

The pie chart will show the allocation of the different task types. E.g. Undone, running, finished ... .

#### Service Controls 

<p align="center"><img alt="Sherlock Gopher" src="https://github.com/ob-algdatii-20ss/SherlockGopher/blob/develop/assets/frontend/controls.png" width="750" height="350"></p>

Also the user is able to control the entire system via the webinterface. As the Image above shows there are various options. E.g. press the stop-button to stop all services, drop the Graph database via the dropit-button or ping the webserver.

<hr>

### Details of a Node - /nodedetails

<p align="center"><img alt="Sherlock Gopher" src="https://github.com/ob-algdatii-20ss/SherlockGopher/blob/develop/assets/frontend/details.png" width="750" height="350"></p>

On this Page you can access advanced information about a node in the database. Just enter the address in the search field and the web interface will show all stored information as the image above shows.

<hr>

### Scripting-API - /scriptingapi

<p align="center"><img alt="Sherlock Gopher" src="https://github.com/ob-algdatii-20ss/SherlockGopher/blob/develop/assets/frontend/apimenu.png" width="650" height="350"></p>

<p align="center"><img alt="Sherlock Gopher" src="https://github.com/ob-algdatii-20ss/SherlockGopher/blob/develop/assets/frontend/apimeta.png" width="850" height="350"></p>

There are a Scripting-API (REST) to work with. So an user will be able to access all data and work with it in a different way then SherlockGopher does. So be creative and build your own tools and find something amazing.
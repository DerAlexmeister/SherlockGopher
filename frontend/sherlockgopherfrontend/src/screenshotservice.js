import React from 'react';
import SearchBar from './searchbar.js';
import ReactPaginate from 'react-paginate';
import axios from 'axios'
import "./screenshotservice.css";

export default class Screenshotservice extends React.Component {

  SCREENSHOT = "http://0.0.0.0:8081/graph/v1/getscreenshots"
  BASICURL = "http://0.0.0.0:8081"

    state = {
        postdata: [],
        currentPage: 0,
        maxPage: 1,
        pageRange: 0,
        sdatamessage: "",
        hasSdataError: false
    }

    constructor(props) {
        super(props);
        this.handlePageClick = this
            .handlePageClick
            .bind(this);
    }

    receivedData() {
        axios.get(this.SCREENSHOT + "/" + this.state.currentPage).then(res => {
            try {
                const data = res.data.map;
                const postData = data.map(pd => (<React.Fragment>
                    <p>{pd.imageurl}</p>
                    <img src={this.BASICURL + "/static/images/" + pd.imagepath}/>
                </React.Fragment>))
                
                this.setState({
                    sdatamessage: "No error",
                    hasSdataError: false,
                    postData,
                    currentPage: res.data.currentpage,
                    maxPage: res.data.maxpage,
                    pageRange: res.data.pagerange,
                })
            } catch (error) {
                console.log(error)
                this.setState({
                    sdatamessage: "For some Reason an Error occured. Cannot process response.",
                    hasMetaError: true,
                })  
            }
        }).catch(error => {
            console.log(error)
            this.setState({
                sdatamessage: "For some Reason an Error occured. Is the Webserver up?",
                hasSdataError: true,
            })
        })
    }

    handlePageClick = (e) => {
        const selectedPage = e.selected;
        this.setState({
            currentPage: selectedPage,
        }, () => {
            this.receivedData()
        });
    };

    componentDidMount(){
        this.interval = setInterval(() => {
            try {
                this.receivedData()
            } catch(error) {
                console.log(error)
            }
        }, 5000)
    }

    componentWillUnmount() {
        clearInterval(this.interval)
    }

    render() {
        const {
            sdatamessage,
            hasSdataError
        } = this.state
        return(
            <div>
                <SearchBar></SearchBar>
                <div class="body">
                    <p>Screenshots of all visited Websites</p>
                        { 
                            hasSdataError ? 
                                <div class="alert alert-danger">
                                    {sdatamessage}
                                </div>
                                :
                                <div>
                                    <hr></hr>
                                    <br></br>
                                    {this.state.postData}
                                    <ReactPaginate
                                        previousLabel={"<<"}
                                        nextLabel={">>"}
                                        breakLabel={"..."}
                                        breakClassName={"break-me"}
                                        pageCount={this.state.maxPage}
                                        marginPagesDisplayed={1}
                                        pageRangeDisplayed={this.state.pageRange}
                                        onPageChange={this.handlePageClick}
                                        containerClassName={"pagination"}
                                        subContainerClassName={"pages pagination"}
                                        activeClassName={"active"} 
                                    /> 
                                </div>     
                        }            
                </div>
            </div>
        )
    }
}
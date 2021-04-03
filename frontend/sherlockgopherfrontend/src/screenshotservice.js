import React from 'react';
import SearchBar from './searchbar.js';
import ReactPaginate from 'react-paginate';
import axios from 'axios'
import "./screenshotservice.css";

export default class Screenshotservice extends React.Component {

  SCREENSHOT = "http://localhost:8081/graph/v1/getscreenshots"

    state = {
        data: [],
        currentPage: 0,
        maxPage: 0,
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
        axios.get(this.SCREENSHOT).then(res => {
            try {
                const data = res.data.map;
                const postData = data.map(pd => <React.Fragment>
                    <p>{pd.imageurl}</p>
                    <img src={pd.imagepath} alt=""/>
                </React.Fragment>)
                this.setState({
                    sdatamessage: "No error",
                    hasSdataError: false,
                    data: postData,
                    currentPage: res.data.currentpage,
                    maxPage: res.data.maxpage,
                    pageRange: res.data.pagerange,
                })
            } catch (error) {
                this.setState({
                    sdatamessage: "For some Reason an Error occured. Cannot process response.",
                    hasMetaError: true,
                })  
            }
        }).catch(error => {
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

    componentDidMount() {
        this.receivedData()
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
                                    {this.state.sliceData}
                                    <ReactPaginate
                                        previousLabel={"<<"}
                                        nextLabel={">>"}
                                        breakLabel={"..."}
                                        breakClassName={"break-me"}
                                        pageCount={this.state.maxPage}
                                        marginPagesDisplayed={0}
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
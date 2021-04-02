import React from 'react';
import SearchBar from './searchbar.js';
import ReactPaginate from 'react-paginate';
import axios from 'axios'
import "./screenshotservice.css";

export default class Screenshotservice extends React.Component {

  SCREENSHOT = "http://localhost:8081/getscreenshots/:page"

    state = {
        slicedata: [],
        sdatamessage: "",
        hasSdataError: false,
        offset: 0,
        data: [],
        perPage: 25,
        currentPage: 0,
        pageRage: 0
    }

    constructor(props) {
        super(props);
        this.handlePageClick = this
            .handlePageClick
            .bind(this);
    }

    receivedData() {
        axios.get(this.SCREENSHOT).then(res => {
            const data = res.data;
            try {
                const slice = data.slice(this.state.offset, this.state.offset + this.state.perPage)
                const tmpPageCount = Math.ceil(data.length / this.state.perPage);
                const tmpPageRange = 0;
                if (tmpPageCount > 4) {
                    tmpPageRange = 5;
                }
                const postData = slice.map(pd => <React.Fragment>
                    <p>{pd.imageurl}</p>
                    <img src={pd.imagepath} alt=""/>
                </React.Fragment>)
                this.setState({
                    sdatamessage: "No error",
                    hasSdataError: false,
                    pageCount: tmpPageCount,
                    slicedata: postData,
                    pageRage: tmpPageRange,
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
        const offset = selectedPage * this.state.perPage;

        this.setState({
            currentPage: selectedPage,
            offset: offset
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
                                        pageCount={this.state.pageCount}
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
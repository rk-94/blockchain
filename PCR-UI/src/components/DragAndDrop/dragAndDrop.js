import React,{Component} from 'react';
import './DragAndDropContainers.css';
import '../SearchOutputComponent/SearchOutputComponent';
import config from '../../config.json'
import img from '../../images/logo.png';
import logout from '../../images/logout.png';

class DragDrop extends Component{

    constructor(){
        super();
        this.state = {
            fieldsTaken : [],
            hashValue : '',
            finalKeys:[],
            fields : [],
            setFields :[],
            initialFields : [],
            isDataPresent : false,
            username :'',
            fhirUrl : '',
            jsonData : [],
            totalFhirResponse : {}
        }

   
        this.handleChange = this.handleChange.bind(this);
    }

    componentDidMount = () => {
        if(this.props.location.state === undefined){
            this.props.history.push({ 
                pathname : '/'
            });
        }else{
            
            if(this.props.location.state.setFields.length >0){
                this.setState({isDataPresent : true})
                this.setState({fields : this.props.location.state.initialFields})
                this.props.location.state.setFields.forEach((eachSet)=>{
                    this.props.location.state.initialFields.forEach((eachField)=>{
                        if(eachField.keyName === eachSet.keyName){
                            eachField.category = "used"
                        }
                    })
                })                
                this.setState({initialFields : this.props.location.state.initialFields, hashValue : this.props.location.state.hashValue});
            }else{
                            
                this.setState({username: this.props.location.state.username});
                
            }
            
        }    
    }

    handleChange(event) {
        this.setState({hashValue: event.target.value});
    }

    onSubmit = (event) =>{
        event.preventDefault();
        this.fetchHashValid();
    }

    fetchURL = (fhirUrl) => {
          
        let configParams = {
            method: 'GET',
            headers: {
              'Accept': 'application/json+fhir'
            },
          }
        fetch(fhirUrl, configParams)
            .then(response =>  response.json())
            .then((response) => {
                
                this.setState({ jsonData : response.entry , totalFhirResponse: response});
                var finalKeys =[];
                this.state.jsonData.forEach((eachEntry)=>{
                    finalKeys.push({"keyName" : eachEntry.resource.code.text, "category" : "unused", "value" :eachEntry.resource.text.div });
                    this.setState({fields: finalKeys, initialFields: finalKeys,isDataPresent : true});
                });
                

            });
        
        
    } 

    fetchData() {
        let configParams = {
            method: 'GET',
            headers: {
            'authorization': 'Bearer '+config.authToken,
            'content-Type': 'application/json'
            },
        }
      // eslint-disable-next-line
        fetch('http://'+config.hostIP+':'+config.port+'' + '/channels/'+config.channelName+'/chaincodes/'+config.chaincodeName+'?peer='+config.peerName+'&fcn=queryCustom&args=%5B%22%7B%5C%22selector%5C%22:%7B%5C%22_rev%5C%22:%5C%22'+this.state.hashValue+'%5C%22,%5C%22payerId%5C%22:%5C%22'+this.state.username+'%5C%22%7D%7D%22%5D', configParams)
            .then(response =>  response.json() )
            .then(response => {
                if (JSON.stringify(response) === '[]'){
                  
                } else {
                   
                    const fhirUrl = response.Record.fhirUrl;
                    this.fetchURL(fhirUrl);
                    this.setState({  fhirUrl: response.Record.fhirUrl })
                }
            } )
            .catch(err => console.log(err))
 
    }


    fetchHashValid(){
        let configParams = {
            method: 'GET',
            headers: {
              'authorization': 'Bearer '+config.authToken,
              'content-Type': 'application/json'
            },
          }
        // eslint-disable-next-line
        fetch('http://'+config.hostIP+':'+config.port+'' + '/channels/'+config.channelName+'/chaincodes/'+config.chaincodeName+'?peer='+config.peerName+'&fcn=isValid&args=%5B%22hash%22,%22'+this.state.hashValue+'%22%5D', configParams)
            .then(response =>  response.text() )
            .then(response => {
                if (response.length === 0 && response[0] !== 'E'){  
                    this.fetchData();
                    
                } else {
                    this.setState({isDataPresent : false})
                    alert("Invalid Hash key")
                }
            } )
            .catch(err => console.log(err))
        
    }


    onDragOver = (event) =>{
        event.preventDefault();
    }

    onDragStart = (ev,keyName) => {
        
        ev.dataTransfer.setData("keyName", keyName)
    }

    onDrop = (e,category) => {
        var fieldsList = [];
        if(category == "selectAll"){
            fieldsList = this.state.fields.filter((field)=> {               
                field.category = "used";               
                return field
            })
        } else if (category == "removeAll"){
            fieldsList = this.state.fields.filter((field)=> {
                field.category = "unused";              
                return field
            })
        }else{
            let keyName = e.dataTransfer.getData("keyName");
            fieldsList = this.state.fields.filter((field)=> {
                if(field.keyName === keyName){
                    field.category = category;
                }
                return field
            })
        }
        

        this.setState({fields: fieldsList});
    }

    onclick = () => {
        var fieldsTaken = [];
        var hashValue =  this.props.location.state.hashValue;
        this.state.fields.filter((field) => {
            if(field.category === "used")
                fieldsTaken.push({keyName : field.keyName, value : field.value})
        });
        this.setState({fieldsTaken : fieldsTaken});
        this.props.history.push({ 
            pathname : '/hashCode',state : {fieldsTaken :fieldsTaken, hashValue : this.state.hashValue, initialFields : this.state.fields}
        });
        
    }

    getResponseData = () =>{
        
        if(this.state.isDataPresent === true){
            var fields = {
                used : [],
                unused :[]
            }

            this.state.fields.forEach((field)=>{
                fields[field.category].push(
                    <div key = {field.keyName}>
                        <span className="circles" draggable = 'true' onDragStart = {(ev) => this.onDragStart(ev, field.keyName, field.value)} ></span>      
                        <span className = "text-style">{field.keyName}</span>          
                    </div>
                )
            })

            return(   <div className = "containers">  
                        <div className = "input-fields">
                            <span className = "text-draggable">FIELDS</span>
                            <div title ="Drag this for selection" className = "draggable" onDragOver = {(event) =>this.onDragOver(event)} onDrop = {(e) => this.onDrop(e, "unused")}>
                            {fields.unused}
                            </div>
                        </div>
                        <div className = "operations">
                                <span>
                                    <button className = "select-all" onClick ={(e) => this.onDrop(e, "selectAll")} > Select All </button>
                                    <button className = "remove-all" onClick ={(e) => this.onDrop(e, "removeAll")}> Remove All </button>
                                    <button type = "submit"  className = "searchBox"  onClick = {this.onclick} disabled = {!fields.used.length >0}>Search</button>
                                </span>
                        </div>
                        <div className = "output-fields">
                             <span className = "text-droppable"> OUTPUT FIELDS</span> 
                            <div title ="Drag this to unselect"className = "droppable" onDragOver = {(event) =>this.onDragOver(event)} onDrop = {(e) => this.onDrop(e, "used")}>
                                {fields.used}
                            </div>
                        </div>
                        {/* <div>                        
                            <button type = "submit"  className = "searchBox"  onClick = {this.onclick} disabled = {!fields.used.length >0}>Search</button>
                        </div> */}
                    </div>
            )
        }
    }

    goTologinPage = () => {
        this.props.history.push({ 
            pathname : '/'
        });
    }

    render(){
        

        return(
            <div className = "dragdrop">

                    <div className = "drag-drop-header">
                        <img src= {img} alt="Cerner Logo" style={{width: '200px' , float : 'left', height: '65px'}}/>
                        <span className = "header-text"> PAYER CHART REVIEW </span>
                        <img src= {logout} alt="Logout" style = {{float : 'right', marginTop :'1%', cursor : 'pointer'}} onClick = {this.goTologinPage}></img>
                    </div>

                    <form className = "hash-placeholder">
                         <div>
                            <input type = "text" className = "hashBox" value = {this.state.hashValue} onChange={this.handleChange}  placeholder = "Enter Hash" required/>                 
                                    
                            <button type = "submit"  className = "submitButton" onClick = {(event) => this.onSubmit(event)} disabled = {!this.state.hashValue} >Submit</button>
                                    
                        </div>
                    </form>
                    <div>
                        {this.getResponseData()}
                    </div>
                    

            </div>
            

        );
    }
}

export default  DragDrop;
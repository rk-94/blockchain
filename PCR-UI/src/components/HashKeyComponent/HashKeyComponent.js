import React,{Component} from 'react';
import './HashKeyComponent.css';
import img from '../../images/logo.png';
import home from '../../images/home.png'
import logout from '../../images/logout.png';
class HaskKeyComponent extends Component{

    state = {
        haskKey : '',
        fieldsValues : [],
        displayFields : false,
        selectedFieldValue : '',
        isSelectedField : false,
        selectedField :''
    }

    componentDidMount = () => {
        if(this.props.location.state === undefined){
            this.props.history.push({ 
                pathname : '/'
            });
        }else{
            var fieldsUsed = this.props.location.state.fieldsTaken;
            this.setState({fieldsValues: fieldsUsed, hashKey : this.props.location.state.hashValue});
            
        }
    }

    getData = () => {
        this.setState({displayFields : true});

    }
    
    getFieldValue = (event, fieldName) => {
        
        this.state.fieldsValues.forEach((field)=>{
            if(field.keyName === fieldName){
                
                this.setState({selectedFieldValue : field.value, isSelectedField : true, selectedField : field.keyName});

            }
        })
    }

    outputValue = () => {
        
        if(this.state.isSelectedField === true){
            return(
                <div>
                    <div className = "value-conatiner">
                        <span className = "select-field" >{this.state.selectedField}</span>
                        <div className = "text-placement" dangerouslySetInnerHTML = {{__html: this.state.selectedFieldValue}}/>
                    </div>
                </div>
                
            )
        }
    }

    goHome = () =>{
        this.props.history.push({ pathname: "/responseData", state :{setFields :this.state.fieldsValues, initialFields : this.props.location.state.initialFields, hashValue : this.props.location.state.hashValue}});
    }

    goTologinPage = () => {
        this.props.history.push({ 
            pathname : '/'
        });
    }
    
    render(){
        var fieldsValues, selectedValues = [];
        if(this.state.displayFields === true){
            fieldsValues = this.state.fieldsValues.map((field)=>{
                return( <div key = {field.keyName}>
                            <span className="smallCircle" onClick = {(event) => this.getFieldValue(event, field.keyName)}/> 
                            <span className = "text-placement">{field.keyName}</span>          
                             
                        </div>   
                    )            
            });
        }

        selectedValues = this.state.fieldsValues;
        
        return(
            
            <div className = "hash-main">                 
                {/* <div>
                    <span className="biggerCircle" onClick = {this.getData}>{this.state.hashKey}</span>
                </div>
                <div className = "circles-conatiner">
                    {fieldsValues}
                </div>
                <div>
                    {this.outputValue()}
                </div>
                 */}

                

                <div className = "hash-header">
                        <img src= {img} alt="Cerner Logo" style={{width: '200px' , float : 'left', height: '65px'}}/>
                        <img src= {home} alt="Home icon" style = {{ width: '50px' ,float : 'left',height: '50px',cursor :'pointer', marginTop :'1%'}} onClick = {this.goHome}/>
                        <span className = "header-text"> PAYER CHART REVIEW </span>
                        <img src= {logout} alt="Logout" style = {{float : 'right', marginTop :'1%', cursor : 'pointer'}} onClick = {this.goTologinPage}></img>
                </div>

                 <div>
                     <span className = "text-styling"> ADDITIONAL CLINICAL DATA </span>
                 </div>
                <div className = "wrapper">
                    
                    {selectedValues.map((eachValue,index)=>{
                        return(
                            <div key = {eachValue.keyName} className ="col-sm-6">
                                <div className = "rectangle-holders" >
                                    <span className = "select-field" >{eachValue.keyName}</span>
                                    <div className = "text-placement" dangerouslySetInnerHTML = {{__html: eachValue.value}}/>
                                </div>
                            </div>  
                        )
                    }) }
                        
                </div>
            </div>
        )
    }
}

export default HaskKeyComponent;
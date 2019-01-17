import React, { Component } from 'react';
import {
  Route,
  BrowserRouter,
  Switch
} from "react-router-dom";
import URLComponent from './components/URLComponent/URLComponent';
import DragDrop from './components/DragAndDrop/dragAndDrop';
import HashKeyComponent from './components/HashKeyComponent/HashKeyComponent';
import LoginComponent from './components/LoginComponent/LoginComponent';

class App extends Component {

  render() {
    
    return (
      <BrowserRouter>
        <div>
          <Switch>
            <Route exact  path="/" component={LoginComponent}/> 
            <Route exact  path="/chaincode" component={URLComponent}/> 
            <Route  exact path="/responseData"  component={DragDrop}/>
            <Route  exact path="/hashCode"  component={HashKeyComponent}/>
          </Switch>
          
        </div>
      </BrowserRouter>
    );

  }
}

export default App;

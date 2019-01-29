import React, { Component } from 'react';
import { BrowserRouter, Switch, Route } from 'react-router-dom';
import Button from 'terra-button/lib/Button';
import Card from 'terra-card/lib/Card';
import ApplicationMenuName from 'terra-application-name/lib/ApplicationMenuName';
import Image from 'terra-image';
import img from './cerner.png';

import Input from 'terra-form-input';


import Main from './main';

import ApplicationHeaderLayout from 'terra-application-header-layout';
import ApplicationHeaderName from 'terra-application-name/lib/ApplicationHeaderName';
import Placeholder from 'terra-application-name/lib/terra-dev-site/doc/common/Placeholder';

class App extends Component {
  render() {
    const size = { height: '50px', width: '100%' };
    const ApplicationHeaderNameStandard = () => (
      <ApplicationHeaderLayout
        style={size}
        logo={(
          <div style={{ backgroundColor: 'green' }}>
            <ApplicationHeaderName
              title="App-Name"
              accessory={<Image variant="rounded" src="https://github.com/cerner/terra-core/raw/master/terra.png" height="26px" width="26px" isFluid />}
            />
          </div>
        )}
        extensions={<Placeholder text="Extensions" width="100px" />}
        navigation={<Placeholder text="Content" />}
        utilities={<Placeholder text="Utiltities" width="100px" />}
      />
    );
    return (
      <div style={{ border: '1px solid lightGray', backgroundColor: '#2481ca', width: '100%', height: '50px', position: 'relative', }} >
        <ApplicationMenuName title="Payer Chart Review" accessory={<Image src={img} height="80px" width="80px" isFluid />} />
      </div>
    );
  }
}

export default App;

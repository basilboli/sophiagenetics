import React, { Component } from "react";
import './App.css';
import { BrowserRouter as Router, Link, Route } from 'react-router-dom';
import HistogramChart from './HistogramChart';

class App extends Component {

  state = {
    clients: null
  }

  // fetch all necessary data
  componentDidMount() {
    fetch('http://localhost:8080/clients')
      .then(res => res.json())
      .then(clients => {
        console.log("Clients:", clients)
        this.setState({ clients })
      })
  }

  render() {
    const { clients } = this.state
    return (
      <Router>
        <Root>
          <Sidebar>
            <SidebarItem key='home'>
              <Link to='/'>
                About
          </Link>
            </SidebarItem>
            <SidebarItem key='clients'>
              <Link to='/clients'>
                Choose a client
          </Link>
            </SidebarItem>
          </Sidebar>
          <Main>
            <Route exact={true} path="/" render={() => (
              <div>
                <h1>Task: Data Visualisation</h1>
                
                Made by <a href="https://basilboli.keybase.pub/">Vasyl Vaskul</a><br/>
                <a href="https://github.com/basiboli/sophiagenetics">Link to the code</a>
              </div>
            )} />
            <Route exact={true} path="/clients" render={() => (
              <div className="App-links">
                {clients ? (
                  clients.map(client => (
                    <Link className="App-link" to={`/g/${client.salesforceId}`}>
                      {client.salesforceId || '[no description]'}
                    </Link>
                  ))
                ) : (
                    <div> Loading ... </div>
                  )}
              </div>
            )} />
            {clients && (
              <Route path="/g/:salesforceId" render={({ match }) => (
                <Client client={clients.find(c => c.salesforceId == match.params.salesforceId)} />
              )} />
            )}
          </Main>
        </Root>
      </Router>
    );
  }
}


const Client = ({ client }) => {
  let allUsages = [...client.actualUsage, ...client.predictedUsage];
  // we calculate max of actualUsage and predictedUsage to make graph y axis consistent
  let max = Math.max(...allUsages.map(item => item));

  return (
    <div className="App">
      <div><span class="App-title-header">Saleforce id:</span> {client.salesforceId || 'No description'}</div>
      <div><span class="App-title-header">Owner:</span>{client.owner || 'No description'}</div>
      <div><span class="App-title-header">Country:</span>{client.country || 'No description'}</div>
      <div><span class="App-title-header">Manager:</span>{client.manager || 'No description'}</div>
      <HistogramChart data={client.actualUsage} max={max} title="Actual usage" />
      <HistogramChart data={client.predictedUsage} max={max} title="Predicted usage" />
    </div>
  );
}

const Root = (props) => (
  <div style={{
    display: 'flex'
  }} {...props} />
)

const Sidebar = (props) => (
  <div style={{
    width: '15vw',
    height: '100vh',
    overflow: 'auto',
  }} {...props} />
)

const SidebarItem = (props) => (
  <div style={{
    whiteSpace: 'nowrap',
    textOverflow: 'elipsis',
    overflow: 'hidden',
    padding: '5px 10px',
  }} {...props} />
)

const Main = (props) => (
  <div style={{
    flex: 1,
    height: '200vh',
    overflow: 'auto',
  }}>
    <div style={{ padding: '0px' }} {...props} />
  </div>
)

export default App;
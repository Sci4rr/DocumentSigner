import React from 'react';
import { Provider } from 'react-redux';
import store from './store';
import DocumentUploader from './components/DocumentUploader';
import DocumentViewer from './components/DocumentViewer';
import SignaturePad from './components/SignaturePad';
import UserDashboard from './components/UserDashboard';
import './App.css';

class App extends React.Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
  }

  render() {
    return (
      <Provider store={store}>
        <div className="App">
          <header className="App-header">
            <h1>DocuSigner</h1>
          </header>
          <main>
            <UserDashboard />
            <DocumentUploader />
            <DocumentViewer />
            <SignaturePad />
          </main>
        </div>
      </Provider>
    );
  }
}

export default App;
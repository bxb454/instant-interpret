import React from 'react';
import logo from './logo.svg';
import './App.css';
import Header from './Header';
import ChatContext from './ChatContext';
import ChatMsg from './ChatMsg';
import ChatHistory from './ChatHistory';


function App() {
  return (
    <div className="App">
    <Header></Header>
    <ChatMsg/>
    <ChatHistory/>
    </div>
  );
}

export default App;

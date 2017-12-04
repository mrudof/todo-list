import React from 'react';
import ReactDOM from 'react-dom';
import Home from './components/home';
// import {Provider} from 'react-redux';
import store from './redux/store';
ReactDOM.render(
    <Home />,
  document.getElementById('root')
);
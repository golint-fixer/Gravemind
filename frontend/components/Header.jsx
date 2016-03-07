import React from 'react';
import { Link } from 'react-router';

export default class Header extends React.Component {
  render() {
    return (
      <div className="header">
        <div className="container">
          <Link className="brand" to="/">
            <img src={require('assets/images/logo.svg')} width="40" height="40" alt="Gravemind" />
            <span>Gravemind</span>
          </Link>
          <div className="nav">
            <Link to="/login">Log in</Link>
          </div>
        </div>
      </div>
    );
  }
}

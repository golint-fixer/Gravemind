import React from 'react';
import { Link } from 'react-router';

export default class Footer extends React.Component {
  render() {
    return (
      <div className="footer">
        <div className="container">
          <Link to="/">Home</Link>
          <span>·</span>
          <Link to="/about">About</Link>
          <span>·</span>
          <Link to="/contact">Contact</Link>
          <span>·</span>
          <Link to="/privacy">Privacy</Link>
        </div>
      </div>
    );
  }
}

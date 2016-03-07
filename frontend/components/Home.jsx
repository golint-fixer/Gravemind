import React from 'react';
import { Link } from 'react-router';

export default class Home extends React.Component {
  render() {
    return (
      <div className="home">
        <div className="overlay">
          <div className="container">
            <h1 className="title">A free, highly available bot of your own design</h1>
            <div className="row">
              <div className="developer">
                <h1>Developers</h1>
                <p>
                  Easily contribute plugins for any broadcaster to use,
                  without worrying about scaling or networking.
                </p>
                <Link to="/contribute" className="button">Contribute</Link>
              </div>
              <div className="spacer"></div>
              <div className="broadcaster">
                <h1>Broadcasters</h1>
                <p>
                  Select exactly the functionality you want from a pool of community created features.
                </p>
                <Link to="/setup" className="button">Set Up</Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

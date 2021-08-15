import React, { useState } from 'react';

const GettingStarted = (props) => {
  const [listener, setListener] = useState("");

  const handleSubmit = (evt) => {
    evt.preventDefault();
    props.history.push("/_/"+listener);
  }

  return (
    <>
      <div className="hero">
        Demystify Webhooks With <i>Whdbg</i>.
      </div>
      <div className="panel">
        <form onSubmit={handleSubmit}>
          <label>
          <div className="row">
            <input
              className="listener"
              type="text"
              style={{
                width:280,
                textAlign: 'right',
                paddingLeft: 12,
              }}
              placeholder="Name Your Listener"
              value={listener}
              onChange={e => setListener(e.target.value)}
            />
            <input type="text"
            className="listener"
            style={{
              width:250,
              textAlign: 'left',
              color:'#afafaf',
              paddingRight: 12,
            }}
            disabled="disabled"
            value=".whdbg.dev" />
            </div>
          </label><br />
          <input type="submit" value="Go" />
        </form>
      </div>
    </>
  );
}

export default GettingStarted;

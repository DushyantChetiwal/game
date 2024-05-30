import '../styles/Start.css'
import React, { useContext } from 'react';
import { SocketContext } from '../SocketContext';

const Start = ({onStart}) => {
  const { socket, initializeSocket } = useContext(SocketContext);

  const enterMatchmaking = () => {
    initializeSocket()
    onStart()
  }
    return (
      <div className="Start">
        <div>Rock-Paper-Scissor</div>
        <div onClick={enterMatchmaking}>Start Game</div>
      </div>
    );
  }
  
  export default Start;
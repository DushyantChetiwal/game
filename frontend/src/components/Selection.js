import '../styles/Selection.css'
import React, { useEffect, useState, useContext } from 'react';
import { SocketContext } from '../SocketContext';

const Selection  = ({onChoice}) => {
  const { socket, initializeSocket,closeSocket } = useContext(SocketContext);
  

  const sendMessage = (message) => {
    if (socket) {
      socket.send( message );
      onChoice();
    }
  };

    return (
      <div className="Selection">
        <button onClick = {() => sendMessage("rock")}>Rock</button>
        <button onClick = {() => sendMessage("paper")}>Paper</button>
        <button onClick = {() => sendMessage("scissor")}>Scissor</button>
      </div>
    );
  }
  
  export default Selection;
import '../styles/Matchmaking.css'
import React, { useEffect, useState, useContext } from 'react';
import { SocketContext } from '../SocketContext';

const Matchmaking  = ({onMatch}) => {
  const { socket, initializeSocket, closeSocket } = useContext(SocketContext);
  
  useEffect(() => {
    if (socket) {
      socket.onmessage = (event) => {
        const message = event.data;
        if(message==="Please enter your choice rock, paper or scissor:"){
          onMatch()
        }
        console.log('Received message:', message);
      };

      socket.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
    }
  }, [socket]);

    return (
      <div className="Matchmaking">
        Finding a match...
      </div>
    );
  }
  
  export default Matchmaking;
import '../styles/Result.css'
import React, { useEffect, useState, useContext } from 'react';
import { SocketContext } from '../SocketContext';

const Result = ({onRestart}) =>{
  const [Outcome, setOutcome] = useState("")
  const { socket, initializeSocket, closeSocket } = useContext(SocketContext);

  useEffect(() => {
    if (socket) {
      socket.onmessage = (event) => {
        const message = event.data;
        setOutcome(message)
        console.log('Received message:', message);
        closeSocket()
      };

      socket.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
    }
  }, [socket, closeSocket]);

    return (
      <div className="Result">
        {Outcome != "" && 
        <div>
          {Outcome}
          <button onClick={onRestart}>Restart</button>
        </div>
        } 

        
      </div>
    );
  }
  
  export default Result;
import React, { createContext, useState, useEffect } from 'react';
import SocketService from './socketService';

const SocketContext = createContext();

const SocketProvider = ({ children }) => {
  const [socket, setSocket] = useState(null);

  const initializeSocket = () => {
    if (!socket) {
      const socketInstance = SocketService.getSocket();
      setSocket(socketInstance);

      socketInstance.onclose = () => {
        setSocket(null);
      };
    }
  };

  const closeSocket = () => {
    if (socket) {
      SocketService.closeSocket();
      setSocket(null);
    }
  };

  return (
    <SocketContext.Provider value={{ socket, initializeSocket, closeSocket }}>
      {children}
    </SocketContext.Provider>
  );
};

export { SocketContext, SocketProvider };

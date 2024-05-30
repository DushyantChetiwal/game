import './App.css';
import React, { useState } from 'react';
import Start from './components/Start'
import Matchmaking from './components/Matchmaking';
import Selection from './components/Selection';
import Result from './components/Result';



function App() {
  const [count, setCount] = useState(0);

  const increment = () => {
    setCount((count + 1)%4);
  };

  return (
    <div className='App'>
      {count===0 && <div><Start onStart = {increment}/></div>}
      {count===1 && <div><Matchmaking onMatch = {increment}/></div>}
      {count===2 && <div><Selection onChoice = {increment}/></div>}
      {count===3 && <div><Result onRestart = {increment}/></div>}
    </div>
  );
}

export default App;

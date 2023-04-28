import React, { useState } from 'react';
import ComplexNavbar from './NavBar';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = (event) => {
    event.preventDefault();
      const requestData = {
        username: username,
        password: password
      };
      console.log(requestData);
    
      const options = {
        method: 'POST',
        body: JSON.stringify(requestData) // Convertir a cadena JSON
      };
  
    
      fetch('http://52.91.77.62/login', options)
        .then(response => response.json())
        .catch(err => console.error(err));
  };

  return (
    <form onSubmit={handleSubmit}>
      <ComplexNavbar username={username} password={password} />
      <label>
        Username:
        <input
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
      </label>
      <br />
      <label>
        Password:
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </label>
      <br />
      <button type="submit">Submit</button>
    </form>
  );
};

export default Login;
